package main

import (
	"context"
	"fmt"
	"os"

	"ire.com/restgwdemo/gateway"
	"ire.com/restgwdemo/grpcsvc"
	"ire.com/slog"
)

const (
	GRPCADDR = "localhost:9090"
	HTTPADDR = ":8080"
)

func main() {
	ctx := context.Background()

	chErr := runServers(ctx)

	select {
	case err := <-chErr:
		slog.Error(err)
		os.Exit(1)

	case <-ctx.Done():
	}
}

func runServers(ctx context.Context) <-chan error {
	ch := make(chan error, 2)

	go func() {
		if err := grpcsvc.Run(ctx, "tcp", GRPCADDR); err != nil {
			ch <- fmt.Errorf("cannot run grpc service: %v", err)
		}
	}()

	go func() {
		opts := gateway.Options{
			Addr: HTTPADDR,
			GRPCServer: gateway.Endpoint{
				Network: "tcp",
				Addr:    GRPCADDR,
			},
		}
		if err := gateway.Run(ctx, opts); err != nil {
			ch <- fmt.Errorf("cannot run gateway service: %v", err)
		}
	}()

	return ch
}
