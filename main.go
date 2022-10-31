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

func main() {
	log.SetFlags(0) // no timestamp
	log.SetPrefix(os.Args[0] + ": ")

	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("supply github <owner>/<repo>")
	}
	repo := flag.Args()[0]

	urls, err := assets.GetDownloadUrls(repo)
	if err != nil {
		log.Fatalf("get donwload URLs for release assets: %v", err)
	}

	var checksumFiles []string
	var count struct {
		mu    sync.Mutex
		files int
	}
	var infoFmt = "%-30s"

	fmt.Printf(infoFmt, "downloading release files ")
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
	fmt.Printf("OK (%d)\n", count.files)

	fmt.Printf(infoFmt, "verifying checksums ")
	var verifiedFiles int
	for _, c := range checksumFiles {
		n, err := verifyChecksums(c)
		if err != nil {
			log.Printf("%s: %v", c, err)
			continue
		}
		verifiedFiles += n
	}
	fmt.Printf("OK (%d)\n", verifiedFiles)
}
