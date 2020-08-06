package box

// box gets information about the Vagrant box "default" directly from VirtualBox
// This package gives most up-to-date information about the server box

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/paulusrobin/go-memory-cache/memory-cache"

	"naksu/config"
	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
	vbm "naksu/box/vboxmanage"
)

const (
	boxName           = "NaksuAbittiKTP"
	boxOSType         = "Debian"
	boxCPUs           = 2
	boxMemory         = 4096  // RAM in megs
	boxFinalImageSize = 56909 // VDI disk size in megs
	boxType           = "digabi/ktp-qa"
	boxVersion        = "SERVER7108X v69"
)

var vBoxResponseCache memory_cache.Cache
var vBoxManageStarted int64

func callRunVBoxManage(args []string) (string, error) {
	// There is an ongoing VBoxManage call (break free after 240 loops)
	// This locking avoids executing multiple instances of VBoxManage at the same time. Calling
	// VBoxManage simulaneously tends to cause E_ACCESSDENIED errors from VBoxManage.
	tryCounter := 0
	for (vBoxManageStarted != 0) && (tryCounter < 240) {
		time.Sleep(500 * time.Millisecond)
		tryCounter++
		log.Debug(fmt.Sprintf("callRunVBoxManage is waiting VBoxManage to exit (race condition lock count %d)", tryCounter))
	}

	vBoxManageStarted = time.Now().Unix()
	vBoxManageOutput, err := mebroutines.RunVBoxManage(args)
	vBoxManageStarted = 0

	return vBoxManageOutput, err
}

func ensureVBoxResponseCacheInitialised() {
	var err error

	if vBoxResponseCache == nil {
		vBoxResponseCache, err = memory_cache.New()
		if err != nil {
			log.Debug(fmt.Sprintf("Fatal error: Failed to initialise memory cache: %v", err))
			panic(err)
		}
	}
}

func resetVBoxResponseCache() {
	vBoxResponseCache = nil
	ensureVBoxResponseCacheInitialised()
}

// getVMInfoRegexp returns result of the given vmRegexp from the current VBoxManage showvminfo
// output. This function gets the output either from the cache or calls getVBoxManageOutput()
func getVMInfoRegexp(vmRegexp string) string {
	var rawVMInfo string

	ensureVBoxResponseCacheInitialised()

	rawVMInfoInterface, err := vBoxResponseCache.Get("showvminfo")
	if err != nil {
		rawVMInfo, err = callRunVBoxManage([]string{"showvminfo", "--machinereadable", boxName})
		if err != nil {
			log.Debug(fmt.Sprintf("Could not get VM info: %v", err))
			rawVMInfo = ""
		}

		errCache := vBoxResponseCache.Set("showvminfo", rawVMInfo, constants.VBoxManageCacheTimeout)
		if errCache != nil {
			log.Debug(fmt.Sprintf("Could not store VM info to cache: %v", errCache))
		}
	} else {
		rawVMInfo = fmt.Sprintf("%v", rawVMInfoInterface)
	}

	// Extract server name
	pattern := regexp.MustCompile(vmRegexp)
	result := pattern.FindStringSubmatch(rawVMInfo)

	if len(result) > 1 {
		return result[1]
	}

	return ""
}

// CreateNewBox creates new VM using the given imagePath
func CreateNewBox(ddImagePath string) error {
	_, err := callRunVBoxManage([]string{"createvm", "--name", boxName, "--register"})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not create machine: %v", err))
		return err
	}

	_, err = callRunVBoxManage([]string{
		"modifyvm", boxName,
		"--pae", "on",
		"--cpus", fmt.Sprintf("%d", boxCPUs),
		"--memory", fmt.Sprintf("%d", boxMemory),
		"--acpi", "on",
		"--ioapic", "on",
		"--ostype", boxOSType,
		"--firmware", "efi",
		"--audio", "none",
	})

	if err != nil {
		log.Debug(fmt.Sprintf("Failed to set CPUs, memory, etc: %v", err))
		return err
	}

	_, err = callRunVBoxManage([]string{
		"guestproperty", "set", boxName,
		"boxType", boxType,
	})
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to SetGuestProperty: %v", err))
		return err
	}

	_, err = callRunVBoxManage([]string{
		"guestproperty", "set", boxName,
		"boxVersion", boxVersion,
	})
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to SetGuestProperty: %v", err))
		return err
	}

	vdiImagePath := filepath.Join(mebroutines.GetHomeDirectory(), "ktp", "naksu_ktp_disk.vdi")
	if mebroutines.ExistsFile(vdiImagePath) {
		log.Debug(fmt.Sprintf("Existing VDI file %s already exists", vdiImagePath))
		return fmt.Errorf("existing vdi file %s already exists", vdiImagePath)
	}

	err = createNewBoxMedia(ddImagePath, vdiImagePath)
	if err != nil {
		return err
	}

	_, err = callRunVBoxManage([]string{
		"sharedfolder", "add", boxName,
		"--name", "media_usb1",
		"--hostpath", mebroutines.GetMebshareDirectory(),
	})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not set shared folder: %v", err))
		return err
	}

	resetVBoxResponseCache()

	return nil
}

