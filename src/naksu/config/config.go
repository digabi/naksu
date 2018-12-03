package config

import (
	"fmt"
	"naksu/mebroutines"
	"strconv"

	"github.com/go-ini/ini"
)

var cfg *ini.File

func setIfMissing(section string, key string, defaultValue string) {
	if !cfg.Section(section).HasKey(key) {
		cfg.Section(section).Key(key).SetValue(defaultValue)
	}
}

type defaultValue struct {
	section, key, value string
}

var defaults []defaultValue

func initDefaults() {
	defaults = []defaultValue{
		defaultValue{"common", "iniVersion", strconv.FormatInt(1, 10)},
		defaultValue{"common", "language", "fi"},
		defaultValue{"selfupdate", "disabled", strconv.FormatBool(false)},
		defaultValue{"selfupdate", "useBeta", strconv.FormatBool(false)},
	}
}

func fillDefaults() {
	for _, defaultValue := range defaults {
		setIfMissing(defaultValue.section, defaultValue.key, defaultValue.value)
	}
}

func getBoolean(section string, key string) bool {
	value, err := getValue("selfupdate", "disabled").Bool()
	if err != nil {
		mebroutines.LogDebug(fmt.Sprintf("Parsing key %s / %s as bool failed", section, key))
		panic(fmt.Sprintf("Invalid boolean configuration flag! section: %s, key: %s", section, key))
	}
	return value
}

func getValue(section string, key string) *ini.Key {
	return cfg.Section(section).Key(key)
}

func setValue(section string, key string, value string) {
	cfg.Section(section).Key(key).SetValue(value)
}

// Load or initialize configuration to empty object
func Load() {
	initDefaults()
	var err error
	cfg, err = ini.Load("naksu.ini")
	if err != nil {
		mebroutines.LogDebug("naksu.ini not found, setting up empty config with defaults")
		cfg = ini.Empty()
	}
	fillDefaults()
}

// Save configuration to disk
func Save() {
	cfg.SaveTo("naksu.ini")
}

// GetLanguage returns user language preference. defaults to fi
func GetLanguage() string {
	return getValue("common", "language").String()
}

// SetLanguage stores user language preference
func SetLanguage(language string) {
	setValue("common", "language", language)
}

// IsSelfUpdateDisabled returns true, if self-update functionality should be disabled
func IsSelfUpdateDisabled() bool {
	return getBoolean("selfupdate", "disabled")
}

// SetReleaseChannel sets the release channel to selected value. Default is "release"
func SetReleaseChannel(channel string) {
	setValue("selfupdate", "channel", channel)
}

// GetReleaseChannel returns selected release channel. Default is "release"
func GetReleaseChannel() string {
	return getValue("selfupdate", "channel").String()
}
