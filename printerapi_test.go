package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	keystoreclient "github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/printer/proto"
)

func InitTestServer() *Server {
	s := Init()
	s.pretend = true
	s.SkipLog = true
	s.SkipIssue = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.KSclient.Save(context.Background(), KEY, &pb.Config{})

	return s
}

func TestPrint(t *testing.T) {
	server := InitTestServer()

	err := server.readyToPrint(context.Background())
	if err != nil {
		t.Errorf("Bad load: %v", err)
	}

	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "inwhitelist"})
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello2", Origin: "inwhitelist"})

	list, err := server.List(context.Background(), &pb.ListRequest{})
	if err != nil {
		t.Fatalf("Bad call: %v", err)
	}

	if len(list.GetQueue()) != 2 {
		t.Errorf("Bad queue: %v", list)
	}

	server.drainQueue()
	if server.prints != 2 {
		t.Errorf("Unable to print: %v", server.prints)
	}
}

func TestPrintFailOnLoop(t *testing.T) {
	server := InitTestServer()
	server.pretendret = fmt.Errorf("Built to fail")
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "inwhitelist"})

	err := server.readyToPrint(context.Background())
	if err != nil {
		t.Errorf("Bad load: %v", err)
	}

	// Let some prints go through
	time.Sleep(time.Second * 2)

	if server.prints != 0 {
		t.Errorf("Wrong number of prints recorded: %v", server.prints)
	}
}

func TestClear(t *testing.T) {
	server := InitTestServer()
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "inwhitelist"})
	server.Clear(context.Background(), &pb.ClearRequest{})

	if server.prints != 0 {
		t.Errorf("We've recorded %v prints, despite not processing", server.prints)
	}

	err := server.readyToPrint(context.Background())
	if err != nil {
		t.Errorf("Bad load: %v", err)
	}

	err = server.readyToPrint(context.Background())
	if err != nil {
		t.Errorf("Bad load: %v", err)
	}

	server.drainQueue()

	if server.prints != 0 {
		t.Errorf("Wrong number of prints recorded: %v", server.prints)
	}
}

func TestClearSingle(t *testing.T) {
	server := InitTestServer()
	server.Print(context.Background(), &pb.PrintRequest{Text: "hello", Origin: "inwhitelist"})
	res, _ := server.Print(context.Background(), &pb.PrintRequest{Text: "hello2", Origin: "inwhitelist"})

	if res.GetUid() == 0 {
		t.Fatalf("Print element has no UID: %v", res)
	}

	server.Clear(context.Background(), &pb.ClearRequest{Uid: res.GetUid()})

	if server.prints != 0 {
		t.Errorf("We've recorded %v prints, despite not processing", server.prints)
	}

	err := server.readyToPrint(context.Background())
	if err != nil {
		t.Errorf("Bad load: %v", err)
	}

	err = server.readyToPrint(context.Background())
	if err != nil {
		t.Errorf("Bad load: %v", err)
	}

	time.Sleep(time.Second * 5)

	server.drainQueue()

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

func TestFails(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	_, err := s.Print(context.Background(), nil)
	if err == nil {
		t.Errorf("Print did not fail")
	}

	_, err = s.Clear(context.Background(), nil)
	if err == nil {
		t.Errorf("Print did not fail")
	}

	_, err = s.List(context.Background(), nil)
	if err == nil {
		t.Errorf("Print did not fail")
	}

}
