package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"queuekeeper/qs"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type message struct {
	Msg   string `json:"message"`
	Delay int    `json:"delay"`
}

type postJson struct {
	Msgs   []message `json:"messages"`
	Cnt    int       `json:"count"`
	Action string    `json:"method"`
}

type answerJson struct {
	Msgs    []string `json:"messages"`
	Success bool     `json:"success"`
	Error   string   `json:"error"`
}

func resolveJson(jsonString []byte) (postJson, error) {
	result := postJson{}
	err := json.Unmarshal(jsonString, &result)
	if err == nil {
		return result, nil
	}
	return result, err
}

func makeQueueItem(msg message) *qs.QueueItem {
	qi := qs.NewQueueItem(msg.Msg, int64(msg.Delay))
	return qi
}

func postQueueHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	body, err := extractBodyAsByteArray(req)
	if nil != err {
		http.Error(w, err.Error(), 400)
		return
	}
	in, err := resolveJson(body)
	if nil != err {
		http.Error(w, err.Error(), 400)
	}
	if "" == in.Action {
		http.Error(w, "Method must be 'get' or 'put'", 400)
	}

	//	var answer answerJson
	switch strings.ToLower(in.Action) {
	case "get":
		return
	case "put":
		return
	default:
		http.Error(w, fmt.Sprintf("Unknown method %s", in.Action), 400)
		return
	}

}
