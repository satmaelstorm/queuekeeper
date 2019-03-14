package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

func getFromQueueHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	queueName := ps.ByName("queue")
	q, err := qm.GetQueue(queueName)
	if err != nil {
		io.WriteString(w, "Queue "+queueName+" not found")
		logger.log(QK_LOG_LEVEL_WARNING, "Queue "+queueName+" not found")
		http.NotFound(w, req)
		return
	}
	msg, err := q.Get()
	if err != nil {
		http.Error(w, err.Error(), 400)
		logger.log(QK_LOG_LEVEL_ERROR, "Get from queue error: "+err.Error())
		return
	}
	logger.log(QK_LOG_LEVEL_TRACE, "Get from queue: "+msg.String())
	io.WriteString(w, msg.String())
}

func putToQueueHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	queueName := ps.ByName("queue")
	q, err := qm.GetQueue(queueName)
	if err != nil {
		io.WriteString(w, "Queue "+queueName+" not found")
		logger.log(QK_LOG_LEVEL_WARNING, "Queue "+queueName+" not found")
		http.NotFound(w, req)
		return
	}
	body, err := extractBody(req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	preparedMessage := message{Msg: body, Delay: -1}
	qi := makeQueueItem(preparedMessage)
	msg, err := q.Put(qi)
	if err != nil {
		http.Error(w, err.Error(), 400)
		logger.log(QK_LOG_LEVEL_ERROR, "Put in queue error: "+err.Error())
	}
	logger.log(QK_LOG_LEVEL_TRACE, "Put in queue: "+msg.String())
	io.WriteString(w, msg.String())
}

func adminReloadQueueConfigHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	adminReloadQueueConfigMutex.Lock()
	defer adminReloadQueueConfigMutex.Unlock()
	logger.log(QK_LOG_LEVEL_INFO, "Reload queue config")
	readQueuesConfigs(qm, conf)
	io.WriteString(w, qm.String())
}

func healthRoute(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	logger.log(QK_LOG_LEVEL_INFO, "Read health info")
	h := health()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonBA, err := json.Marshal(h)
	if err != nil {
		http.Error(w, err.Error(), 400)
		logger.log(QK_LOG_LEVEL_ERROR, "json error: "+err.Error())
		return
	}
	jsonStr := string(jsonBA[:])
	io.WriteString(w, jsonStr)
}
