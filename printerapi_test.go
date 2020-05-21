package main

import (
	"context"
	"fmt"
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

	list, err := server.List(context.Background(), &pb.ListRequest{})
	if err != nil {
		t.Fatalf("Bad call: %v", err)
	}
	if len(list.GetQueue()) != 1 {
		t.Errorf("Bad queue: %v", list)
	}

	server.processPrints(context.Background())

	if server.prints != 1 {
		t.Errorf("Unable to print")
	}
}

func TestPrintFailOnLoop(t *testing.T) {
	server := InitTestServer()
	server.pretendret = fmt.Errorf("Built to fail")
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "inwhitelist"})

	err := server.processPrints(context.Background())

	if err == nil {
		t.Errorf("Did not fail")
	}
}

func TestClear(t *testing.T) {
	server := InitTestServer()
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "inwhitelist"})
	server.Clear(context.Background(), &pb.ClearRequest{})

	if server.prints != 0 {
		t.Errorf("We've recorded %v prints, despite not processing", server.prints)
	}

	server.Print(context.Background(), &pb.PrintRequest{Text: "hello there", Origin: "inwhitelist"})
	server.processPrints(context.Background())

	if server.prints != 1 {
		t.Errorf("Wrong number of prints recorded: %v", server.prints)
	}
}

func TestPrintFail(t *testing.T) {
	server := InitTestServer()
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "notinwhitelist"})

	if server.prints > 0 {
		t.Errorf("Unwhitelisted origin was printed")
	}
}
