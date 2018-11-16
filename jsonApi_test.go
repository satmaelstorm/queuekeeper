package main

import (
	"testing"
)

func TestResolveJsonGet(t *testing.T) {
	j := `{"count": 2,"method": "get"}`
	v, err := resolveJson([]byte(j))
	if err != nil {
		t.Error(err)
	}
	if v.Action != "get" {
		t.Error("method must be 'get', but got %s", v.Action)
	}
	if v.Cnt != 2 {
		t.Error("count must be 2, but got %d", v.Cnt)
	}
}

func TestResolveJsonPut(t *testing.T) {
	j := `{"method": "put", "messages": [{"message": "test", "delay": 0},{"message": "test2", "delay": 300}, {"message": "test3", "delay": 300}]}`
	v, err := resolveJson([]byte(j))
	if err != nil {
		t.Error(err)
	}
	if v.Action != "put" {
		t.Error("method must be 'put', but got %s", v.Action)
	}
	if len(v.Msgs) != 3 {
		t.Error("number of messages must be 3, %d got", len(v.Msgs))
	}
	if v.Msgs[0].Msg != "test" {
		t.Error("first message must be test, but %s got", v.Msgs[0].Msg)
	}
	if v.Msgs[2].Delay != 300 {
		t.Error("third message delay must be 300, but %d got", v.Msgs[2].Delay)
	}

}
