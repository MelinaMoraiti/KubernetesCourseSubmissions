package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const logFile = "/usr/src/app/files/"

func logHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests.
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	file, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Log file not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to read log file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to send log file", http.StatusInternalServerError)
	}
}

func main() {
	// Load .env file (optional if it doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}
	http.HandleFunc("/logs", logHandler)

	addr := ":" + port

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}