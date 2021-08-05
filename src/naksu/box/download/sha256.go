package download

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"regexp"
)

// CleanSHA256ChecksumString takes sha256 checksum file content as a string and
// returns the hash part (the first 64 "word" characters). If there is no valid
// match, an empty string is returned
func CleanSHA256ChecksumString(checksumString string) string {
	re := regexp.MustCompile(`^\w{64}`)
	return re.FindString(checksumString)
}

// GetSHA256ChecksumFromFile reads given file, calculates it SHA256 hash
// and returns it as a string
func GetSHA256ChecksumFromFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error while opening file to calculate sha256 from '%s': %v", filePath, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("error while reading file to calculate sha256 from '%s': %v", filePath, err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
