package backup

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"naksu/box"
	"naksu/constants"
	"naksu/host"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/ui/progress"
	"naksu/xlate"

	humanize "github.com/dustin/go-humanize"
)

var generalErrorString = xlate.GetRaw("Backup failed: %v")

// MakeBackup creates virtual machine backup to path
func MakeBackup(backupPath string) error {
	err := ensureBoxInstalledAndNotRunning()
	if err != nil {
		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, err)
	}

	err = host.CheckFreeDisk(constants.LowDiskLimit, []string{filepath.Dir(backupPath)})
	var lowDiskSizeError *host.LowDiskSizeError
	if errors.As(err, &lowDiskSizeError) {
		mebroutines.ShowTranslatedWarningMessage("Your free disk size is getting low (%s). If backup process fails please consider freeing some disk space.", humanize.Bytes(lowDiskSizeError.LowSize))
	} else if err != nil {
		mebroutines.ShowTranslatedErrorMessage("Failed to calculate free disk size: %v", err)
	}

	progress.TranslateAndSetMessage("Checking existing file...")
	if mebroutines.ExistsFile(backupPath) {
		mebroutines.ShowTranslatedErrorMessage("File %s already exists", backupPath)

		return errors.New("backup file already exists")
	}

	// Check if path_backup is writeable
	progress.TranslateAndSetMessage("Checking backup path...")
	err = mebroutines.CreateFile(backupPath)
	if err != nil {
		mebroutines.ShowTranslatedErrorMessage("Could not write test backup file %s. Try another location.", backupPath)

		return fmt.Errorf("could not write test backup file: %w", err)
	}

	err = os.Remove(backupPath)
	if err != nil {
		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, fmt.Errorf("removing test backup file returned error code: %w", err))
	}

	// Get disk location
	progress.TranslateAndSetMessage("Getting disk location...")
	diskLocation := box.GetDiskLocation()
	log.Debug("Disk location: %s", diskLocation)
	if diskLocation == "" {
		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, errors.New("could not get disk location"))
	}

	progress.TranslateAndSetMessage("Checking for FAT32 filesystem...")
	err = checkForFATFilesystem(backupPath, diskLocation)
	if err != nil {
		mebroutines.ShowTranslatedErrorMessage("The backup file is too large for a FAT32 filesystem. Please reformat the backup disk as exFAT.")

		return fmt.Errorf("backup file too large for fat32 filesystem: %w", err)
	}

	// Make clone to path_backup
	progress.TranslateAndSetMessage("Please wait, writing backup...")
	err = box.WriteDiskClone(backupPath)
	if err != nil {
		return mebroutines.ShowTranslatedErrorMessageAndPassError(generalErrorString, fmt.Errorf("failed to make clone: %w", err))
	}

	return nil
}

func ensureBoxInstalledAndNotRunning() error {
	isInstalled, err := box.Installed()
	if err != nil {
		return fmt.Errorf("could not back up server as we could not detect whether vm is installed: %w", err)
	}

	if !isInstalled {
		return errors.New("no server has been installed")
	}

	isRunning, err := box.Running()
	if err != nil {
		return fmt.Errorf("could not back up server as we could not detect whether vm is running: %w", err)
	}

	if isRunning {
		return errors.New("could not back up server as the server is running")
	}

	return nil
}

func checkForFATFilesystem(backupPath string, vmDiskLocation string) error {
	// Check VM disk size
	mediumSizeMB, err := box.MediumSizeOnDisk(vmDiskLocation)

	// If we can't get medium size, we'll just ignore the error and continue.
	if err != nil {
		log.Error("Error getting VirtualBox medium size: %s", err)

		return nil
	}

	// FAT32 is enough to store this backup, so we don't need to check the filesystem.
	if mediumSizeMB < 4*1024 {
		return nil
	}

	// If there is an error checking whether the backup medium has a FAT32
	// filesystem, we'll just allow the user to continue. The user will
	// see an error eventually, if the backup disk actually is FAT32.
	isFAT32, err := isFAT32(backupPath)
	if err != nil {
		log.Error("Error checking if the backup medium has a FAT filesystem: %s", err)

		return nil
	}

	if isFAT32 {
		return errors.New("backup too large for FAT32")
	}

	return nil
}

// GetBackupFilename returns generated filename
func GetBackupFilename(timestamp time.Time) string {
	// Don't get confused by fixed date, this is correct. see: https://golang.org/src/time/format.go
	return timestamp.Format("2006-01-02_15-04-05.vmdk")
}
