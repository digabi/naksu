package box

// box gets information about the currently installed VM

import (
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strconv"

	semver "github.com/blang/semver/v4"

	"naksu/box/vboxmanage"
	"naksu/config"
	"naksu/constants"
	"naksu/host"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/xlate"
)

const (
	boxName           = "NaksuAbittiKTP"
	boxOSType         = "Debian"
	boxFinalImageSize = 55 * 1024 // VDI disk size in megs
	boxSnapshotName   = "Installed"
)

func calculateBoxCPUs() int {
	calculatedCores := host.GetCPUCoreCount() - 1

	if calculatedCores <= 2 {
		return 2
	}

	return calculatedCores
}

func calculateBoxMemory() (uint64, error) {
	hostMemory, err := host.GetMemory()

	if err != nil {
		return 0, fmt.Errorf("could not read system memory: %v", err)
	}

	freeVMMemory := uint64(math.Round(float64(hostMemory) * 0.74))
	lowVMMemoryLimit := uint64(math.Round((8192 - 1024) * 0.74))

	if freeVMMemory < lowVMMemoryLimit {
		return 0, fmt.Errorf("allocated vm memory %d is less than required minimum memory limit %d", freeVMMemory, lowVMMemoryLimit)
	}

	return freeVMMemory, nil
}

