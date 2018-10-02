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
  "xlate"
  "progress"
)

const version = "1.5.0"
const low_disk_limit = 5000000 // 5 Gb
// Test URL for checking network connection
const URL_TEST = "http://static.abitti.fi/usbimg/qa/latest.txt"

var is_debug bool

func doSelfUpdate() bool {
  v := semver.MustParse(version)

  if (mebroutines.Is_debug()) {
    selfupdate.EnableLog()
  }

  latest, err := selfupdate.UpdateSelf(v, "digabi/naksu")
  if err != nil {
    mebroutines.Message_warning(fmt.Sprintf(xlate.Get("Naksu update failed. Maybe you don't have network connection?\n\nError: %s"), err))
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
  // Set default UI language
  xlate.SetLanguage("fi")

  // Process command line parameters
  flag.BoolVar(&is_debug, "debug", false, "Turn debugging on")
  flag.Parse()

  mebroutines.Set_debug(is_debug)

  // Determine/set path for debug log
  log_path := mebroutines.Get_vagrant_directory()
  if mebroutines.ExistsDir(log_path) {
    mebroutines.Set_debug_filename(log_path + string(os.PathSeparator) + "naksu_lastlog.txt")
  } else {
    mebroutines.Set_debug_filename(os.TempDir() + string(os.PathSeparator)+ "naksu_lastlog.txt")
  }

  mebroutines.Message_debug(fmt.Sprintf("This is Naksu %s. Hello world!", version))

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

  // Get list of backup locations (as there is not SaveAs/directory dialog in libui)
  // We do this before starting GUI to avoid "cannot change thread mode" in Windows WMI call
  backup_media := backup.Get_backup_media()

  // UI (main menu)

  err := ui.Main(func () {
    // Define main window
    button_start_server := ui.NewButton(xlate.Get("Start Stickless Exam Server"))
    button_get_server := ui.NewButton("Install or update Abitti Stickless Exam Server")
    button_switch_server := ui.NewButton("Install or update Stickless Matriculation Exam Server")
    button_make_backup := ui.NewButton("Make Stickless Exam Server Backup")
    button_mebshare := ui.NewButton("Open virtual USB stick (ktp-jako)")

    combobox_lang := ui.NewCombobox()
    combobox_lang.Append("Suomeksi")
    combobox_lang.Append("På svenska")
    combobox_lang.Append("In English")
    combobox_lang.SetSelected(0)

    label_box := ui.NewLabel("")
    label_status := ui.NewLabel("")

    checkbox_advanced := ui.NewCheckbox("")

    // Box version and language selection dropdown
    box_basic_upper := ui.NewHorizontalBox()
    box_basic_upper.SetPadded(true)
    box_basic_upper.Append(label_box, true)
    box_basic_upper.Append(combobox_lang, false)

    box_basic := ui.NewVerticalBox()
    box_basic.SetPadded(true)
    box_basic.Append(box_basic_upper, true)
    box_basic.Append(label_status, true)
    box_basic.Append(button_start_server, true)
    box_basic.Append(button_mebshare, true)
    box_basic.Append(checkbox_advanced, true)

    box_advanced := ui.NewVerticalBox()
    box_advanced.SetPadded(true)
    box_advanced.Append(button_get_server, true)
    box_advanced.Append(button_switch_server, true)
    box_advanced.Append(button_make_backup, true)

    box := ui.NewVerticalBox()
    box.Append(box_basic, false)
    box.Append(box_advanced, false)

    window := ui.NewWindow(fmt.Sprintf("naksu %s", version), 1, 1, false)

    mebroutines.Set_main_window(window)
    progress.Set_label_object(label_status)

    // Run auto-update
    if doSelfUpdate() {
      mebroutines.Message_warning("naksu has been automatically updated. Please restart naksu.")
      os.Exit(0)
    }

    // Define command channel & goroutine for disabling/enabling main UI buttons
    main_ui_status := make(chan string)
    main_ui_netupdate := time.NewTicker(5 * time.Second)
    go func() {
      last_status := ""
      for {
        select {
        case <- main_ui_netupdate.C:
          if last_status == "enable" {
            // Require network connection for install/update

            ui.QueueMain(func () {
              if install.If_http_get(URL_TEST) {
                button_get_server.Enable()
                button_switch_server.Enable()
              } else {
                button_get_server.Disable()
                button_switch_server.Disable()
              }
            })
          }
        case new_status := <- main_ui_status:
          mebroutines.Message_debug(fmt.Sprintf("main_ui_status: %s", new_status))
          // Got new status
          if new_status == "enable" {
            mebroutines.Message_debug("enable ui")

            ui.QueueMain(func() {
              combobox_lang.Enable()

              button_start_server.Enable()
              button_mebshare.Enable()

              // Require network connection for install/update
              if install.If_http_get(URL_TEST) {
                button_get_server.Enable()
                button_switch_server.Enable()
              } else {
                button_get_server.Disable()
                button_switch_server.Disable()
              }
              button_make_backup.Enable()
            })

            last_status = new_status
          }
          if new_status == "disable" {
            mebroutines.Message_debug("disable ui")

            ui.QueueMain(func () {
              combobox_lang.Disable()

              button_start_server.Disable()
              button_mebshare.Enable()

              button_get_server.Disable()
              button_switch_server.Disable()
              button_make_backup.Disable()
            })

            last_status = new_status
          }
        }
      }
    }()

    main_ui_status <- "enable"

    window.SetMargined(true)
		window.SetChild(box)

    // Advanced group is hidden by default
    box_advanced.Hide()

    // Define Backup SaveAs window/dialog
    backup_label := ui.NewLabel("Please select target path")

    backup_combobox := ui.NewCombobox()
    // Refresh media selection
    backup_media_path := backup_combobox_populate(backup_media, backup_combobox)

    backup_button_save := ui.NewButton("Save")
    backup_button_cancel := ui.NewButton("Cancel")

    backup_box := ui.NewVerticalBox()
    backup_box.SetPadded(true)
    backup_box.Append(backup_label, false)
    backup_box.Append(backup_combobox, false)
    backup_box.Append(backup_button_save, false)
    backup_box.Append(backup_button_cancel, false)

    backup_window := ui.NewWindow("", 1, 1, false)

    backup_window.SetMargined(true)
    backup_window.SetChild(backup_box)

    // (Re)write UI labels
    rewrite_ui_labels := func () {
      button_start_server.SetText(xlate.Get("Start Stickless Exam Server"))
      button_get_server.SetText(xlate.Get("Install or update Abitti Stickless Exam Server"))
      button_switch_server.SetText(xlate.Get("Install or update Stickless Matriculation Exam Server"))
      button_make_backup.SetText(xlate.Get("Make Stickless Exam Server Backup"))
      button_mebshare.SetText(xlate.Get("Open virtual USB stick (ktp-jako)"))

      label_box.SetText(fmt.Sprintf(xlate.Get("Current version: %s"), mebroutines.Get_vagrantbox_version()))

      checkbox_advanced.SetText(xlate.Get("Show management features"))

      backup_window.SetTitle(xlate.Get("naksu: SaveTo"))
      backup_label.SetText(xlate.Get("Please select target path"))
      backup_button_save.SetText(xlate.Get("Save"))
      backup_button_cancel.SetText(xlate.Get("Cancel"))
    }

    // Set UI labels with default language
    rewrite_ui_labels()

    // Define language selection action main window
    combobox_lang.OnSelected(func(*ui.Combobox) {
      switch combobox_lang.Selected() {
      case 0: xlate.SetLanguage("fi")
      case 1: xlate.SetLanguage("sv")
      case 2: xlate.SetLanguage("en")
      }

      rewrite_ui_labels()
    })

    // Show/hide advanced features
    checkbox_advanced.OnToggled(func(*ui.Checkbox) {
      switch checkbox_advanced.Checked() {
      case true: {
        box_advanced.Show()
      }
      case false: {
        box_advanced.Hide()
      }
      }
    })

    // Define actions for main window
    button_start_server.OnClicked(func(*ui.Button) {
      go func () {
        main_ui_status <- "disable"
        start.Do_start_server()
        main_ui_status <- "enable"
        progress.Set_message("")
      }()
    })

    button_get_server.OnClicked(func(*ui.Button) {
      ch_free_disk := make(chan int)
      ch_disk_low_popup := make(chan bool)

      // Check free disk
      // Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
      go func () {
        free_disk := 0
        if mebroutines.ExistsDir(mebroutines.Get_vagrant_directory()) {
          free_disk,_ = mebroutines.Get_disk_free(mebroutines.Get_vagrant_directory())
        } else {
          free_disk,_ = mebroutines.Get_disk_free(mebroutines.Get_home_directory())
        }
        ch_free_disk <- free_disk
      }()

      go func () {
        free_disk := <- ch_free_disk
        if (free_disk != -1 && free_disk < low_disk_limit) {
          mebroutines.Message_warning("Your free disk size is getting low. If update/install process fails please consider freeing some disk space.")
        }

        ch_disk_low_popup <- true
      }()

      go func () {
        // Wait until disk low popup has been processed
        <- ch_disk_low_popup

        go func () {
          main_ui_status <- "disable"
          install.Do_get_server("")
          rewrite_ui_labels()
          main_ui_status <- "enable"
          progress.Set_message("")
        }()
      }()
    })

    button_switch_server.OnClicked(func(*ui.Button) {
      ch_free_disk := make(chan int)
      ch_disk_low_popup := make(chan bool)
      ch_path_new_vagrantfile := make(chan string)

      // Check free disk
      // Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
      go func () {
        free_disk := 0
        if mebroutines.ExistsDir(mebroutines.Get_vagrant_directory()) {
          free_disk,_ = mebroutines.Get_disk_free(mebroutines.Get_vagrant_directory())
        } else {
          free_disk,_ = mebroutines.Get_disk_free(mebroutines.Get_home_directory())
        }
        ch_free_disk <- free_disk
      }()

      go func () {
        free_disk := <- ch_free_disk
        if (free_disk != -1 && free_disk < low_disk_limit) {
          mebroutines.Message_warning("Your free disk size is getting low. If update/install process fails please consider freeing some disk space.")
        }

        ch_disk_low_popup <- true
      }()

      go func () {
        // Wait until free disk check has been carried out
        <- ch_disk_low_popup

        ui.QueueMain(func () {
          path_new_vagrantfile := ui.OpenFile(window)
          ch_path_new_vagrantfile <- path_new_vagrantfile
        })
      }()

      go func () {
        // Wait until you have path_new_vagrantfile
        path_new_vagrantfile := <- ch_path_new_vagrantfile

        // Path to ~/ktp/Vagrantfile
        path_the_vagrantfile := mebroutines.Get_vagrant_directory()+string(os.PathSeparator)+"Vagrantfile"

        if path_new_vagrantfile == "" {
          mebroutines.Message_error(xlate.Get("Did not get a path for a new Vagrantfile"))
        } else if path_new_vagrantfile == path_the_vagrantfile {
          mebroutines.Message_error(xlate.Get("Please place the new Exam Vagrantfile to another location (e.g. desktop or home directory)"))
        } else {
          go func () {
            main_ui_status <- "disable"
            install.Do_get_server(path_new_vagrantfile)
            rewrite_ui_labels()
            main_ui_status <- "enable"
            progress.Set_message("")
          }()
        }
      }()
    })

    button_make_backup.OnClicked(func(*ui.Button) {
      main_ui_status <- "disable"
      backup_window.Show()
    })

    button_mebshare.OnClicked(func(*ui.Button) {
      mebroutines.Message_debug("Opening MEB share (~/ktp-jako)")
      mebroutines.Open_meb_share()
    })

    // Define actions for SaveAs window/dialog
    backup_button_save.OnClicked(func(*ui.Button) {
      path_backup := fmt.Sprintf("%s%s%s", backup_media_path[backup_combobox.Selected()], string(os.PathSeparator), backup.Get_backup_filename(time.Now()))

      ch_free_disk := make(chan int)
      ch_disk_low_popup := make(chan bool)

      // Check free disk
      // Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
      go func () {
        free_disk,_ := mebroutines.Get_disk_free(fmt.Sprintf("%s%s", backup_media_path[backup_combobox.Selected()], string(os.PathSeparator)))
        ch_free_disk <- free_disk
      }()

      go func () {
        free_disk := <- ch_free_disk
        if (free_disk != -1 && free_disk < low_disk_limit) {
          mebroutines.Message_warning("Your free disk size is getting low. If backup process fails please consider freeing some disk space.")
        }
        ch_disk_low_popup <- true
      }()

      go func () {
        <- ch_disk_low_popup

        go func () {
          backup_window.Hide()
          backup.Do_make_backup(path_backup)
          main_ui_status <- "enable"
        }()
      }()
    })

    backup_button_cancel.OnClicked(func(*ui.Button) {
      backup_window.Hide()
      main_ui_status <- "enable"
    })

    window.OnClosing(func(*ui.Window) bool {
      mebroutines.Message_debug("User exists through window exit")
      ui.Quit()
      return true
    })

    window.Show()
    main_ui_status <- "enable"
    backup_window.Hide()

    // Make sure we have vagrant
  	if (! mebroutines.If_found_vagrant()) {
  		mebroutines.Message_error(xlate.Get("Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?"))
  	}

  	// Make sure we have VBoxManage
  	if (! mebroutines.If_found_vboxmanage()) {
  		mebroutines.Message_error(xlate.Get("Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?"))
  	}

    // Check if home directory contains non-american characters which may cause problems to vagrant
    if (mebroutines.If_intl_chars_in_path(mebroutines.Get_home_directory())) {
      mebroutines.Message_warning(fmt.Sprintf(xlate.Get("Your home directory path (%s) contains characters which may cause problems to Vagrant."), mebroutines.Get_home_directory()))
    }

  })

  if err != nil {
    panic(err)
  }

  mebroutines.Message_debug("Exiting GUI loop")
}

func backup_combobox_populate (backup_media map[string]string, combobox *ui.Combobox) []string {
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
