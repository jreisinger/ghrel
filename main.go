package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"text/tabwriter"

	"github.com/jreisinger/ghrel/asset"
)

var pattern = flag.String("p", "", "only assets matching shell `pattern`")
var list = flag.Bool("l", false, "list assets")

func main() {
	flag.Usage = func() {
		desc := "List or download assets (files) of the latest release from a GitHub repository."
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

	assets, err := asset.GetDownloadUrls(repo)
	if err != nil {
		log.Fatalf("get donwload URLs for release assets: %v", err)
	}

	if *list {
		const format = "%v\t%v\t%v\t%v\n"
		tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintf(tw, format, "Asset", "Updated", "Size", "Download count")
		fmt.Fprintf(tw, format, "-----", "-------", "----", "--------------")
		for _, a := range assets {
			fn, _ := fileName(a.BrowserDownloadUrl)
			if *pattern != "" {
				if matched, _ := filepath.Match(*pattern, fn); !matched {
					continue
				}
			}
			fmt.Fprintf(tw, format, fn, a.UpdatedAt.Format("2006-01-02"), a.Size, a.DownloadCount)
		}
		tw.Flush()
		return
	}

	var checksumsFiles []string
	var count struct {
		mu    sync.Mutex
		files int
	}

	var wg sync.WaitGroup
	for _, a := range assets {
		wg.Add(1)
		go func(asset asset.Asset) {
			defer wg.Done()

			file, err := fileName(asset.BrowserDownloadUrl)
			if err != nil {
				log.Printf("extract filename from URL: %v", err)
				return
			}

			if isChecksumsFile(file) {
				checksumsFiles = append(checksumsFiles, file)
			} else if *pattern != "" {
				matched, err := filepath.Match(*pattern, file)
				if err != nil {
					log.Fatal(err)
				}
				if !matched {
					return
				}
			}

			if err := download(asset.BrowserDownloadUrl, file); err != nil {
				log.Printf("download %s: %v", asset.BrowserDownloadUrl, err)
				return
			}

			count.mu.Lock()
			count.files++
			count.mu.Unlock()
		}(a)
	}
	wg.Wait()

	fmt.Printf("downloaded\t%d (+ %d checksums file)\n",
		count.files-len(checksumsFiles), len(checksumsFiles))

	var verifiedFiles int
	for _, c := range checksumsFiles {
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
	fmt.Printf("verified\t%d\n", verifiedFiles)
}
