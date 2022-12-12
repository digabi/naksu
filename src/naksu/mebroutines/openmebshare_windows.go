package mebroutines

import (
	"fmt"

	"naksu/log"
)

// OpenMebShare executes command that opens meb share directory
func OpenMebShare() {
	mebSharePath := GetMebshareDirectory()

	log.Debug("MEB share directory: %s", mebSharePath)

	if !ExistsDir(mebSharePath) {
		ShowTranslatedWarningMessage("Cannot open MEB share directory since it does not exist")

		return
	}

	runParams := []string{"explorer", mebSharePath}

	// For some not-obvious reason Run_get_output() results err
	output, err := RunAndGetOutput(runParams, true)

	if err != nil {
		errStr := fmt.Sprintf("%v", err)
		// Opening explorer results exit code 1
		if errStr != "exit status 1" {
			ShowTranslatedWarningMessage("Could not open MEB share directory")
			log.Warning("Could not open MEB share directory: %v", err)
		}
	}

	log.Debug("MEB share directory open output:")
	log.Debug(output)
}
