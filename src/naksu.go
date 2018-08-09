package main

// required by selfupdate (needs context)
// +build go1.7

import (
  "os"
  "fmt"
  "flag"
  "time"
  "strings"

  "github.com/blang/semver"
  "github.com/rhysd/go-github-selfupdate/selfupdate"
  "github.com/andlabs/ui"
  "github.com/kardianos/osext"

  "mebroutines"
  "mebroutines/install"
  "mebroutines/start"
  "mebroutines/backup"
)

const version = "1.2.0"

var is_debug bool

func doSelfUpdate() bool {
  v := semver.MustParse(version)

  if (mebroutines.Is_debug()) {
    selfupdate.EnableLog()
  }

  latest, err := selfupdate.UpdateSelf(v, "digabi/naksu")
  if err != nil {
    mebroutines.Message_warning(fmt.Sprintf("Binary update failed: %s", err))
    return false
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

  // Check whether we have a terminal (restart with x-terminal-emulator, if missing)
  if (! mebroutines.ExistsStdin()) {
    path_to_me, _ := osext.Executable()
    command_args := []string{"x-terminal-emulator", "-e", path_to_me}

    mebroutines.Message_debug(fmt.Sprintf("No stdin, restarting with terminal: %s", strings.Join(command_args, " ")))
    _, _ = mebroutines.Run_get_output(command_args)
    mebroutines.Message_debug(fmt.Sprintf("No stdin, returned from %s", strings.Join(command_args, " ")))

    // Normal termination
    os.Exit(0)
  }

  // UI (main menu)

  err := ui.Main(func () {
    // Define main window
    button_start_server := ui.NewButton("Start Stickless Exam Server")
    button_get_server := ui.NewButton("Install or update Abitti Stickless Exam Server")
    button_switch_server := ui.NewButton("Install or update Stickless Matriculation Exam Server")
    button_make_backup := ui.NewButton("Make Stickless Exam Server Backup")
    button_exit := ui.NewButton("Exit")

    group_common := ui.NewGroup("Basic Functions")
    group_common.SetMargined(true)
    box_common := ui.NewVerticalBox()
    box_common.Append(button_start_server, false)
    box_common.Append(button_exit, false)
    group_common.SetChild(box_common)

    group_abitti := ui.NewGroup("Abitti")
    group_abitti.SetMargined(true)
    box_abitti := ui.NewVerticalBox()
    box_abitti.Append(button_get_server, false)
    group_abitti.SetChild(box_abitti)

    group_matric := ui.NewGroup("Matriculation Exam")
    group_matric.SetMargined(true)
    box_matric := ui.NewVerticalBox()
    box_matric.Append(button_switch_server, false)
    box_matric.Append(button_make_backup, false)
    group_matric.SetChild(box_matric)

    box := ui.NewVerticalBox()
    box.Append(group_common, false)
    box.Append(group_abitti, false)
    box.Append(group_matric, false)

    window := ui.NewWindow(fmt.Sprintf("naksu %s", version), 1, 1, false)

    mebroutines.Set_main_window(window)

    // Run auto-update
    if doSelfUpdate() {
      mebroutines.Message_warning("naksu has been automatically updated. Please restart naksu.")
      os.Exit(0)
    }

    window.SetMargined(true)
		window.SetChild(box)

    // Define Backup SaveAs window/dialog
    backup_label := ui.NewLabel("Please select target path")

    backup_combobox := ui.NewCombobox()
    // Refresh media selection
    backup_media_path := backup_combobox_populate(backup_combobox)

    backup_button_save := ui.NewButton("Save")
    backup_button_cancel := ui.NewButton("Cancel")

    backup_box := ui.NewVerticalBox()
    backup_box.Append(backup_label, false)
    backup_box.Append(backup_combobox, false)
    backup_box.Append(backup_button_save, false)
    backup_box.Append(backup_button_cancel, false)

    backup_window := ui.NewWindow("naksu: SaveTo", 1, 1, false)

    backup_window.SetMargined(true)
    backup_window.SetChild(backup_box)

    // Define actions for main window
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
        install.Do_get_server("")
        os.Exit(0)
      })
    })

    button_switch_server.OnClicked(func(*ui.Button) {
      path_new_vagrantfile := ui.OpenFile(window)

      if path_new_vagrantfile == "" {
        mebroutines.Message_error("Did not get a path for a new Vagrantfile")
      }

      window.Hide()
      ui.QueueMain(func () {
        install.Do_get_server(path_new_vagrantfile)
        os.Exit(0)
      })
    })

    button_make_backup.OnClicked(func(*ui.Button) {
      window.Hide()
      backup_window.Show()
    })

    button_exit.OnClicked(func(*ui.Button) {
      mebroutines.Message_debug("Exiting by user request")
      os.Exit(0)
    })

    // Define actions for SaveAs window/dialog
    backup_button_save.OnClicked(func(*ui.Button) {
      backup_window.Hide()

      path_backup := fmt.Sprintf("%s%s%s", backup_media_path[backup_combobox.Selected()], string(os.PathSeparator), backup.Get_backup_filename(time.Now()))

      ui.QueueMain(func () {
        backup.Do_make_backup(path_backup)
        os.Exit(0)
      })
    })

    backup_button_cancel.OnClicked(func(*ui.Button) {
      backup_window.Hide()
      window.Show()
    })

    window.Show()
    backup_window.Hide()

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

func backup_combobox_populate (combobox *ui.Combobox) []string {
  // Get list of backup locations (as there is not SaveAs/directory dialog in libui)
  backup_media := backup.Get_backup_media()

  // Collect all paths to this slice
  media_path := make([]string, len(backup_media))
  media_path_n := 0

  for this_path := range backup_media {
    combobox.Append(fmt.Sprintf("%s [%s]", backup_media[this_path], this_path))

    media_path[media_path_n] = this_path
    media_path_n++
  }

  return media_path
}
