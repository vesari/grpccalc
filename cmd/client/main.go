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

func add(ctx context.Context, logger *log.Logger) error {
	if len(os.Args) != 4 {
		return fmt.Errorf("you have to provide 3 arguments")
	}
	operand1, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		return fmt.Errorf("unable to convert first operand to an integer")
	}
	operand2, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		return fmt.Errorf("unable to convert second operand to an integer")
	}

	c, conn, err := newClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	r, err := c.Add(ctx, &pb.AddRequest{Number1: operand1, Number2: operand2})
	if err != nil {
		return fmt.Errorf("could not add: %w", err)
	}

	logger.Printf("Result: %d", r.GetValue())
	return nil
}

func newClient() (pb.CalcClient, io.Closer, error) {
	portStr := strings.TrimSpace(os.Getenv("PORT"))
	if portStr == "" {
		portStr = "50051"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid env var PORT: %q", portStr)
	}
	host := strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		host = "localhost"
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %w", err)
	}
	c := pb.NewCalcClient(conn)

	return c, conn, nil
}

func multiplyF(ctx context.Context, logger *log.Logger) error {
	if len(os.Args) != 4 {
		return fmt.Errorf("you have to provide 3 arguments")
	}
	operand1, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		return fmt.Errorf("unable to convert first operand to a float")
	}
	operand2, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		return fmt.Errorf("unable to convert second operand to a float")
	}

	c, conn, err := newClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	r, err := c.MultiplyF(ctx, &pb.MultiplyFRequest{Number1: operand1, Number2: operand2})
	if err != nil {
		return fmt.Errorf("could not multiply: %w", err)
	}

	logger.Printf("Result: %.2f", r.GetValue())
	return nil
}

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	err := realMain(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func realMain(logger *log.Logger) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("not enough arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	op := os.Args[1]
	switch op {
	case "add":
		if err := add(ctx, logger); err != nil {
			return err
		}
	case "multiplyF":
		if err := multiplyF(ctx, logger); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid operator %q", op)
	}

	return nil
}
