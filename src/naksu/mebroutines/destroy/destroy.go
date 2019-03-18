package destroy

import (
	"regexp"

	"naksu/mebroutines"
	"naksu/progress"
	"naksu/xlate"
)

// Server destroys existing exam server
func Server() {
	// chdir ~/ktp
	if !mebroutines.ChdirVagrantDirectory() {
		mebroutines.ShowErrorMessage("Could not change to vagrant directory ~/ktp")
	}

	// Start VM
	progress.TranslateAndSetMessage("Removing exams. This takes a while.")
	destroyRunParams := []string{mebroutines.GetVagrantPath(), "destroy", "-f"}
	destroyOutput, destroyErr := mebroutines.RunAndGetOutput(destroyRunParams, false)

	if destroyErr == nil {
		reBoxExists, errBoxExists := regexp.MatchString("Destroying VM and associated drives", destroyOutput)
		reBoxNotCreated, errBoxNotCreated := regexp.MatchString("VM not created", destroyOutput)

		if errBoxExists == nil && reBoxExists {
			mebroutines.LogDebug("Destroy complete. There was an existing box which has been destroyed.")
			progress.TranslateAndSetMessage("Exams were removed successfully.")
			return
		}

		if errBoxNotCreated == nil && reBoxNotCreated {
			mebroutines.LogDebug("Destroy completed. There was no existing box but the destroy process finished without errors.")
			progress.TranslateAndSetMessage("Exams were removed successfully.")
			return
		}
	}

	mebroutines.ShowWarningMessage(xlate.Get("Failed to remove exams."))
	progress.TranslateAndSetMessage("Failed to remove exams.")
}
