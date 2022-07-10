package main

import (
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/printer/proto"
)

// Print performs a print
func (s *Server) Print(ctx context.Context, req *pb.PrintRequest) (*pb.PrintResponse, error) {
	s.printlock.Lock()
	defer s.printlock.Unlock()
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	req.Id = time.Now().UnixNano()
	config.Requests = append(config.Requests, req)

	err = s.save(ctx, config)

	if err == nil {
		s.printq <- req
	}

	return &pb.PrintResponse{Uid: req.Id}, err
}

// Clear clears all the backlog
func (s *Server) Clear(ctx context.Context, req *pb.ClearRequest) (*pb.ClearResponse, error) {
	s.printlock.Lock()
	defer s.printlock.Unlock()
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetUid() > 0 {
		rs := []*pb.PrintRequest{}
		for _, pr := range config.Requests {
			if pr.Id != req.GetUid() {
				rs = append(rs, pr)
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
