package mebroutines

import "naksu/log"

// ExistsStdin returns always true on windows
func ExistsStdin() bool {
	log.LogDebug("Windows is always expected to have STDIN")
	return true
}
