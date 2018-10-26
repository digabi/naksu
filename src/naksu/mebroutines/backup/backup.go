package backup

import (
	"errors"
	"fmt"
	"io/ioutil"
	"naksu/mebroutines"
	"naksu/progress"
	"naksu/xlate"
	"os"
	"regexp"
	"time"
)

// MakeBackup creates virtual maching backup to path
func MakeBackup(backupPath string) {
	progress.TranslateAndSetMessage("Check that there is no existing backup file")
	if mebroutines.ExistsFile(backupPath) {
		mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("File %s already exists"), backupPath))
	}

	// Check if path_backup is writeable
	progress.TranslateAndSetMessage("Check that backup path is writeable")
	wrErr := mebroutines.CreateFile(backupPath)
	if wrErr != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Could not write backup file %s. Try another location."), backupPath))
		return
	}
	err := os.Remove(backupPath)
	if err != nil {
		mebroutines.LogDebug("Backup remove returned error code")
	}

	// Get box
	progress.TranslateAndSetMessage("Getting vagrantbox ID")
	boxID := getVagrantBoxID()
	mebroutines.LogDebug(fmt.Sprintf("Vagrantbox ID: %s", boxID))

	// Get disk UUID
	progress.TranslateAndSetMessage("Getting disk UUID")
	diskUUID := getDiskUUID(boxID)
	mebroutines.LogDebug(fmt.Sprintf("Disk UUID: %s", diskUUID))

	// Make clone to path_backup
	progress.TranslateAndSetMessage("Making backup. This takes a while.")
	cloneErr := makeClone(diskUUID, backupPath)
	if cloneErr != nil {
		progress.TranslateAndSetMessage("Backup failed.")
		return
	}

	// Close backup media (detach it from VirtualBox disk management)
	progress.TranslateAndSetMessage("Detaching backup disk from disk management")
	deleteClone(backupPath)

	progress.SetMessage(fmt.Sprintf(xlate.Get("Backup can be found at %s"), backupPath))
	mebroutines.ShowInfoMessage(fmt.Sprintf(xlate.Get("Backup has been made to %s"), backupPath))
}

func getVagrantBoxID() string {
	vagrantPath := mebroutines.GetVagrantDirectory()

	pathID := vagrantPath + string(os.PathSeparator) + ".vagrant" + string(os.PathSeparator) + "machines" + string(os.PathSeparator) + "default" + string(os.PathSeparator) + "virtualbox" + string(os.PathSeparator) + "id"

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
		return errors.New("Backup failed")
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
