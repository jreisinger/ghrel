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

var (
	l = flag.Bool("l", false, "just list assets, don't download them")
	p = flag.String("p", "", "assets matching shell `pattern` (doesn't apply to checksum files)")
	v = flag.Bool("v", false, "be verbose")
)

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
		if *v {
			asset.Table(assets)
		} else {
			for _, a := range assets {
				fmt.Println(a.Name)
			}
		}
		os.Exit(0)
	}

	// Download and checksum assets.
	var wg sync.WaitGroup
	for i := range assets {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			if err := asset.Download(assets[i]); err != nil {
				log.Printf("download asset %v", err)
				return
			}

			sum, err := checksum.Sha256(assets[i].Name)
			if err != nil {
				log.Printf("checksum asset: %v", err)
				return
			}
			assets[i].Checksum = sum
		}(i)
	}
	wg.Wait()

	// Print download statistics.
	if *v {
		nFiles, nCheckumsFiles := asset.Count(assets)
		fmt.Printf("downloaded\t%d + %d checksum file(s)\n", nFiles, nCheckumsFiles)
	}

	var pairs []checksum.Pair

	// Get checksum/filename pairs from assets that are checksum files.
	for _, a := range assets {
		if !a.IsChecksumFile {
			continue
		}

		cs, err := checksum.Parse(a.Name)
		if err != nil {
			log.Printf("get checksums from %s: %v", a.Name, err)
			continue
		}
		pairs = append(pairs, cs...)
	}

	// Verify checksums.
	var verifiedFiles int
Asset:
	for _, a := range assets {
		if a.IsChecksumFile {
			continue
		}

		for _, c := range pairs {
			if a.Name == c.Filename {
				if a.Checksum == c.Checksum {
					verifiedFiles++
				} else {
					log.Printf("%s not verified, has bad checksum %s", a.Name, c.Checksum)
				}
				continue Asset
			}
		}
		log.Printf("%s not verified, has no checksum", a.Name)
	}
	if *v {
		fmt.Printf("verified\t%d\n", verifiedFiles)
	}
}
