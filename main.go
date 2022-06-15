package main

import (
	"fmt"
	"log"
	"os"
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
