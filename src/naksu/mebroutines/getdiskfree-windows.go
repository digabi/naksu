//go:build windows
// +build windows

package mebroutines

import (
	"errors"
	"fmt"
	"regexp"

	"naksu/log"

	"github.com/yusufpapurcu/wmi"
)

func getDiskFreeWindows(path string) (uint64, error) {
	patternDisk := regexp.MustCompile(`^(\w\:)`)
	patternResult := patternDisk.FindStringSubmatch(path)

	const minimumDriveLetterLength = 2
	if len(patternResult) < minimumDriveLetterLength {
		log.Error("Could not detect drive letter from path: %s", path)

		return 0, errors.New("could not detect drive letter")
	}

	diskletter := patternResult[1]
	// gosec complains here "SQL string formatting" but this can be safely turned off
	/* #nosec */
	wmiQuery := fmt.Sprintf(`WHERE DeviceID="%s"`, diskletter)

	result := make(chan []Win32_LogicalDisk)

	// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
	go func() {
		var dst []Win32_LogicalDisk
		/* #nosec */
		query := wmi.CreateQuery(&dst, wmiQuery)
		err := wmi.Query(query, &dst)
		if err != nil {
			log.Error("getDiskFreeWindows() could not make WMI query (%s): %s", wmiQuery, fmt.Sprint(err))
			result <- []Win32_LogicalDisk{}
		} else {
			result <- dst
		}
	}()

	return ExtractDiskFreeWindows(<-result)
}
