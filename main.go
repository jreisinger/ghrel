package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jreisinger/ghrel/asset"
	"github.com/jreisinger/ghrel/checksum"
)

var l = flag.Bool("l", false, "just list assets")
var p = flag.String("p", "", "only assets matching shell `pattern` (doesn't apply to checksums files)")

func main() {
	flag.Usage = func() {
		desc := "Download and verify assets (files) of the latest release from a GitHub repository."
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n\n%s [flags] <owner>/<repo>\n", desc, os.Args[0])
		flag.PrintDefaults()
	}

	// Parse CLI arguments.
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	repo := flag.Args()[0]

	// Set CLI-style logging.
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	assets, err := asset.Get(repo, p)
	if err != nil {
		log.Fatalf("get release assets: %v", err)
	}

	if *l {
		asset.Print(assets)
		os.Exit(0)
	}

	// Download assets.
	var wg sync.WaitGroup
	for _, a := range assets {
		wg.Add(1)
		go func(a asset.Asset) {
			defer wg.Done()
			if err := a.Download(); err != nil {
				log.Printf("download %s: %v", a.BrowserDownloadUrl, err)
				return
			}
		}(a)
	}
	wg.Wait()

	// Print download statistics.
	nFiles, nCheckumsFiles := asset.Count(assets)
	fmt.Printf("downloaded\t%d (+ %d checksums files)\n", nFiles, nCheckumsFiles)

	checksums, err := checksum.Get(assets)
	if err != nil {
		log.Fatalf("get checksums: %v", err)
	}

	// Verify checksums.
	var verifiedFiles int
Asset:
	for _, a := range assets {
		if a.IsChecksumsFile {
			continue
		}
		for _, c := range checksums {
			if a.Name == c.Name {
				ok, err := c.Verify()
				if err != nil {
					log.Printf("verifying: %v", err)
					continue
				}
				if ok {
					verifiedFiles++
				} else {
					log.Printf("%s not verified (bad checksum)", a.Name)
				}
				continue Asset
			}
		}
		log.Printf("%s not verified (no checksum)", a.Name)
	}
	fmt.Printf("verified\t%d\n", verifiedFiles)
}
