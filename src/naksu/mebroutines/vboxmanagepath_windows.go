package mebroutines

import (
	"os"
)

func getVBoxManagePath() string {
	var path = "VBoxManage"
	if os.Getenv("VBOXMANAGEPATH") != "" {
		path = os.Getenv("VBOXMANAGEPATH")
	} else {
		var pathVirtualbox = os.Getenv("VBOX_MSI_INSTALL_PATH")
		if pathVirtualbox != "" {
			path = pathVirtualbox + string(os.PathSeparator) + path
		}
	}

	return path
}
