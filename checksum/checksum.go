package checksum

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/jreisinger/ghrel/asset"
)

// Checksum represents a line from a checksums file. The line looks like:
// ba47c83b6038dda089dd1410b9e97d1de7e4adea7620c856f9b74a782048e272  checkip_0.45.1_linux_amd64.tar.gz
type Checksum struct {
	Checksum []byte // in hex
	Name     string // filename
}

// Get extracts checksums from the checksums files.
func Get(assets []asset.Asset) ([]Checksum, error) {
	var checksums []Checksum
	for _, a := range assets {
		if !a.IsChecksumsFile {
			continue
		}

		file, err := os.Open(a.Name)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		cs, err := parseChecksumsLines(b)
		if err != nil {
			return nil, err
		}
		checksums = append(checksums, cs...)
	}
	return checksums, nil
}

// parseChecksumsLines parses bytes from a checksums file that look like:
// 712f37d14687e10ae0425bf7e5d0faf17c49f9476b8bb6a96f2a3f91b0436db2  checkip_0.45.1_linux_armv6.tar.gz
// ba47c83b6038dda089dd1410b9e97d1de7e4adea7620c856f9b74a782048e272  checkip_0.45.1_linux_amd64.tar.gz
func parseChecksumsLines(b []byte) ([]Checksum, error) {
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
			Name:     fields[1],
		})
	}
	return checksums, nil
}

// Verify computes SHA256 checksum of the filename from c and compares it to the
// checksum from c.
func (c Checksum) Verify() (ok bool, err error) {
	file, err := os.Open(c.Name)
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
