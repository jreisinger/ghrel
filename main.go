package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sync"
)

func main() {
	log.SetFlags(0) // no timestamp
	log.SetPrefix(os.Args[0] + ": ")

	if len(os.Args[1:]) != 1 {
		log.Fatal("supply github <owner>/<repo>")
	}
	repo := os.Args[1]

	urls, err := getDownloadUrls(repo)
	if err != nil {
		log.Fatal(err)
	}

	var checksumFiles []string
	var infoFmt = "%-30s"

	fmt.Printf(infoFmt, "downloading release files ")
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			file, _ := fileName(url)
			if err := download(url); err != nil {
				log.Print(err)
				return
			}
			if isChecksumsFile(file) {
				checksumFiles = append(checksumFiles, file)
			}
		}(url)
	}
	wg.Wait()
	fmt.Printf("OK\n")

	fmt.Printf(infoFmt, "verifying checksums ")
	for _, c := range checksumFiles {
		if err := verifyChecksums(c); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("OK\n")
}

func verifyChecksums(checksumsFile string) error {
	cmd := exec.Command("sha256sum", "-c", checksumsFile)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}

var checksumsFile = regexp.MustCompile(`(?i)checksum`)

func isChecksumsFile(filename string) bool {
	return checksumsFile.MatchString(filename)
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

var api_url = "https://api.github.com"

// getDownloadUrls returns URLs for downloading assets from the latest repo release.
func getDownloadUrls(repo string) ([]string, error) {
	url := api_url + "/repos/" + repo + "/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
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
