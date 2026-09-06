package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
    "github.com/joho/godotenv"
)
const counterFile = "/usr/src/app/files/counter.txt"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func fetchPings(pingsURL string) string {
	resp, err := http.Get(pingsURL)
	if err != nil {
		return "Error fetching pings"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response"
	}

	return string(body)
}

func main() {
    if err := godotenv.Load(); err != nil {
        fmt.Println("Warning: .env file not found")
    }

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
		"@#!$%^&*()-+=_~`}{][|\\><?"

	// Read port from environment variable (default 8080)
	port := os.Getenv("PORT")
    ping_pong_host := os.Getenv("PING_PONG_APP_HOST")
	// Generated once at startup
	randomString := generateRandomString(36, charset)

	// HTTP endpoint
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        pingsURL := ping_pong_host + "/pings"
        pings := fetchPings(pingsURL)
        timestamp := time.Now().UTC().Format(time.RFC3339Nano)
		fmt.Fprintf(
			w,
			"%s: %s.\n%s",
			timestamp,
			randomString,
			pings,
		)
    })
	// Start HTTP server
	go func() {
		fmt.Printf("HTTP server listening on :%s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			panic(err)
		}
	}()

	// Continue logging every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

    for {
        //now := time.Now().UTC().Format(time.RFC3339Nano)

        //fmt.Printf("%s: %s.", now, randomString)

        <-ticker.C
    }
}