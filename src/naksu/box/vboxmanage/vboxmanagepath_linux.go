package vboxmanage

import (
	"os"
)

func getVBoxManagePath() string {
	var path = "VBoxManage"
	if os.Getenv("VBOXMANAGEPATH") != "" {
		path = os.Getenv("VBOXMANAGEPATH")
	}

	return path
}
