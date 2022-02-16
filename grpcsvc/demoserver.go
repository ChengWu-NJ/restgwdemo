package grpcsvc

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"google.golang.org/protobuf/types/known/emptypb"
	"ire.com/restgwdemo/pb"
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

	defer func() {
		<-ctx.Done()

		db.Close()
	}()

	if err := _creatTables(db); err != nil {
		return nil, err
	}

	return &demoServer{
		ctx: ctx,
		DB:  db,
	}, nil
}

func _creatTables(db *pg.DB) error {
	if err := db.CreateTable(&pb.StorNode{}, &orm.CreateTableOptions{IfNotExists: true}); err != nil {
		return err
	}
	//TODO: create others ...

	return nil
}

func (s *demoServer) EnableStorNode(ctx context.Context, req *pb.EnableStorNodeRequest) (
	*emptypb.Empty, error) {

	//name is unique, so there is just one
	if r, err := s.DB.QueryOne(&pb.StorNode{}, `select * from stor_nodes where name = ?`,
		req.StorNode.Name); err != nil {

		//there is none, do insert
		if r.RowsAffected() == -1 {
			slog.Debug(`do insert...`)

			//2. oterwise create(/insert) it
			return nil, s.DB.Insert(req.StorNode)
		}

		//there is multiple records
		return nil, err
	}

	//1. update it if exists
	return nil, s.DB.Update(req.StorNode)
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
