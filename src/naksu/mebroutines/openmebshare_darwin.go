package mebroutines

import (
	"naksu/log"
)

// OpenMebShare opens file explorer with meb share path
func OpenMebShare() {
	mebSharePath := GetMebshareDirectory()

	log.Debug("MEB share directory: %s", mebSharePath)

	if !ExistsDir(mebSharePath) {
		ShowTranslatedWarningMessage("Cannot open MEB share directory since it does not exist")
		return
	}

	runParams := []string{"open", mebSharePath}

	output, err := RunAndGetOutput(runParams, true)

	if err != nil {
		ShowTranslatedWarningMessage("Could not open MEB share directory")
	}

	log.Debug("MEB share directory open output:")
	log.Debug(output)
}
