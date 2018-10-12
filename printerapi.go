package main

import (
	"time"

	pb "github.com/brotherlogic/printer/proto"
	"golang.org/x/net/context"
)

// Print performs a print
func (s *Server) Print(ctx context.Context, req *pb.PrintRequest) (*pb.PrintResponse, error) {
	err := s.localPrint(req.Text, req.Lines, time.Now())
	return &pb.PrintResponse{}, err
}
