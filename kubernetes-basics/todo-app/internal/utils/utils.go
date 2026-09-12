package utils
import (
	"io"
	"net/http"
	"os"


	"time"
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