package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jreisinger/ghrel/asset"
	"github.com/jreisinger/ghrel/checksums"
)

var l = flag.Bool("l", false, "list assets")
var p = flag.String("p", "", "only assets matching shell `pattern` (doesn't apply to checksums files)")

func main() {
	flag.Usage = func() {
		desc := "Download or list assets (files) of the latest release from a GitHub repository."
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

	assets, err := asset.Get(repo, p)
	if err != nil {
		log.Fatalf("get download URLs for release assets: %v", err)
	}

	if *l {
		asset.Print(assets)
		os.Exit(0)
	}

	var wg sync.WaitGroup
	for _, a := range assets {
		wg.Add(1)
		go func(asset asset.Asset) {
			defer wg.Done()
			if err := asset.Download(); err != nil {
				log.Printf("download %s: %v", asset.BrowserDownloadUrl, err)
				return
			}
		}(a)
	}
	wg.Wait()

	var nFiles, nChecksumsFiles int
	for _, a := range assets {
		switch {
		case a.IsChecksumsFile:
			nChecksumsFiles++
		default:
			nFiles++
		}
	}
	fmt.Printf("downloaded\t%d (+ %d checksums files)\n", nFiles, nChecksumsFiles)

	var cs []checksums.Checksum
	for _, a := range assets {
		if a.IsChecksumsFile {
			c, err := checksums.Extract(a.Name)
			if err != nil {
				log.Fatalf("extracting checksums: %v", err)
			}
			cs = append(cs, c...)
		}
	}
	var verifiedFiles int
Asset:
	for _, a := range assets {
		if a.IsChecksumsFile {
			continue
		}
		for _, c := range cs {
			if a.Name == c.Filename {
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
