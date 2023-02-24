package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/jreisinger/ghrel/assets"
)

var shellPattern = flag.String("p", "", "donwload only files matching shell `pattern`")
var onlyList = flag.Bool("l", false, "only list files that would be downloaded")

func main() {
	flag.Usage = func() {
		desc := "Download assets (files) of the latest release from a GitHub repository."
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n\n%s [flags] <owner>/<repo>\n", desc, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	repo := flag.Args()[0]

	log.SetFlags(0) // no timestamp
	log.SetPrefix(os.Args[0] + ": ")

	urls, err := assets.GetDownloadUrls(repo)
	if err != nil {
		log.Fatalf("get donwload URLs for release assets: %v", err)
	}

	var checksumFiles []string
	var count struct {
		mu    sync.Mutex
		files int
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			file, err := fileName(url)
			if err != nil {
				log.Printf("extract filename from URL: %v", err)
				return
			}

			if isChecksumsFile(file) {
				checksumFiles = append(checksumFiles, file)
			} else if *shellPattern != "" {
				matched, err := filepath.Match(*shellPattern, file)
				if err != nil {
					log.Fatal(err)
				}
				if !matched {
					return
				}
			}

			if *onlyList {
				fmt.Println(file)
				return
			}

			if err := download(url, file); err != nil {
				log.Printf("download %s: %v", url, err)
				return
			}

			count.mu.Lock()
			count.files++
			count.mu.Unlock()
		}(url)
	}
	wg.Wait()

	if !*onlyList {
		fmt.Printf("downloaded %d file(s)\n", count.files)

		var verifiedFiles int
		for _, c := range checksumFiles {
			checksums, err := verifyChecksums(c)
			if err != nil {
				log.Fatalf("%s: %v", c, err)
			}
			for _, c := range checksums {
				if !c.verified {
					log.Printf("%s not verified", c.filename)
				} else {
					verifiedFiles += 1
				}
			}
		}
		fmt.Printf("verified %d file(s)\n", verifiedFiles)
	}
}
