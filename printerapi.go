package main

import "golang.org/x/net/context"
import pb "github.com/brotherlogic/printer/proto"

// Print performs a print
func (s *Server) Print(ctx context.Context, req *pb.PrintRequest) (*pb.PrintResponse, error) {
	err := s.localPrint(req.Text)
	return nil, err
}
