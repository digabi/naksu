package install

import (
	"fmt"
	"os"

	"naksu/box"
	"naksu/cloud"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
	"naksu/xlate"
)

// NewServerAbitti downloads and creates new Abitti server using the given image path
func NewServerAbitti() {
	isInstalled, errInstalled := box.Installed()
	if errInstalled != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Could not install server as we could not detect whether existing VM is installed: %v", errInstalled))
		return
	}

	if isInstalled {
		mebroutines.ShowErrorMessage("Please remove existing server before installing a new one.")
		return
	}

	_, _, errDir := ensureNaksuDirectoriesExist()

	if errDir != nil {
		log.Debug(fmt.Sprintf("Failed to ensure Naksu directories exist: %v", errDir))
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Could not create directory: %v"), errDir))
		return
	}

	newImagePath, errTemp := mebroutines.GetTempFilename()
	if errTemp != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Failed to create temporary file: %v", errTemp))
		return
	}

	errRemove := os.Remove(newImagePath)
	if errRemove != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf("Failed to remove raw image file %s: %v", newImagePath, errRemove))
	}

	progress.TranslateAndSetMessage("Getting Image from the Cloud")
	errGet := cloud.GetAbittiImage(newImagePath, progress.TranslateAndSetMessage)

	if errGet != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Failed to get new VM image: %v", errGet))
		return
	}

	progress.TranslateAndSetMessage("Creating New VM")
	errCreate := box.CreateNewBox(newImagePath)

	if errCreate != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Failed to create new VM: %v", errCreate))
		errRemove := os.Remove(newImagePath)

		if errRemove != nil {
			log.Debug(fmt.Sprintf("Failed to remove image file %s: %v", newImagePath, errRemove))
		}

		return
	}

	progress.SetMessage("Removing temporary raw image file")
	errRemove = os.Remove(newImagePath)

	if errRemove != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf("Failed to remove raw image file %s: %v", newImagePath, errRemove))
	}

	progress.SetMessage("New VM was created")
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
