package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/vesari/grpccalc/grpccalc"
	"google.golang.org/grpc"
)

func add(ctx context.Context) {
	if len(os.Args) != 4 {
		log.Fatal("You have to provide 3 arguments")
	}
	operand1, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		log.Fatal("Unable to convert first operand to an integer")
	}
	operand2, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		log.Fatal("Unable to convert second operand to an integer")
	}
	c, conn := newClient()
	defer conn.Close()
	r, err := c.Add(ctx, &pb.AddRequest{Number1: operand1, Number2: operand2})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Result: %d", r.GetValue())
}

func newClient() (pb.CalcClient, io.Closer) {
	portStr := strings.TrimSpace(os.Getenv("PORT"))
	if portStr == "" {
		portStr = "50051"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid env var PORT: %q", portStr)
	}
	host := strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		host = "localhost"
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewCalcClient(conn)

	return c, conn
}

func multiplyF(ctx context.Context) {
	if len(os.Args) != 4 {
		log.Fatal("You have to provide 3 arguments")
	}
	operand1, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		log.Fatal("Unable to convert first operand to a float")
	}
	operand2, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		log.Fatal("Unable to convert second operand to a float")
	}
	c, conn := newClient()
	defer conn.Close()
	r, err := c.MultiplyF(ctx, &pb.MultiplyFRequest{Number1: operand1, Number2: operand2})
	if err != nil {
		log.Fatalf("could not multiply: %v", err)
	}
	log.Printf("Result: %.2f", r.GetValue())
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	op := os.Args[1]
	switch op {
	case "add":
		add(ctx)
	case "multiplyF":
		multiplyF(ctx)
	default:
		log.Fatalf("Invalid operator %q", op)
	}
}
