package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	Error   []string `json:"error"`
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

func jsonResponse(w http.ResponseWriter, answer answerJson, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonBA, err := json.Marshal(answer)
	if err != nil {
		http.Error(w, err.Error(), 400)
		logger.log(QK_LOG_LEVEL_ERROR, "json error: "+err.Error())
		return
	}
	jsonStr := string(jsonBA[:])
	io.WriteString(w, jsonStr)
}

func postQueueHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	var answer answerJson

	body, err := extractBodyAsByteArray(req)
	if nil != err {
		answer.Error = append(answer.Error, err.Error())
		jsonResponse(w, answer, 400)
		return
	}
	in, err := resolveJson(body)
	if nil != err {
		answer.Error = append(answer.Error, err.Error())
		jsonResponse(w, answer, 400)
		logger.log(QK_LOG_LEVEL_WARNING, "json error: "+err.Error())
		return
	}
	if "" == in.Action {
		answer.Error = append(answer.Error, err.Error())
		logger.log(QK_LOG_LEVEL_WARNING, "json error: "+err.Error())
		jsonResponse(w, answer, 400)
		return
	}
	queueName := ps.ByName("queue")
	q, err := qm.GetQueue(queueName)
	if err != nil {
		answer.Error = append(answer.Error, "Queue "+queueName+" not found")
		jsonResponse(w, answer, 404)
		logger.log(QK_LOG_LEVEL_WARNING, "Queue "+queueName+" not found")
		return
	}
	code := 200
	switch strings.ToLower(in.Action) {
	case "get":
		answer = getJson(q, in)
	case "put":
		answer = putJson(q, in)
	default:
		answer.Error = append(answer.Error, fmt.Sprintf("Unknown method %s", in.Action))
		logger.log(QK_LOG_LEVEL_WARNING, "Unknown method "+in.Action)
		code = 400
	}
	jsonResponse(w, answer, code)
}

func putJson(q qs.ICommonQueue, in postJson) answerJson {
	ans := answerJson{}
	if len(in.Msgs) < 1 {
		ans.Error = append(ans.Error, "No messages")
		logger.log(QK_LOG_LEVEL_INFO, "putJson: No messages")
		return ans
	}
	for _, v := range in.Msgs {
		qi := makeQueueItem(v)
		qi, _ = q.Put(qi)
	}
	ans.Success = true
	return ans
}

func getJson(q qs.ICommonQueue, in postJson) answerJson {
	ans := answerJson{}
	if in.Cnt < 1 {
		ans.Error = append(ans.Error, "No messages")
		logger.log(QK_LOG_LEVEL_INFO, "getJson: No messages")
		return ans
	}
	for i := 0; i < in.Cnt; i++ {
		msg, err := q.Get()
		if nil != err {
			break
		}
		ans.Msgs = append(ans.Msgs, msg.String())
	}
	ans.Success = true
	return ans
}
