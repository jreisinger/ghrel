package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

// download downloads file from the url. It creates or truncates the file.
func download(url, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// fileName extracts filename from the URL.
func fileName(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	_, file := path.Split(u.Path)
	return file, nil
}
