package main

import (
	"log"

	"todo-app/internal/server"
    "time"
	"github.com/joho/godotenv"
	"todo-app/internal/utils"
)

func main() {
	_ = godotenv.Load()
    // Define cache parameters for Picsum image
	imagePath :=  os.Getenv("STORED_IMAGES_PATH") + "current.jpg"
	imageURL := os.Getenv("IMAGE_URL")

	cacheDurationString := os.Getenv("IMAGE_CACHE_DURATION")
	cacheDuration, err := time.ParseDuration(cacheDurationString)

	if err := utils.CacheImage(imagePath, imageURL, cacheDuration); err != nil {
    	log.Printf("Warning: Failed to initial cache image: %v", err)
    }

    // 2. Start a background ticker to update the image periodically
    go func() {
    		ticker := time.NewTicker(cacheDuration)
    		defer ticker.Stop()

    		for range ticker.C {
    			if err := utils.CacheImage(imagePath, imageURL, cacheDuration); err != nil {
    				log.Printf("Warning: Failed to update cached image: %v", err)
    			} else {
    				log.Println("Successfully updated cached image")
    			}
    		}
    	}()

	srv := server.NewServer()
	log.Printf("Server listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}