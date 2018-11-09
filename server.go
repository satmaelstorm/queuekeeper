package main

import (
	"fmt"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"
	_ "os"
	_ "os/signal"
	"runtime"

	"queuekeeper/qs"

	_ "github.com/braintree/manners"
	"github.com/julienschmidt/httprouter"
)

var qm *qs.QueueManager

func extractBody(req *http.Request) string {
	bodyA, _ := ioutil.ReadAll(req.Body)
	return string(bodyA[:])

}

func getFromQueueHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	io.WriteString(w, "hello, "+ps.ByName("queue")+"!\n")
}

func putToQueueHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	io.WriteString(w, "hello, "+ps.ByName("queue")+"!\n")
	fmt.Fprintf(w, "%v", extractBody(req))
	req.Body.Close()
}

func main() {
	//	signalChannel := make(chan os.Signal)
	//	signal.Notify(signalChannel, os.Interrupt, os.Kill)
	//	go func(ch <-chan os.Signal) {
	//		<-ch
	//		manners.Close()
	//	}(signalChannel)

	//	router := httprouter.New()
	//	router.GET("/:queue", getFromQueueHandler)
	//	router.POST("/:queue", putToQueueHandler)
	//	err := manners.ListenAndServe(":12345", router)
	//	if err != nil {
	//		log.Fatal("ListenAndServe: ", err)
	//	}
	qm := qs.NewQueueManager()
	conf := readGlobalConfig()
	qm = readQueuesConfigs(qm, conf)
	runtime.GOMAXPROCS(conf.maxWorkers)

	q1, _ := qm.GetQueue("test")
	q2, _ := qm.GetQueue("test")

	q1.Put("test1")
	q1.Put("test1")
	q1.Put("test2")

	val, _ := q2.Get()
	fmt.Printf("%s\n", val)
	val, _ = q2.Get()
	fmt.Printf("%s\n", val)
	q2.Put("test3")
	val, _ = q2.Get()
	fmt.Printf("%s\n", val)
}
