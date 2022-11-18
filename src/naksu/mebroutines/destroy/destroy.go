package destroy

import (
	"errors"
	"fmt"

	"naksu/box"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
	"naksu/xlate"
)

var generalErrorString = xlate.GetRaw("Failed to remove exams: %v")

// Server destroys existing exam server by restoring the fresh snapshot.
func Server() error {
	isInstalled, err := box.Installed()
	if err != nil {
		log.Debug("Could not start destroying server as we could not detect whether existing VM is installed: %v", err)

		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, errors.New("could not detect whether there is an existing vm installed"))
	}

	if !isInstalled {
		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, errors.New("there is no vm installed"))
	}

	isRunning, err := box.Running()
	if err != nil {
		log.Debug("Could not start destroying server as we could not detect whether existing VM is running: %v", err)

		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, errors.New("could not detect whether there is existing vm running"))
	}

	if isRunning {
		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, errors.New("the vm is running, please stop it first"))
	}

	progress.TranslateAndSetMessage("Removing exams. This takes a while.")

	err = box.RestoreSnapshot()
	if err != nil {
		log.Debug("Could not destroy VM / restore initial snapshot: %v", err)

		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, fmt.Errorf("could not restore snapshot: %w", err))
	}

	progress.SetMessage("")

	return nil
}
