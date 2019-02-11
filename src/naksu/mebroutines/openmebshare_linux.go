package mebroutines

import (
	"fmt"
)

// OpenMebShare command executes command that opens file browser to meb share directory
func OpenMebShare() {
	mebSharePath := GetMebshareDirectory()

	LogDebug(fmt.Sprintf("MEB share directory: %s", mebSharePath))

	if !ExistsDir(mebSharePath) {
		ShowWarningMessage("Cannot open MEB share directory since it does not exist")
		return
	}

	runParams := []string{"xdg-open", mebSharePath}

	output, err := RunAndGetOutput(runParams, false)

	if err != nil {
		ShowWarningMessage("Could not open MEB share directory")
	}

	LogDebug("MEB share directory open output:")
	LogDebug(output)
}
