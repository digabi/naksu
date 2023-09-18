package host

import (
	"fmt"

	"naksu/log"

	"github.com/yusufpapurcu/wmi"
)

// isHyperVOptionalFeature returns true if Hyper-V features are installed
// Currently this is used only for logging
func isHyperVOptionalFeature() bool {
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
		log.Debug("isHyperVOptionalFeature() could not query WMI")
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

			log.Debug("Windows Hyper-V Optional Feature found: %s (%s)", thisName, thisCaption)
			// We're not returning this value as there might be a number of Hyper-V -related features found
			isRunning = true
		}
	}

	return isRunning
}

// isHypervisorPresent returns true if Win32_ComputerSystem.HypervisorPresent is true
// it is true at least when Windows Hypervisor feature or Windows Virtualization Platform feature is on
func isHypervisorPresent() bool {
	// This struct must be named with an underscore, otherwise it is not recognised
	// and results "Invalid class" exception.
	type Win32_ComputerSystem struct { //nolint
		HypervisorPresent bool
	}

	var dst []Win32_ComputerSystem
	query := wmi.CreateQuery(&dst, "")
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Debug("isHypervisorPresent() could not query WMI")
		log.Debug(fmt.Sprint(err))

		return false
	}

	for thisEntry := range dst {
		if dst[thisEntry].HypervisorPresent {
			log.Debug("Windows Hypervisor was detected")

			return true
		}
	}

	log.Debug("Windows Hypervisor was not detected")

	return false
}

// IsHyperV returns true if Hypervisor is present
func IsHyperV() bool {
	// We call this just to log the status of optional features
	_ = isHyperVOptionalFeature()

	// This is the value we want to return
	return isHypervisorPresent()
}
