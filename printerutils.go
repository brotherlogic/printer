package main

import (
	"time"

	"github.com/brotherlogic/goserver/utils"
	"golang.org/x/net/context"

	pb "github.com/brotherlogic/printer/proto"
)

func (s *Server) printQueue() {
	for val := range s.printq {
		ctx, cancel := utils.ManualContext("printqueue", "printqueue", time.Minute, true)

		_, err := s.load(ctx)
		var t time.Duration

		if err == nil {
			t, err = s.processPrint(ctx, val)

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
	t, err := s.localPrint(req.Text, req.Lines, time.Now())

	if err != nil {
		return t, err
	}

	return t, s.dequeue(ctx, req)
}
