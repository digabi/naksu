package box

// box gets information about the Vagrant box "default" directly from VirtualBox
// This package gives most up-to-date information about the server box

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"naksu/boxversion"
	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/xlate"
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

	if !mebroutines.ExistsFile(pathID) {
		return ""
	}

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
	return mebroutines.RunVBoxManage([]string{"showvminfo", "--machinereadable", boxID}, false)
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

func getVagrantfileData() (string, string) {
	pathVagrantfile := filepath.Join(mebroutines.GetVagrantDirectory(), "Vagrantfile")

	if !mebroutines.ExistsFile(pathVagrantfile) {
		log.Debug(fmt.Sprintf("There is no Vagrantfile '%s' so no box type/version can't be identified", pathVagrantfile))
		return "", ""
	}

	/* #nosec */
	fileContent, err := ioutil.ReadFile(pathVagrantfile)
	if err != nil {
		log.Debug(fmt.Sprintf("Could not read Vagrantfile from '%s': %v", pathVagrantfile, err))
		return "", ""
	}

	fileContentString := string(fileContent)

	boxType, boxVersion, err := boxversion.GetVagrantVersionDetails(fileContentString)
	if err != nil {
		return "", ""
	}

	return boxType, boxVersion
}

// GetType returns the box type (e.g. "digabi/ktp-qa") of the current VM
// If there is no current VM installed, get the value from ~/ktp/Vagrantfile
func GetType() string {
	result := getVMInfoRegexp("description=\"(.*?)\"")

	if result == "" {
		boxType, _ := getVagrantfileData()
		if boxType != "" {
			log.Debug(fmt.Sprintf("Got box type '%s' from Vagrantfile as VM did not give any value", boxType))
			result = boxType
		}
	}

	return result
}

// GetVersion returns the version string (e.g. "SERVER7108X v69") of the current VM
// If there is no current VM installed, get the value from ~/ktp/Vagrantfile
func GetVersion() string {
	result := getVMInfoRegexp("name=\"(.*?)\"")

	if result == "" {
		_, boxVersion := getVagrantfileData()
		if boxVersion != "" {
			log.Debug(fmt.Sprintf("Got box version '%s' from Vagrantfile as VM did not give any value", boxVersion))
			result = boxVersion
		}
	}

	return result
}

// GetDiskUUID returns the VirtualBox UUID for the image of the current VM
func GetDiskUUID() string {
	return getVMInfoRegexp("\"SATA Controller-ImageUUID-0-0\"=\"(.*?)\"")
}

// GetDiskLocation returns the full path of the current VM disk image.
func GetDiskLocation() string {
	return getVMInfoRegexp("\"SATA Controller-0-0\"=\"(.*)\"")
}

// MediumSizeOnDisk returns the size of the current VM disk image on disk
// (= the expected size of a VM backup) in megabytes.
func MediumSizeOnDisk(location string) (uint64, error) {
	// According to documentation, showmediuminfo should also accept a disk uuid
	// as a parameter, but that doesn't seem to be the case. To be safe, we'll
	// use the location of the disk instead.
	mediumInfo := mebroutines.RunVBoxManage([]string{"showmediuminfo", location}, false)
	sizeOnDiskRE := regexp.MustCompile(`Size on disk:\s+(\d+)\s+MBytes`)
	result := sizeOnDiskRE.FindStringSubmatch(mediumInfo)
	if len(result) > 1 {
		size := result[1]
		return strconv.ParseUint(size, 10, 64)
	}
	return 0, errors.New("failed to get medium size: no regex matches")
}
