package install

import (
	"fmt"
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
	// Create ~/ktp if missing
	progress.TranslateAndSetMessage("Creating ~/ktp")
	var ktpPath = createKtpDir()
	log.Debug(fmt.Sprintf("ktpPath is %s", ktpPath))

	// Create ~/ktp-jako if missing
	progress.TranslateAndSetMessage("Creating ~/ktp-jako")
	var ktpJakoPath = createKtpJakoDir()
	log.Debug(fmt.Sprintf("ktpJakoPath is %s", ktpJakoPath))

	var vagrantfilePath = filepath.Join(ktpPath, "Vagrantfile")

	// If no given Vagrantfile try to download one
	if newVagrantfilePath == "" {
		// Download Vagrantfile (Abitti)
		abittiVagrantfilePath := vagrantfilePath + ".abitti"

		progress.TranslateAndSetMessage("Downloading Abitti Vagrantfile")
		errDownload := network.DownloadFile(constants.AbittiVagrantURL, abittiVagrantfilePath)
		if errDownload != nil {
			log.Debug(fmt.Sprintf("Download failed: %v", errDownload))
			mebroutines.ShowWarningMessage(xlate.Get("Could not update Abitti stickless server. Please check your network connection."))
			return
		}

		newVagrantfilePath = abittiVagrantfilePath
	}

	// chdir ~/ktp
	if !mebroutines.ChdirVagrantDirectory() {
		mebroutines.ShowErrorMessage("Could not change to vagrant directory ~/ktp")
	}

	// If there is ~/ktp/Vagrantfile
	if mebroutines.ExistsFile(vagrantfilePath) {
		// Destroy current VM
		progress.TranslateAndSetMessage("Destroying existing server")
		destroyRunParams := []string{"destroy", "-f"}
		mebroutines.RunVagrant(destroyRunParams)

		removeVagrantfile(vagrantfilePath)
	}

	progress.TranslateAndSetMessage("Copying Vagrantfile")
	err := mebroutines.CopyFile(newVagrantfilePath, vagrantfilePath)

	if err != nil {
		log.Debug(fmt.Sprintf("Copying failed, error: %v", err))
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Error while copying new Vagrantfile: %d"), err))
	}

	progress.TranslateAndSetMessage("Installing/updating VM: box update")
	updateRunParams := []string{"box", "update"}
	mebroutines.RunVagrant(updateRunParams)

	progress.TranslateAndSetMessage("Installing/updating VM: box prune")
	pruneRunParams := []string{"box", "prune"}
	mebroutines.RunVagrant(pruneRunParams)

	progress.TranslateAndSetMessage("Downloading stickless server and starting it for the first time. This takes a long time...\n\nIf the server fails to start please try to start it again from the Naksu main menu.")
	upRunParams := []string{"up"}
	mebroutines.RunVagrant(upRunParams)
}

func createKtpDir() string {
	var ktpPath = mebroutines.GetVagrantDirectory()

	if !mebroutines.ExistsDir(ktpPath) {
		if mebroutines.CreateDir(ktpPath) != nil {
			mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Could not create ~/ktp to %s"), ktpPath))
		}
	}

	return ktpPath
}

func createKtpJakoDir() string {
	var ktpJakoPath = mebroutines.GetMebshareDirectory()

	if !mebroutines.ExistsDir(ktpJakoPath) {
		if mebroutines.CreateDir(ktpJakoPath) != nil {
			mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Could not create ~/ktp-jako to %s"), ktpJakoPath))
		}
	}

	return ktpJakoPath
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
