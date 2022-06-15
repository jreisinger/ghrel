package assets

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDownloadUrls(t *testing.T) {
	t.Run("given valid response, download URLs and no error is returned", func(t *testing.T) {
		testUrl := startMockApiServer(t)
		setMockApiUrl(t, testUrl)

		urls, err := GetDownloadUrls("jreisinger/checkip")
		require.NoError(t, err)
		assert.Equal(t, 6, len(urls)) // number of release assets for jreisinger/checkip repo
	})

}

// --- test helpers ---

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
