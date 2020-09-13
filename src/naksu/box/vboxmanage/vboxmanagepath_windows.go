package vboxmanage

import (
	"os"
	"path/filepath"
)

func getVBoxManagePath() string {
	var path = "VBoxManage"
	if os.Getenv("VBOXMANAGEPATH") != "" {
		path = os.Getenv("VBOXMANAGEPATH")
	} else {
		var pathVirtualbox = os.Getenv("VBOX_MSI_INSTALL_PATH")
		if pathVirtualbox != "" {
			path = filepath.Join(pathVirtualbox, path)
		}
	}

	return path
}
