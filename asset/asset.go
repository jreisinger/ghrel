package asset

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var GitHubApiUrl = "https://api.github.com"

type Asset struct {
	BrowserDownloadUrl string    `json:"browser_download_url"`
	UpdatedAt          time.Time `json:"updated_at"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
}

// GetDownloadUrls returns URLs for downloading assets from the latest repo
// release. Repo is a GitHub repository in the form <user>/<repo>.
func GetDownloadUrls(repo string) ([]Asset, error) {
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
	r := struct{ Assets []Asset }{}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return r.Assets, nil
}
