// required by selfupdate (needs context)
// +build go1.7

package main

import (
	"fmt"
	"naksu/mebroutines"
	"naksu/xlate"
	"naksu/mebroutines/install"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// RunSelfUpdate executes self-update
func RunSelfUpdate() {
	// Run auto-update
	if doSelfUpdate() {
		mebroutines.ShowWarningMessage("naksu has been automatically updated. Please restart naksu.")
		os.Exit(0)
	}
}

func doSelfUpdate() bool {
	v := semver.MustParse(version)

	if mebroutines.IsDebug() {
		selfupdate.EnableLog()
	}

	// Test network connection here with a timeout
	if !install.TestHTTPGet(URLTest, URLTestTimeout) {
		mebroutines.ShowWarningMessage(xlate.Get("Naksu could not check for updates as there is no network connection."))
		return false
	}

	latest, err := selfupdate.UpdateSelf(v, "digabi/naksu")
	if err != nil {
		mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Naksu update failed. Maybe you don't have network connection?\n\nError: %s"), err))
		return false
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		mebroutines.LogDebug(fmt.Sprintf("Current binary is the latest version: %s", version))
		return false
	}
	mebroutines.LogDebug(fmt.Sprintf("Successfully updated to version: %s", latest.Version))
	return true
	//log.Println("Release note:\n", latest.ReleaseNotes)
}
