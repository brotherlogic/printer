package main

import (
	"fmt"

	pb "github.com/brotherlogic/printer/proto"
	"golang.org/x/net/context"
)

// Print performs a print
func (s *Server) Print(ctx context.Context, req *pb.PrintRequest) (*pb.PrintResponse, error) {
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
	s.save(ctx)
	return &pb.PrintResponse{}, nil
}
