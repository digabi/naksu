package start

import (
	"errors"
	"fmt"
	"io/ioutil"
	"naksu/box"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
	"os"
	"path"
	"regexp"
	"strings"
)

// Server starts exam server by running vagrant
func Server() {
	cleanUpTrashVMDirectories()

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

// cleanUpTrashVMDirectories tries to find and delete leftover VM directories that only contain one .vbox file and nothing else
func cleanUpTrashVMDirectories() {
	defaultVMDirectory, err := virtualBoxDefaultVMDirectory()
	if err != nil {
		log.Debug(fmt.Sprintf("Error searching for trash VM directories (get default vm dir): %v", err))
		return
	}

	entriesInDefaultVMDir, err := ioutil.ReadDir(defaultVMDirectory)
	if err != nil {
		log.Debug(fmt.Sprintf("Error searching for trash VM directories (list default vm dir %s): %v", defaultVMDirectory, err))
		return
	}

	for _, entryInDefaultVMDirRoot := range entriesInDefaultVMDir {
		if !entryInDefaultVMDirRoot.IsDir() {
			continue
		}

		fullPathToPotentialTrashVMDir := path.Join(defaultVMDirectory, entryInDefaultVMDirRoot.Name())
		entriesInSubDir, err := ioutil.ReadDir(fullPathToPotentialTrashVMDir)
		if err != nil {
			log.Debug(fmt.Sprintf("Error searching for trash VM directories (listing '%s'): %v", fullPathToPotentialTrashVMDir, err))
			return
		}
		if len(entriesInSubDir) == 1 && !entriesInSubDir[0].IsDir() && strings.HasSuffix(entriesInSubDir[0].Name(), ".vbox") {
			log.Debug(fmt.Sprintf("Removing trash VM dir %s", fullPathToPotentialTrashVMDir))
			err := os.RemoveAll(fullPathToPotentialTrashVMDir)
			if err != nil {
				log.Debug(fmt.Sprintf("Error removing trash VM dir %s", fullPathToPotentialTrashVMDir))
			}
		}
	}
}

func virtualBoxDefaultVMDirectory() (string, error) {
	systemProperties, err := mebroutines.RunVBoxManage([]string{"list", "systemproperties"})

	if err != nil {
		log.Debug("Failing to list system properties is not a fatal error, continuing normally")
	}

	defaultVMDirectoryRE := regexp.MustCompile(`Default machine folder:\s+(\S.*)`)
	result := defaultVMDirectoryRE.FindStringSubmatch(systemProperties)
	if len(result) > 1 {
		return strings.TrimSpace(result[1]), nil
	}
	return "", errors.New("failed to get defaultVMDirectory: no regex matches")
}
