package main

import (
	"time"

	"golang.org/x/net/context"
)

func (s *Server) processPrints(ctx context.Context) error {
	if len(s.config.Requests) > 0 {
		req := s.config.Requests[0]
		err := s.localPrint(req.Text, req.Lines, time.Now())

		if err != nil {
			return err
		}

		s.config.Requests = s.config.Requests[1:]
		Backlog.Set(float64(len(s.config.Requests)))
		s.save(ctx)
	}

	return nil
}
