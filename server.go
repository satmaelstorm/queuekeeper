package main

import (
	"fmt"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"

	"queuekeeper/qs"

	_ "github.com/bmizerany/pat"
)

func extractBody(req *http.Request) string {
	bodyA, _ := ioutil.ReadAll(req.Body)
	return string(bodyA[:])

}

func getFromQueueHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, "+req.URL.Query().Get(":queue")+"!\n")
}

func putToQueueHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, "+req.URL.Query().Get(":queue")+"!\n")
	fmt.Fprintf(w, "%v", extractBody(req))
	req.Body.Close()
}

func main() {
	//	m := pat.New()
	//	m.Get("/:queue", http.HandlerFunc(getFromQueueHandler))
	//	m.Post("/:queue", http.HandlerFunc(putToQueueHandler))

	//	http.Handle("/", m)
	//	err := http.ListenAndServe(":12345", nil)
	//	if err != nil {
	//		log.Fatal("ListenAndServe: ", err)
	//	}
	q := qs.NewQueue()
	q.Put("1")
	q.Put("2")
	q.Put("3")
	v, _ := q.Get()
	fmt.Printf("%s\n", v)
	v, _ = q.Get()
	fmt.Printf("%s\n", v)
	v, _ = q.Get()
	fmt.Printf("%s\n", v)
}
