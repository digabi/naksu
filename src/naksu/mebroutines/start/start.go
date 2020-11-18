package start

import (
	"fmt"
	"naksu/box"
	"naksu/box/vboxmanage"
	"naksu/mebroutines"
	"naksu/ui/progress"
)

// Server starts exam server
func Server() {
	vboxmanage.CleanUpTrashVMDirectories()

	isInstalled, errInstalled := box.Installed()
	if errInstalled != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not start server as we could not detect whether existing VM is installed: %v", errInstalled))
		return
	}

	if !isInstalled {
		mebroutines.ShowErrorMessage("No server has been installed.")
		return
	}

	isRunning, errRunning := box.Running()
	if errRunning != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not start server as we could not detect whether existing VM is running: %v", errRunning))
		return
	}

	if isRunning {
		mebroutines.ShowErrorMessage("The server is already running.")
		return
	}

	err := box.StartCurrentBox()
	if err != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not start VM: %v", err))
		return
	}

	progress.SetMessage("Virtual machine was started")
}
