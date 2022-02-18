package grpcsvc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"ire.com/pg"
	"ire.com/pg/orm"
	"ire.com/restgwdemo/pb"
	"ire.com/restgwdemo/pgutils"
	"ire.com/slog"
)

// Implements of EchoServiceServer

type demoServer struct {
	ctx context.Context
	pb.DemoServer

	DB *pg.DB
}

func newDemoServer(ctx context.Context) (pb.DemoServer, error) {
	db := pg.Connect(&pg.Options{
		User:                  "postgres",
		Password:              "4y7sV96vA9wv46VR",
		Database:              "postgres",
		Addr:                  `localhost:5432`,
		RetryStatementTimeout: true,
		MaxRetries:            4,
		MinRetryBackoff:       250 * time.Millisecond,
	})

	go func() {
		defer db.Close()

		<-ctx.Done()
	}()

	if err := _createSchema(db); err != nil {
		return nil, err
	}

	return &demoServer{
		ctx: ctx,
		DB:  db,
	}, nil
}

func _createSchema(db *pg.DB) error {
	models := []interface{}{
		(*pb.StorNode)(nil),
		//TODO: add others...
	}

	for _, model := range models {
		opts := &orm.CreateTableOptions{
			IfNotExists: true,
		}

		q := db.Model(model)

		if sqlstr, err := pgutils.CreateTableQueryString(q, opts); err != nil {
			return err
		} else {
			slog.Debugf(`got sql:[%s]`, sqlstr)
		}

		err := q.CreateTable(opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *demoServer) UpsertStorNode(ctx context.Context, req *pb.UpsertStorNodeRequest) (
	*emptypb.Empty, error) {

	//name is unique, so there is just one
	if req == nil || req.StorNode == nil || req.StorNode.Name == "" {
		return nil, fmt.Errorf(`got an empty request or name`)
	}

	err := pgutils.Upsert(s.DB, `name`, req.StorNode)

	return nil, err
}

func (s *demoServer) Healthz(ctx context.Context, emp *emptypb.Empty) (*pb.HealthzResponse, error) {
	return &pb.HealthzResponse{
		State: "ok",
		Htime: timestamppb.Now(),
	}, nil
}

/*
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
*/
