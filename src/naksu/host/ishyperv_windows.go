package host

import (
	"fmt"

	"naksu/log"

	"github.com/StackExchange/wmi"
)

// IsHyperV returns true if Hyper-V Hearbeat Service is running
func IsHyperV() bool {
	// This struct must be named with an underscore, otherwise it is not recognised
	// and results "Invalid class" exception.
	type Win32_OptionalFeature struct { //nolint
		Caption      *string
		Name         *string
		InstallState uint32
	}

	var dst []Win32_OptionalFeature
	query := wmi.CreateQuery(&dst, "WHERE Name LIKE '%hyper-v%' AND InstallState=1")
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Debug("IsHyperV() could not query WMI")
		log.Debug(fmt.Sprint(err))
		return false
	}

	isRunning := false

	for thisService := range dst {
		if dst[thisService].InstallState == 1 {
			thisName := "N/A"
			thisCaption := "N/A"

			if dst[thisService].Name != nil {
				thisName = *dst[thisService].Name
			}

			if dst[thisService].Caption != nil {
				thisCaption = *dst[thisService].Caption
			}

			log.Debug(fmt.Sprintf("Windows Hyper-V Optional Feature found: %s (%s)", thisName, thisCaption))
			isRunning = true
		}
	}

	return isRunning
}
