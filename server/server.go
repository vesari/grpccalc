package server

import (
	"context"

	pb "github.com/vesari/grpccalc/grpccalc"
)

// Server is used to implement grpccalc.CalcServer.
type Server struct {
	pb.UnimplementedCalcServer
}

// Add implements grpccalc.CalcServer
func (s *Server) Add(ctx context.Context, in *pb.AddRequest) (*pb.ValueReply, error) {
	result := in.Number1 + in.Number2
	return &pb.ValueReply{Value: result}, nil
}

// MultiplyF implements grpccalc.CalcServer
func (s *Server) MultiplyF(ctx context.Context, in *pb.MultiplyFRequest) (*pb.ValueFReply, error) {
	result := in.Number1 * in.Number2
	return &pb.ValueFReply{Value: result}, nil
}
