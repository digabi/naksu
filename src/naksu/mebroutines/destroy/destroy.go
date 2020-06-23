package destroy

import (
	"errors"
	"fmt"
	"regexp"

	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
)

// Server destroys existing exam server. The errors are reported upstream.
func Server() error {
	// chdir ~/ktp
	if !mebroutines.ChdirVagrantDirectory() {
		log.Debug("Could not change to vagrant directory ~/ktp")
		return errors.New("could not chmod ~/ktp")
	}

	// Start VM
	progress.TranslateAndSetMessage("Removing exams. This takes a while.")
	destroyRunParams := []string{mebroutines.GetVagrantPath(), "destroy", "-f"}
	destroyOutput, err := mebroutines.RunAndGetOutput(destroyRunParams)

	if err == nil {
		reBoxExists, errBoxExists := regexp.MatchString("Destroying VM and associated drives", destroyOutput)
		reBoxNotCreated, errBoxNotCreated := regexp.MatchString("VM not created", destroyOutput)

		if errBoxExists == nil && reBoxExists {
			log.Debug("Destroy complete. There was an existing box which has been destroyed.")
			return nil
		}

		if errBoxNotCreated == nil && reBoxNotCreated {
			log.Debug("Destroy completed. There was no existing box but the destroy process finished without errors.")
			return nil
		}

		err = errors.New("destroy was executed but success was not confirmed")
	}

	log.Debug(fmt.Sprintf("Could not remove exams (%v). vagrant destroy says:\n%s", err, destroyOutput))

	return err
}
