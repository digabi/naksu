package backup

import (
	"fmt"
	"os"
	"path/filepath"

	"naksu/log"
	"naksu/mebroutines"
	"naksu/xlate"
)

// GetBackupMedia returns backup media path
func GetBackupMedia() map[string]string {
	media := getBackupMediaLinux()

	// Add some entries from environment variables
	if os.Getenv("HOME") != "" {
		media[os.Getenv("HOME")] = xlate.Get("Home directory")

		// Try ~/Desktop
		desktopPath := filepath.Join(os.Getenv("HOME"), "Desktop")
		if mebroutines.ExistsDir(desktopPath) {
			media[desktopPath] = xlate.Get("Desktop")
		}

		// Try ~/desktop
		desktopPath = filepath.Join(os.Getenv("HOME"), "desktop")
		if mebroutines.ExistsDir(desktopPath) {
			media[desktopPath] = xlate.Get("Desktop")
		}
	}
	if os.TempDir() != "" {
		media[os.TempDir()] = xlate.Get("Temporary files")
	}

	return media
}

func listBlockDevices() (*LsblkOutput, error) {
	runParams := []string{"lsblk", "-J", "-o", "NAME,FSTYPE,MOUNTPOINT,VENDOR,MODEL,HOTPLUG"}

	lsblkJSON, lsblkErr := mebroutines.RunAndGetOutput(runParams, true)

	log.Debug("lsblk says:")
	log.Debug(lsblkJSON)

	if lsblkErr != nil {
		log.Debug("Failed to run lsblk")
		return &LsblkOutput{}, lsblkErr
	}

	output, jsonErr := ParseLsblkJSON(lsblkJSON)
	if jsonErr != nil {
		log.Debug("Unable to unmarshal lsblk response:")
		log.Debug(fmt.Sprintf("%s", jsonErr))
		return &LsblkOutput{}, lsblkErr
	}

	return output, nil
}

// isFAT32 returns true if the filesystem of the drive
// pointed to by backupPath is FAT32.
func isFAT32(backupPath string) (bool, error) {
	lsblk, err := listBlockDevices()
	if err != nil {
		return false, err
	}

	mountPoint := filepath.Dir(backupPath)
	device := findBlockDevice(lsblk.BlockDevices, mountPoint)
	if device != nil {
		return device.FileSystem == "vfat", nil
	}

	return false, nil
}

func findBlockDevice(blockDevices []BlockDevice, mountPoint string) *BlockDevice {
	for _, device := range blockDevices {
		if mountPoint == device.MountPoint {
			return &device
		}

		if len(device.Children) > 0 {
			childDevice := findBlockDevice(device.Children, mountPoint)
			if childDevice != nil {
				return childDevice
			}
		}
	}

	return nil
}

func getBackupMediaLinux() map[string]string {
	lsblk, err := listBlockDevices()
	if err != nil {
		// Return empty set of media
		return map[string]string{}
	}

	return lsblk.GetRemovableDisks()
}
