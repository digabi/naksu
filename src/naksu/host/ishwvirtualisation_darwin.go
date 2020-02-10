package host

import (
	"naksu/log"
)

// IsHWVirtualisation returns true if hardware virtualisation support (VT-x or AMD-V)
// is available. FIXME: Not impmemented for Darwin - returns always true
func IsHWVirtualisation() bool {
	log.Debug("Warning: Detection of hardware virtualisation is not implemented for Darwin.")

	return true
}
