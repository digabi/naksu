package main

import (
	"fmt"
	"os"

	"naksu/config"
	"naksu/host"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/xlate"

	flags "github.com/jessevdk/go-flags"

	_ "github.com/andlabs/ui/winmanifest"
)

const thisNaksuVersion = "2.0.8"

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
		{"KTP directory (~/ktp)", mebroutines.GetKtpDirectory()},
		{"MEB share directory (~/ktp-jako)", mebroutines.GetMebshareDirectory()},
		{"VirtualBox hidden directory (~/.VirtualBox)", mebroutines.GetVirtualBoxHiddenDirectory()},
		{"VirtualBox VMs directory (~/VirtualBox VMs)", mebroutines.GetVirtualBoxVMsDirectory()},
	}

	for _, thisDir := range listDirs {
		if mebroutines.ExistsDir(thisDir.dirPath) {
			log.Debug("%s: %s [Directory exists]", thisDir.dirName, thisDir.dirPath)
		} else {
			log.Debug("%s: %s [Directory does not exist]", thisDir.dirName, thisDir.dirPath)
		}
	}
}

func logHardwareDetails() {
	log.Debug("---Hardware data dump (start)\n%s\n---Hardware data dump (end)", host.GetHwLog())
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
		fmt.Printf("Naksu version is %v\n", thisNaksuVersion)
		os.Exit(0)
	})

	handleOptionalArgument("self-update", parser, func(opt *flags.Option) {
		log.Debug("Self-update: %v", opt.Value())
		if opt.Value() == "disabled" {
			config.SetSelfUpdateDisabled(true)
		} else {
			config.SetSelfUpdateDisabled(false)
		}
	})

	log.SetDebug(isDebug)

	// Determine/set path for debug log
	log.SetDebugFilename(log.GetNewDebugFilename())

	log.Action("This is Naksu %s. Hello world!", thisNaksuVersion)

	logDirectoryPaths()

	logHardwareDetails()

	var err = RunUI()

	if err != nil {
		panic(err)
	}

	log.Action("Exiting GUI loop")
}
