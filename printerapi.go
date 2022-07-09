package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/printer/proto"
)

// Print performs a print
func (s *Server) Print(ctx context.Context, req *pb.PrintRequest) (*pb.PrintResponse, error) {
	//Don't print it we're out of paper
	if s.outOfPaper {
		return nil, fmt.Errorf("Out of paper")
	}

	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	req.Id = time.Now().UnixNano()
	config.Requests = append(config.Requests, req)
	s.printq <- req
	return &pb.PrintResponse{Uid: req.Id}, s.save(ctx, config)
}

// Clear clears all the backlog
func (s *Server) Clear(ctx context.Context, req *pb.ClearRequest) (*pb.ClearResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetUid() > 0 {
		rs := []*pb.PrintRequest{}
		for _, req := range config.Requests {
			if req.Id != req.GetId() {
				rs = append(rs, req)
			}
		}
		config.Requests = rs
	} else {
		config.Requests = nil
	}
	return &pb.ClearResponse{}, s.save(ctx, config)
}

//List lists the backlog
func (s *Server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListResponse{Queue: config.Requests}, nil
}
