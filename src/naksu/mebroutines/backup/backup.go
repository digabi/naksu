package backup

import (
	"errors"
	"fmt"
	"io/ioutil"
	"naksu/mebroutines"
	"naksu/progress"
	"naksu/xlate"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// MakeBackup creates virtual maching backup to path
func MakeBackup(backupPath string) error {
	progress.TranslateAndSetMessage("Checking existing file...")
	if mebroutines.ExistsFile(backupPath) {
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("File %s already exists"), backupPath))
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
		mebroutines.LogDebug("Backup remove returned error code")
		return errors.New("removing backup returned error code")
	}

	// Get box
	progress.TranslateAndSetMessage("Getting vagrantbox ID...")
	boxID := getVagrantBoxID()
	mebroutines.LogDebug(fmt.Sprintf("Vagrantbox ID: %s", boxID))

	// Get disk UUID
	progress.TranslateAndSetMessage("Getting disk UUID...")
	diskUUID := getDiskUUID(boxID)
	mebroutines.LogDebug(fmt.Sprintf("Disk UUID: %s", diskUUID))

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

func getVagrantBoxID() string {
	vagrantPath := mebroutines.GetVagrantDirectory()

	pathID := filepath.Join(vagrantPath, ".vagrant", "machines", "default", "virtualbox", "id")

	/* #nosec */
	fileContent, err := ioutil.ReadFile(pathID)
	if err != nil {
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Could not get vagrantbox ID: %d"), err))
	}

	return string(fileContent)
}

func getDiskUUID(boxID string) string {
	vBoxManageOutput := mebroutines.RunVBoxManage([]string{"showvminfo", "-machinereadable", boxID})

	// Extract server disk image path
	pattern := regexp.MustCompile("\"SATA Controller-ImageUUID-0-0\"=\"(.*?)\"")
	result := pattern.FindStringSubmatch(vBoxManageOutput)

	if len(result) > 1 {
		return result[1]
	}

	// No match
	mebroutines.LogDebug(vBoxManageOutput)
	mebroutines.ShowErrorMessage(xlate.Get("Could not make backup: failed to get disk UUID"))

	return ""
}

func makeClone(diskUUID string, backupPath string) error {
	vBoxManageOutput := mebroutines.RunVBoxManage([]string{"clonemedium", diskUUID, backupPath})

	// Check whether clone was successful or not
	matched, errRe := regexp.MatchString("Clone medium created in format 'VMDK'", vBoxManageOutput)
	if errRe != nil || !matched {
		// Failure
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Could not back up disk %s to %s"), diskUUID, backupPath))
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
