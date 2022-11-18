package start

import (
	"errors"
	"fmt"

	"naksu/box"
	"naksu/box/vboxmanage"
	"naksu/mebroutines"
	"naksu/xlate"
)

var generalErrorString = xlate.GetRaw("Failed to start server: %v")

// Server starts the exam server
func Server() error {
	vboxmanage.CleanUpTrashVMDirectories()

	isInstalled, err := box.Installed()
	if err != nil {
		mebroutines.ShowTranslatedErrorMessage("Could not start server as we could not detect whether existing VM is installed: %v", err)

		return fmt.Errorf("could not detect whether an existing vm is installed: %w", err)
	}

	if !isInstalled {
		mebroutines.ShowTranslatedErrorMessage("No server has been installed.")

		return errors.New("no server has been installed")
	}

	isRunning, err := box.Running()
	if err != nil {
		mebroutines.ShowTranslatedErrorMessage("Could not start server as we could not detect whether existing VM is running: %v", err)

		return fmt.Errorf("could not detect whether the server is running: %w", err)
	}

	if isRunning {
		mebroutines.ShowTranslatedErrorMessage("The server is already running.")

		return errors.New("the server is already running")
	}

	err = box.StartCurrentBox()
	if err != nil {
		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, err)
	}

	return nil
}
