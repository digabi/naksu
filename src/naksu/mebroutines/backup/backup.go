package backup

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"naksu/box"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/progress"
	"naksu/xlate"
)

// MakeBackup creates virtual maching backup to path
func MakeBackup(backupPath string) error {
	progress.TranslateAndSetMessage("Checking existing file...")
	if mebroutines.ExistsFile(backupPath) {
		mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("File %s already exists"), backupPath))
		return errors.New("backup file already exists")
	}

	// Check if path_backup is writeable
	progress.TranslateAndSetMessage("Checking backup path...")
	wrErr := mebroutines.CreateFile(backupPath)
	if wrErr != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Could not write backup file %s. Try another location."), backupPath))
		return errors.New("could not write backup file")
	}

	remErr := os.Remove(backupPath)
	if remErr != nil {
		log.Debug("Backup remove returned error code")
		return errors.New("removing backup returned error code")
	}

	// Get disk UUID
	progress.TranslateAndSetMessage("Getting disk UUID...")
	diskUUID := box.GetDiskUUID()
	log.Debug(fmt.Sprintf("Disk UUID: %s", diskUUID))
	if diskUUID == "" {
		mebroutines.ShowWarningMessage(xlate.Get("Could not make backup: failed to get disk UUID"))
		return errors.New("could not get disk uuid")
	}

	// Make clone to path_backup
	progress.TranslateAndSetMessage("Please wait, writing backup...")
	cloneErr := makeClone(diskUUID, backupPath)
	if cloneErr != nil {
		return errors.New("writing backup failed")
	}

	// Close backup media (detach it from VirtualBox disk management)
	progress.TranslateAndSetMessage("Detaching backup disk image...")
	deleteClone(backupPath)

	return nil
}

func makeClone(diskUUID string, backupPath string) error {
	vBoxManageOutput := mebroutines.RunVBoxManage([]string{"clonemedium", diskUUID, backupPath})

	// Check whether clone was successful or not
	matched, errRe := regexp.MatchString("Clone medium created in format 'VMDK'", vBoxManageOutput)
	if errRe != nil || !matched {
		// Failure
		mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Could not back up disk %s to %s"), diskUUID, backupPath))
		return errors.New("backup failed")
	}

	return nil
}

func deleteClone(backupPath string) {
	_ = mebroutines.RunVBoxManage([]string{"closemedium", backupPath})
}

// GetBackupFilename returns generated filename
func GetBackupFilename(timestamp time.Time) string {
	// Don't get confused by fixed date, this is correct. see: https://golang.org/src/time/format.go
	return timestamp.Format("2006-01-02_15-04-05.vmdk")
}
