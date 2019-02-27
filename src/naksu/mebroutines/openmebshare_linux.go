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

	// Try to open MEB share folder with any of these utils
	// Hopefully we have at least one of them installed!
	openers := [3]string {"xdg-open", "gnome-open", "nautilus"}

	for _, thisOpener := range openers {
		runParams := []string{thisOpener, mebSharePath}
		output, err := RunAndGetOutput(runParams, false)

		if err == nil {
			LogDebug("MEB share directory open output:")
			LogDebug(output)

			return
		}
	}

	ShowWarningMessage("Could not open MEB share directory")
}
