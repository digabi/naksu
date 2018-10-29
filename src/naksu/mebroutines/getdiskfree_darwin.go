package mebroutines

import (
	"fmt"
	"regexp"
	"strconv"
)

// GetDiskFree returns disk space usage as float
func GetDiskFree(path string) (int, error) {
	runParams := []string{"df", "--output=avail", path}

	output, err := RunAndGetOutput(runParams)

	if err != nil {
		return -1, err
	}

	// Extract server disk image path
	pattern := regexp.MustCompile("(\\d+)")
	result := pattern.FindStringSubmatch(output)

	if len(result) > 1 {
		floatResult, _ := strconv.ParseFloat(result[1], 10)
		LogDebug(fmt.Sprintf("Disk free for path %s: %d", path, int(floatResult)))
		return int(floatResult), nil
	}

	return -1, nil
}
