package main

import (
    "context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
)

var (
	numOfRequests int
	mu      sync.Mutex
	db      *pgx.Conn
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
func connectToDB(databaseURL string) error {
	var err error

	db, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Loads the current counter value from the database.
func loadCounterFromDB() error {
	err := db.QueryRow(
		context.Background(),
		"SELECT value FROM counter WHERE id = 1",
	).Scan(&numOfRequests)

	if err != nil {
		return fmt.Errorf("failed to load counter: %w", err)
	}

	return nil
}

// Saves the provided counter value to the database.
func saveCounterToDB(counter int) error {
	_, err := db.Exec(
		context.Background(),
		"UPDATE counter SET value = $1 WHERE id = 1",
		counter,
	)

	if err != nil {
		return fmt.Errorf("failed to save counter: %w", err)
	}

	return nil
}
func main() {
    if err := godotenv.Load(); err != nil {
        fmt.Println("Warning: .env file not found")
    }

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("DATABASE_URL is not set")
		return
	}

	// Connect to PostgreSQL.
	if err := connectToDB(databaseURL); err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close(context.Background())

	// Load the persisted counter.
	if err := loadCounterFromDB(); err != nil {
		fmt.Println(err)
		return
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
		if err := saveCounterToDB(numOfRequests); err != nil {
			http.Error(w, "Failed to save counter", http.StatusInternalServerError)
			return
		}
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