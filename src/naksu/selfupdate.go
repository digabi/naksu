// required by selfupdate (needs context)
// +build go1.7

package main

import (
	"fmt"

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
func RunSelfUpdate() {
	// Run auto-update
	if doReleaseSelfUpdate() {
		mebroutines.ShowTranslatedInfoMessage("Naksu has been automatically updated. Please restart Naksu.")
	}
	if WarnUserAboutStaleVersionIfUpdateDisabled() {
		mebroutines.ShowTranslatedInfoMessage("Naksu has update available, but your version of Naksu has updates disabled. Please update or ask your administrator to update Naksu.")
	}
}

func doReleaseSelfUpdate() bool {
	progress.TranslateAndSetMessage("Checking for new versions of Naksu...")

	v := semver.MustParse(version)

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
		if err != nil {
			log.Debug(fmt.Sprintf("Version check failed: %s", err))
			return false
		}
		if found && latest.Version.GT(v) {
			isOutOfDate = true
		}
		return false
	}

	latest, err := selfupdate.UpdateSelf(v, "digabi/naksu")
	progress.SetMessage("")
	if err != nil {
		mebroutines.ShowTranslatedWarningMessage("Naksu update failed. Maybe you don't have network connection?\n\nError: %s", err)
		return false
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		log.Debug(fmt.Sprintf("Current binary is the latest version: %s", version))
		return false
	}
	log.Debug(fmt.Sprintf("Successfully updated to version: %s", latest.Version))
	return true
	//log.Println("Release note:\n", latest.ReleaseNotes)
}

// WarnUserAboutStaleVersionIfUpdateDisabled tells us if we should warn user that they are running old version if self-update is disabled. This is very corner-case check
// for environments that prefer distributing naksu via AD or other centralized management environment
func WarnUserAboutStaleVersionIfUpdateDisabled() bool {
	return config.IsSelfUpdateDisabled() && isOutOfDate
}
