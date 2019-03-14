package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"

	"queuekeeper/qs"

	"github.com/braintree/manners"
)

var qm *qs.QueueManager
var conf configuration
var logger queueKeeperLogger
var adminReloadQueueConfigMutex sync.Mutex

func extractBody(req *http.Request) (string, error) {
	bodyA, err := ioutil.ReadAll(req.Body)
	if nil != err {
		logger.log(QK_LOG_LEVEL_ERROR, "Extract body error: "+err.Error())
		return "", err
	}
	return string(bodyA[:]), nil

}

func extractBodyAsByteArray(req *http.Request) ([]byte, error) {
	bodyA, err := ioutil.ReadAll(req.Body)
	if nil != err {
		logger.log(QK_LOG_LEVEL_ERROR, "Extract body as byte array error: "+err.Error())
		return nil, err
	}
	return bodyA, nil
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
