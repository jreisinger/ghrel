package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// verifyChecksums reads checksums from checksumsFile and checks them. It
// implements "sha256sum -c" functionality.
func verifyChecksums(checksumsFile string) (okfiles int, err error) {
	// cmd := exec.Command("sha256sum", "-c", checksumsFile)
	// return cmd.Run()

	checksums, err := extractChecksums(checksumsFile)
	if err != nil {
		return 0, err
	}

	for _, c := range checksums {
		if *shellPattern != "" {
			if matched, _ := filepath.Match(*shellPattern, c.filename); !matched {
				continue
			}
		}

		ok, err := verifyFileChecksum(c)
		if err != nil {
			return okfiles, err
		}
		if !ok {
			return okfiles, fmt.Errorf("%s is not ok", c.filename)
		}
		okfiles++
	}

	return okfiles, nil
}

// checksum represents a line from a checksums file.
type checksum struct {
	checksum []byte // in hex
	filename string
}

// extractChecksums parses a checksumsFile and extracts checksums from it.
func extractChecksums(checksumsFile string) ([]checksum, error) {
	file, err := os.Open(checksumsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var checksums []checksum

	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		c, err := hex.DecodeString(fields[0])
		if err != nil {
			return nil, err
		}
		checksums = append(checksums, checksum{c, fields[1]})
	}

	return checksums, nil
}

// verifyFileChecksum computes SHA256 checksum of the file from c and compares
// it to the checksum from c.
func verifyFileChecksum(c checksum) (ok bool, err error) {
	file, err := os.Open(c.filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}
	sum := hash.Sum(nil)

	return bytes.Equal(c.checksum, sum), nil
}

type checksumFilePattern struct {
	pattern    string
	singleFile bool // checksum file holds checksum only for single file
}

var checksumFilePatterns = []checksumFilePattern{
	// ghrel_0.3.1_checksums.txt
	{`(?i)check.?sum`, false},
	// brave-browser-nightly-1.47.27-linux-amd64.zip.sha256
	{`\.sha256$`, true},
}

// isChecksumsFile tells whether filename looks like a file containing checksums.
func isChecksumsFile(filename string) bool {
	for _, c := range checksumFilePatterns {
		if c.singleFile && *shellPattern != "" {
			if matched, _ := filepath.Match(*shellPattern, filename); !matched {
				continue
			}
		}
		if regexp.MustCompile(c.pattern).MatchString(filename) {
			return true
		}
	}
	return false
}