// createNewBoxMedia sets media to new box
// this is called by CreateNewBox()
func createNewBoxMedia(ddImagePath string, vdiImagePath string) error {
	_, err := callRunVBoxManage([]string{
		"convertfromraw", ddImagePath, vdiImagePath,
		"--format", "VDI",
	})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not convert image file %s to VDI file: %v", ddImagePath, err))
		return err
	}

	_, err = callRunVBoxManage([]string{
		"modifyhd", vdiImagePath,
		"--resize", fmt.Sprintf("%d", boxFinalImageSize),
	})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not enlarge VDI file %s: %v", vdiImagePath, err))
		return err
	}

	_, err = callRunVBoxManage([]string{
		"storagectl", boxName,
		"--add", "sata",
		"--name", "SATA Controller",
	})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not add SATA controller to VM: %v", err))
		return err
	}

	_, err = callRunVBoxManage([]string{
		"storageattach", boxName,
		"--storagectl", "SATA Controller",
		"--port", "0",
		"--device", "0",
		"--type", "hdd",
		"--medium", vdiImagePath,
	})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not attach media to SATA controller: %v", err))
		return err
	}

	return nil
}

// StartCurrentBox starts currently installed VM
func StartCurrentBox() error {
	_, err := callRunVBoxManage([]string{"modifyvm", boxName, "--nic1", "bridged"})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not start VM, as we could not set VM network mode to bridged: %v", err))
		return err
	}

	_, err = callRunVBoxManage([]string{"modifyvm", boxName, "--bridgeadapter1", config.GetExtNic()})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not start VM, as we could not set %s to VM bridged adapter: %v", config.GetExtNic(), err))
		return err
	}

	_, err = callRunVBoxManage([]string{"modifyvm", boxName, "--nictype1", config.GetNic()})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not start VM, as we could not set VM network type: %v", err))
		return err
	}

	_, err = callRunVBoxManage([]string{"startvm", boxName, "--type", "gui"})
	if err != nil {
		log.Debug(fmt.Sprintf("Could not start VM: %v", err))
		return err
	}

	log.Debug("VM was started successfully")

	return nil
}

func getBoxProperty(property string) string {
	ensureVBoxResponseCacheInitialised()

	propertyValue := ""

	propertyValueInterface, errCache := vBoxResponseCache.Get(property)
	if errCache != nil {
		output, errVBoxManage := callRunVBoxManage([]string{"guestproperty", "get", boxName, property})
		if errVBoxManage != nil {
			log.Debug(fmt.Sprintf("Could not get VM guest property: %v", errVBoxManage))
			output = ""
		}

		propRegexp := regexp.MustCompile(`Value: (.+)`)
		propMatches := propRegexp.FindStringSubmatch(output)
		if len(propMatches) == 2 {
			propertyValue = propMatches[1]
		}

		errCacheSet := vBoxResponseCache.Set(property, propertyValue, constants.VBoxManageCacheTimeout)
		if errCacheSet == nil {
			log.Debug(fmt.Sprintf("Stored VM guest property %s to cache: %s", property, propertyValue))
		} else {
			log.Debug(fmt.Sprintf("Could not store VM guest property %s, value %s to cache: %v", property, propertyValue, errCacheSet))
		}
	} else {
		propertyValue = fmt.Sprintf("%v", propertyValueInterface)
		log.Debug(fmt.Sprintf("Got VM guest property %s from cache: %s", property, propertyValue))
	}

	return propertyValue
}

// GetType returns the box type (e.g. "digabi/ktp-qa") of the current VM
func GetType() string {
	return getBoxProperty("boxType")
}

// GetVersion returns the version string (e.g. "SERVER7108X v69") of the current VM
func GetVersion() string {
	return getBoxProperty("boxVersion")
}

// GetDiskUUID returns the VirtualBox UUID for the image of the current VM
func GetDiskUUID() string {
	return vbm.GetVMInfoRegexp(boxName, "\"SATA Controller-ImageUUID-0-0\"=\"(.*?)\"")
}

// GetDiskLocation returns the full path of the current VM disk image.
func GetDiskLocation() string {
	return vbm.GetVMInfoRegexp(boxName, "\"SATA Controller-0-0\"=\"(.*)\"")
}

// GetLogDir returns the full path of VirtualBox log directory
func GetLogDir() string {
	return vbm.GetVMInfoRegexp(boxName, "LogFldr=\"(.*)\"")
}

// MediumSizeOnDisk returns the size of the current VM disk image on disk
// (= the expected size of a VM backup) in megabytes.
func MediumSizeOnDisk(location string) (uint64, error) {
	// According to documentation, showmediuminfo should also accept a disk uuid
	// as a parameter, but that doesn't seem to be the case. To be safe, we'll
	// use the location of the disk instead.

	mediumInfo, err := callRunVBoxManage([]string{"showmediuminfo", location})

	if err != nil {
		log.Debug(fmt.Sprintf("Could not get medium info to calculate its size: %v", err))
		return 0, errors.New("failed to get medium size: could not execute vboxmanage")
	}

	sizeOnDiskRE := regexp.MustCompile(`Size on disk:\s+(\d+)\s+MBytes`)
	result := sizeOnDiskRE.FindStringSubmatch(mediumInfo)
	if len(result) > 1 {
		size := result[1]
		return strconv.ParseUint(size, 10, 64)
	}
	return 0, errors.New("failed to get medium size: no regex matches")
}
