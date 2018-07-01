package main

// required by selfupdate (needs context)
// +build go1.7

import (
  "os"
  "fmt"
  "flag"

  "github.com/blang/semver"
  "github.com/rhysd/go-github-selfupdate/selfupdate"
  "github.com/andlabs/ui"

  "mebroutines"
  "mebroutines/install"
  "mebroutines/start"
)

const version = "1.0.0"

var is_debug bool

func doSelfUpdate() bool {
  v := semver.MustParse(version)

  if (mebroutines.Is_debug()) {
    selfupdate.EnableLog()
  }

  latest, err := selfupdate.UpdateSelf(v, "digabi/naksu")
  if err != nil {
    mebroutines.Message_warning(fmt.Sprintf("Binary update failed: %s", err))
    return true
  }
  if latest.Version.Equals(v) {
    // latest version is the same as current version. It means current binary is up to date.
    mebroutines.Message_debug(fmt.Sprintf("Current binary is the latest version: %s", version))
    return false
  } else {
    mebroutines.Message_debug(fmt.Sprintf("Successfully updated to version: %s", latest.Version))
    return true
    //log.Println("Release note:\n", latest.ReleaseNotes)
  }
}

func main() {
  // Process command line parameters
  flag.BoolVar(&is_debug, "debug", false, "Turn debugging on")
  flag.Parse()

  mebroutines.Set_debug(is_debug)

  // UI (main menu)
  if doSelfUpdate() {
    mebroutines.Message_warning("naksu has been automatically updated. Please restart naksu.")
    os.Exit(0)
  }

  err := ui.Main(func () {
    button_start_server := ui.NewButton("Start Stickless Exam Server")
    button_get_server := ui.NewButton("Install new or update existing Stickless Exam Server")
    button_exit := ui.NewButton("Exit")

    box := ui.NewVerticalBox()
    box.Append(button_start_server, false)
    box.Append(button_get_server, false)
    box.Append(button_exit, false)

    window := ui.NewWindow(fmt.Sprintf("naksu %s", version), 1, 1, false)

    mebroutines.Set_main_window(window)

    window.SetMargined(true)
		window.SetChild(box)

    button_start_server.OnClicked(func(*ui.Button) {
      window.Hide()
      ui.QueueMain(func () {
        start.Do_start_server()
        os.Exit(0)
      })
    })

    button_get_server.OnClicked(func(*ui.Button) {
      window.Hide()
      ui.QueueMain(func () {
        install.Do_get_server()
        os.Exit(0)
      })
    })

    button_exit.OnClicked(func(*ui.Button) {
      mebroutines.Message_debug("Exiting by user request")
      os.Exit(0)
    })

    window.Show()

    // Make sure we have vagrant
  	if (! mebroutines.If_found_vagrant()) {
  		mebroutines.Message_error("Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?")
  	}

  	// Make sure we have VBoxManage
  	if (! mebroutines.If_found_vboxmanage()) {
  		mebroutines.Message_error("Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?")
  	}
  })

  if err != nil {
    panic(err)
  }
}
