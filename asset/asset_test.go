package asset

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilename(t *testing.T) {
	tests := []struct {
		asset Asset
		want  string
	}{
		{Asset{}, ""},
		{Asset{BrowserDownloadUrl: "httpx://wrongurl.com"}, ""},
		{Asset{BrowserDownloadUrl: "https://github.com/jgm/pandoc/releases/download/2.18/pandoc-2.18-1-amd64.deb"}, "pandoc-2.18-1-amd64.deb"},
	}
	for _, test := range tests {
		got, err := test.asset.filename()
		if err != nil {
			t.Error(err)
		}
		if got != test.want {
			t.Errorf("got %s, want %s", got, test.want)
		}
	}
}

func TestGetDownloadUrls(t *testing.T) {
	t.Run("given valid response, download URLs and no error is returned", func(t *testing.T) {
		testUrl := startMockApiServer(t)
		setMockApiUrl(t, testUrl)

		urls, err := Get("jreisinger/checkip", getPointer(""))
		require.NoError(t, err)
		assert.Equal(t, 6, len(urls)) // number of release assets for jreisinger/checkip repo
	})

}

// --- test helpers ---

func getPointer(s string) *string {
	return &s
}

// startMockApiServer returns URL of an HTTP server that mocks GitHub API.
// Getting that URL you will return testdata/gh_api_response.json.
func startMockApiServer(t *testing.T) string {
	b, err := ioutil.ReadFile(filepath.Join("testdata", "gh_api_response.json"))
	require.NoError(t, err)

	handlerFn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
	server := httptest.NewServer(handlerFn)
	t.Cleanup(func() {
		server.Close()
	})
	return server.URL
}

// setMockApiUrl temporarily sets api_url variable to testUrl.
func setMockApiUrl(t *testing.T, testUrl string) {
	origUrl := GitHubApiUrl
	GitHubApiUrl = testUrl
	t.Cleanup(func() {
		GitHubApiUrl = origUrl
	})
}
