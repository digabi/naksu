package vboxmanage

import (
	"github.com/mitchellh/go-homedir"
	"path/filepath"
)

func getVirtualBoxConfigPath() string {
	homeDir, errHome := homedir.Dir()
	if errHome != nil {
		panic("Could not get home directory")
	}
	return filepath.Join(homeDir, "Library", "VirtualBox", "VirtualBox.xml")
}
