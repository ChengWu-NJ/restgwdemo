package grpcsvc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"ire.com/restgwdemo/pb"
	"ire.com/slog"
)

// Run starts the demo gRPC service.
// "network" and "address" are passed to net.Listen.
func Run(ctx context.Context, network, address string) error {
	slog.Infof("grpc starting listening at %s", address)

	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			slog.Errorf("Failed to close %s %s: %v", network, address, err)
		}
	}()

	s := grpc.NewServer()
	pb.RegisterDemoServer(s, newDemoServer(ctx))

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()
	return s.Serve(l)
}

/*
func RunInProcessGateway(ctx context.Context, addr string, opts ...runtime.ServeMuxOption) error {
	mux := runtime.NewServeMux(opts...)

	if err := pb.RegisterDemoHandlerServer(ctx, mux, newDemoServer(ctx)); err != nil {
		return err
	}

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		slog.Infof("Shutting down the http gateway server")
		if err := s.Shutdown(context.Background()); err != nil {
			slog.Errorf("Failed to shutdown http gateway server: %v", err)
		}
	}()

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		slog.Errorf("Failed to listen and serve: %v", err)
		return err
	}
	return nil

}
*/
