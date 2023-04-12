package checksum

import (
	"testing"
)

var tests = []struct {
	fileName     string
	fileChecksum string // calculated using sha256sum
	checksumFile string // contains checksum of file and optionally file name
}{
	{
		fileName:     "testdata/file-to-checksum",
		fileChecksum: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		checksumFile: "testdata/checksum-file",
	},
	{
		fileName:     "testdata/file-to-checksum",
		fileChecksum: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		checksumFile: "testdata/file-to-checksum.sha256",
	},
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		checksums, err := Parse(test.checksumFile)
		if err != nil {
			t.Error(err)
		}

		for _, c := range checksums {
			got := c.Checksum
			want := test.fileChecksum
			if got != want {
				t.Errorf("wrong checksum: got %q, want %q", got, want)
			}

			got = c.Filename
			want = test.fileName
			if got != want {
				t.Errorf("wrong filename: got %q, want %q", got, want)
			}
		}
	}
}

func TestTrimSuffix(t *testing.T) {
	tests := []struct {
		filename string
		want     string
	}{
		{"", ""},
		{"starship-aarch64-apple-darwin.tar.gz.sha256", "starship-aarch64-apple-darwin.tar.gz"},
	}
	for _, test := range tests {
		got := trimSuffix(test.filename)
		if got != test.want {
			t.Errorf("trimSuffix(%s) = %s, want %s", test.filename, got, test.want)
		}
	}
}
