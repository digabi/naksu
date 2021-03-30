package remove

import (
	"errors"
	"fmt"

	"naksu/box"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
)

// Server removes all directories related to VirtualBox
// Returns:
// - bool: If true, the error has already been communicated to the user
// - error: Error object
func Server() (bool, error) {
	isRunning, err := box.Running()

	switch {
	case err != nil:
		mebroutines.ShowWarningMessage(fmt.Sprintf("We could not detect whether existing VM is running: %v, but continued removing the server as you requested.", err))
	case isRunning:
		mebroutines.ShowWarningMessage("There is a server appears to be running but we remove it as you requested.")
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
		return false, errors.New("could not chdir to home directory")
	}

	progress.TranslateAndSetMessage("Deleting ~/.VirtualBox")
	mebroutines.RemoveDirAndLogErrors(mebroutines.GetVirtualBoxHiddenDirectory())

	progress.TranslateAndSetMessage("Deleting ~/VirtualBox VMs")
	err = mebroutines.RemoveDir(mebroutines.GetVirtualBoxVMsDirectory())
	if err != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf("Failed to remove directory %s: %v", mebroutines.GetVirtualBoxVMsDirectory(), err))
		return true, err
	}

	return false, nil
}
