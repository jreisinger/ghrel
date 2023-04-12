package checksums

import (
	"encoding/hex"
	"testing"
)

func TestVerify(t *testing.T) {
	cs, err := hex.DecodeString("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9")
	if err != nil {
		t.Errorf("decode hex string: %v", err)
	}
	checksum := Checksum{Checksum: cs, Filename: "testdata/file-to-checksum"}
	ok, err := checksum.Verify()
	if err != nil {
		t.Errorf("Verify failed: %v", err)
	}
	if !ok {
		t.Errorf("Verify %v, want true", ok)
	}
}

func TestExtractAndVerify(t *testing.T) {
	checksums, err := Extract("testdata/checksum-ok.txt")
	if err != nil {
		t.Errorf("Extract failed: %v", err)
	}
	ok, err := checksums[0].Verify()
	if err != nil {
		t.Errorf("Verify failed: %v", err)
	}
	if !ok {
		t.Error("testdata/checksum-ok.txt not verified")
	}

	checksums, err = Extract("testdata/checksum-notok.txt")
	if err != nil {
		t.Errorf("Extract failed: %v", err)
	}
	ok, err = checksums[0].Verify()
	if err != nil {
		t.Errorf("Verify failed: %v", err)
	}
	if ok {
		t.Error("testdata/checksum-notok.txt verified")
	}
}
