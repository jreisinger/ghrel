package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
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
		ok, err := verifyChecksums(c)
		if err != nil {
			log.Fatal(err)
		}
		if !ok {
			fmt.Printf("NOT OK\n")
			return
		}
	}
	fmt.Printf("OK\n")
}

func verifyChecksums(checksumsFile string) (ok bool, err error) {
	// cmd := exec.Command("sha256sum", "-c", checksumsFile)
	// return cmd.Run()

	checksums, err := extractChecksums(checksumsFile)
	if err != nil {
		return false, err
	}

	for _, c := range checksums {
		ok, err := verifyFileChecksum(c)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

type checksum struct {
	checksum []byte // in hex
	filename string
}

func extractChecksums(checksumsFile string) ([]checksum, error) {
	file, err := os.Open(checksumsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var checksums []checksum

	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		c, err := hex.DecodeString(fields[0])
		if err != nil {
			return nil, err
		}
		checksums = append(checksums, checksum{c, fields[1]})
	}

	return checksums, nil
}

func verifyFileChecksum(c checksum) (ok bool, err error) {
	file, err := os.Open(c.filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}
	sum := hash.Sum(nil)

	return bytes.Equal(c.checksum, sum), nil
}

var checksumsFile = regexp.MustCompile(`(?i)check.?sum`)

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
