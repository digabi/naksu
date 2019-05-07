package mebroutines

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strconv"

	"naksu/log"
)

// Win32_LogicalDisk is a struct used to query Windows WMI
// (Windows Management Instrumentation)
// The struct must be named with an underscore, otherwise it is not recognised
// and results "Invalid class" exception.
type Win32_LogicalDisk struct {
	Size      uint64
	FreeSpace uint64
	DeviceID  string
}

// GetDiskFree returns available disk space
func GetDiskFree(path string) (uint64, error) {
	var diskFree uint64
	var diskError error

	switch runtime.GOOS {
	case "darwin":
		diskFree, diskError = getDiskFreeDarwin(path)
	case "linux":
		diskFree, diskError = getDiskFreeLinux(path)
	case "windows":
		diskFree, diskError = getDiskFreeWindows(path)
	default:
		diskFree = 0
		diskError = errors.New("unknown execution environment")
	}

	return diskFree, diskError
}

func getDiskFreeDarwin(path string) (uint64, error) {
	// FIXME: By default Darwin df gives number of 512-byte blocks
	// GetDiskFree should return number of free disk in BYTES
	runParams := []string{"df", path}

	output, err := RunAndGetOutput(runParams, true)

	if err != nil {
		return 0, err
	}

	return ExtractDiskFreeDarwin(output)
}

// ExtractDiskFreeDarwin extracts free disk space from a given df output
func ExtractDiskFreeDarwin(dfOutput string) (uint64, error) {
	// Extract server disk image path using tail-hooked regexp
	pattern := regexp.MustCompile(`([0-9]+)\s+[0-9]+%\s+[0-9]+\s+[0-9]+\s+[0-9]+%\s+[a-zA-Z0-9/]+$`)
	result := pattern.FindStringSubmatch(dfOutput)

	if len(result) > 1 {
		floatResult, err := strconv.ParseFloat(result[1], 64)
		if err == nil {
			// df gives available disk space in 1K blocks
			intResult := uint64(floatResult) * 512
			log.Debug(fmt.Sprintf("ExtractDiskFreeDarwin: %d", intResult))
			return intResult, nil
		}
	}

	log.Debug("ExtractDiskFreeDarwin failed to parse df output")
	log.Debug(dfOutput)

	return 0, errors.New("could not extract free disk size from darwin df output")
}

func getDiskFreeLinux(path string) (uint64, error) {
	runParams := []string{"df", "--output=avail", path}

	output, err := RunAndGetOutput(runParams, true)

	if err != nil {
		return 0, err
	}

	return ExtractDiskFreeLinux(output)
}

// ExtractDiskFreeLinux extracts free disk space from a given df output
func ExtractDiskFreeLinux(dfOutput string) (uint64, error) {
	// Extract server disk image path
	pattern := regexp.MustCompile(`(\d+)`)
	result := pattern.FindStringSubmatch(dfOutput)

	if len(result) > 1 {
		floatResult, err := strconv.ParseFloat(result[1], 64)
		if err == nil {
			// df gives available disk space in 1K blocks
			intResult := uint64(floatResult) * 1024
			log.Debug(fmt.Sprintf("ExtractDiskFreeLinux: %d", intResult))
			return intResult, nil
		}
	}

	log.Debug("ExtractDiskFreeLinux failed to parse df output")
	log.Debug(dfOutput)

	return 0, errors.New("could not extract free disk size from linux df output")
}

// ExtractDiskFreeWindows extracts free disk space from a given WMI query result slice
func ExtractDiskFreeWindows(wmiData []Win32_LogicalDisk) (uint64, error) {
	if len(wmiData) > 0 {
		freeSpace := wmiData[0].FreeSpace
		log.Debug(fmt.Sprintf("ExtractDiskFreeWindows: %d", freeSpace))
		return freeSpace, nil
	}

	return 0, errors.New("could not extract free disk size from wmi data")
}
