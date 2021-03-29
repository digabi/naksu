package destroy

import (
	"errors"
	"fmt"

	"naksu/box"
	"naksu/log"
	"naksu/ui/progress"
)

// Server destroys existing exam server by restoring the fresh snapshot.
// Returns:
// - bool: If true, the error has already been communicated to the user
// - error: Error object
func Server() (bool, error) {
	isInstalled, err := box.Installed()
	if err != nil {
		log.Debug("Could not start destroying server as we could not detect whether existing VM is installed: %v", err)
		return false, errors.New("could not detect whether there is an existing vm installed")
	}

	if !isInstalled {
		return false, errors.New("there is no vm installed")
	}

	isRunning, err := box.Running()
	if err != nil {
		log.Debug("Could not start destroying server as we could not detect whether existing VM is running: %v", err)
		return false, errors.New("could not detect whether there is existing vm running")
	}

	if isRunning {
		return false, errors.New("the vm is running, please stop it first")
	}

	progress.TranslateAndSetMessage("Removing exams. This takes a while.")

	err = box.RestoreSnapshot()
	if err != nil {
		log.Debug("Could not destroy VM / restore initial snapshot: %v", err)
		return false, fmt.Errorf("could not restore snapshot: %v", err)
	}

	progress.SetMessage("")

	return false, nil
}
