package main

import (
	"context"
	"log"
	"time"

	pb "github.com/vesari/grpccalc/grpccalc"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCalcClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Add(ctx, &pb.AddRequest{Number1: 1, Number2: 45})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Result: %d", r.GetValue())
}
