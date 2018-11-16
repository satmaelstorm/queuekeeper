package main

import (
	"testing"
)

func TestResolveJson(t *testing.T) {
	j := `{"count": 2,"method": "get"}`
	bj := []byte(j)
	v, _ := resolveJson(bj)
	if v.Action != "get" {
		t.Error("method must be 'get', but got %s", v.Action)
	}
}
