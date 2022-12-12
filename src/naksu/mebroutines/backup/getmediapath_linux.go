package backup

import (
	"os"
	"path/filepath"

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

// isFAT32 returns true if the filesystem of the drive
// pointed to by backupPath is FAT32.
func isFAT32(backupPath string) (bool, error) {
	lsblk, err := ListBlockDevices()
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
			correctDevice := device

			return &correctDevice
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
	lsblk, err := ListBlockDevices()
	if err != nil {
		// Return empty set of media
		return map[string]string{}
	}

	return lsblk.GetRemovableDisks()
}
