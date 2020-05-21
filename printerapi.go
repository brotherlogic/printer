package main

import (
	"fmt"

	pb "github.com/brotherlogic/printer/proto"
	"golang.org/x/net/context"
)

// Print performs a print
func (s *Server) Print(ctx context.Context, req *pb.PrintRequest) (*pb.PrintResponse, error) {
	s.Log(fmt.Sprintf("RECEIVED PRINT REQUEST"))
	found := false
	for _, whitelisted := range s.whitelist {
		if req.Origin == whitelisted {
			found = true
		}
	}

	if !found {
		return &pb.PrintResponse{}, fmt.Errorf("Origin is not in the whitelist")
	}

	s.config.Requests = append(s.config.Requests, req)
	Backlog.Set(float64(len(s.config.Requests)))
	s.save(ctx)
	return &pb.PrintResponse{}, nil
}

// Clear clears all the backlog
func (s *Server) Clear(ctx context.Context, req *pb.ClearRequest) (*pb.ClearResponse, error) {
	s.config.Requests = nil
	s.save(ctx)
	return &pb.ClearResponse{}, nil
}

//List lists the backlog
func (s *Server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	return &pb.ListResponse{Queue: s.config.Requests}, nil
}
