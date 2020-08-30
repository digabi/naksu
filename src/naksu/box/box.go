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

	vbm "naksu/box/vboxmanage"
	"naksu/config"
	"naksu/log"
	"naksu/mebroutines"
)

const (
	boxName           = "NaksuAbittiKTP"
	boxOSType         = "Debian"
	boxCPUs           = 2
	boxMemory         = 4096  // RAM in megs
	boxFinalImageSize = 56909 // VDI disk size in megs
	boxType           = "digabi/ktp-qa"
	boxVersion        = "SERVER7108X v69"
	boxSnapshotName   = "Installed"
)

// CreateNewBox creates new VM using the given imagePath
func CreateNewBox(ddImagePath string) error {
	vdiImagePath := filepath.Join(mebroutines.GetHomeDirectory(), "ktp", "naksu_ktp_disk.vdi")
	if mebroutines.ExistsFile(vdiImagePath) {
		log.Debug(fmt.Sprintf("Existing VDI file %s already exists", vdiImagePath))
		return fmt.Errorf("existing vdi file %s already exists", vdiImagePath)
	}

	createCommands := []vbm.VBoxCommand{
		{"convertfromraw", ddImagePath, vdiImagePath, "--format", "VDI"},
		{"modifyhd", vdiImagePath, "--resize", fmt.Sprintf("%d", boxFinalImageSize)},
		{"createvm", "--name", boxName, "--register"},
		{
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
		{
			"guestproperty", "set", boxName,
			"boxType", boxType,
		},
		{
			"guestproperty", "set", boxName,
			"boxVersion", boxVersion,
		},
		{
			"sharedfolder", "add", boxName,
			"--name", "media_usb1",
			"--hostpath", mebroutines.GetMebshareDirectory(),
		},
		{
			"storagectl", boxName,
			"--add", "sata",
			"--name", "SATA Controller",
		},
		{
			"storageattach", boxName,
			"--storagectl", "SATA Controller",
			"--port", "0",
			"--device", "0",
			"--type", "hdd",
			"--medium", vdiImagePath,
		},
	}

	v6_1, _ := semver.Make("6.1.0")
	vBoxVersion, errVBoxManageVersion := vbm.GetVBoxManageVersion()
	if errVBoxManageVersion != nil {
		log.Debug(fmt.Sprintf("Could not get VBoxManage version: %v", errVBoxManageVersion))
		return errVBoxManageVersion
	}

	if vBoxVersion.LT(v6_1) {
		// Override bidirectional clipboard for VirtualBox pre-6.1
		createCommands = append(createCommands, vbm.VBoxCommand{"modifyvm", boxName, "--clipboard", "bidirectional"})
	} else {
		// Set bidirectional clipboard for VirtualBox 6.1 and later
		createCommands = append(createCommands, vbm.VBoxCommand{"modifyvm", boxName, "--clipboard-mode", "bidirectional"})
	}

	createCommands = append(createCommands, vbm.VBoxCommand{"snapshot", boxName, "take", boxSnapshotName})

	err := vbm.MultipleCallRunVBoxManage(createCommands)
	if err != nil {
		return err
	}

	vbm.ResetVBoxResponseCache()

	return nil
}

// StartCurrentBox starts currently installed VM
func StartCurrentBox() error {
	startCommands := []vbm.VBoxCommand{
		{"modifyvm", boxName, "--nic1", "bridged"},
		{"modifyvm", boxName, "--bridgeadapter1", config.GetExtNic()},
		{"modifyvm", boxName, "--nictype1", config.GetNic()},
		{"startvm", boxName, "--type", "gui"},
	}

	return vbm.MultipleCallRunVBoxManage(startCommands)
}

// RestoreSnapshot returns installed VM to fresh state (to the snapshot taken just after the install)
func RestoreSnapshot() error {
	restoreCommands := []vbm.VBoxCommand{
		{"snapshot", boxName, "restore", boxSnapshotName},
	}

	return vbm.MultipleCallRunVBoxManage(restoreCommands)
}

// WriteDiskClone creates a disk clone of the first disk of the current VM
func WriteDiskClone(clonePath string) error {
	diskUUID := getDiskUUID()
	if diskUUID == "" {
		return fmt.Errorf("could not get disk uuid")
	}

	vBoxManageOutput, err := vbm.CallRunVBoxManage(vbm.VBoxCommand{"clonemedium", diskUUID, clonePath, "--format", "VMDK"})

	if err != nil {
		return err
	}

	// Check whether clone was successful or not
	matched, errRe := regexp.MatchString("Clone medium created in format 'VMDK'", vBoxManageOutput)
	if errRe != nil || !matched {
		// Failure
		log.Debug("VBoxManage output does not report successful clone in format 'VMDK'")
		return errors.New("could not get correct response from vboxmanage")
	}

	// Detach media from VirtualBox disk management
	_, errCloseMedium := vbm.CallRunVBoxManage(vbm.VBoxCommand{"closemedium", clonePath})
	return errCloseMedium
}

// Installed returns true if we have box installed, otherwise false
func Installed() (bool, error) {
	isInstalled, err := vbm.Installed(boxName)

	if err == nil {
		log.Debug(fmt.Sprintf("Server '%s' installed: %t", boxName, isInstalled))
	} else {
		log.Debug(fmt.Sprintf("box.Installed() could not detect whether VM is installed: %v", err))
	}

	return isInstalled, nil
}

// Running returns true if current box is installed and running, otherwise false
func Running() (bool, error) {
	isInstalled, errInstalled := Installed()
	isRunning, errRunning := vbm.Running(boxName)

	isInstalledAndRunning := isInstalled && isRunning

	if errInstalled == nil {
		log.Debug(fmt.Sprintf("box.Running(): VM installed: %t", isInstalled))
	} else {
		log.Debug(fmt.Sprintf("box.Running() could not detect whether VM is installed: %v", errInstalled))
		return isInstalledAndRunning, errInstalled
	}

	if errRunning == nil {
		log.Debug(fmt.Sprintf("box.Running(): VM running: %t", isRunning))
	} else {
		log.Debug(fmt.Sprintf("box.Running() could not detect whether VM is running: %v", errRunning))
		return isInstalledAndRunning, errRunning
	}

	log.Debug(fmt.Sprintf("box.Running(): Server '%s' is installed and running: %t", boxName, isInstalledAndRunning))

	return isInstalledAndRunning, nil
}

// GetType returns the box type (e.g. "digabi/ktp-qa") of the current VM
func GetType() string {
	return vbm.GetBoxProperty(boxName, "boxType")
}

// GetVersion returns the version string (e.g. "SERVER7108X v69") of the current VM
func GetVersion() string {
	return vbm.GetBoxProperty(boxName, "boxVersion")
}

// getDiskUUID returns the VirtualBox UUID for the image of the current VM
func getDiskUUID() string {
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
