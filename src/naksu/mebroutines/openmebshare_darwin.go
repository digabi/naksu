package mebroutines

import (
	"fmt"
)

func Open_meb_share() {
	meb_share_path := GetMebshareDirectory()

	LogDebug(fmt.Sprintf("MEB share directory: %s", meb_share_path))

	if !ExistsDir(meb_share_path) {
		ShowWarningMessage("Cannot open MEB share directory since it does not exist")
		return
	}

	run_params := []string{"open", meb_share_path}

	output, err := RunAndGetOutput(run_params)

	if err != nil {
		ShowWarningMessage("Could not open MEB share directory")
	}

	LogDebug("MEB share directory open output:")
	LogDebug(output)
}
