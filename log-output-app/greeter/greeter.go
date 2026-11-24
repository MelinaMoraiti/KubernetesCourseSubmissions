package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/hello" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "greetings: %s\n", os.Getenv("HELLO_GREETING_MESSAGE"))
	})

	fmt.Printf("Listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}