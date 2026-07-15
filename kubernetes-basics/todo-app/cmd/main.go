package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	imagePath     = "/usr/src/app/images/current.jpg"
	cacheDuration = 10 * time.Minute
)

func downloadImage() error {

	resp, err := http.Get("https://picsum.photos/1200")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func ensureImage() error {
	info, err := os.Stat(imagePath)

	// Download if it doesn't exist
	if os.IsNotExist(err) {
		return downloadImage()
	}

	if err != nil {
		return err
	}

	// Refresh after 10 minutes
	if time.Since(info.ModTime()) > cacheDuration {
		return downloadImage()
	}

	return nil
}

func main() {
	_ = godotenv.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// Serve images folder
	r.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("/usr/src/app/images"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := ensureImage()
		if err != nil {
			http.Error(w, "Could not get image", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>To-Do App</title>
		</head>
		<body>
			<h2>To-Do List</h2>
			<img src="/images/current.jpg" width="600">
		</body>
		</html>
		`)
	})

	http.ListenAndServe("0.0.0.0:"+port, r)
}