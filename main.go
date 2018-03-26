package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

var url = "http://ae35bcc3f30d911e88672066e78bb0cd-49118c41d965be0f.elb.us-west-2.amazonaws.com"
var successes uint64
var failures uint64

func main() {

	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range ticker.C {
			resp, err := http.Get(url)
			if err != nil || resp.StatusCode != 200 {
				atomic.AddUint64(&failures, 1)
				fmt.Println(0)
				continue
			}
			atomic.AddUint64(&successes, 1)
			fmt.Println(1)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", statsHandler)
	log.Fatal(http.ListenAndServe(":8000", r))

}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	status := fmt.Sprintf("successes: %d, failures: %d", atomic.LoadUint64(&successes), atomic.LoadUint64(&failures))
	w.Write([]byte(status))
}
