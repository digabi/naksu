package install

import (
	"fmt"

	"naksu/box"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
	"naksu/xlate"
)

// NewServer creates new server using the given image path
func NewServer(newImagePath string) {
	_, _, errDir := ensureNaksuDirectoriesExist()

	if errDir != nil {
		log.Debug(fmt.Sprintf("Failed to ensure Naksu directories exist: %v", errDir))
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Could not create directory: %v"), errDir))
		return
	}

	if !mebroutines.ExistsFile(newImagePath) {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Image file %s does not exist", newImagePath))
		return
	}

	progress.TranslateAndSetMessage("Installing new server")
	errCreate := box.CreateNewBox(newImagePath)

	if errCreate != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf("Failed to create new VM: %v", errCreate))
		return
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
