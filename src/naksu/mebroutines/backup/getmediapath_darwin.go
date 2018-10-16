package backup

import (
	"encoding/json"
	"fmt"
	"naksu/mebroutines"
	"naksu/xlate"
	"os"
)

// GetBackupMedia returns map of backup medias
func GetBackupMedia() map[string]string {
	media := getBackupMediaLinux()

	// Add some entries from environment variables
	if os.Getenv("HOME") != "" {
		media[os.Getenv("HOME")] = xlate.Get("Home directory")

		// Try ~/Desktop
		desktopPath := os.Getenv("HOME") + string(os.PathSeparator) + "Desktop"
		if mebroutines.ExistsDir(desktopPath) {
			media[desktopPath] = xlate.Get("Desktop")
		}

		// Try ~/desktop
		desktopPath = os.Getenv("HOME") + string(os.PathSeparator) + "desktop"
		if mebroutines.ExistsDir(desktopPath) {
			media[desktopPath] = xlate.Get("Desktop")
		}
	}
	if os.TempDir() != "" {
		media[os.TempDir()] = xlate.Get("Temporary files")
	}

	return media
}

func getBackupMediaLinux() map[string]string {
	var media = map[string]string{}

	runParams := []string{"lsblk", "-J", "-o", "NAME,FSTYPE,MOUNTPOINT,VENDOR,MODEL,HOTPLUG"}

	lsblkJSON, lsblkErr := mebroutines.RunAndGetOutput(runParams)

	mebroutines.LogDebug("lsblk says:")
	mebroutines.LogDebug(lsblkJSON)

	if lsblkErr != nil {
		mebroutines.LogDebug("Failed to run lsblk")
		// Return empty set of media
		return media
	}

	var jsonData map[string]interface{}

	jsonErr := json.Unmarshal([]byte(lsblkJSON), &jsonData)
	if jsonErr != nil {
		mebroutines.LogDebug("Unable on decode lsblk response:")
		mebroutines.LogDebug(fmt.Sprintf("%s", jsonErr))
		// Return empty set of media
		return media
	}

	blockdevices := jsonData["blockdevices"].([]interface{})

	media = getRemovableDisks(blockdevices)

	return media
}

func getRemovableDisks(blockdevices []interface{}) map[string]string {
	var media = map[string]string{}
	//media_n := 0

	if blockdevices == nil {
		return media
	}

	for blockdeviceIndex := range blockdevices {
		//fmt.Println(blockdevices[blockdevice_n])
		blockdevice := blockdevices[blockdeviceIndex].(map[string]interface{})
		if deviceFieldString(blockdevice["hotplug"]) == "1" && blockdevice["children"] != nil {
			children := blockdevice["children"].([]interface{})

			for childIndex := range children {
				child := children[childIndex].(map[string]interface{})

				mountpoint := deviceFieldString(child["mountpoint"])
				if mountpoint != "" {
					media[mountpoint] = fmt.Sprintf("%s, %s", deviceFieldString(blockdevice["vendor"]), deviceFieldString(blockdevice["model"]))
				}
			}
		}
	}

	return media
}

func deviceFieldString(field interface{}) string {
	if field == nil {
		return ""
	}

	return field.(string)
}
