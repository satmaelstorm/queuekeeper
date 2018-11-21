package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"

	"queuekeeper/qs"

	"github.com/braintree/manners"
	"github.com/julienschmidt/httprouter"
)

var qm *qs.QueueManager
var conf configuartion
var logger queueKeeperLogger
var adminReloadQueueConfigMutex sync.Mutex

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
		http.Error(w, err.Error(), 400)
		return
	}
	preparedMessage := message{Msg: body, Delay: -1}
	qi := makeQueueItem(preparedMessage)
	msg, err := q.Put(qi)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	io.WriteString(w, msg.String())
}

func adminReloadQueueConfigHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()
	adminReloadQueueConfigMutex.Lock()
	defer adminReloadQueueConfigMutex.Unlock()
	readQueuesConfigs(qm, conf)
	io.WriteString(w, qm.String())
}

func main() {
	qm = qs.NewQueueManager()
	conf = readGlobalConfig()
	logger = initLogger(conf.logConf)
	logger.log(QK_LOG_LEVEL_INFO, fmt.Sprintf("Read configuration: %v", conf))
	qm = readQueuesConfigs(qm, conf)
	runtime.GOMAXPROCS(conf.maxWorkers)

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	go func(ch <-chan os.Signal) {
		<-ch
		manners.Close()
	}(signalChannel)

	router := route()

	err := manners.ListenAndServe(":"+strconv.FormatInt(int64(conf.httpPort), 10), router)
	if err != nil {
		logger.log(QK_LOG_LEVEL_CRITICAL, "ListenAndServe: "+err.Error())
	}
}
