package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"github.com/joho/godotenv"
)

var (
	numOfRequests int
	mu      sync.Mutex
)

const filename = "/counter.txt"

// Opens (or creates) the counter file.
func createFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
}

// Writes formatted text to the file.
func writeToFile(file *os.File, format string, a ...interface{}) {
	// Clear previous contents
	if err := file.Truncate(0); err != nil {
		fmt.Printf("Error truncating file: %v\n", err)
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("Error seeking file: %v\n", err)
		return
	}

	if _, err := fmt.Fprintf(file, format, a...); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	file.Sync()
}

func main() {
    if err := godotenv.Load(); err != nil {
        fmt.Println("Warning: .env file not found")
    }

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Open counter file once
	counterFileHandle, err := createFile(filename)
	if err != nil {
		fmt.Printf("Failed to open counter file: %v\n", err)
		return
	}
	defer counterFileHandle.Close()
    numOfRequests := 0
    writeToFile(counterFileHandle, "%d\n", numOfRequests)
	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pingpong" {
			http.NotFound(w, r)
			return
		}

		mu.Lock()

		numOfRequests++

		// Save updated counter
		writeToFile(counterFileHandle, "%d\n", numOfRequests)

		mu.Unlock()

		fmt.Fprintf(w, "Ping / Pongs: %d", numOfRequests)
	})
	http.HandleFunc("/pings", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pings" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Ping / Pongs: %d", numOfRequests)
	})
	fmt.Printf("Listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}