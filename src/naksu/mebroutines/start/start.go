package start

import (
	"errors"
	"fmt"

	"naksu/box"
	"naksu/box/vboxmanage"
	"naksu/mebroutines"
)

// Server starts the exam server
// Returns:
// - bool: If true, the error has already been communicated to the user
// - error: Error object
func Server() (bool, error) {
	vboxmanage.CleanUpTrashVMDirectories()

	isInstalled, err := box.Installed()
	if err != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not start server as we could not detect whether existing VM is installed: %v", err))
		return true, fmt.Errorf("could not detect whether an existing vm is installed: %v", err)
	}

	if !isInstalled {
		mebroutines.ShowErrorMessage("No server has been installed.")
		return true, errors.New("no server has been installed")
	}

	isRunning, err := box.Running()
	if err != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not start server as we could not detect whether existing VM is running: %v", err))
		return true, fmt.Errorf("could not detect whether the server is running: %v", err)
	}

	if isRunning {
		mebroutines.ShowErrorMessage("The server is already running.")
		return true, errors.New("the server is already running")
	}

	err = box.StartCurrentBox()
	if err != nil {
		return false, err
	}

	return false, nil
}
