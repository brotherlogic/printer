package main

import (
	"fmt"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brotherlogic/printer/proto"
)

func (s *Server) printQueue() {
	for val := range s.printq {
		ctx, cancel := utils.ManualContext("printqueue", "printqueue", time.Minute, true)

		_, err := s.load(ctx)
		var t time.Duration

		if err == nil {
			for i := 0; i < len(val.GetLines()); i++ {
				t, err = s.processPrint(ctx, val)
				if err != nil && status.Convert(err).Code() != codes.Unavailable {
					s.RaiseIssue("Unable to print", fmt.Sprintf("Cannot print: %v", err))
					break
				}
				val.LinePointer++
			}

			s.processPrint(ctx, nil)

			time.Sleep(t)
		}
		cancel()

		if err != nil {
			s.printq <- val
		}

		Backlog.Set(float64(len(s.printq)))
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
	Backlog.Set(float64(len(config.Requests)))
	return s.save(ctx, config)
}

func (s *Server) processPrint(ctx context.Context, req *pb.PrintRequest) (time.Duration, error) {
	if req != nil {
		t, err := s.localPrint(req.Text, req.Lines[req.LinePointer], time.Now())

		if err != nil {
			return t, err
		}

		return t, s.dequeue(ctx, req)
	} else {
		s.localPrint("", "", time.Now())
		s.localPrint("", "", time.Now())
		return time.Second, nil
	}
}
