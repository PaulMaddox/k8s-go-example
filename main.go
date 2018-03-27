package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

var url = os.Getenv("URL")
var successes uint64
var errors = []string{}
var errorsToDisplay = 10

func main() {

	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range ticker.C {

			timeout := time.Duration(5 * time.Second)
			client := http.Client{Timeout: timeout}
			resp, err := client.Get(url)

			if err != nil {
				errors = append(errors, fmt.Sprintf("ERROR: %s", err))
				if len(errors) > errorsToDisplay {
					errors = errors[:errorsToDisplay]
				}
				fmt.Println(0)
				continue
			}

			if resp.StatusCode != 200 {
				errors = append(errors, fmt.Sprintf("ERROR: non-200 status code (%d)", resp.StatusCode))
				if len(errors) > errorsToDisplay {
					errors = errors[:errorsToDisplay]
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
