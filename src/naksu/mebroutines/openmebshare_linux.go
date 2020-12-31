package mebroutines

import (
	"fmt"

	"naksu/log"
)

// OpenMebShare command executes command that opens file browser to meb share directory
func OpenMebShare() {
	mebSharePath := GetMebshareDirectory()

	log.Debug(fmt.Sprintf("MEB share directory: %s", mebSharePath))

	if !ExistsDir(mebSharePath) {
		ShowWarningMessage("Cannot open MEB share directory since it does not exist")
		return
	}

	// Try to open MEB share folder with any of these utils
	// Hopefully we have at least one of them installed!
	// We do the opening in a goroutine to avoid any lags to the UI
	
	openers := [3]string{"xdg-open", "gnome-open", "nautilus"}

	go func (){
		for _, thisOpener := range openers {
			runParams := []string{thisOpener, mebSharePath}
			output, err := RunAndGetOutput(runParams, true)

			if err == nil {
				log.Debug("MEB share directory open output:")
				log.Debug(output)

				return
			}
		}

		ShowWarningMessage("Could not open MEB share directory")
	}()
}
