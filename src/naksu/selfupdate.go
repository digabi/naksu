// required by selfupdate (needs context)
//go:build go1.7
// +build go1.7

package main

import (
	"naksu/config"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/network"
	"naksu/ui/progress"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

var isOutOfDate bool

// RunSelfUpdate executes self-update
func RunSelfUpdate(currentVersionStr string) {
	// Run auto-update
	if doReleaseSelfUpdate(currentVersionStr) {
		mebroutines.ShowTranslatedInfoMessage("Naksu has been automatically updated. Please restart Naksu.")
	}
	if WarnUserAboutStaleVersionIfUpdateDisabled() {
		mebroutines.ShowTranslatedInfoMessage("Naksu has update available, but your version of Naksu has updates disabled. Please update or ask your administrator to update Naksu.")
	}
}

func doReleaseSelfUpdate(currentVersionStr string) bool {
	progress.TranslateAndSetMessage("Checking for new versions of Naksu...")

	currentVersion := semver.MustParse(currentVersionStr)

	if log.IsDebug() {
		selfupdate.EnableLog()
	}

	// Test network connection here with a timeout
	if !network.CheckIfNetworkAvailable() {
		progress.TranslateAndSetMessage("Naksu self-update needs network connection")

		return false
	}

	// If self-update is disabled, do a version check nevertheless and store information for user warning
	if config.IsSelfUpdateDisabled() {
		latest, found, err := selfupdate.DetectLatest("digabi/naksu")
		progress.SetMessage("")
		if err != nil {
			log.Error("Version check failed: %s", err)

			return false
		}
		if found && latest.Version.GT(currentVersion) {
			isOutOfDate = true
		}

		return false
	}

	latest, err := selfupdate.UpdateSelf(currentVersion, "digabi/naksu")
	progress.SetMessage("")
	if err != nil {
		mebroutines.ShowTranslatedWarningMessage("Naksu update failed. Maybe you don't have network connection?\n\nError: %s", err)

		return false
	}
	if latest.Version.Equals(currentVersion) {
		// latest version is the same as current version. It means current binary is up to date.
		log.Debug("Current binary is the latest version: %s", currentVersionStr)

		return false
	}
	log.Debug("Successfully updated to version: %s", latest.Version)

	return true
}

// WarnUserAboutStaleVersionIfUpdateDisabled tells us if we should warn user that they are running old version if self-update is disabled. This is very corner-case check
// for environments that prefer distributing naksu via AD or other centralized management environment
func WarnUserAboutStaleVersionIfUpdateDisabled() bool {
	return config.IsSelfUpdateDisabled() && isOutOfDate
}
