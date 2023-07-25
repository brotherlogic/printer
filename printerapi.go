package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/printer/proto"
)

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	data, err := exec.Command("/usr/bin/lsusb").Output()
	if err != nil {
		return nil, fmt.Errorf("error listing usb components: %v", err)
	}
	s.CtxLog(ctx, fmt.Sprintf("USBRES: %v", string(data)))
	if strings.Contains(string(data), "TSP100II") {
		return &pb.PingResponse{}, nil
	}

	return nil, status.Errorf(codes.FailedPrecondition, "printer is not available")
}

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

// List lists the backlog
func (s *Server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	config, err := s.load(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListResponse{Queue: config.Requests}, nil
}
