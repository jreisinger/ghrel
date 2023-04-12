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

func TestGet(t *testing.T) {
	testUrl := startMockApiServer(t)
	setMockApiUrl(t, testUrl)
	t.Run("get assets of jreisinger/checkip repo", func(t *testing.T) {
		assets, err := Get("jreisinger/checkip", stringPointer(""))
		require.NoError(t, err)
		assert.Equal(t, 6, len(assets))
	})
	t.Run("get assets of jreisinger/checkip repo with pattern", func(t *testing.T) {
		assets, err := Get("jreisinger/checkip", stringPointer("*linux*"))
		require.NoError(t, err)
		assert.Equal(t, 3, len(assets))
	})
}

func TestCount(t *testing.T) {
	tests := []struct {
		assets         []Asset
		nFiles         int
		nChecksumFiles int
	}{
		{
			assets:         []Asset{},
			nFiles:         0,
			nChecksumFiles: 0,
		},
		{
			assets: []Asset{
				{IsChecksumFile: false},
				{IsChecksumFile: false},
			},
			nFiles:         2,
			nChecksumFiles: 0,
		},
		{
			assets: []Asset{
				{IsChecksumFile: false},
				{IsChecksumFile: true},
			},
			nFiles:         1,
			nChecksumFiles: 1,
		},
	}
	for _, test := range tests {
		nFiles, nChecksumFiles := Count(test.assets)
		if nFiles != test.nFiles {
			t.Errorf("nFiles = %d, want %d", nFiles, test.nFiles)
		}
		if nChecksumFiles != test.nChecksumFiles {
			t.Errorf("nChecksumFiles = %d, want %d", nChecksumFiles, test.nChecksumFiles)
		}
	}
}

func TestIsChecksumFile(t *testing.T) {
	tests := []struct {
		filename       string
		isChecksumFile bool
	}{
		{
			filename:       "",
			isChecksumFile: false,
		},
		{
			filename:       "ghrel_0.6.2_linux_armv6.tar.gz",
			isChecksumFile: false,
		},
		{
			filename:       "checksums.txt",
			isChecksumFile: true,
		},
		{
			filename:       "brave-v1.50.114-darwin-arm64-symbols.zip.sha256",
			isChecksumFile: true,
		},
	}
	for _, test := range tests {
		ok := isChecksumFile(test.filename)
		if ok != test.isChecksumFile {
			t.Errorf("isChecksumFile(%s) = %t, want %t", test.filename, ok, test.isChecksumFile)
		}
	}
}

// --- test helpers ---

func stringPointer(s string) *string {
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
