package mebroutines

import (
	"fmt"
	"os"
)

// ExistsStdin checks if app has stdin
func ExistsStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		LogDebug(fmt.Sprintf("Checking status for Stdin results error: %s", fmt.Sprint(err)))
		return false
	}

	stat := fmt.Sprintf("%v", fi.Mode())

	LogDebug(fmt.Sprintf("STDIN stat is: %s", stat))

	if stat == "Dcrw-rw-rw-" {
		// No stdin
		LogDebug("No STDIN detected - naksu is executed without a terminal")
		return false
	}

	LogDebug("STDIN detected - naksu is executed from a terminal")

	return true
}
