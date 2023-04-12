package checksums

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"
)

// from `sha256sum testdata/file-to-checksum`
const helloWorld = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

func TestVerify(t *testing.T) {
	cs, err := hex.DecodeString(helloWorld)
	if err != nil {
		t.Errorf("decode hex string: %v", err)
	}
	checksum := Checksum{Checksum: cs, Name: "testdata/file-to-checksum"}
	ok, err := checksum.Verify()
	if err != nil {
		t.Errorf("Verify failed: %v", err)
	}
	if !ok {
		t.Errorf("Verify %v, want true", ok)
	}
}

func TestGetChecksums(t *testing.T) {
	file, err := os.Open("testdata/checksum-file")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	cs, err := parseChecksumsLines(b)
	if err != nil {
		t.Errorf("get checksums: %v", err)
	}
	for _, c := range cs {
		got := fmt.Sprintf("%x", c.Checksum)
		want := helloWorld
		if got != want {
			t.Errorf("wrong checksum: got %s, want %s", got, want)
		}
	}
}
