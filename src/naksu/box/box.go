package box

// box gets information about the Vagrant box "default" directly from VirtualBox
// This package gives most up-to-date information about the server box

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	semver "github.com/blang/semver/v4"

	"naksu/config"
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


// CreateNewBox creates new VM using the given imagePath
func CreateNewBox(ddImagePath string) error {
	vdiImagePath := filepath.Join(mebroutines.GetHomeDirectory(), "ktp", "naksu_ktp_disk.vdi")
	if mebroutines.ExistsFile(vdiImagePath) {
		log.Debug(fmt.Sprintf("Existing VDI file %s already exists", vdiImagePath))
		return fmt.Errorf("existing vdi file %s already exists", vdiImagePath)
	}

	createCommands := []vbm.VBoxCommand{
		vbm.VBoxCommand{ "convertfromraw", ddImagePath, vdiImagePath, "--format", "VDI" },
		vbm.VBoxCommand{ "modifyhd", vdiImagePath, "--resize", fmt.Sprintf("%d", boxFinalImageSize) },
		vbm.VBoxCommand{ "createvm", "--name", boxName, "--register" },
		vbm.VBoxCommand{
			"modifyvm", boxName,
			"--pae", "on",
			"--cpus", fmt.Sprintf("%d", boxCPUs),
			"--memory", fmt.Sprintf("%d", boxMemory),
			"--acpi", "on",
			"--ioapic", "on",
			"--ostype", boxOSType,
			"--firmware", "efi",
			"--audio", "none",
		},
		vbm.VBoxCommand{
			"guestproperty", "set", boxName,
			"boxType", boxType,
		},
		vbm.VBoxCommand{
			"guestproperty", "set", boxName,
			"boxVersion", boxVersion,
		},
		vbm.VBoxCommand{
			"sharedfolder", "add", boxName,
			"--name", "media_usb1",
			"--hostpath", mebroutines.GetMebshareDirectory(),
		},
		vbm.VBoxCommand{
			"storagectl", boxName,
			"--add", "sata",
			"--name", "SATA Controller",
		},
		vbm.VBoxCommand{
			"storageattach", boxName,
			"--storagectl", "SATA Controller",
			"--port", "0",
			"--device", "0",
			"--type", "hdd",
			"--medium", vdiImagePath,
		},
	}

	err := vbm.MultipleCallRunVBoxManage(createCommands)
	if err != nil {
		return err
	}

	v6_1, _ := semver.Make("6.1.0")
	vBoxVersion, errVBoxManageVersion := vbm.GetVBoxManageVersion()
	if errVBoxManageVersion != nil {
		log.Debug(fmt.Sprintf("Could not get VBoxManage version: %v", errVBoxManageVersion))
		return errVBoxManageVersion
	}

	// Set bidirectional clipboard for VirtualBox 6.1 and later
	clipboardCommands := []vbm.VBoxCommand{
		vbm.VBoxCommand{"modifyvm", boxName, "--clipboard-mode", "bidirectional"},
	}

	if vBoxVersion.LT(v6_1) {
		// Override bidirectional clipboard for VirtualBox pre-6.1
		clipboardCommands = []vbm.VBoxCommand{
			vbm.VBoxCommand{"modifyvm", boxName, "--clipboard", "bidirectional"},
		}
	}

	err = vbm.MultipleCallRunVBoxManage(clipboardCommands)
	if err != nil {
		return err
	}

	vbm.ResetVBoxResponseCache()

	return nil
}

// StartCurrentBox starts currently installed VM
func StartCurrentBox() error {
	startCommands := []vbm.VBoxCommand{
		vbm.VBoxCommand{ "modifyvm", boxName, "--nic1", "bridged" },
		vbm.VBoxCommand{ "modifyvm", boxName, "--bridgeadapter1", config.GetExtNic() },
		vbm.VBoxCommand{ "modifyvm", boxName, "--nictype1", config.GetNic() },
		vbm.VBoxCommand{ "startvm", boxName, "--type", "gui" },
	}

	return vbm.MultipleCallRunVBoxManage(startCommands)
}

// GetType returns the box type (e.g. "digabi/ktp-qa") of the current VM
func GetType() string {
	return vbm.GetBoxProperty(boxName, "boxType")
}

// GetVersion returns the version string (e.g. "SERVER7108X v69") of the current VM
func GetVersion() string {
	return vbm.GetBoxProperty(boxName, "boxVersion")
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

	mediumInfo, err := vbm.CallRunVBoxManage([]string{"showmediuminfo", location})

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
