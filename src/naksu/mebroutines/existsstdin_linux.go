package mebroutines

import (
	"fmt"
	"os"

	"naksu/log"
)

// ExistsStdin returns true if stdin is available
func ExistsStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.LogDebug(fmt.Sprintf("Checking status for Stdin results error: %s", fmt.Sprint(err)))
		return false
	}

	stat := fmt.Sprintf("%v", fi.Mode())

	log.LogDebug(fmt.Sprintf("STDIN stat is: %s", stat))

	if stat == "Dcrw-rw-rw-" {
		// No stdin
		log.LogDebug("No STDIN detected - naksu is executed without a terminal")
		return false
	}

	log.LogDebug("STDIN detected - naksu is executed from a terminal")

	return true
}
