package checksum

import (
	"io"
	"os"
	"testing"
)

var test = struct {
	fileName      string
	fileChecksum  string // calculated using sha256sum
	checksumsFile string // contains checksum of file and file name
}{
	fileName:      "testdata/file-to-checksum",
	fileChecksum:  "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
	checksumsFile: "testdata/checksum-file", // only one line here
}

func TestParseChecksumsLines(t *testing.T) {
	file, err := os.Open(test.checksumsFile)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	cs := parseChecksumLines(b)
	for _, c := range cs {
		got := c.Checksum
		want := test.fileChecksum
		if got != want {
			t.Errorf("wrong checksum: got %s, want %s", got, want)
		}
	}
}
