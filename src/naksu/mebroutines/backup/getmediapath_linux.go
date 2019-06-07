package backup

import (
	"encoding/json"
	"fmt"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/xlate"
	"os"
	"path/filepath"
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

func getBackupMediaLinux() map[string]string {
	var media = map[string]string{}

	runParams := []string{"lsblk", "-J", "-o", "NAME,FSTYPE,MOUNTPOINT,VENDOR,MODEL,HOTPLUG"}

	lsblkJSON, lsblkErr := mebroutines.RunAndGetOutput(runParams, true)

	log.Debug("lsblk says:")
	log.Debug(lsblkJSON)

	if lsblkErr != nil {
		log.Debug("Failed to run lsblk")
		// Return empty set of media
		return media
	}

	var jsonData map[string]interface{}

	jsonErr := json.Unmarshal([]byte(lsblkJSON), &jsonData)
	if jsonErr != nil {
		log.Debug("Unable on decode lsblk response:")
		log.Debug(fmt.Sprintf("%s", jsonErr))
		// Return empty set of media
		return media
	}

	blockdevices := jsonData["blockdevices"].([]interface{})

	media = getRemovableDisks(blockdevices)

	return media
}

func getRemovableDisks(blockdevices []interface{}) map[string]string {
	media := map[string]string{}

	if blockdevices == nil {
		return media
	}

	for blockdeviceIndex := range blockdevices {
		thisBlockdevice := blockdevices[blockdeviceIndex].(map[string]interface{})
		if deviceFieldString(thisBlockdevice["hotplug"]) == "1" && thisBlockdevice["children"] != nil {
			thisChildren := thisBlockdevice["children"].([]interface{})

			for thisChildIndex := range thisChildren {
				thisChild := thisChildren[thisChildIndex].(map[string]interface{})

				thisMountpoint := deviceFieldString(thisChild["mountpoint"])
				if thisMountpoint != "" {
					media[thisMountpoint] = fmt.Sprintf("%s, %s", deviceFieldString(thisBlockdevice["vendor"]), deviceFieldString(thisBlockdevice["model"]))
				}
			}
		}
	}

	return media
}

func deviceFieldString(thisField interface{}) string {
	if thisField == nil {
		return ""
	}
	switch v := thisField.(type) {
	case bool:
		fieldBool := thisField.(bool)
		if fieldBool {
			return "1"
		}
		return "0"
	case string:
		return thisField.(string)
	default:
		log.Debug("Fail on getmediapath.deviceFieldString")
		log.Debug(fmt.Sprintf("unexpected type %T", v))
	}
	return thisField.(string)
}
