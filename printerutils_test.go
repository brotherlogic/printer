package main

import (
	"testing"

	"golang.org/x/net/context"
)

func TestBadDequeue(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	err := s.dequeue(context.Background(), nil)
	if err == nil {
		t.Errorf("Bad dequeue did not fail")
	}
}

func TestBadProcessPrint(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	_, err := s.processPrint(context.Background(), nil)
	if err != nil {
		t.Errorf("Bad process did not fail")
	}
}
