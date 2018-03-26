package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

// var url = "http://ae35bcc3f30d911e88672066e78bb0cd-49118c41d965be0f.elb.us-west-2.amazonaws.com"
var url = "http://www.bbc.co.uk"
var successes uint64
var errors = []string{}

func main() {

	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range ticker.C {

			timeout := time.Duration(5 * time.Second)
			client := http.Client{Timeout: timeout}
			resp, err := client.Get(url)

			if err != nil {
				errors = append(errors, fmt.Sprintf("ERROR: %s", err))
				if len(errors) > 50 {
					errors = errors[:50]
				}
				fmt.Println(0)
				continue
			}

			if resp.StatusCode != 200 {
				errors = append(errors, fmt.Sprintf("ERROR: non-200 status code (%d)", resp.StatusCode))
				if len(errors) > 50 {
					errors = errors[:50]
				}
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
	status := strings.Join(errors, "\n")
	status += fmt.Sprintf("\nsuccesses: %d, failures: %d\n", atomic.LoadUint64(&successes), len(errors))
	w.Write([]byte(status))
}
