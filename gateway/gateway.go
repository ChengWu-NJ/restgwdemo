package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/felixge/httpsnoop"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"ire.com/restgwdemo/pb"

	"ire.com/slog"
)

const (
	EXEMPT_SECONDS_AFTER_CHECK_PASSWORD = 60
)

var rightAuthenticationStrings = make(map[string]time.Time)

//please go to "THE CUSTOMIZED PART" to registry grpc service

// Endpoint describes a gRPC endpoint
type Endpoint struct {
	Network, Addr string
}

// Options is a set of options to be passed to Run
type Options struct {
	// Addr is the address to listen
	Addr string

	// GRPCServer defines an endpoint of a gRPC service
	GRPCServer Endpoint

	// Mux is a list of options to be passed to the gRPC-Gateway multiplexer
	Mux []gwruntime.ServeMuxOption
}

// Run starts a HTTP server and blocks while running if successful.
// The server will be shutdown when "ctx" is canceled.
func Run(ctx context.Context, opts Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := dial(ctx, opts.GRPCServer.Network, opts.GRPCServer.Addr)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			slog.Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzServer(conn))

	gw, err := newGateway(ctx, conn, opts.Mux)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	s := &http.Server{
		Addr: opts.Addr,
		// add authentication and logger for restful access
		Handler: withAuthAndLogger(mux),
	}
	go func() {
		<-ctx.Done()
		slog.Infof("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			slog.Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	slog.Infof("http gateway starting listening at %s", opts.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		slog.Errorf("Failed to listen and serve: %v", err)
		return err
	}
	return nil
}

// newGateway returns a new gateway server which translates HTTP into gRPC.
func newGateway(ctx context.Context, conn *grpc.ClientConn, opts []gwruntime.ServeMuxOption) (http.Handler, error) {
	// add Authentication ...
	opts = append(opts,
		gwruntime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),
	)

	mux := gwruntime.NewServeMux(opts...)

	for _, f := range []func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error{
		pb.RegisterDemoHandler,
		//THE CUSTOMIZED PART
		//register other handler if any
	} {
		if err := f(ctx, mux, conn); err != nil {
			return nil, err
		}
	}
	return mux, nil
}

func dial(ctx context.Context, network, addr string) (*grpc.ClientConn, error) {
	switch network {
	case "tcp":
		return dialTCP(ctx, addr)
	case "unix":
		return dialUnix(ctx, addr)
	default:
		return nil, fmt.Errorf("unsupported network type %q", network)
	}
}

// dialTCP creates a client connection via TCP.
// "addr" must be a valid TCP address with a port number.
func dialTCP(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// dialUnix creates a client connection via a unix domain socket.
// "addr" must be a valid path to the socket.
func dialUnix(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	d := func(ctx context.Context, addr string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "unix", addr)
	}
	return grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(d))
}

// healthzServer returns a simple health handler which returns ok.
func healthzServer(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}

func withAuthAndLogger(handler http.Handler) http.Handler {
	// the create a handler
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if !canExemptCheckPassword(request) {
			username, password, ok := request.BasicAuth()
			if !ok || !checkUsernameAndPassword(username, password) {
				writer.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
				writer.WriteHeader(401)
				_, _ = writer.Write([]byte("Unauthorised.\n"))
				return
			}

			rightAuthenticationStrings[request.Header.Get("Authorization")+string(net.ParseIP(request.RemoteAddr))] = time.Now()
		}

		// pass the handler to httpsnoop to get http status and latency
		m := httpsnoop.CaptureMetrics(handler, writer, request)

		// printing exracted data
		slog.Infof("from:[%s], http[%d]-- %s -- %s\n", request.RemoteAddr, m.Code, m.Duration, request.URL.Path)
	})
}

// canExemptCheckPassword is for reducing to invoke checkUsernameAndPassword
// which will actually invoke a backend authentication service
func canExemptCheckPassword(request *http.Request) bool {
	authStr := request.Header.Get("Authorization") + string(net.ParseIP(request.RemoteAddr))
	checkTime, ok := rightAuthenticationStrings[authStr]
	if !ok {
		return false
	}

	if checkTime.Add(EXEMPT_SECONDS_AFTER_CHECK_PASSWORD * time.Second).Before(time.Now()) {
		return false
	}

	return true
}

func checkUsernameAndPassword(username, password string) bool {
	// TODO...
	slog.Info("invoke some authentication service...")
	return username == "testuser" && password == "testpass"
}
