package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"seckilling-practice-project/common"
	"time"

	pb "seckilling-practice-project/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn := common.GetGrpcClientConn(address)
	defer conn.Close()
	c := pb.NewGetOneServiceClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetOne(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %t", r.GetValue())
}
