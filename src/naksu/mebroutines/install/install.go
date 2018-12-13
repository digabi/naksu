package install

import (
	"errors"
	"fmt"
	"io"
	"naksu/mebroutines"
	"naksu/progress"
	"naksu/xlate"
	"net/http"
	"os"
)

// GetServer downloads vagrantfile and starts server
func GetServer(newVagrantfilePath string) {
	const VagrantURL = "http://static.abitti.fi/usbimg/qa/vagrant/Vagrantfile"

	// Create ~/ktp if missing
	progress.TranslateAndSetMessage("Creating ~/ktp")
	var ktpPath = createKtpDir()
	mebroutines.LogDebug(fmt.Sprintf("ktpPath is %s", ktpPath))

	// Create ~/ktp-jako if missing
	progress.TranslateAndSetMessage("Creating ~/ktp-jako")
	var ktpJakoPath = createKtpJakoDir()
	mebroutines.LogDebug(fmt.Sprintf("ktpJakoPath is %s", ktpJakoPath))

	var vagrantfilePath = ktpPath + string(os.PathSeparator) + "Vagrantfile"

	// If no given Vagrantfile try to download one
	if newVagrantfilePath == "" {
		// Download Vagrantfile (Abitti)
		abittiVagrantfilePath := vagrantfilePath + ".abitti"

		progress.TranslateAndSetMessage("Downloading Abitti Vagrantfile")
		errDownload := downloadFile(VagrantURL, abittiVagrantfilePath)
		if errDownload != nil {
			mebroutines.LogDebug(fmt.Sprintf("Download failed: %v", errDownload))
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
		mebroutines.LogDebug(fmt.Sprintf("Copying failed, error: %v", err))
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

func downloadFile(url string, filepath string) error {
	mebroutines.LogDebug(fmt.Sprintf("Starting download from URL %s to file %s", url, filepath))

	out, err1 := os.Create(filepath)
	if err1 != nil {
		return errors.New("Failed to create file")
	}
	defer mebroutines.Close(out)

	/* #nosec */
	resp, err2 := http.Get(url)
	if err2 != nil {
		return errors.New("Failed to retrieve file")
	}
	defer mebroutines.Close(resp.Body)

	_, err3 := io.Copy(out, resp.Body)
	if err3 != nil {
		return errors.New("Failed to copy body")
	}

	mebroutines.LogDebug(fmt.Sprintf("Finished download from URL %s to file %s", url, filepath))
	return nil
}
