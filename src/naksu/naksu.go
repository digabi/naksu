package main

import (
	"fmt"
	"naksu/boxversion"
	"naksu/config"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/xlate"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/kardianos/osext"
)

const version = "1.8.1"

var isDebug bool

// Options contains command line options
type Options struct {
	IsDebug    bool   `short:"D" long:"debug" description:"Turn debugging on" optional:"true"`
	Version    bool   `short:"v" long:"version" description:"Print naksu version" optional:"true"`
	SelfUpdate string `long:"self-update" choice:"enabled" choice:"disabled" description:"Control self-update behaviour. Naksu will always warn if your version is out-of-date. This flag will store the setting to ini-file." optional:"true"`
}

var options Options

func handleOptionalArgument(longName string, parser *flags.Parser, function func(option *flags.Option)) {
	opt := parser.FindOptionByLongName(longName)
	if opt != nil && opt.IsSet() {
		function(opt)
	}
}

func logDirectoryPaths() {
	listDirs := []struct {
		dirName string
		dirPath string
	}{
		{"Home directory (~)", mebroutines.GetHomeDirectory()},
		{"Vagrant directory (~/ktp)", mebroutines.GetVagrantDirectory()},
		{"MEB share directory (~/ktp-jako)", mebroutines.GetMebshareDirectory()},
		{"Vagrant internal settings directory (~/vagrant.d)", mebroutines.GetVagrantdDirectory()},
		{"VirtualBox hidden directory (~/.VirtualBox)", mebroutines.GetVirtualBoxHiddenDirectory()},
		{"VirtualBox VMs directory (~/VirtualBox VMs)", mebroutines.GetVirtualBoxVMsDirectory()},
	}

	for _, thisDir := range listDirs {
		if mebroutines.ExistsDir(thisDir.dirPath) {
			log.LogDebug(fmt.Sprintf("%s: %s [Directory exists]", thisDir.dirName, thisDir.dirPath))
		} else {
			log.LogDebug(fmt.Sprintf("%s: %s [Directory does not exist]", thisDir.dirName, thisDir.dirPath))
		}
	}
}

func main() {
	// Load configuration if it exists
	config.Load()

	// Set default UI language
	xlate.SetLanguage(config.GetLanguage())

	var parser = flags.NewParser(&options, flags.Default)
	_, parseErr := parser.Parse()

	if flags.WroteHelp(parseErr) {
		os.Exit(0)
	} else if parseErr != nil {
		panic(parseErr)
	}

	handleOptionalArgument("debug", parser, func(opt *flags.Option) {
		isDebug = true
	})

	handleOptionalArgument("version", parser, func(opt *flags.Option) {
		fmt.Printf("Naksu version is %v\n", version)
		os.Exit(0)
	})

	handleOptionalArgument("self-update", parser, func(opt *flags.Option) {
		log.LogDebug(fmt.Sprintf("Self-update: %v", opt.Value()))
		if opt.Value() == "disabled" {
			config.SetSelfUpdateDisabled(true)
		} else {
			config.SetSelfUpdateDisabled(false)
		}
	})

	RunSelfUpdate()

	log.SetDebug(isDebug)

	// Determine/set path for debug log
	log.SetDebugFilename(log.GetNewDebugFilename())

	log.LogDebug(fmt.Sprintf("This is Naksu %s. Hello world!", version))

	logDirectoryPaths()

	log.LogDebug(fmt.Sprintf("Currently installed box: %s", boxversion.GetVagrantFileVersion("")))

	// Check whether we have a terminal (restart with x-terminal-emulator, if missing)
	if !mebroutines.ExistsStdin() {
		pathToMe, err := osext.Executable()
		if err != nil {
			log.LogDebug("Failed to get executable path")
		}
		commandArgs := []string{"x-terminal-emulator", "-e", pathToMe}

		log.LogDebug(fmt.Sprintf("No stdin, restarting with terminal: %s", strings.Join(commandArgs, " ")))
		_, err = mebroutines.RunAndGetOutput(commandArgs, true)
		if err != nil {
			log.LogDebug(fmt.Sprintf("Failed to restart with terminal: %d", err))
		}
		log.LogDebug(fmt.Sprintf("No stdin, returned from %s", strings.Join(commandArgs, " ")))

		// Normal termination
		os.Exit(0)
	}

	var err = RunUI()

	if err != nil {
		panic(err)
	}

	log.LogDebug("Exiting GUI loop")
}
