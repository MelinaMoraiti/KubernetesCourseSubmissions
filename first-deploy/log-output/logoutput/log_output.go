package main

import (
	"fmt"
	"math/rand"
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

// Opens (or creates) the log file.
func createLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
}

// Writes formatted text to the file.
func writeToFile(file *os.File, format string, a ...interface{}) {
	if _, err := fmt.Fprintf(file, format, a...); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
	// Optional: ensure data is flushed immediately.
	file.Sync()
}

func main() {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
		"-"

	randomString := generateRandomString(36, charset)

	// Open log file
	logFile, err := createLogFile("/usr/src/app/files/")
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Print and write startup message
	writeToFile(logFile, "Generated string: %s\n", randomString)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		now := time.Now().Format(time.RFC3339)

		// Write same output to file
		writeToFile(logFile, "[%s] %s\n", now, randomString)

		<-ticker.C
	}
}