package checksums

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

// Checksum represents a line from a checksums file.
type Checksum struct {
	Checksum []byte // in hex
	Filename string
}

// Extract parses a checksumsFile and extracts checksums from it.
func Extract(checksumsFile string) ([]Checksum, error) {
	file, err := os.Open(checksumsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var checksums []Checksum
	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		c, err := hex.DecodeString(fields[0])
		if err != nil {
			return nil, err
		}
		checksums = append(checksums, Checksum{
			Checksum: c,
			Filename: fields[1],
		})
	}
	return checksums, nil
}

// Verify computes SHA256 checksum of the filename from c and compares it to the
// checksum from c.
func (c Checksum) Verify() (ok bool, err error) {
	file, err := os.Open(c.Filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}
	sum := hash.Sum(nil)

	return bytes.Equal(c.Checksum, sum), nil
}
