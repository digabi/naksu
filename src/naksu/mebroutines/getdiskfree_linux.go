package mebroutines

import (
	"fmt"
	"regexp"
	"strconv"
)

// GetDiskFree returns available disk space
func GetDiskFree(path string) (int, error) {
	runParams := []string{"df", "--output=avail", path}

	output, err := RunAndGetOutput(runParams)

	if err != nil {
		return -1, err
	}

	// Extract server disk image path
	pattern := regexp.MustCompile(`(\d+)`)
	result := pattern.FindStringSubmatch(output)

	if len(result) > 1 {
		resultFloat, err := strconv.ParseFloat(result[1], 10)
		if err != nil {
			LogDebug(fmt.Sprintf("result parsing failed for string %s", result[1]))
		}
		LogDebug(fmt.Sprintf("Disk free for path %s: %d", path, int(resultFloat)))
		return int(resultFloat), nil
	}

	return -1, nil
}
