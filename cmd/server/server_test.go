package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/vesari/grpccalc/grpccalc"
	"github.com/vesari/grpccalc/server"
	"google.golang.org/grpc"
)

func TestAdd(t *testing.T) {
	p := startServer(t)
	c := newClient(t, p)

	t.Log("Starting tests")
	t.Run("1 and 2", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second) //timeout after 1 sec
		t.Cleanup(cancel)                                                     //this is to cancel the context after the test if necessary
		r, err := c.Add(ctx, &pb.AddRequest{Number1: 1, Number2: 2})
		require.NoError(t, err)

		assert.Equal(t, int64(3), r.Value) //making sure it's int64
	})
}

func TestMultiplyF(t *testing.T) {
	p := startServer(t)
	c := newClient(t, p)

	t.Log("Starting tests")
	t.Run("1.2 and 2.1", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		t.Cleanup(cancel)
		r, err := c.MultiplyF(ctx, &pb.MultiplyFRequest{Number1: 1.2, Number2: 2.1})
		require.NoError(t, err)

		assert.Equal(t, 2.52, r.Value)
	})
}

func startServer(t *testing.T) int {
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	elems := strings.Split(lis.Addr().String(), ":")
	pStr := elems[len(elems)-1]
	p, err := strconv.Atoi(pStr)
	require.NoError(t, err)
	s := grpc.NewServer()
	pb.RegisterCalcServer(s, &server.Server{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("Server failed: %s", err)
		}
	}()

	t.Cleanup(func() {
		t.Log("Stopping server")
		s.Stop()
	})

	return p
}

func newClient(t *testing.T, port int) pb.CalcClient {
	address := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn.Close()
	})
	c := pb.NewCalcClient(conn)

	return c
}
