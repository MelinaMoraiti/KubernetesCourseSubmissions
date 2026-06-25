package main

import (
	"fmt"
	"math/rand"
    "time"
)
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int, charset string) (string) {
    b := make([]byte, length)
	for i := range b {
      b[i] = charset[seededRand.Intn(len(charset))]
    }
    return string(b)
}

func main() {
	// Generate a random string once at startup.
	const charset = "abcdefghijklmnopqrstuvwxyz" +
    "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
    "@#!$%^&*()-+=_~`}{][|\\><?"
	randomString:= generateRandomString(36, charset)


	fmt.Printf("Generated string: %s\n", randomString)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		now := time.Now().Format(time.RFC3339)
		fmt.Printf("[%s] %s\n", now, randomString)

		<-ticker.C
	}
}