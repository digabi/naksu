package config

import (
	"fmt"
	"path/filepath"
	"strconv"

	"naksu/constants"
	"naksu/log"

	"github.com/go-ini/ini"
	"github.com/mitchellh/go-homedir"
)

var cfg *ini.File

func getIniFilePath() string {
	homeDir, errHome := homedir.Dir()
	if errHome != nil {
		panic("Could not get home directory")
	}
	return filepath.Join(homeDir, "naksu.ini")
}

func setIfMissing(section string, key string, defaultValue string) {
	if !cfg.Section(section).HasKey(key) {
		cfg.Section(section).Key(key).SetValue(defaultValue)
	}
}

type defaultValue struct {
	section, key, value string
}

var defaults = []defaultValue{
	defaultValue{"common", "iniVersion", strconv.FormatInt(1, 10)},
	defaultValue{"common", "language", constants.AvailableLangs[0].ConfigValue},
	defaultValue{"selfupdate", "disabled", strconv.FormatBool(false)},
	defaultValue{"environment", "nic", constants.AvailableNics[0].ConfigValue},
	defaultValue{"environment", "extnic", ""},
}

func fillDefaults() {
	for _, defaultValue := range defaults {
		setIfMissing(defaultValue.section, defaultValue.key, defaultValue.value)
	}
}

func getDefault(section string, key string) string {
	for _, defaultValue := range defaults {
		if defaultValue.section == section && defaultValue.key == key {
			return defaultValue.value
		}
	}
	panic(fmt.Sprintf("Default for %v / %v is not defined!", section, key))
}

func getIniKey(section string, key string) *ini.Key {
	return cfg.Section(section).Key(key)
}

func getBoolean(section string, key string) bool {
	value, err := getIniKey(section, key).Bool()
	if err != nil {
		log.Debug(fmt.Sprintf("Parsing key %s / %s as bool failed", section, key))
		defaultValue := getDefault(section, key)
		value, err = strconv.ParseBool(defaultValue)
		if err != nil {
			panic(fmt.Sprintf("Default boolean parsing for %v / %v (%v) failed to parse to boolean!", section, key, defaultValue))
		}
		setValue(section, key, defaultValue)
	}
	return value
}

func getString(section string, key string) string {
	return getIniKey(section, key).String()
}

func setValue(section string, key string, value string) {
	log.Debug(fmt.Sprintf("Setting new configuration: section %s, key: %s, value: %s", section, key, value))
	cfg.Section(section).Key(key).SetValue(value)
	save()
}

// Load or initialize configuration to empty object
func Load() {
	naksuIniPath := getIniFilePath()

	var err error
	cfg, err = ini.Load(naksuIniPath)
	if err != nil {
		log.Debug(fmt.Sprintf("%s not found, setting up empty config with defaults", naksuIniPath))
		cfg = ini.Empty()
	}
	fillDefaults()
	save()
}

func validateStringChoice(section string, key string, choices []constants.AvailableSelection) string {
	value := getString(section, key)

	id := constants.GetAvailableSelectionID(value, choices, -1)

	if id >= 0 {
		return value
	}
	defaultValue := getDefault(section, key)
	log.Debug(fmt.Sprintf("Correcting malformed ini-key %v / %v to default value %v", section, key, defaultValue))
	setValue(section, key, defaultValue)
	return defaultValue
}

// Save configuration to disk
func save() {
	naksuIniPath := getIniFilePath()

	err := cfg.SaveTo(naksuIniPath)
	if err != nil {
		log.Debug(fmt.Sprintf("%s save failed: %v", naksuIniPath, err))
	}
}

// GetLanguage returns user language preference. defaults to fi
func GetLanguage() string {
	return validateStringChoice("common", "language", constants.AvailableLangs)
}

// SetLanguage stores user language preference
func SetLanguage(language string) {
	if constants.GetAvailableSelectionID(language, constants.AvailableLangs, -1) < 0 {
		setValue("common", "language", getDefault("common", "language"))
	} else {
		setValue("common", "language", language)
	}
}

// IsSelfUpdateDisabled returns true, if self-update functionality should be disabled
func IsSelfUpdateDisabled() bool {
	return getBoolean("selfupdate", "disabled")
}

// SetSelfUpdateDisabled sets the state of self-update functionality
func SetSelfUpdateDisabled(isSelfUpdateDisabled bool) {
	setValue("selfupdate", "disabled", strconv.FormatBool(isSelfUpdateDisabled))
}

// GetNic returns vagrant NIC value. Defaults to "virtio"
func GetNic() string {
	return validateStringChoice("environment", "nic", constants.AvailableNics)
}

// SetNic sets the state of vagrant NIC value
func SetNic(nic string) {
	if constants.GetAvailableSelectionID(nic, constants.AvailableNics, -1) < 0 {
		setValue("environment", "nic", getDefault("environment", "nic"))
	} else {
		setValue("environment", "nic", nic)
	}
}

// GetExtNic returns vagrant EXTNIC value. If empty no external network connection
// should be passed to Vagrant via EXTNIC environment variable.
func GetExtNic() string {
	// Since there are no pre-set selection of variables we dont use validateStringChoice() here
	return getString("environment", "extnic")
}

// SetExtNic sets the state of vagrant EXTNIC value
func SetExtNic(nic string) {
	setValue("environment", "extnic", nic)
}
