package main

import (
	"log"

	"todo-app/internal/server"
    "time"
	"todo-app/internal/utils"
	"todo-app/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

    // Define cache parameters for Picsum image
	imagePath :=  cfg.StoredImagesPath + "current.jpg"

	if err := utils.CacheImage(imagePath, cfg.ImageURL, cfg.ImageCacheDuration); err != nil {
    	log.Printf("Warning: Failed to initial cache image: %v", err)
    }

    // 2. Start a background ticker to update the image periodically
    go func() {
    		ticker := time.NewTicker(cfg.ImageCacheDuration)
    		defer ticker.Stop()

    		for range ticker.C {
    			if err := utils.CacheImage(imagePath, cfg.ImageURL, cfg.ImageCacheDuration); err != nil {
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