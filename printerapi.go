package main

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/printer/proto"
	pqpb "github.com/brotherlogic/printqueue/proto"
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
	s.RaiseIssue("Bad printer", fmt.Sprintf("%v is trying to use print directly - migrate them to printqueue", req.GetOrigin()))

	// Reflect this over to the printqueue
	conn, err := grpc.Dial("print.brotherlogic-backend.com:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pqpb.NewPrintServiceClient(conn)
	urgency := pqpb.Urgency_URGENCY_REGULAR
	if req.GetOverride() {
		urgency = pqpb.Urgency_URGENCY_IMMEDIATE
	}
	_, err = client.Print(ctx, &pqpb.PrintRequest{
		Lines:       req.GetLines(),
		Urgency:     urgency,
		Destination: pqpb.Destination_DESTINATION_RECEIPT,
		Origin:      req.GetOrigin(),
		Fanout:      pqpb.Fanout_FANOUT_ONE,
	})
	if err != nil {
		return nil, err
	}

	return &pb.PrintResponse{}, err
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
