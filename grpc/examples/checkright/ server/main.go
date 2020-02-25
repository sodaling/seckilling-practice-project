package main

import (
	"context"
	"log"
	"net"
	"seckilling-practice-project/common"
	pb "seckilling-practice-project/grpc"
)

var port = ":50052"
var accessControl *common.AccessControl

type server struct {
	pb.UnimplementedCheckRightServiceServer
}

//实现CheckOneServiceServer
func (s *server) CheckRight(ctx context.Context, req *pb.Uid) (*pb.IsOk, error) {
	uid := req.GetValue()
	isOk := accessControl.GetDataFromMap(uid)
	return &pb.IsOk{Value: isOk}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	defer lis.Close()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	accessControl = common.GetAccessControl()
	s := common.GetGrpcServer()
	pb.RegisterCheckRightServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
