package box

// box gets information about the currently installed VM

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"

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
	boxName                 = "NaksuAbittiKTP"
	boxOSType               = "Debian"
	boxFinalImageSize       = 55 * 1024 // VDI disk size in megs
	boxVRamSize             = 24        // Video RAM size in megs
	boxSnapshotName         = "Installed"
	boxMinimumNumberOfCores = 2
	boxMemorySizePercentage = 0.74        // 0.74 = box RAM size will be 74% of the host RAM size
	boxLowMemoryLimit       = 8192 - 1024 // 8G minus 1G for display adapter
)

var boxInstalled bool

func calculateBoxCPUs() (int, error) {
	detectedCores, err := host.GetCPUCoreCount()
	if err != nil {
		return 0, err
	}

	calculatedCores := detectedCores - 1

	if calculatedCores <= boxMinimumNumberOfCores {
		return boxMinimumNumberOfCores, nil
	}

	return calculatedCores, nil
}

func calculateBoxMemory() (uint64, error) {
	hostMemory, err := host.GetMemory()

	if err != nil {
		return 0, fmt.Errorf("could not read system memory: %w", err)
	}

	freeVMMemory := uint64(math.Round(float64(hostMemory) * boxMemorySizePercentage))
	lowVMMemoryLimit := uint64(math.Round(float64(boxLowMemoryLimit) * boxMemorySizePercentage))

	if freeVMMemory < lowVMMemoryLimit {
		return 0, fmt.Errorf("allocated vm memory %d is less than required minimum memory limit %d", freeVMMemory, lowVMMemoryLimit)
	}

	return freeVMMemory, nil
}

func getCreateNewBoxBasicCommands(boxName string, boxType string, boxVersion string, calculatedBoxCPUs int, calculatedBoxMemory uint64) []vboxmanage.VBoxCommand {
	createCommands := []vboxmanage.VBoxCommand{
		{"convertfromraw", mebroutines.GetImagePath(), mebroutines.GetVDIImagePath(), "--format", "VDI"},
		{"modifyhd", mebroutines.GetVDIImagePath(), "--resize", fmt.Sprintf("%d", boxFinalImageSize)},
		{"createvm", "--name", boxName, "--register"},
		{
			"modifyvm", boxName,
			"--pae", "on",
			"--cpus", fmt.Sprintf("%d", calculatedBoxCPUs),
			"--memory", fmt.Sprintf("%d", calculatedBoxMemory),
			"--vram", fmt.Sprintf("%d", boxVRamSize),
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
			"--medium", mebroutines.GetVDIImagePath(),
		},
		{
			"setextradata", boxName,
			"GUI/RestrictedCloseActions",
			"SaveState,PowerOffRestoringSnapshot",
		},
	}

	return createCommands
}

func getCreateNewBoxClipboadCommand(vBoxVersion semver.Version) (vboxmanage.VBoxCommand, error) {
	v6_1String := "6.1.0"
	v6_1, err := semver.Make(v6_1String)
	if err != nil {
		return nil, fmt.Errorf("hard-coded version string %s could not be converted to sematic version object", v6_1String)
	}

	// Defaults to 6.1 or newer
	clipboardCommand := vboxmanage.VBoxCommand{"modifyvm", boxName, "--clipboard-mode", "bidirectional"}

	if vBoxVersion.LT(v6_1) {
		clipboardCommand = vboxmanage.VBoxCommand{"modifyvm", boxName, "--clipboard", "bidirectional"}
	}

	return clipboardCommand, nil
}

