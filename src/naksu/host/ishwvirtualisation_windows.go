package host

import (
	"fmt"
	"syscall"

	"naksu/log"
)

// IsHWVirtualisation returns true if hardware virtualisation support is
// detected by the OS. Return true if all is good
func IsHWVirtualisation() bool {
	var mod = syscall.NewLazyDLL("kernel32.dll")
	var proc = mod.NewProc("IsProcessorFeaturePresent")
	var PF_VIRT_FIRMWARE_ENABLED = 21 //nolint

	ret, _, err := proc.Call(
		uintptr(PF_VIRT_FIRMWARE_ENABLED),
	)

	if (err != nil) {
		log.Debug(fmt.Sprintf("Error while making call to kernel32.dll, IsProcessorFeaturePresent: %v", err))
		return false
	}
	
	log.Debug(fmt.Sprintf("Kernel IsProcessorFeaturePresent returns %d", ret))

	if ret > 0 {
		log.Debug("Hardware virtualisation support is enabled in BIOS (Windows)")
		return true
	}

	log.Debug("Hardware virtualisation support is disabled in BIOS (Windows)")
	return false
}
