package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"

	"queuekeeper/qs"

	"github.com/braintree/manners"
	"github.com/julienschmidt/httprouter"
)

var qm *qs.QueueManager
var conf configuartion

func extractBody(req *http.Request) (string, error) {
	bodyA, err := ioutil.ReadAll(req.Body)
	if nil != err {
		return "", err
	}
	return string(bodyA[:]), nil

}

func extractBodyAsByteArray(req *http.Request) ([]byte, error) {
	bodyA, err := ioutil.ReadAll(req.Body)
	if nil != err {
		return nil, err
	}
	return bodyA, nil
}

func getFromQueueHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	queueName := ps.ByName("queue")
	q, err := qm.GetQueue(queueName)
	if err != nil {
		io.WriteString(w, "Queue "+queueName+" not found")
		http.NotFound(w, req)
		return
	}
	msg, err := q.Get()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	io.WriteString(w, msg.String())
}

func putToQueueHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	queueName := ps.ByName("queue")
	q, err := qm.GetQueue(queueName)
	if err != nil {
		io.WriteString(w, "Queue "+queueName+" not found")
		http.NotFound(w, req)
		return
	}
	body, err := extractBody(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	msg, err := q.Put(body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	io.WriteString(w, msg.String())
}

func adminReloadQueueConfigHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	readQueuesConfigs(qm, conf)
}

func main() {
	qm = qs.NewQueueManager()
	conf = readGlobalConfig()
	fmt.Printf("Read configuration: %v\n", conf)
	qm = readQueuesConfigs(qm, conf)
	runtime.GOMAXPROCS(conf.maxWorkers)

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	go func(ch <-chan os.Signal) {
		<-ch
		manners.Close()
	}(signalChannel)

	router := httprouter.New()
	router.GET("/q/:queue", getFromQueueHandler)
	router.PUT("/q/:queue", putToQueueHandler)
	router.GET("/admin/reload/queues", adminReloadQueueConfigHandler)
	err := manners.ListenAndServe(":"+strconv.FormatInt(int64(conf.httpPort), 10), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
