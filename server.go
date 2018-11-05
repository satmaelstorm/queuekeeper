package main

import (
	"fmt"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"
	_ "os"
	_ "os/signal"

	"queuekeeper/qs"

	_ "github.com/braintree/manners"
	"github.com/julienschmidt/httprouter"
)

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
	q := qs.NewDedupQueue(qs.NewQueueFlags())
	q.Put("3")
	q.Put("1")
	q.Put("1")
	q.Put("2")
	q.Put("2")
	q.Put("1")
	v, _ := q.Get()
	fmt.Printf("%s\n", v)
	v, _ = q.Get()
	fmt.Printf("%s\n", v)
	v, _ = q.Get()
	fmt.Printf("%s\n", v)
	fmt.Printf("%v\n", readGlobalConfig())
}
