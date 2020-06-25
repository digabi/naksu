package remove

import (
	"errors"
	"fmt"

	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
)

// Server removes all directories related to Vagrant and VirtualBox
func Server() error {
	var err error

	// Chdir to home directory to avoid problems with Windows where deleting
	// a directory where a
	progress.TranslateAndSetMessage("Chdir ~")
	if !mebroutines.ChdirHomeDirectory() {
		return errors.New("could not chdir to home directory")
	}

	progress.TranslateAndSetMessage("Deleting ~/.vagrant.d")
	err = mebroutines.RemoveDir(mebroutines.GetVagrantdDirectory())
	if err != nil {
		deleteFailed(mebroutines.GetVagrantdDirectory(), err)
		return err
	}

	progress.TranslateAndSetMessage("Deleting ~/.VirtualBox")
	err = mebroutines.RemoveDir(mebroutines.GetVirtualBoxHiddenDirectory())
	if err != nil {
		deleteFailed(mebroutines.GetVirtualBoxHiddenDirectory(), err)
		return err
	}

	// Close current debug file in case it is located in ~/ktp
	log.SetDebugFilename("-")

	progress.TranslateAndSetMessage("Deleting ~/ktp")
	err = mebroutines.RemoveDir(mebroutines.GetVagrantDirectory())
	if err != nil {
		deleteFailed(mebroutines.GetVagrantDirectory(), err)
		return err
	}

	// Set new filename for debug log (we probably just deleted previous one above)
	log.SetDebugFilename(log.GetNewDebugFilename())

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
