package download

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"regexp"

	"naksu/xlate"
)

// CleanSHA256ChecksumString takes sha256 checksum file content as a string and
// returns the hash part (the first 64 "word" characters). If there is no valid
// match, an empty string is returned
func CleanSHA256ChecksumString(checksumString string) string {
	re := regexp.MustCompile(`^[0-9a-f]{64}`)

	return re.FindString(checksumString)
}

// GetSHA256ChecksumFromFile reads given file, calculates it SHA256 hash
// and returns it as a string
func GetSHA256ChecksumFromFile(filePath string, progressCallbackFn func(string, int)) (string, error) {
	checksumFile, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error while opening file to calculate sha256 from '%s': %w", filePath, err)
	}
	defer checksumFile.Close()

	// To get the file size
	fileStat, err := checksumFile.Stat()
	if err != nil {
		return "", fmt.Errorf("error while trying to get file info for '%s': %w", filePath, err)
	}

	progressCallbackFn(xlate.Get("Checking disk image..."), 0)

	counter := &writeCounter{
		ProgressCallbackFn: progressCallbackFn,
		FileSize:           uint64(fileStat.Size()),
		ProgressString:     xlate.GetRaw("Checking disk image..."),
		Total:              0,
	}
	checksumCalculator := sha256.New()
	if _, err = io.Copy(checksumCalculator, io.TeeReader(checksumFile, counter)); err != nil {
		return "", fmt.Errorf("error while reading file to calculate sha256 from '%s': %w", filePath, err)
	}

	const progressCheckFinished = 100
	progressCallbackFn(xlate.Get("Disk image checked"), progressCheckFinished)

	return fmt.Sprintf("%x", checksumCalculator.Sum(nil)), nil
}
