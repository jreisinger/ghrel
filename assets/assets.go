package assets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var GitHubApiUrl = "https://api.github.com"

// GetDownloadUrls returns URLs for downloading assets from the latest repo
// release. Repo is a GitHub repository in the form <user>/<repo>.
func GetDownloadUrls(repo string) ([]string, error) {
	url := GitHubApiUrl + "/repos/" + repo + "/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := struct {
		Assets []struct {
			BrowserDownloadUrl string `json:"browser_download_url"`
		}
	}{}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	var urls []string
	for _, a := range r.Assets {
		urls = append(urls, a.BrowserDownloadUrl)
	}
	return urls, nil
}
