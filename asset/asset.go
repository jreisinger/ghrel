package asset

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"text/tabwriter"
	"time"
)

var GitHubApiUrl = "https://api.github.com" // changed during testing

// Asset represents a GitHub API release asset.
type Asset struct {
	BrowserDownloadUrl string    `json:"browser_download_url"`
	Name               string    `json:"name"` // filename
	IsChecksumsFile    bool      `json:"-"`
	UpdatedAt          time.Time `json:"updated_at"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
}

// Get queries GitHub API for assets of the latest repo release whose
// (file)names match the shell pattern, if that is not empty. Pattern matching
// does not apply to checksums files. Repo is a GitHub repository in the form
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
		if isChecksumsFile(a.Name) {
			a.IsChecksumsFile = true
		} else if *pattern != "" {
			if matched, _ := filepath.Match(*pattern, a.Name); !matched {
				continue
			}
		}
		assets = append(assets, a)
	}

	return assets, nil
}

// filename extracts filename from the Asset's BrowserDownloadUrl.
func (a Asset) filename() (string, error) {
	u, err := url.Parse(a.BrowserDownloadUrl)
	if err != nil {
		return "", err
	}
	_, file := path.Split(u.Path)
	return file, nil
}

// Download downloads asset from BrowserDownloadUrl to a file who's name is
// extracted from BrowserDownloadUrl. It creates or truncates the file.
func (a Asset) Download() error {
	file, err := a.filename()
	if err != nil {
		return err
	}
	f, err := os.Create(file)
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
	fmt.Fprintf(tw, format, "Asset", "Checksums file", "Updated", "Size", "Download count")
	fmt.Fprintf(tw, format, "-----", "--------------", "-------", "----", "--------------")
	for _, a := range assets {
		fmt.Fprintf(tw, format, a.Name, a.IsChecksumsFile, a.UpdatedAt.Format("2006-01-02"), a.Size, a.DownloadCount)
	}
	tw.Flush()
}

func Count(assets []Asset) (nFiles, nChecksumsFiles int) {
	for _, a := range assets {
		switch {
		case a.IsChecksumsFile:
			nChecksumsFiles++
		default:
			nFiles++
		}
	}
	return
}

var checksumsFilePatterns = []string{
	`(?i)check.?sum`, // ghrel_0.3.1_checksums.txt
	`\.sha256$`,      // brave-browser-nightly-1.47.27-linux-amd64.zip.sha256
}

// isChecksumsFile tells whether asset looks like a file containing checksum(s).
func isChecksumsFile(filename string) bool {
	for _, c := range checksumsFilePatterns {
		if regexp.MustCompile(c).MatchString(filename) {
			return true
		}
	}
	return false
}
