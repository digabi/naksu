package mebroutines

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/StackExchange/wmi"
)

// GetDiskFree returns free disk space amount
func GetDiskFree(path string) (int, error) {
	patternDisk := regexp.MustCompile(`^(\w\:)`)
	patternResult := patternDisk.FindStringSubmatch(path)

	// This struct must be named with an underscore, otherwise it is not recognised
	// and results "Invalid class" exception.
	type Win32_LogicalDisk struct {
		Size      int
		FreeSpace int
		DeviceID  string
	}

	if len(patternResult) < 2 {
		LogDebug(fmt.Sprintf("Could not detect drive letter from path: %s", path))

		return -1, errors.New("could not detect drive letter")
	}

	diskletter := patternResult[1]
	wmiQuery := fmt.Sprintf("WHERE DeviceID=\"%s\"", diskletter)

	var dst []Win32_LogicalDisk
	/* #nosec */
	query := wmi.CreateQuery(&dst, wmiQuery)
	err := wmi.Query(query, &dst)
	if err != nil {
		LogDebug(fmt.Sprintf("GetDiskFree() could not make WMI query (%s): %s", wmiQuery, fmt.Sprint(err)))
		return -1, errors.New("could not detect free disk size as it could not query wmi")
	}

	if len(dst) > 0 {
		freeSpace := dst[0].FreeSpace / 1000
		LogDebug(fmt.Sprintf("Disk free for path %s: %d", path, freeSpace))
		return freeSpace, nil
	}

	return -1, errors.New("could not detect free disk size")
}
