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

const version = "1.6.1"
const lowDiskLimit = 50 * 1024 * 1024 // 50 Gb

var isDebug bool

func main() {
	// Load configuration if it exists
	config.Load()
	// Set default UI language
	xlate.SetLanguage(config.GetLanguage())


	// Process command line parameters
	flag.BoolVar(&isDebug, "debug", false, "Turn debugging on")
	flag.Parse()

	RunSelfUpdate()

	mebroutines.SetDebug(isDebug)

	// Determine/set path for debug log
	mebroutines.SetDebugFilename(mebroutines.GetNewDebugFilename())

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

	config.Save()

	mebroutines.LogDebug("Exiting GUI loop")
}
