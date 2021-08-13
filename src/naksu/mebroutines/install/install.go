package install

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"naksu/box"
	"naksu/box/download"
	"naksu/constants"
	"naksu/host"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
	"naksu/xlate"

	humanize "github.com/dustin/go-humanize"
)

// newServer downloads and creates new Abitti or Exam server using the given image URL
func newServer(boxType string, imageURL string, versionURL string) error {
	version, err := download.GetAvailableVersion(versionURL)
	switch fmt.Sprintf("%v", err) {
	case "<nil>":
	case "404":
		mebroutines.ShowTranslatedErrorMessage("Please check the install passphrase")
		return errors.New("wrong passphrase entered (got 404)")
	default:
		mebroutines.ShowTranslatedErrorMessage("Could not get version string for a new server: %v", err)
		return fmt.Errorf("error from server: %v", err)
	}

	// Clean message
	progress.SetMessage("")

	// Initialize dialog
	progressDialog := progress.TranslateAndShowProgressDialog("Preparing...")

	updateProgressFunc := func(message string, value int) {
		//fmt.Println(message)
		progress.UpdateProgressDialog(progressDialog, value, &message)
	}

	// Check prerequisites
	if ensureServerIsNotRunningAndDoesNotExist() != nil || ensureDiskIsReady(&progressDialog) != nil {
		progress.CloseProgressDialog(progressDialog)
		return errors.New("server exists or disk is not ready")
	}

	updateProgressFunc(xlate.GetRaw("Getting Image from the Cloud"), 100*(1/3))
	err = download.GetServerImage(imageURL, updateProgressFunc)
	if err != nil {
		progress.CloseProgressDialog(progressDialog)
		mebroutines.ShowTranslatedErrorMessage("Failed to get new VM image: %v", err)
		return fmt.Errorf("downloading image failed: %v", err)
	}

	updateProgressFunc(xlate.GetRaw("Creating New VM"), 100*(2/3))
	err = box.CreateNewBox(boxType, version)

	if err != nil {
		mebroutines.ShowTranslatedErrorMessage("Failed to create new VM: %v", err)

		removeErr := os.Remove(mebroutines.GetImagePath())
		if removeErr != nil {
			log.Debug("Failed to remove image file %s: %v", mebroutines.GetImagePath(), removeErr)
		}
		progress.CloseProgressDialog(progressDialog)
		return fmt.Errorf("failed to create new vm: %v", err)
	}

	updateProgressFunc(xlate.GetRaw("Removing temporary raw image file"), 100)
	err = os.Remove(mebroutines.GetImagePath())

	if err != nil {
		progress.CloseProgressDialog(progressDialog)
		mebroutines.ShowTranslatedWarningMessage("Failed to remove raw image file %s: %v", mebroutines.GetImagePath(), err)
	}
	progress.CloseProgressDialog(progressDialog)
	return nil
}

// NewAbittiServer downloads and installs a new Abitti server
func NewAbittiServer() error {
	return newServer(constants.AbittiBoxType, constants.AbittiEtcherURL, constants.AbittiVersionURL)
}

// NewExamServer downloads and installs a new exam server
func NewExamServer(passphrase string) error {
	passphraseHash := getPassphraseHash(passphrase)
	imageURL := getExamURL(constants.MatriculationExamEtcherURL, passphraseHash)
	versionURL := getExamURL(constants.MatriculationExamVersionURL, passphraseHash)

	return newServer(constants.MatriculationExamBoxType, imageURL, versionURL)
}

func ensureServerIsNotRunningAndDoesNotExist() error {
	isRunning, errRunning := box.Running()
	if errRunning != nil {
		mebroutines.ShowTranslatedErrorMessage("Could not install server as we could not detect whether existing VM is running: %v", errRunning)
		return errRunning
	}

	if isRunning {
		mebroutines.ShowTranslatedErrorMessage("Please stop the current server before installing a new one")
		return errors.New("please stop the current server before installing a new one")
	}

	isInstalled, errInstalled := box.Installed()
	if errInstalled != nil {
		mebroutines.ShowTranslatedErrorMessage("Could not install server as we could not detect whether existing VM is installed: %v", errInstalled)
		return errInstalled
	}

	if isInstalled {
		errRemove := box.RemoveCurrentBox()
		if errRemove != nil {
			mebroutines.ShowTranslatedWarningMessage("Could not remove current VM before installing new one: %v", errRemove)
		}
	}

	return nil
}

func ensureDiskIsReady(dialog *progress.Dialog) error {
	err := ensureNaksuDirectoriesExist(dialog)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to ensure Naksu directories exist: %v", err))
		mebroutines.ShowTranslatedErrorMessage("Could not create directory: %v", err)
		return err
	}

	err = ensureFreeDisk()
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to ensure we have enough free disk: %v", err))
		mebroutines.ShowTranslatedErrorMessage("Could not calculate free disk size: %v", err)
		return err
	}

	return nil
}

func ensureNaksuDirectoriesExist(dialog *progress.Dialog) error {
	// Create ~/ktp if missing
	if dialog != nil {
		progress.TranslateAndUpdateProgressDialogWithMessage(*dialog, 1, "Creating ~/ktp")
	} else {
		progress.TranslateAndSetMessage("Creating ~/ktp")
	}
	ktpPath, errKtpPath := createKtpDir()

	if errKtpPath != nil {
		return fmt.Errorf("could not create ktp (%s): %v", ktpPath, errKtpPath)
	}

	log.Debug(fmt.Sprintf("ktpPath is %s", ktpPath))

	// Create ~/ktp-jako if missing
	if dialog != nil {
		progress.TranslateAndUpdateProgressDialogWithMessage(*dialog, 1, "Creating ~/ktp-jako")
	} else {
		progress.TranslateAndSetMessage("Creating ~/ktp-jako")
	}
	ktpJakoPath, errKtpJakoPath := createKtpJakoDir()

	if errKtpJakoPath != nil {
		return fmt.Errorf("could not create ktp-jako (%s): %v", ktpJakoPath, errKtpJakoPath)
	}

	log.Debug(fmt.Sprintf("ktpJakoPath is %s", ktpJakoPath))

	return nil
}

func ensureFreeDisk() error {
	err := host.CheckFreeDisk(constants.LowDiskLimit, []string{mebroutines.GetKtpDirectory(), mebroutines.GetVirtualBoxHiddenDirectory(), mebroutines.GetVirtualBoxVMsDirectory()})

	if err != nil {
		if err, ok := err.(*host.LowDiskSizeError); ok {
			// We just inform the user instead of returning an error
			mebroutines.ShowTranslatedWarningMessage("Your free disk size is getting low (%s)", humanize.Bytes(err.LowSize))
		} else {
			mebroutines.ShowTranslatedErrorMessage("Failed to calculate free disk space: %v", err)
			return err
		}
	}

	return nil
}

func createKtpDir() (string, error) {
	var ktpPath = mebroutines.GetKtpDirectory()

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

func getPassphraseHash(passphrase string) string {
	o := sha256.New()
	_, err := io.WriteString(o, passphrase)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", o.Sum(nil))
}

func getExamURL(url string, passphraseHash string) string {
	re := regexp.MustCompile(`###PASSPHRASEHASH###`)
	return re.ReplaceAllString(url, passphraseHash)
}
