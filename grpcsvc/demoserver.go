package grpcsvc

import (
	"context"
	"fmt"
	"io"

	"ire.com/restgwdemo/pb"
	"ire.com/slog"
)

// Implements of EchoServiceServer

type demoServer struct {
	ctx context.Context
	pb.DemoServer
}

func newDemoServer(ctx context.Context) pb.DemoServer {
	return &demoServer{
		ctx: ctx,
	}
}

func (s *demoServer) UnaryDemo(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	if req == nil || req.InMsg == "" {
		return nil, fmt.Errorf(`got an requset with empty message`)
	}

	return &pb.Response{OutMsg: fmt.Sprintf(`I got your message:[%s]`, req.InMsg)}, nil
}

func (s *demoServer) BulkUpload(stream pb.Demo_BulkUploadServer) error {
	ctx := stream.Context()

	msgs := make([]string, 0)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-s.ctx.Done():
			return s.ctx.Err()

		default:
		}

		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if req == nil || req.InMsg == "" {
			continue
		}

		msgs = append(msgs, req.InMsg)
	}

	slog.Infof(`BulkUpload got %v`, msgs)

	return nil
}

func (s *demoServer) BulkDownload(req *pb.Request, stream pb.Demo_BulkDownloadServer) error {
	if req == nil || req.InMsg == "" {
		return fmt.Errorf(`got an requset with empty message`)
	}

	ctx := stream.Context()
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-s.ctx.Done():
			return s.ctx.Err()

		default:
		}

		if err := stream.Send(&pb.Response{
			OutMsg: fmt.Sprintf(`No.%d response to %s`, i, req.InMsg),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *demoServer) DoubleStream(stream pb.Demo_DoubleStreamServer) error {
	ctx := stream.Context()

	i := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-s.ctx.Done():
			return s.ctx.Err()

		default:
		}

		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if req == nil || req.InMsg == "" {
			continue
		}

		_ = stream.Send(&pb.Response{OutMsg: fmt.Sprintf(`I got your No.%d message:[%s]`,
			i, req.InMsg)})
		i++
	}

	return nil
}
