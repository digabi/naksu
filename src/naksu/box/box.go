package box

// box gets information about the Vagrant box "default" directly from VirtualBox
// This package gives most up-to-date information about the server box

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

	"naksu/constants"
	"naksu/mebroutines"
	"naksu/xlate"
	"naksu/log"
)

// cache for VBoxMange showvminfo --machinereadable
// See getVMInfoRegexp
type cacheShowVMInfoType struct {
	output          string
	outputTimestamp int64
	updateStarted   int64
}

var cacheShowVMInfo cacheShowVMInfoType

func getVagrantBoxID() string {
	vagrantPath := mebroutines.GetVagrantDirectory()

	pathID := filepath.Join(vagrantPath, ".vagrant", "machines", "default", "virtualbox", "id")

	/* #nosec */
	fileContent, err := ioutil.ReadFile(pathID)
	if err != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Could not get vagrantbox ID: %d"), err))
		return ""
	}

	return string(fileContent)
}

// SetCacheShowVMInfo sets cacheShowVMInfo content
// It is exported for unit tests
func SetCacheShowVMInfo(newShowVMInfo string) {
	cacheShowVMInfo.output = newShowVMInfo
	cacheShowVMInfo.outputTimestamp = time.Now().Unix()
	cacheShowVMInfo.updateStarted = 0
}

// getVBoxManageOutput executes "VBoxManage showvminfo".
func getVBoxManageOutput() string {
	boxID := getVagrantBoxID()
	if boxID == "" {
		return ""
	}

	cacheShowVMInfo.updateStarted = time.Now().Unix()
	return mebroutines.RunVBoxManage([]string{"showvminfo", "--machinereadable", boxID})
}

// getVMInfoRegexp returns result of the given vmRegexp from the current VBoxManage showvminfo
// output. This function gets the output either from the cache or calls getVBoxManageOutput()
func getVMInfoRegexp(vmRegexp string) string {
	var vBoxManageOutput string

	// There is a avail version fetch going on (break free after 240 loops)
	// This locking avoids executing multiple instances of VBoxManage at the same time. Calling
	// VBoxManage simulaneously tends to cause E_ACCESSDENIED errors from VBoxManage.
	tryCounter := 0
	for (cacheShowVMInfo.updateStarted != 0) && (tryCounter < 240) {
		time.Sleep(500 * time.Millisecond)
		tryCounter++
		log.Debug(fmt.Sprintf("getVMIInfoRegexp is waiting 'VBoxManage showvminfo' to exit (race condition lock count %d)", tryCounter))
	}

	if cacheShowVMInfo.outputTimestamp < (time.Now().Unix() - constants.VBoxManageCacheTimeout) {
		// Cache is too old or not set

		vBoxManageOutput = getVBoxManageOutput()
		SetCacheShowVMInfo(vBoxManageOutput)
		log.Debug("getVMIInfoRegexp executed 'VBoxManage showvminfo' and updated cache")
	} else {
		vBoxManageOutput = cacheShowVMInfo.output
		log.Debug("getVMIInfoRegexp got 'VBoxManage showvminfo' output from the cache")
	}

	// Extract server name
	pattern := regexp.MustCompile(vmRegexp)
	result := pattern.FindStringSubmatch(vBoxManageOutput)

	if len(result) > 1 {
		return result[1]
	}

	return ""
}

// GetType returns the box type (e.g. "digabi/ktp-qa") of the current VM
func GetType() string {
	return getVMInfoRegexp("description=\"(.*?)\"")
}

// GetVersion returns the version string (e.g. "SERVER7108X v69") of the current VM
func GetVersion() string {
	return getVMInfoRegexp("name=\"(.*?)\"")
}

// GetDiskUUID returns the VirtualBox UUID for the image of the current VM
func GetDiskUUID() string {
	return getVMInfoRegexp("\"SATA Controller-ImageUUID-0-0\"=\"(.*?)\"")
}
