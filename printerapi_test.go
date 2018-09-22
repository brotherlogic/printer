package main

import (
	"context"
	"testing"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/printer/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.print = false
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")

	return s
}

func TestPrint(t *testing.T) {
	server := InitTestServer()
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello"})
}
