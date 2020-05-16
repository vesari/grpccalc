//go:generate protoc -I../grpccalc/ --go_opt=paths=source_relative --go_out=plugins=grpc:../grpccalc ../grpccalc/grpccalc.proto

// Package main implements a server for the gRPCCalc service.
package main

import (
	"context"
	"log"
	"net"

	pb "github.com/vesari/grpccalc/grpccalc"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement grpccalc.CalcServer.
type server struct {
	pb.UnimplementedCalcServer
}

// Add implements grpccalc.CalcServer
func (s *server) Add(ctx context.Context, in *pb.AddRequest) (*pb.ValueReply, error) {
	result := in.Number1 + in.Number2
	return &pb.ValueReply{Value: result}, nil
}

func (s *server) MultiplyF(ctx context.Context, in *pb.MultiplyFRequest) (*pb.ValueFReply, error) {
	result := in.Number1 * in.Number2
	return &pb.ValueFReply{Value: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalcServer(s, &server{})
	log.Printf("Running the server!")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