// CreateNewBox creates new VM using the given imagePath
func CreateNewBox(boxType string, ddImagePath string, boxVersion string) error {
	vdiImagePath := filepath.Join(mebroutines.GetHomeDirectory(), "ktp", "naksu_ktp_disk.vdi")
	if mebroutines.ExistsFile(vdiImagePath) {
		log.Debug(fmt.Sprintf("VDI file %s already exists", vdiImagePath))
		return fmt.Errorf("vdi file %s already exists", vdiImagePath)
	}

	calculatedBoxCPUs := calculateBoxCPUs()

	calculatedBoxMemory, errMemory := calculateBoxMemory()
	if errMemory != nil {
		return errMemory
	}

	log.Debug(fmt.Sprintf("Calculated new VM specs - CPUs: %d, Memory: %d", calculatedBoxCPUs, calculatedBoxMemory))

	createCommands := []vboxmanage.VBoxCommand{
		{"convertfromraw", ddImagePath, vdiImagePath, "--format", "VDI"},
		{"modifyhd", vdiImagePath, "--resize", fmt.Sprintf("%d", boxFinalImageSize)},
		{"createvm", "--name", boxName, "--register"},
		{
			"modifyvm", boxName,
			"--pae", "on",
			"--cpus", fmt.Sprintf("%d", calculatedBoxCPUs),
			"--memory", fmt.Sprintf("%d", calculatedBoxMemory),
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

	v6_1String := "6.1.0"
	v6_1, err := semver.Make(v6_1String)
	if err != nil {
		return fmt.Errorf("hard-coded version string %s could not be converted to sematic version object", v6_1String)
	}

	vBoxVersion, err := vboxmanage.GetVBoxManageVersion()
	if err != nil {
		log.Debug(fmt.Sprintf("Could not get VBoxManage version: %v", err))
		return err
	}

	if vBoxVersion.LT(v6_1) {
		createCommands = append(createCommands, vboxmanage.VBoxCommand{"modifyvm", boxName, "--clipboard", "bidirectional"})
	} else {
		createCommands = append(createCommands, vboxmanage.VBoxCommand{"modifyvm", boxName, "--clipboard-mode", "bidirectional"})
	}

	createCommands = append(createCommands, vboxmanage.VBoxCommand{"snapshot", boxName, "take", boxSnapshotName})

	err = vboxmanage.RunCommands(createCommands)
	if err != nil {
		return err
	}

	vboxmanage.ResetVBoxResponseCache()

	return nil
}

// StartCurrentBox starts currently installed VM
func StartCurrentBox() error {
	startCommands := []vboxmanage.VBoxCommand{
		{"modifyvm", boxName, "--nic1", "bridged"},
		{"modifyvm", boxName, "--bridgeadapter1", config.GetExtNic()},
		{"modifyvm", boxName, "--nictype1", config.GetNic()},
		{"startvm", boxName, "--type", "gui"},
	}

	return vboxmanage.RunCommands(startCommands)
}

// RestoreSnapshot returns installed VM to fresh state (to the snapshot taken just after the install)
func RestoreSnapshot() error {
	restoreCommands := []vboxmanage.VBoxCommand{
		{"snapshot", boxName, "restore", boxSnapshotName},
	}

	return vboxmanage.RunCommands(restoreCommands)
}

// RemoveCurrentBox deletes currently installed VM
func RemoveCurrentBox() error {
	removeCommands := []vboxmanage.VBoxCommand{
		{"unregistervm", boxName, "--delete"},
	}

	return vboxmanage.RunCommands(removeCommands)
}

// WriteDiskClone creates a disk clone of the first disk of the current VM
func WriteDiskClone(clonePath string) error {
	diskUUID := getDiskUUID()
	if diskUUID == "" {
		return fmt.Errorf("could not get disk uuid")
	}

	vBoxManageOutput, err := vboxmanage.RunCommand(vboxmanage.VBoxCommand{"clonemedium", diskUUID, clonePath, "--format", "VMDK"})

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
	_, errCloseMedium := vboxmanage.RunCommand(vboxmanage.VBoxCommand{"closemedium", clonePath})
	return errCloseMedium
}

// Installed returns true if we have box installed, otherwise false
func Installed() (bool, error) {
	isInstalled, err := vboxmanage.Installed(boxName)

	if err == nil {
		log.Debug(fmt.Sprintf("Server '%s' installed: %t", boxName, isInstalled))
	} else {
		log.Debug(fmt.Sprintf("box.Installed() could not detect whether VM is installed: %v", err))
	}

	return isInstalled, err
}

func Running() (bool, error) {
	isRunning, err := vboxmanage.Running(boxName)

	if err == nil {
		log.Debug(fmt.Sprintf("Server '%s' running: %t", boxName, isRunning))
	} else {
		log.Debug(fmt.Sprintf("box.Running() could not detect whether VM is running: %v", err))
	}

	return isRunning, err
}

// GetType returns the box type (e.g. "digabi/ktp-qa") of the current VM
func GetType() string {
	return vboxmanage.GetBoxProperty(boxName, "boxType")
}

// GetTypeLegend returns an user-readable type legend of the current VM
func GetTypeLegend() string {
	if TypeIsAbitti() {
		return xlate.Get("Abitti server")
	}

	if TypeIsMatriculationExam() {
		return xlate.Get("Matric Exam server")
	}

	// Unknown box type
	log.Debug(fmt.Sprintf("Warning: We have a type string '%s' which does not resolve to Abitti/Matriculation box type (GetTypeLegend)", GetType()))
	return "-"
}

// TypeIsAbitti returns true if currently installed box is Abitti box
func TypeIsAbitti() bool {
	boxType := GetType()

	return (boxType == constants.AbittiBoxType)
}

// TypeIsMatriculationExam returns true if currently installed box is Matriculation Exam box
func TypeIsMatriculationExam() bool {
	boxType := GetType()

	return (boxType == constants.ExamBoxType)
}

// GetVersion returns the version string (e.g. "SERVER7108X v69") of the current VM
func GetVersion() string {
	return vboxmanage.GetBoxProperty(boxName, "boxVersion")
}

// getDiskUUID returns the VirtualBox UUID for the image of the current VM
func getDiskUUID() string {
	return vboxmanage.GetVMInfoByRegexp(boxName, "\"SATA Controller-ImageUUID-0-0\"=\"(.*?)\"")
}

// GetDiskLocation returns the full path of the current VM disk image.
func GetDiskLocation() string {
	return vboxmanage.GetVMInfoByRegexp(boxName, "\"SATA Controller-0-0\"=\"(.*)\"")
}

// GetLogDir returns the full path of VirtualBox log directory
func GetLogDir() string {
	return vboxmanage.GetVMInfoByRegexp(boxName, "LogFldr=\"(.*)\"")
}

// MediumSizeOnDisk returns the size of the current VM disk image on disk
// (= the expected size of a VM backup) in megabytes.
func MediumSizeOnDisk(location string) (uint64, error) {
	// According to documentation, showmediuminfo should also accept a disk uuid
	// as a parameter, but that doesn't seem to be the case. To be safe, we'll
	// use the location of the disk instead.

	mediumInfo, err := vboxmanage.RunCommand([]string{"showmediuminfo", location})

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
