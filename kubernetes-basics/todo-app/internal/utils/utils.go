package utils
import (
	"io"
	"net/http"
	"os"
	"encoding/json"
	"fmt"
	"time"
	"todo-app/internal/models"
)

func DownloadImage(url, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func CacheImage(path, url string, maxAge time.Duration) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return DownloadImage(url, path)
	}

	if err != nil {
		return err
	}

	if time.Since(info.ModTime()) > maxAge {
		return DownloadImage(url, path)
	}

	return nil
}

func FetchAllTodosFromJSON(filePath string) ([]models.Todo, error) {
    file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open json file: %w", err)
	}
	defer file.Close()

	var todos []models.Todo
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todos); err != nil {
		return nil, fmt.Errorf("failed to decode json data: %w", err)
	}

	return todos, nil
}