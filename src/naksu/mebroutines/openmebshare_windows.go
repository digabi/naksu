package mebroutines

import (
	"fmt"
)

// OpenMebShare executes command that opens meb share directory
func OpenMebShare() {
	mebSharePath := GetMebshareDirectory()

	LogDebug(fmt.Sprintf("MEB share directory: %s", mebSharePath))

	if !ExistsDir(mebSharePath) {
		ShowWarningMessage("Cannot open MEB share directory since it does not exist")
		return
	}

	runParams := []string{"explorer", mebSharePath}

	// For some not-obvious reason Run_get_output() results err
	output, err := RunAndGetError(runParams)

	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		// Opening explorer results exit code 1
		if errStr != "exit status 1" {
			ShowWarningMessage("Could not open MEB share directory")
			LogDebug(fmt.Sprintf("Could not open MEB share directory: %v", err))
		}
	}

	LogDebug("MEB share directory open output:")
	LogDebug(output)
}
