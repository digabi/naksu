package backup

import (
	"naksu/mebroutines"
	"naksu/xlate"
	"os"
)

// GetBackupMedia returns map of backup medias
func GetBackupMedia() map[string]string {
	media := make(map[string]string)

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
