package install

import (
	"fmt"
	"naksu/mebroutines/start"
	"os"
	"path/filepath"

	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/network"
	"naksu/ui/progress"
	"naksu/xlate"
)

// GetServer downloads vagrantfile and starts server
func GetServer(newVagrantfilePath string) {
	ktpPath, _, errDir := ensureNaksuDirectoriesExist()

	if errDir != nil {
		log.Debug(fmt.Sprintf("Failed to ensure Naksu directories exist: %v", errDir))
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Could not create directory: %v"), errDir))
		return
	}

	var vagrantfilePath = filepath.Join(ktpPath, "Vagrantfile")

	// If no given Vagrantfile try to download one
	if newVagrantfilePath == "" {
		// Download Vagrantfile (Abitti)
		abittiVagrantfilePath := vagrantfilePath + ".abitti"

		progress.TranslateAndSetMessage("Downloading Abitti Vagrantfile")
		errDownload := network.DownloadFile(constants.AbittiVagrantURL, abittiVagrantfilePath)
		if errDownload != nil {
			log.Debug(fmt.Sprintf("Download failed: %v", errDownload))
			mebroutines.ShowErrorMessage(xlate.Get("Could not update Abitti stickless server. Please check your network connection."))
			return
		}

		newVagrantfilePath = abittiVagrantfilePath
	}

	// chdir ~/ktp
	if !mebroutines.ChdirVagrantDirectory() {
		mebroutines.ShowErrorMessage("Could not change to vagrant directory ~/ktp")
		return
	}

	// If there is ~/ktp/Vagrantfile
	if mebroutines.ExistsFile(vagrantfilePath) {
		// Destroy current VM
		progress.TranslateAndSetMessage("Destroying existing server")
		destroyRunParams := []string{"destroy", "-f"}
		errDestroy := mebroutines.RunVagrant(destroyRunParams)

		if errDestroy != nil {
			log.Debug("Failed to destroy the existing server while installing a new server")
			mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Failed to execute %s: %v"), "vagrant destroy -f", errDestroy))
			// This is only a warning as the "vagrant box update" is the critical command here
		}

		removeVagrantfile(vagrantfilePath)
	}

	progress.TranslateAndSetMessage("Copying Vagrantfile")
	errCopy := mebroutines.CopyFile(newVagrantfilePath, vagrantfilePath)

	if errCopy != nil {
		log.Debug(fmt.Sprintf("Failed to copy Vagrantfile, error: %v", errCopy))
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Error while copying new Vagrantfile: %d"), errCopy))
		return
	}

	progress.TranslateAndSetMessage("Installing/updating VM: box update")
	updateRunParams := []string{"box", "update"}
	errUpdate := mebroutines.RunVagrant(updateRunParams)

	if errUpdate != nil {
		log.Debug("Failed to install/update new box when installing a new server")
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Failed to execute %s: %v"), "vagrant box update", errUpdate))
		return
	}

	progress.TranslateAndSetMessage("Installing/updating VM: box prune")
	pruneRunParams := []string{"box", "prune"}
	errPrune := mebroutines.RunVagrant(pruneRunParams)

	if errPrune != nil {
		log.Debug("Failed to prune new box when installing a new server")
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Failed to execute %s: %v"), "vagrant box prune", errPrune))
		return
	}

	progress.TranslateAndSetMessage("Downloading stickless server and starting it for the first time. This takes a long time...\n\nIf the server fails to start please try to start it again from the Naksu main menu.")
	start.Server()
}

func ensureNaksuDirectoriesExist() (string, string, error) {
	// Create ~/ktp if missing
	progress.TranslateAndSetMessage("Creating ~/ktp")
	ktpPath, errKtpPath := createKtpDir()

	if errKtpPath != nil {
		return "", "", fmt.Errorf("could not create ktp (%s): %v", ktpPath, errKtpPath)
	}

	log.Debug(fmt.Sprintf("ktpPath is %s", ktpPath))

	// Create ~/ktp-jako if missing
	progress.TranslateAndSetMessage("Creating ~/ktp-jako")
	ktpJakoPath, errKtpJakoPath := createKtpJakoDir()

	if errKtpJakoPath != nil {
		return "", "", fmt.Errorf("could not create ktp-jako (%s): %v", ktpJakoPath, errKtpJakoPath)
	}

	log.Debug(fmt.Sprintf("ktpJakoPath is %s", ktpJakoPath))

	return ktpPath, ktpJakoPath, nil
}

func createKtpDir() (string, error) {
	var ktpPath = mebroutines.GetVagrantDirectory()

	var err error

	if !mebroutines.ExistsDir(ktpPath) {
		err = mebroutines.CreateDir(ktpPath)
	}

	return ktpPath, err
}

func createKtpJakoDir() (string, error) {
	var ktpJakoPath = mebroutines.GetMebshareDirectory()

	var err error

	if !mebroutines.ExistsDir(ktpJakoPath) {
		err = mebroutines.CreateDir(ktpJakoPath)
	}

	return ktpJakoPath, err
}

func removeVagrantfile(vagrantfilePath string) {
	// Delete Vagrantfile.bak
	if mebroutines.ExistsFile(vagrantfilePath + ".bak") {
		err := os.Remove(vagrantfilePath + ".bak")
		if err != nil {
			mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Failed to delete %s"), vagrantfilePath+".bak"))
		}
	}

	// Rename Vagrantfile to Vagrantfile.bak
	err := os.Rename(vagrantfilePath, vagrantfilePath+".bak")
	if err != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Failed to rename %s to %s"), vagrantfilePath, vagrantfilePath+".bak"))
	}
}
