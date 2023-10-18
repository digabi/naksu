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

const (
	installProgressDownloadingImage = 33
	installProgressCreatingVM       = 66
	installProgressFinished         = 100
)

// newServer downloads and creates new Abitti or Exam server using the given image URL
func newServer(boxType string, imageURL string, versionURL string) error {
	version, err := download.GetAvailableVersion(versionURL)
	switch fmt.Sprintf("%v", err) {
	case "<nil>":
	case "403", "404":
		mebroutines.ShowTranslatedErrorMessage("Please check the install passphrase")

		return fmt.Errorf("wrong passphrase entered (got %w)", err)
	default:
		mebroutines.ShowTranslatedErrorMessage("Could not get version string for a new server: %v", err)

		return fmt.Errorf("error from server: %w", err)
	}

	// Clean message
	progress.SetMessage("")

	// Initialize dialog
	progressDialog := progress.TranslateAndShowProgressDialog("Preparing...")

	// Check prerequisites
	if ensureServerIsNotRunningAndDoesNotExist() != nil || ensureDiskIsReady(&progressDialog) != nil {
		progress.CloseProgressDialog(progressDialog)

		return errors.New("server exists or disk is not ready")
	}

	err = downloadAndInstallVM(&progressDialog, imageURL, boxType, version)
	if err != nil {
		log.Error("Failed to download and install VM: %v", err)

		return err
	}

	progressDialogMessage := xlate.GetRaw("Removing temporary raw image file")
	progress.UpdateProgressDialog(progressDialog, installProgressFinished, &progressDialogMessage)
	err = os.Remove(mebroutines.GetImagePath())

	if err != nil {
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
		log.Error("Failed to ensure Naksu directories exist: %v", err)
		mebroutines.ShowTranslatedErrorMessage("Could not create directory: %v", err)

		return err
	}

	err = ensureFreeDisk()
	if err != nil {
		log.Error("Failed to ensure we have enough free disk: %v", err)
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
		return fmt.Errorf("could not create ktp (%s): %w", ktpPath, errKtpPath)
	}

	log.Debug("ktpPath is %s", ktpPath)

	// Create ~/ktp-jako if missing
	if dialog != nil {
		progress.TranslateAndUpdateProgressDialogWithMessage(*dialog, 1, "Creating ~/ktp-jako")
	} else {
		progress.TranslateAndSetMessage("Creating ~/ktp-jako")
	}
	ktpJakoPath, errKtpJakoPath := createKtpJakoDir()

	if errKtpJakoPath != nil {
		return fmt.Errorf("could not create ktp-jako (%s): %w", ktpJakoPath, errKtpJakoPath)
	}

	log.Debug("ktpJakoPath is %s", ktpJakoPath)

	return nil
}

func ensureFreeDisk() error {
	err := host.CheckFreeDisk(constants.LowDiskLimit, []string{mebroutines.GetKtpDirectory(), mebroutines.GetVirtualBoxHiddenDirectory(), mebroutines.GetVirtualBoxVMsDirectory()})
	var lowDiskSizeError *host.LowDiskSizeError
	if errors.As(err, &lowDiskSizeError) {
		mebroutines.ShowTranslatedWarningMessage("Your free disk size is getting low (%s)", humanize.Bytes(lowDiskSizeError.LowSize))
	} else if err != nil {
		mebroutines.ShowTranslatedErrorMessage("Failed to calculate free disk space: %v", err)

		return err
	}

	return nil
}

func downloadAndInstallVM(progressDialog *progress.Dialog, imageURL string, boxType string, version string) error {
	updateProgressFunc := func(message string, value int) {
		progress.UpdateProgressDialog(*progressDialog, value, &message)
	}

	updateProgressFunc(xlate.GetRaw("Getting Image from the Cloud"), installProgressDownloadingImage)
	err := download.GetServerImage(imageURL, updateProgressFunc)

	if errors.Is(err, download.ErrDownloadedDiskImageCorrupted) {
		progress.CloseProgressDialog(*progressDialog)
		mebroutines.ShowTranslatedErrorMessage("Downloaded image is corrupted. Try again.")

		return fmt.Errorf("downloading image failed (image corrupted): %w", err)
	} else if err != nil {
		progress.CloseProgressDialog(*progressDialog)
		mebroutines.ShowTranslatedErrorMessage("Failed to get new VM image: %v", err)

		return fmt.Errorf("downloading image failed: %w", err)
	}

	updateProgressFunc(xlate.GetRaw("Creating New VM"), installProgressCreatingVM)
	err = box.CreateNewBox(boxType, version)

	if err != nil {
		mebroutines.ShowTranslatedErrorMessage("Failed to create new VM: %v", err)

		removeErr := os.Remove(mebroutines.GetImagePath())
		if removeErr != nil {
			log.Debug("Failed to remove image file %s: %v", mebroutines.GetImagePath(), removeErr)
		}
		progress.CloseProgressDialog(*progressDialog)

		return fmt.Errorf("failed to create new vm: %w", err)
	}

	box.ResetCache()

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
	hashCalculator := sha256.New()
	_, err := io.WriteString(hashCalculator, passphrase)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hashCalculator.Sum(nil))
}

func getExamURL(url string, passphraseHash string) string {
	re := regexp.MustCompile(`###PASSPHRASEHASH###`)

	return re.ReplaceAllString(url, passphraseHash)
}
