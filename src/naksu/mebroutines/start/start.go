package start

import (
	"errors"
	"fmt"
	"io/ioutil"
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

	// chdir ~/ktp
	if !mebroutines.ChdirVagrantDirectory() {
		mebroutines.ShowErrorMessage("Could not change to vagrant directory ~/ktp")
	}

	// Start VM
	progress.TranslateAndSetMessage("Starting Exam server. This takes a while.")
	upRunParams := []string{"up"}
	mebroutines.RunVagrant(upRunParams)
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
		if ! entryInDefaultVMDirRoot.IsDir() {
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
	systemProperties := mebroutines.RunVBoxManage([]string{"list", "systemproperties"}, false)
	defaultVMDirectoryRE := regexp.MustCompile(`Default machine folder:\s+(\S.*)`)
	result := defaultVMDirectoryRE.FindStringSubmatch(systemProperties)
	if len(result) > 1 {
		return strings.TrimSpace(result[1]), nil
	}
	return "", errors.New("failed to get defaultVMDirectory: no regex matches")
}
