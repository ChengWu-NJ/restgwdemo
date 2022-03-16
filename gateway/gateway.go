package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc/credentials/insecure"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"ire.com/restgwdemo/pb"

	"ire.com/slog"
)

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

// TODO ...
// https://github.com/shaj13/go-guardian/blob/master/_examples/token/main.go

var (


	/*
	basicAuthMiddleware = func(handle http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, pass, ok := req.BasicAuth()
			if !ok || !checkBasicAuth(user, pass) {
				http.Error(w, "Unauthorized.", 401)
				return
			}
			handle.ServeHTTP(w, req)
		})
	}*/
)

// Run starts a HTTP server and blocks while running if successful.
// The server will be shutdown when "ctx" is canceled.
func Run(ctx context.Context, opts Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := dial(ctx, opts.GRPCServer.Network, opts.GRPCServer.Addr)
	if err != nil {
		return err
	}
	slog.Debugf(`after dial, conn.GetState() == %s`, conn.GetState())

	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			slog.Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := http.NewServeMux()

	gw, err := newGateway(ctx, conn, opts.Mux)
	if err != nil {
		return err
	}

	//if isBasicAuth {
	//	mux.Handle("/", basicAuthMiddleware(gw))
	//} else {
		mux.Handle("/", gw)
	//}

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
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

func checkBasicAuth(user, pass string) bool {
	//TODO ...
	return true
}
