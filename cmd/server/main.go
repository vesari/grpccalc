//go:generate protoc -I../../grpccalc --go_out ../../grpccalc --go_opt paths=source_relative --go-grpc_out ../../grpccalc --go-grpc_opt paths=source_relative grpccalc.proto

// Package main implements a server for the gRPCCalc service.
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	pb "github.com/vesari/grpccalc/grpccalc"
	"github.com/vesari/grpccalc/server"
	"google.golang.org/grpc"
)

func main() {
	portStr := strings.TrimSpace(os.Getenv("PORT"))
	if portStr == "" {
		portStr = "50051"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Env var PORT has invalid value %q", portStr)
	}
	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalcServer(s, &server.Server{})
	log.Printf("Listening on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
