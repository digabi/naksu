package vboxmanage

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func getVirtualBoxConfigPath() string {
	homeDir, errHome := homedir.Dir()
	if errHome != nil {
		panic("Could not get home directory")
	}

	return filepath.Join(homeDir, ".config", "VirtualBox", "VirtualBox.xml")
}
