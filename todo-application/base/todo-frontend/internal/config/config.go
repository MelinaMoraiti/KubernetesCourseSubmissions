package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload" // loads .env if present
)

type Config struct {
	Port               int
	TodoBackendURL     string
	StoredImagesPath   string
	ImageURL           string
	ImageCacheDuration time.Duration
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("PORT is not set or invalid: %w", err)
	}

	backendURL, err := required("TODO_BACKEND_URL")
	if err != nil {
		return nil, err
	}

	imagesPath, err := required("STORED_IMAGES_PATH")
	if err != nil {
		return nil, err
	}

	imageURL, err := required("IMAGE_URL")
	if err != nil {
		return nil, err
	}

	cacheDuration, err := time.ParseDuration(os.Getenv("IMAGE_CACHE_DURATION"))
	if err != nil {
		return nil, fmt.Errorf("IMAGE_CACHE_DURATION is not set or invalid (e.g. 10m): %w", err)
	}

	return &Config{
		Port:               port,
		TodoBackendURL:     backendURL,
		StoredImagesPath:   imagesPath,
		ImageURL:           imageURL,
		ImageCacheDuration: cacheDuration,
	}, nil
}

func required(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("%s is not set", key)
	}
	return v, nil
}