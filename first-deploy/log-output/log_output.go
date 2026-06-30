package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type Status struct {
	Timestamp    string `json:"timestamp"`
	RandomString string `json:"randomString"`
}

func main() {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
		"@#!$%^&*()-+=_~`}{][|\\><?"

	// Read port from environment variable (default 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Generated once at startup
	randomString := generateRandomString(36, charset)

	fmt.Printf("Generated string: %s\n", randomString)

	// HTTP endpoint
	http.HandleFunc("/current-status", func(w http.ResponseWriter, r *http.Request) {
		status := Status{
			Timestamp:    time.Now().Format(time.RFC3339),
			RandomString: randomString,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
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
		now := time.Now().Format(time.RFC3339)
		fmt.Printf("[%s] %s\n", now, randomString)
		<-ticker.C
	}
}