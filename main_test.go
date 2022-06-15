package main

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyFileChecksum(t *testing.T) {
	cs, err := hex.DecodeString("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9")
	if err != nil {
		t.Errorf("decode hex string: %v", err)
	}
	ok, err := verifyFileChecksum(checksum{cs, "testdata/file-to-checksum"})
	if err != nil {
		t.Errorf("verifyFileChecksum: %v, want nil", err)
	}
	if !ok {
		t.Errorf("verifyFileChecksum: %v, want true", ok)
	}
}

func TestVerifyChecksums(t *testing.T) {
	ok, err := verifyChecksums("testdata/checksum-ok.txt")
	if err != nil {
		t.Errorf("verifyChecksums: %v, want nil", err)
	}
	if !ok {
		t.Errorf("verifyChecksums: %v, want true", ok)
	}
	ok, err = verifyChecksums("testdata/checksum-notok.txt")
	if err != nil {
		t.Errorf("verifyChecksums: %v, want nil", err)
	}
	if ok {
		t.Errorf("verifyChecksums: %v, want false", ok)
	}
}

func TestFilename(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"httpx://wrongurl.com", ""},
		{"", ""},
		{"https://github.com/jgm/pandoc/releases/download/2.18/pandoc-2.18-1-amd64.deb", "pandoc-2.18-1-amd64.deb"},
	}
	for _, test := range tests {
		got, err := fileName(test.url)
		if err != nil {
			t.Error(err)
		}
		if got != test.want {
			t.Errorf("got %s, want %s", got, test.want)
		}
	}
}

func TestGetUrls(t *testing.T) {
	t.Run("given valid response, download URLs and no error is returned", func(t *testing.T) {
		testUrl := startMockApiServer(t)
		setMockApiUrl(t, testUrl)

		urls, err := getDownloadUrls("jreisinger/checkip")
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
	origUrl := api_url
	api_url = testUrl
	t.Cleanup(func() {
		api_url = origUrl
	})
}
