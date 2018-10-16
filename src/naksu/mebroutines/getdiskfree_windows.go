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

	type Win32LogicalDisk struct {
		Size      int
		FreeSpace int
		DeviceID  string
	}

	if len(patternResult) < 2 {
		LogDebug(fmt.Sprintf("Could not detect drive letter from path: %s", path))

		return -1, errors.New("Could not detect drive letter")
	}

	diskletter := patternResult[1]

	var dst []Win32LogicalDisk
	/* #nosec */
	query := wmi.CreateQuery(&dst, fmt.Sprintf("WHERE DeviceID=\"%s\"", diskletter))
	err := wmi.Query(query, &dst)
	if err != nil {
		LogDebug(fmt.Sprintf("Get_disk_free() could not make WMI query: %s", fmt.Sprint(err)))
		return -1, errors.New("Get_disk_free() could not detect free disk size as it could not query WMI")
	}

	if len(dst) > 0 {
		freeSpace := dst[0].FreeSpace / 1000
		LogDebug(fmt.Sprintf("Disk free for path %s: %d", path, freeSpace))
		return freeSpace, nil
	}

	return -1, errors.New("Get_disk_free() could not detect free disk size")
}
