package host

import (
	"errors"

	"golang.org/x/sys/windows"

	"naksu/log"
)

// IsHWVirtualisation returns true if hardware virtualisation support is
// detected by the OS. Return true if all is good
func IsHWVirtualisation() bool {
	var mod = windows.NewLazyDLL("kernel32.dll")
	var proc = mod.NewProc("IsProcessorFeaturePresent")
	var PF_VIRT_FIRMWARE_ENABLED = 21 //nolint

	ret, _, err := proc.Call(uintptr(PF_VIRT_FIRMWARE_ENABLED))

	if !errors.Is(err, windows.ERROR_SUCCESS) {
		log.Error("Error while making call to kernel32.dll, IsProcessorFeaturePresent: %v", err)

		return false
	}

	log.Debug("Kernel IsProcessorFeaturePresent returns %d", ret)

	if ret > 0 {
		log.Debug("Hardware virtualisation support is enabled in BIOS (Windows)")

		return true
	}

	log.Debug("Hardware virtualisation support is disabled in BIOS (Windows)")

	return false
}
