package asset

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/tabwriter"
	"time"
)

var GitHubApiUrl = "https://api.github.com" // changed during testing

// Asset represents a GitHub API release asset with some calculated fields.
type Asset struct {
	BrowserDownloadUrl string    `json:"browser_download_url"`
	Name               string    `json:"name"` // filename
	IsChecksumFile     bool      `json:"-"`
	Checksum           string    `json:"-"` // in hex
	UpdatedAt          time.Time `json:"updated_at"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
}

// Get queries GitHub API for assets of the latest repo release whose name
// matches the shell pattern, if the pattern is not empty. Pattern matching does
// not apply to checksum files. Repo is a GitHub repository in the form
// <user>/<repo>.
func Get(repo string, pattern *string) ([]Asset, error) {
	url := GitHubApiUrl + "/repos/" + repo + "/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s %s", url, resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	api := struct{ Assets []Asset }{}
	if err := json.Unmarshal(b, &api); err != nil {
		return nil, err
	}

	var assets []Asset

	for _, a := range api.Assets {
		if isChecksumFile(a.Name) {
			a.IsChecksumFile = true
		} else if *pattern != "" {
			if matched, _ := filepath.Match(*pattern, a.Name); !matched {
				continue
			}
		}
		assets = append(assets, a)
	}

	return assets, nil
}

// Download downloads asset from a.BrowserDownloadUrl to a file named a.Name. It
// creates or truncates the file.
func Download(a Asset) error {
	f, err := os.Create(a.Name)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(a.BrowserDownloadUrl)
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

// Print prints table of assets.
func Print(assets []Asset) {
	const format = "%v\t%v\t%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Asset", "Checksum file", "Updated", "Size", "Download count")
	fmt.Fprintf(tw, format, "-----", "-------------", "-------", "----", "--------------")
	for _, a := range assets {
		fmt.Fprintf(tw, format, a.Name, a.IsChecksumFile, a.UpdatedAt.Format("2006-01-02"), a.Size, a.DownloadCount)
	}
	tw.Flush()
}

// Count counts files and checksum files.
func Count(assets []Asset) (nFiles, nChecksumFiles int) {
	for _, a := range assets {
		switch {
		case a.IsChecksumFile:
			nChecksumFiles++
		default:
			nFiles++
		}
	}
	return
}

// isChecksumFile tells whether asset looks like a file containing checksum(s).
func isChecksumFile(filename string) bool {
	checksumFiles := []string{
		`(?i)check.?sum`, // ghrel_0.3.1_checksums.txt
		`\.sha256$`,      // brave-browser-nightly-1.47.27-linux-amd64.zip.sha256
	}

	var checksumFilePatterns []*regexp.Regexp

	for _, f := range checksumFiles {
		re := regexp.MustCompile(f)
		checksumFilePatterns = append(checksumFilePatterns, re)
	}

	for _, re := range checksumFilePatterns {
		if re.MatchString(filename) {
			return true
		}
	}
	return false
}
