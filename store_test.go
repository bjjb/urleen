package main

import (
	"testing"
)

func Test_redisStore(t *testing.T) {
	s := &redisStore{url: "redis://localhost:6379"}
	s.open()
	s.ping()
	k := s.put("", "bar")
	if got := s.get(k); got != "bar" {
		t.Errorf("expected %q, got %q", "bar", got)
	}
}

func Test_mapStore(t *testing.T) {
	s := &mapStore{}
	s.open()
	k := s.put("", "bar")
	if got := s.get(k); got != "bar" {
		t.Errorf("expected %q, got %q", "bar", got)
	}
}
