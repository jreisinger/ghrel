package checksum

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Sha256 calculates SHA256 checksum of filename. Checksum is in hex format.
func Sha256(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum), nil
}

// Pair represents a line from a checksum file. The line looks like:
//
//	ba47c83b6038dda089dd1410b9e97d1de7e4adea7620c856f9b74a782048e272  checkip_0.45.1_linux_amd64.tar.gz
type Pair struct {
	Checksum string // in hex
	Filename string
}

// Parse parses a checksum file into checksum/filename pairs.
func Parse(checksumFile string) ([]Pair, error) {
	var checksums []Pair

	file, err := os.Open(checksumFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cs := parseChecksumLines(b)
	if len(cs) == 0 {
		return checksums, fmt.Errorf("no checksums in %s", checksumFile)
	}
	for i := range cs {
		// Since there is no filename field in the checksumFile get the
		// filename from the checksumFile name.
		if cs[i].Filename == "" {
			cs[i].Filename = trimSuffix(checksumFile)
		}
	}
	checksums = append(checksums, cs...)

	return checksums, nil
}

// trimSuffix removes suffix (like .sha256) from the filename.
func trimSuffix(filename string) string {
	suffix := filepath.Ext(filename)
	return strings.TrimSuffix(filename, suffix)
}

// parseChecksumLines parses bytes from a checksums file. The bytes look like:
// 712f37d14687e10ae0425bf7e5d0faf17c49f9476b8bb6a96f2a3f91b0436db2  checkip_0.45.1_linux_armv6.tar.gz
// ba47c83b6038dda089dd1410b9e97d1de7e4adea7620c856f9b74a782048e272  checkip_0.45.1_linux_amd64.tar.gz
func parseChecksumLines(b []byte) []Pair {
	var checksums []Pair
	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		switch len(fields) {
		case 2:
			checksums = append(checksums, Pair{
				Checksum: fields[0],
				Filename: fields[1],
			})
		case 1:
			checksums = append(checksums, Pair{
				Checksum: fields[0],
			})
		}
	}
	return checksums
}
