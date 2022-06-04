package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

func main() {
	log.SetFlags(0) // no timestamp
	log.SetPrefix(os.Args[0] + ": ")

	if len(os.Args[1:]) != 1 {
		log.Fatal("supply github <owner>/<repo> to download latest release assets from")
	}
	repo := os.Args[1]

	urls, err := getUrls(repo)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if err := download(url); err != nil {
				log.Print(err)
			}
		}(url)
	}
	wg.Wait()
}

// download downloads file from the url.
func download(url string) error {
	file, err := fileName(url)
	if err != nil {
		return err
	}
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

// getUrls returns URLs for downloading assets from the latest repo release.
func getUrls(repo string) ([]string, error) {
	api_url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := http.Get(api_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting %s: %s", api_url, resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := struct {
		Assets []struct {
			BrowserDownloadUrl string `json:"browser_download_url"`
		}
	}{}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	var urls []string
	for _, a := range r.Assets {
		urls = append(urls, a.BrowserDownloadUrl)
	}
	return urls, nil
}
