package mebroutines

import (
	"fmt"
)

// OpenMebShare opens file explorer with meb share path
func OpenMebShare() {
	mebSharePath := GetMebshareDirectory()

	LogDebug(fmt.Sprintf("MEB share directory: %s", mebSharePath))

	if !ExistsDir(mebSharePath) {
		ShowWarningMessage("Cannot open MEB share directory since it does not exist")
		return
	}

	runParams := []string{"open", mebSharePath}

	output, err := RunAndGetOutput(runParams)

	if err != nil {
		ShowWarningMessage("Could not open MEB share directory")
	}

	LogDebug("MEB share directory open output:")
	LogDebug(output)
}
