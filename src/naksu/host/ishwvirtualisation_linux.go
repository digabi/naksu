package host

import (
	"naksu/log"
	"naksu/mebroutines"
)

// IsHWVirtualisation returns true if hardware virtualisation support (VT-x or AMD-V)
// is available.
func IsHWVirtualisation() bool {
	// This method has been tested on
	// * Ubuntu 18.04 LTS
	// * Debian 10
	if mebroutines.ExistsCharDevice("/dev/kvm") {
		log.Debug("Hardware virtualisation support is enabled in BIOS (Linux)")

		return true
	}

	log.Debug("Hardware virtualisation support is disabled in BIOS (Linux)")

	return false
}
