package main

import (
	"flag"
	"fmt"
	"naksu/mebroutines"
	"naksu/xlate"
	"os"
	"strings"

	"github.com/kardianos/osext"
)

const version = "1.5.0"
const lowDiskLimit = 50000000 // 50 Gb

// URLTest is testing URL for checking network connection
const URLTest = "http://static.abitti.fi/usbimg/qa/latest.txt"

var isDebug bool

func main() {
	// Set default UI language
	xlate.SetLanguage("fi")

	// Process command line parameters
	flag.BoolVar(&isDebug, "debug", false, "Turn debugging on")
	flag.Parse()

	RunSelfUpdate()

	mebroutines.SetDebug(isDebug)

	// Determine/set path for debug log
	logPath := mebroutines.GetVagrantDirectory()
	if mebroutines.ExistsDir(logPath) {
		mebroutines.SetDebugFilename(logPath + string(os.PathSeparator) + "naksu_lastlog.txt")
	} else {
		mebroutines.SetDebugFilename(os.TempDir() + string(os.PathSeparator) + "naksu_lastlog.txt")
	}

	mebroutines.LogDebug(fmt.Sprintf("This is Naksu %s. Hello world!", version))

	// Check whether we have a terminal (restart with x-terminal-emulator, if missing)
	if !mebroutines.ExistsStdin() {
		pathToMe, err := osext.Executable()
		if err != nil {
			mebroutines.LogDebug("Failed to get executable path")
		}
		commandArgs := []string{"x-terminal-emulator", "-e", pathToMe}

		mebroutines.LogDebug(fmt.Sprintf("No stdin, restarting with terminal: %s", strings.Join(commandArgs, " ")))
		_, err = mebroutines.RunAndGetOutput(commandArgs)
		if err != nil {
			mebroutines.LogDebug(fmt.Sprintf("Failed to restart with terminal: %d", err))
		}
		mebroutines.LogDebug(fmt.Sprintf("No stdin, returned from %s", strings.Join(commandArgs, " ")))

		// Normal termination
		os.Exit(0)
	}

	var err = RunUI()

	if err != nil {
		panic(err)
	}

	mebroutines.LogDebug("Exiting GUI loop")
}
