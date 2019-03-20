package destroy

import (
	"errors"
	"regexp"
	"fmt"

	"naksu/mebroutines"
	"naksu/progress"
)

// Server destroys existing exam server
func Server() error {
	// chdir ~/ktp
	if !mebroutines.ChdirVagrantDirectory() {
		mebroutines.LogDebug("Could not change to vagrant directory ~/ktp")
		return errors.New("could not chmod ~/ktp")
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
			return nil
		}

		if errBoxNotCreated == nil && reBoxNotCreated {
			mebroutines.LogDebug("Destroy completed. There was no existing box but the destroy process finished without errors.")
			return nil
		}
	}

	mebroutines.LogDebug(fmt.Sprintf("Could not remove exams. vagrant destroy says:\n%s", destroyOutput))

	return errors.New("failed to remove exams")
}
