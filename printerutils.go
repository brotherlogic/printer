package main

import (
	"fmt"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/printer/proto"
)

func (s *Server) printQueue() {
	for val := range s.printq {
		goqueue.Set(float64(len(s.printq)))
		ctx, cancel := utils.ManualContext("printqueue", time.Minute)

		config, err := s.load(ctx)
		var t time.Duration

		stillInQueue := false
		for _, entry := range config.GetRequests() {
			if val.Id == entry.Id {
				stillInQueue = true
			}
		}

		s.CtxLog(ctx, fmt.Sprintf("Printing %v -> %v", val, stillInQueue))
		if err == nil && stillInQueue {
			t, err = s.processPrint(ctx, val)
			if err != nil && status.Convert(err).Code() != codes.Unavailable {
				s.RaiseIssue("Unable to print", fmt.Sprintf("Cannot print: %v", err))
			}

			time.Sleep(t)
		} else {
			s.CtxLog(ctx, fmt.Sprintf("%v is not in queue: %v", val, config.GetRequests()))
		}
		cancel()

		if err != nil {
			printErrors.With(prometheus.Labels{"error": fmt.Sprintf("%v", err)}).Inc()
			s.printq <- val
		}

	}

	s.done <- true
}

func (s *Server) dequeue(ctx context.Context, reqrem *pb.PrintRequest) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}
	newList := []*pb.PrintRequest{}
	for _, req := range config.GetRequests() {
		if req.GetId() != reqrem.GetId() {
			newList = append(newList, req)
		}
	}

	config.Requests = newList
	return s.save(ctx, config)
}

func (s *Server) processPrint(ctx context.Context, req *pb.PrintRequest) (time.Duration, error) {
	if req != nil {
		t, err := s.localPrint(ctx, req.Text, req.Lines, time.Now(), req.GetOverride())

		if err != nil {
			return t, err
		}

		return t, s.dequeue(ctx, req)
	} else {
		return time.Second, nil
	}
}
