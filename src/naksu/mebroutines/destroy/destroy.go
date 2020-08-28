package destroy

import (
	"errors"
	"fmt"

	"naksu/log"
	"naksu/box"
	"naksu/ui/progress"
)

// Server destroys existing exam server by restoring the fresh snapshot. The errors are reported upstream.
func Server() error {
	isInstalled, errInstalled := box.Installed()
	if errInstalled != nil {
		log.Debug(fmt.Sprintf("Could not start destoroying server as we could not detect whether existing VM is installed: %v", errInstalled))
		return errors.New("could not detect wheteher there is an existing vm installed")
	}

	if ! isInstalled {
		return errors.New("there is no vm installed")
	}

	isRunning, errRunning := box.Running()
	if errRunning != nil {
		log.Debug(fmt.Sprintf("Could not start destoroying server as we could not detect whether existing VM is running: %v", errRunning))
		return errors.New("could not detect whether there is existing vm running")
	}

	if isRunning {
		return errors.New("the vm is running, please stop it first")
	}

	progress.TranslateAndSetMessage("Removing exams. This takes a while.")

	err := box.RestoreSnapshot()
	if err != nil {
		log.Debug(fmt.Sprintf("Could not destroy VM / restore initial snapshot: %v", err))
	}

	progress.SetMessage("")

	return err
}
