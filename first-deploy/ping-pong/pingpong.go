package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

var (
	counter int
	mu      sync.Mutex
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	    if r.URL.Path != "/" {
            return
        }
		mu.Lock()
		current := counter
		counter++
		mu.Unlock()

		fmt.Fprintf(w, "pong %d", current)
	})

	fmt.Printf("Listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}