// CreateNewBox creates new VM using the given imagePath
func CreateNewBox(boxType string, boxVersion string) error {
	if mebroutines.ExistsFile(mebroutines.GetVDIImagePath()) {
		err := os.Remove(mebroutines.GetVDIImagePath())
		if err != nil {
			return fmt.Errorf("could not remove old vdi file %s: %w", mebroutines.GetVDIImagePath(), err)
		}
		log.Debug("Removed existing VDI file %s", mebroutines.GetVDIImagePath())
	}

	calculatedBoxCPUs, err := calculateBoxCPUs()
	if err != nil {
		return err
	}

	calculatedBoxMemory, errMemory := calculateBoxMemory()
	if errMemory != nil {
		return errMemory
	}

	log.Debug("Calculated new VM specs - CPUs: %d, Memory: %d", calculatedBoxCPUs, calculatedBoxMemory)

	createCommands := getCreateNewBoxBasicCommands(boxName, boxType, boxVersion, calculatedBoxCPUs, calculatedBoxMemory)

	vBoxVersion, err := vboxmanage.GetVBoxManageVersion()
	if err != nil {
		log.Error("Could not get VBoxManage version: %v", err)

		return err
	}

	clipboardCommand, err := getCreateNewBoxClipboadCommand(vBoxVersion)
	if err != nil {
		log.Error("Could not get new box clipboard creation command: %v", err)

		return err
	}
	createCommands = append(createCommands, clipboardCommand)

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

// StartEnvironmentStatusUpdate starts periodically updating given
// environmentStatus.BoxInstalled and .BoxRunning values
func StartEnvironmentStatusUpdate(environmentStatus *constants.EnvironmentStatus, tickerDuration time.Duration) {
	ticker := time.NewTicker(tickerDuration)

	go func() {
		for {
			<-ticker.C

			boxInstalledStatus, err := Installed()
			if err != nil {
				log.Error("Could not query whether VM is installed: %v", err)
			} else {
				boxInstalled = (err == nil) && boxInstalledStatus
				environmentStatus.BoxInstalled = boxInstalled
			}

			boxRunning, boxRunningErr := Running()
			if boxRunningErr != nil {
				log.Error("Could not query whether VM is running: %v", boxRunningErr)
			} else {
				environmentStatus.BoxRunning = (boxRunningErr == nil) && boxRunning
			}
		}
	}()
}

// Installed returns true if we have box installed, otherwise false
func Installed() (bool, error) {
	isInstalled, err := vboxmanage.IsVMInstalled(boxName)

	if err != nil {
		log.Error("box.Installed() could not detect whether VM is installed: %v", err)
	}

	return isInstalled, err
}

func Running() (bool, error) {
	if !boxInstalled {
		return false, nil
	}

	isRunning, err := vboxmanage.IsVMRunning(boxName)

	if err != nil {
		log.Error("box.Running() could not detect whether VM is running: %v", err)
	}

	return isRunning, err
}

// GetType returns the box type (e.g. "digabi/ktp-qa") of the current VM
func GetType() string {
	if !boxInstalled {
		return ""
	}

	return vboxmanage.GetVMProperty(boxName, "boxType")
}

// GetTypeLegend returns an user-readable type legend of the current VM
func GetTypeLegend() string {
	if !boxInstalled {
		return "-"
	}

	if TypeIsAbitti() {
		return xlate.Get("Abitti server")
	}

	if TypeIsMatriculationExam() {
		return xlate.Get("Matric Exam server")
	}

	// Unknown box type
	log.Warning("Warning: We have a type string '%s' which does not resolve to Abitti/Matriculation box type (GetTypeLegend)", GetType())

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

	return (boxType == constants.MatriculationExamBoxType)
}

// GetVersion returns the version string (e.g. "SERVER7108X v69") of the current VM
func GetVersion() string {
	if !boxInstalled {
		return ""
	}

	return vboxmanage.GetVMProperty(boxName, "boxVersion")
}

// getDiskUUID returns the VirtualBox UUID for the image of the current VM
func getDiskUUID() string {
	if !boxInstalled {
		return ""
	}

	return vboxmanage.GetVMInfoByRegexp(boxName, "\"SATA Controller-ImageUUID-0-0\"=\"(.*?)\"")
}

// GetDiskLocation returns the full path of the current VM disk image.
func GetDiskLocation() string {
	if !boxInstalled {
		return ""
	}

	return vboxmanage.GetVMInfoByRegexp(boxName, "\"SATA Controller-0-0\"=\"(.*)\"")
}

// GetLogDir returns the full path of VirtualBox log directory
func GetLogDir() string {
	if !boxInstalled {
		return ""
	}

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
		log.Error("Could not get medium info to calculate its size: %v", err)

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
