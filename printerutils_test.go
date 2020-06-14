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
