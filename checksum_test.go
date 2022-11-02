package main

import (
	"encoding/hex"
	"testing"
)

func TestVerifyFileChecksum(t *testing.T) {
	cs, err := hex.DecodeString("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9")
	if err != nil {
		t.Errorf("decode hex string: %v", err)
	}
	ok, err := verifyFileChecksum(checksum{cs, "testdata/file-to-checksum", false})
	if err != nil {
		t.Errorf("verifyFileChecksum: %v, want nil", err)
	}
	if !ok {
		t.Errorf("verifyFileChecksum: %v, want true", ok)
	}
}

func TestVerifyChecksums(t *testing.T) {
	checksums, err := verifyChecksums("testdata/checksum-ok.txt")
	if err != nil {
		t.Errorf("verifyChecksums: %v, want nil", err)
	}
	if !checksums[0].verified {
		t.Error("testdata/checksum-ok.txt not verified")
	}

	checksums, err = verifyChecksums("testdata/checksum-notok.txt")
	if err != nil {
		t.Errorf("verifyChecksums: %v, want nil", err)
	}
	if checksums[0].verified {
		t.Error("testdata/checksum-notok.txt verified")
	}
}
