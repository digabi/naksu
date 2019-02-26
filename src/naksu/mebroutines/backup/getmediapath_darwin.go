package backup

import (
	"path/filepath"
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
