package remove

import (
	"errors"

	"naksu/box"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
	"naksu/xlate"
)

var generalErrorString = xlate.GetRaw("Error while removing server: %v")

// Server removes all directories related to VirtualBox
func Server() error {
	isRunning, err := box.Running()

	switch {
	case err != nil:
		mebroutines.ShowTranslatedWarningMessage("We could not detect whether existing VM is running: %v, but continued removing the server as you requested.", err)
	case isRunning:
		mebroutines.ShowTranslatedWarningMessage("The server appears to be running but we remove it as you requested.")
	}

	// Remove current box to syncronise running VirtualBox GUI
	err = box.RemoveCurrentBox()
	if err != nil {
		log.Debug("Got error when removed current box before removing server: %v", err)
	}

	// Chdir to home directory to avoid problems with Windows where deleting
	// a directory where the process is running
	progress.TranslateAndSetMessage("Chdir ~")
	if !mebroutines.ChdirHomeDirectory() {
		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, errors.New("could not chdir to home directory"))
	}

	progress.TranslateAndSetMessage("Deleting ~/.VirtualBox")
	mebroutines.RemoveDirAndLogErrors(mebroutines.GetVirtualBoxHiddenDirectory())

	progress.TranslateAndSetMessage("Deleting ~/VirtualBox VMs")
	err = mebroutines.RemoveDir(mebroutines.GetVirtualBoxVMsDirectory())
	if err != nil {
		mebroutines.ShowTranslatedWarningMessage("Failed to remove directory %s: %v", mebroutines.GetVirtualBoxVMsDirectory(), err)

		return err
	}

	box.ResetCache()

	return nil
}
