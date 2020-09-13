package remove

import (
	"errors"
	"fmt"

	"naksu/box"
	"naksu/mebroutines"
	"naksu/ui/progress"
)

// Server removes all directories related to VirtualBox
func Server() error {
	isRunning, errRunning := box.Running()

	switch {
	case errRunning != nil:
		mebroutines.ShowWarningMessage(fmt.Sprintf("We could not detect whether existing VM is running: %v, but continued removing the server as you requested.", errRunning))
	case isRunning:
		mebroutines.ShowWarningMessage("There is a server appears to be running but we remove it as you requested.")
	}

	var err error

	// Chdir to home directory to avoid problems with Windows where deleting
	// a directory where the process is running
	progress.TranslateAndSetMessage("Chdir ~")
	if !mebroutines.ChdirHomeDirectory() {
		return errors.New("could not chdir to home directory")
	}

	progress.TranslateAndSetMessage("Deleting ~/.VirtualBox")
	err = mebroutines.RemoveDir(mebroutines.GetVirtualBoxHiddenDirectory())
	if err != nil {
		deleteFailed(mebroutines.GetVirtualBoxHiddenDirectory(), err)
		return err
	}

	progress.TranslateAndSetMessage("Deleting ~/VirtualBox VMs")
	err = mebroutines.RemoveDir(mebroutines.GetVirtualBoxVMsDirectory())
	if err != nil {
		deleteFailed(mebroutines.GetVirtualBoxVMsDirectory(), err)
		return err
	}

	return nil
}

// deleteFailed gives user an error message
func deleteFailed(failedPath string, err error) {
	mebroutines.ShowWarningMessage(fmt.Sprintf("Failed to remove directory %s: %v", failedPath, err))
}
