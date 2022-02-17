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
	ds, err := newDemoServer(ctx)
	if err != nil {
		return err
	}

	pb.RegisterDemoServer(s, ds)

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()

	slog.Debug(`grpc starting to serve...`)
	return s.Serve(l)
}
