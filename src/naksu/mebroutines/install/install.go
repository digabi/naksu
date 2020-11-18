package install

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"naksu/box"
	"naksu/cloud"
	"naksu/constants"
	"naksu/host"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
	"naksu/xlate"

	humanize "github.com/dustin/go-humanize"
)

// newServer downloads and creates new Abitti or Exam server using the given image URL
func newServer(boxType string, imageURL string, versionURL string) {
	version, errVersion := cloud.GetAvailableVersion(versionURL)
	if errVersion != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not get version string for a new server: %v", errVersion))
		return
	}

	errRI := ensureServerIsNotRunningAndDoesNotExist()
	if errRI != nil {
		return
	}

	errDir := ensureNaksuDirectoriesExist()
	if errDir != nil {
		log.Debug(fmt.Sprintf("Failed to ensure Naksu directories exist: %v", errDir))
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Could not create directory: %v"), errDir))
		return
	}

	newImagePath, errNewImagePath := getTempFilePath()
	if errNewImagePath != nil {
		return
	}

	progress.TranslateAndSetMessage("Getting Image from the Cloud")
	errGet := cloud.GetServerImage(imageURL, newImagePath, progress.TranslateAndSetMessage)

	if errGet != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Failed to get new VM image: %v", errGet))
		return
	}

	progress.TranslateAndSetMessage("Creating New VM")
	errCreate := box.CreateNewBox(boxType, newImagePath, version)

	if errCreate != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Failed to create new VM: %v", errCreate))
		errRemove := os.Remove(newImagePath)

		if errRemove != nil {
			log.Debug(fmt.Sprintf("Failed to remove image file %s: %v", newImagePath, errRemove))
		}

		return
	}

	progress.SetMessage("Removing temporary raw image file")
	errRemove := os.Remove(newImagePath)

	if errRemove != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf("Failed to remove raw image file %s: %v", newImagePath, errRemove))
	}

	progress.SetMessage("New VM was created")
}

// NewAbittiServer downloads and installs a new Abitti server
func NewAbittiServer() {
	newServer(constants.AbittiBoxType, constants.AbittiEtcherURL, constants.AbittiVersionURL)
}

func NewExamServer(passphrase string) {
	passphraseHash := getPassphraseHash(passphrase)
	imageURL := getExamURL(constants.MatriculationExamEtcherURL, passphraseHash)
	versionURL := getExamURL(constants.MatriculationExamVersionURL, passphraseHash)

	newServer(constants.MatriculationExamBoxType, imageURL, versionURL)
}

func ensureServerIsNotRunningAndDoesNotExist() error {
	isRunning, errRunning := box.Running()
	if errRunning != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not install server as we could not detect whether existing VM is running: %v", errRunning))
		return errRunning
	}

	if isRunning {
		mebroutines.ShowErrorMessage("Please stop the current server before installing a new one")
		return errors.New("please stop the current server before installing a new one")
	}

	isInstalled, errInstalled := box.Installed()
	if errInstalled != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not install server as we could not detect whether existing VM is installed: %v", errInstalled))
		return errInstalled
	}

	if isInstalled {
		errRemove := box.RemoveCurrentBox()
		if errRemove != nil {
			mebroutines.ShowWarningMessage(fmt.Sprintf("Could not remove current VM before installing new one: %v", errRemove))
		}
	}

	return nil
}

func ensureNaksuDirectoriesExist() error {
	// Create ~/ktp if missing
	progress.TranslateAndSetMessage("Creating ~/ktp")
	ktpPath, errKtpPath := createKtpDir()

	if errKtpPath != nil {
		return fmt.Errorf("could not create ktp (%s): %v", ktpPath, errKtpPath)
	}

	log.Debug(fmt.Sprintf("ktpPath is %s", ktpPath))

	// Create ~/ktp-jako if missing
	progress.TranslateAndSetMessage("Creating ~/ktp-jako")
	ktpJakoPath, errKtpJakoPath := createKtpJakoDir()

	if errKtpJakoPath != nil {
		return fmt.Errorf("could not create ktp-jako (%s): %v", ktpJakoPath, errKtpJakoPath)
	}

	log.Debug(fmt.Sprintf("ktpJakoPath is %s", ktpJakoPath))

	return nil
}

func getTempFilePath() (string, error) {
	newImagePath, errTemp := mebroutines.GetTempFilename()
	if errTemp != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Failed to create temporary file: %v", errTemp))
		return "", errTemp
	}

	errRemove := os.Remove(newImagePath)
	if errRemove != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf("Failed to remove raw image file %s: %v", newImagePath, errRemove))
	}

	freeSize, errDiskFree := host.CheckFreeDisk(constants.LowDiskLimit, []string{filepath.Dir(newImagePath), mebroutines.GetKtpDirectory(), mebroutines.GetVirtualBoxHiddenDirectory(), mebroutines.GetVirtualBoxVMsDirectory()})

	switch {
	case errDiskFree != nil && strings.HasPrefix(fmt.Sprintf("%v", errDiskFree), "low:"):
		mebroutines.ShowTranslatedWarningMessage("Your free disk size is getting low (%s)", humanize.Bytes(freeSize))
		// We just inform the user instead of returning an error
	case errDiskFree != nil:
		mebroutines.ShowErrorMessage(fmt.Sprintf("Failed to query free disk space: %v", errDiskFree))
		return "", errDiskFree
	}

	return newImagePath, nil
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
