package mebroutines

import (
	"fmt"
	"os"

	"naksu/log"
)

// ExistsStdin checks if app has stdin
func ExistsStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Debug(fmt.Sprintf("Checking status for Stdin results error: %s", fmt.Sprint(err)))
		return false
	}

	stat := fmt.Sprintf("%v", fi.Mode())

	log.Debug(fmt.Sprintf("STDIN stat is: %s", stat))

	if stat == "Dcrw-rw-rw-" {
		// No stdin
		log.Debug("No STDIN detected - naksu is executed without a terminal")
		return false
	}

	log.Debug("STDIN detected - naksu is executed from a terminal")

	return true
}
