package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Fetches JSON data from an url and deserializes response.
func FetchJson[T any](url string) (*T, http.Header, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.Header, fmt.Errorf("http response status %s", resp.Status)
	}

	var result T
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, resp.Header, err
	}

	return &result, resp.Header, nil
}

func DownloadFile(url string, directory string) (string, error) {
	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Get file name
	var filename string
	cd := resp.Header.Get("Content-Disposition")
	if start := strings.Index(cd, "filename="); start != -1 {
		start += len("filename=")
		end := strings.Index(cd[start:], ";")
		if end == -1 {
			end = len(cd)
		} else {
			end += start
		}
		filename = strings.Trim(cd[start:end], " \"")
	} else {
		filename = "DownloadedFile"
	}

	// Get total file size
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return "", err
	}

	file := filepath.Join(directory, filename)
	out, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(progressWriter(out, size), resp.Body); err != nil {
		return "", err
	}

	return file, nil
}
