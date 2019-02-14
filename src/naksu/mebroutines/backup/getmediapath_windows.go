package backup

import (
	"fmt"
	"naksu/mebroutines"
	"naksu/xlate"
	"os"

	"github.com/StackExchange/wmi"
)

// GetBackupMedia returns the backup media path
func GetBackupMedia() map[string]string {
	media := getBackupMediaWindows()

	if os.Getenv("TEMP") != "" {
		media[os.Getenv("TEMP")] = xlate.Get("Temporary files")
	}
	if os.Getenv("USERPROFILE") != "" {
		media[os.Getenv("USERPROFILE")] = xlate.Get("Profile directory")

		desktopPath := os.Getenv("USERPROFILE") + string(os.PathSeparator) + "Desktop"
		if mebroutines.ExistsDir(desktopPath) {
			media[desktopPath] = xlate.Get("Desktop")
		}
	}

	return media
}

func getBackupMediaWindows() map[string]string {
	// This struct must be named with an underscore, otherwise it is not recognised
	// and results "Invalid class" exception.
	type Win32_LogicalDisk struct {
		DeviceID    string
		DriveType   int
		Description string
		VolumeName  string
	}

	var media = map[string]string{}

	var dst []Win32_LogicalDisk
	query := wmi.CreateQuery(&dst, "WHERE DriveType=2 OR DriveType=3")
	err := wmi.Query(query, &dst)
	if err != nil {
		mebroutines.LogDebug("getBackupMediaWindows() could not detect removable/hard drives as it could not query WMI")
		mebroutines.LogDebug(fmt.Sprint(err))
		return media
	}

	for thisDrive := range dst {
		// We have either hard or removable drive
		thisPath := fmt.Sprintf("%s%s", dst[thisDrive].DeviceID, string(os.PathSeparator))
		media[thisPath] = fmt.Sprintf("%s, %s", dst[thisDrive].VolumeName, dst[thisDrive].Description)
	}

	return media
}
