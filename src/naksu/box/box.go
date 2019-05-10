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
)

// cache for VBoxMange showvminfo --machinereadable
// See getVMInfoRegexp
type cacheShowVMInfoType struct {
	output    string
	timestamp int64
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
	cacheShowVMInfo.timestamp = time.Now().Unix()
}

func getVMInfoRegexp(vmRegexp string) string {
	var vBoxManageOutput string

	// Get showvminfo from cache/by executing VBoxManage
	// Running VBoxManage showvminfo too often makes VBoxManage to exit with code 1:
	//   VBoxManage: error: The object is not ready
	//   VBoxManage: error: Details: code E_ACCESSDENIED (0x80070005), component SessionMachine, interface IMachine, callee nsISupports"
	if cacheShowVMInfo.timestamp < (time.Now().Unix() - constants.VBoxManageCacheTimeout) {
		boxID := getVagrantBoxID()
		if boxID == "" {
			return ""
		}

		vBoxManageOutput = mebroutines.RunVBoxManage([]string{"showvminfo", "--machinereadable", boxID})
		SetCacheShowVMInfo(vBoxManageOutput)
	} else {
		vBoxManageOutput = cacheShowVMInfo.output
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
