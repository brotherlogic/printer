package main

import (
	"context"
	"testing"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/printer/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.pretend = true
	s.SkipLog = true
	s.whitelist = []string{"inwhitelist"}
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")

	return s
}

func TestPrint(t *testing.T) {
	server := InitTestServer()
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "inwhitelist"})

	server.processPrints(context.Background())

	if server.prints != 1 {
		t.Errorf("Unable to print")
	}
}

func TestPrintFail(t *testing.T) {
	server := InitTestServer()
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "notinwhitelist"})

	if server.prints > 0 {
		t.Errorf("Unwhitelisted origin was printed")
	}
}
