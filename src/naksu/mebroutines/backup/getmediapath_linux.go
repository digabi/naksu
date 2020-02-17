package backup

import (
	"encoding/json"
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

type lsblkOutput struct {
	BlockDevices []blockDevice `json:"blockdevices"`
}

type blockDevice struct {
	Name       string        `json:"name"`
	FileSystem string        `json:"fstype"`
	MountPoint string        `json:"mountpoint"`
	Vendor     string        `json:"vendor"`
	Model      string        `json:"model"`
	HotPlug    bool          `json:"hotplug"`
	Children   []blockDevice `json:"children"`
}

func listBlockDevices() ([]blockDevice, error) {
	runParams := []string{"lsblk", "-J", "-o", "NAME,FSTYPE,MOUNTPOINT,VENDOR,MODEL,HOTPLUG"}

	lsblkJSON, lsblkErr := mebroutines.RunAndGetOutput(runParams, true)

	log.Debug("lsblk says:")
	log.Debug(lsblkJSON)

	if lsblkErr != nil {
		log.Debug("Failed to run lsblk")
		return []blockDevice{}, lsblkErr
	}

	var jsonData lsblkOutput

	jsonErr := json.Unmarshal([]byte(lsblkJSON), &jsonData)
	if jsonErr != nil {
		log.Debug("Unable to unmarshal lsblk response:")
		log.Debug(fmt.Sprintf("%s", jsonErr))
		return []blockDevice{}, lsblkErr
	}

	return jsonData.BlockDevices, nil
}

func findBlockDevice(blockDevices []blockDevice, mountPoint string) *blockDevice {
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
	blockDevices, err := listBlockDevices()
	if err != nil {
		// Return empty set of media
		return map[string]string{}
	}

	return getRemovableDisks(blockDevices)
}

func getRemovableDisks(blockdevices []blockDevice) map[string]string {
	media := map[string]string{}

	if blockdevices == nil {
		return media
	}

	for blockdeviceIndex := range blockdevices {
		thisBlockdevice := blockdevices[blockdeviceIndex]
		if thisBlockdevice.HotPlug && thisBlockdevice.Children != nil {
			thisChildren := thisBlockdevice.Children

			for thisChildIndex := range thisChildren {
				thisChild := thisChildren[thisChildIndex]

				thisMountpoint := thisChild.MountPoint
				if thisMountpoint != "" {
					media[thisMountpoint] = fmt.Sprintf("%s, %s", thisBlockdevice.Vendor, thisBlockdevice.Model)
				}
			}
		}
	}

	return media
}
