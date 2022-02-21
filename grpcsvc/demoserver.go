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
		(*pb.VG)(nil),
		//TODO: add others...
	}

	for _, model := range models {
		opts := &orm.CreateTableOptions{
			IfNotExists: true,
			FKConstraints: true,
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

	if req == nil || req.StorNode == nil || req.StorNode.Name == "" {
		return &emptypb.Empty{}, fmt.Errorf(`got an empty request or name`)
	}

	err := pgutils.Upsert(s.DB, `name`, req.StorNode)

	return &emptypb.Empty{}, err
}

func (s *demoServer) GetStorNodeByName(ctx context.Context, req *pb.GetStorNodeByNameRequest) (
	*pb.GetStorNodeByNameResponse, error) {

	if req == nil || req.Name == "" {
		return &pb.GetStorNodeByNameResponse{}, fmt.Errorf(`got an empty request or name`)
	}

	storNode := &pb.StorNode{
		Name: req.Name,
	}

	err := pgutils.SelectOneByKey(s.DB, storNode)

	return &pb.GetStorNodeByNameResponse{StorNode: storNode}, err
}

func (s *demoServer) DelStorNodeByName(ctx context.Context, req *pb.DelStorNodeByNameRequest) (
	*emptypb.Empty, error) {

	if req == nil || req.Name == "" {
		return &emptypb.Empty{}, fmt.Errorf(`got an empty request or name`)
	}	

	storNode := &pb.StorNode{
		Name: req.Name,
	}

	err := pgutils.DeleteOneByKey(s.DB, storNode)

	return  &emptypb.Empty{}, err

}

func (s *demoServer) Healthz(ctx context.Context, emp *emptypb.Empty) (*pb.HealthzResponse, error) {
	return &pb.HealthzResponse{
		State: "ok",
		Htime: timestamppb.Now(),
	}, nil
}

func (s *demoServer) UpsertVG(ctx context.Context, req *pb.UpsertVGRequest) (*emptypb.Empty, error) {
	if req == nil || req.Vg == nil || req.Vg.VgId == "" {
		return &emptypb.Empty{}, fmt.Errorf(`got an empty request or vgId`)
	}

	if req.Vg.StorNodeName == "" {
		return &emptypb.Empty{}, fmt.Errorf(`got an empty storNodeName`)
	}

	err := pgutils.Upsert(s.DB, `vg_id`, req.Vg)

	return &emptypb.Empty{}, err
}


func (s *demoServer) GetVGById(ctx context.Context, req *pb.GetVGByIdRequest) (
	*pb.GetVGByIdResponse, error) {

	if req == nil || req.VgId == "" {
		return &pb.GetVGByIdResponse{}, fmt.Errorf(`got an empty request or vgId`)
	}

	vg := &pb.VG{
		VgId: req.VgId,
	}

	err := pgutils.SelectOneByKey(s.DB, vg)

	return &pb.GetVGByIdResponse{Vg: vg}, err	
}

func (s *demoServer) DelVGById(ctx context.Context, req *pb.DelVGByIdRequest) (*emptypb.Empty, error) {

	if req == nil || req.VgId == "" {
		return &emptypb.Empty{}, fmt.Errorf(`got an empty request or vgId`)
	}

	vg := &pb.VG{
		VgId: req.VgId,
	}

	err := pgutils.DeleteOneByKey(s.DB, vg)

	return &emptypb.Empty{}, err
}