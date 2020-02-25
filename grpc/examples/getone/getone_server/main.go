package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"net"
	"seckilling-practice-project/common"
	pb "seckilling-practice-project/grpc"
	"sync"
)

var port = ":50051"
var ProductSum int64 = 10000
var Sum int64 = 0
var mutex sync.Mutex

type server struct {
	pb.UnimplementedGetOneServiceServer
}

//实现GetOneServiceServer
func (s *server) GetOne(ctx context.Context, req *empty.Empty) (*pb.Result, error) {
	if GetOneProduct() {
		fmt.Println(Sum)
		return &pb.Result{Value: true}, nil
	} else {
		fmt.Println(Sum)
		return &pb.Result{Value: false}, nil
	}
}

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	if Sum < ProductSum {
		Sum += 1

		return true
	}
	return false
}

func main() {
	lis, err := net.Listen("tcp", port)
	defer lis.Close()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := common.GetGrpcServer()
	pb.RegisterGetOneServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
