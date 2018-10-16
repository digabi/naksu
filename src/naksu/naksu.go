// required by selfupdate (needs context)
// +build go1.7

package main

import (
	"flag"
	"fmt"
	"naksu/mebroutines"
	"naksu/mebroutines/backup"
	"naksu/mebroutines/install"
	"naksu/mebroutines/start"
	"naksu/progress"
	"naksu/xlate"
	"os"
	"strings"
	"time"

	"github.com/andlabs/ui"
	"github.com/blang/semver"
	"github.com/kardianos/osext"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const version = "1.5.0"
const lowDiskLimit = 5000000 // 5 Gb

// URLTest is testing URL for checking network connection
const URLTest = "http://static.abitti.fi/usbimg/qa/latest.txt"

var isDebug bool

func doSelfUpdate() bool {
	v := semver.MustParse(version)

	if mebroutines.IsDebug() {
		selfupdate.EnableLog()
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

func main() {
	// Set default UI language
	xlate.SetLanguage("fi")

	// Process command line parameters
	flag.BoolVar(&isDebug, "debug", false, "Turn debugging on")
	flag.Parse()

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
		pathToMe, _ := osext.Executable()
		commandArgs := []string{"x-terminal-emulator", "-e", pathToMe}

		mebroutines.LogDebug(fmt.Sprintf("No stdin, restarting with terminal: %s", strings.Join(commandArgs, " ")))
		_, _ = mebroutines.RunAndGetOutput(commandArgs)
		mebroutines.LogDebug(fmt.Sprintf("No stdin, returned from %s", strings.Join(commandArgs, " ")))

		// Normal termination
		os.Exit(0)
	}

	// Get list of backup locations (as there is not SaveAs/directory dialog in libui)
	// We do this before starting GUI to avoid "cannot change thread mode" in Windows WMI call
	backupMedia := backup.GetBackupMedia()

	// UI (main menu)

	err := ui.Main(func() {
		// Define main window
		buttonStartServer := ui.NewButton(xlate.Get("Start Stickless Exam Server"))
		buttonGetServer := ui.NewButton("Install or update Abitti Stickless Exam Server")
		buttonSwitchServer := ui.NewButton("Install or update Stickless Matriculation Exam Server")
		buttonMakeBackup := ui.NewButton("Make Stickless Exam Server Backup")
		buttonMebshare := ui.NewButton("Open virtual USB stick (ktp-jako)")

		comboboxLang := ui.NewCombobox()
		comboboxLang.Append("Suomeksi")
		comboboxLang.Append("På svenska")
		comboboxLang.Append("In English")
		comboboxLang.SetSelected(0)

		labelBox := ui.NewLabel("")
		labelStatus := ui.NewLabel("")

		checkboxAdvanced := ui.NewCheckbox("")

		// Box version and language selection dropdown
		boxBasicUpper := ui.NewHorizontalBox()
		boxBasicUpper.SetPadded(true)
		boxBasicUpper.Append(labelBox, true)
		boxBasicUpper.Append(comboboxLang, false)

		boxBasic := ui.NewVerticalBox()
		boxBasic.SetPadded(true)
		boxBasic.Append(boxBasicUpper, true)
		boxBasic.Append(labelStatus, true)
		boxBasic.Append(buttonStartServer, true)
		boxBasic.Append(buttonMebshare, true)
		boxBasic.Append(checkboxAdvanced, true)

		boxAdvanced := ui.NewVerticalBox()
		boxAdvanced.SetPadded(true)
		boxAdvanced.Append(buttonGetServer, true)
		boxAdvanced.Append(buttonSwitchServer, true)
		boxAdvanced.Append(buttonMakeBackup, true)

		box := ui.NewVerticalBox()
		box.Append(boxBasic, false)
		box.Append(boxAdvanced, false)

		window := ui.NewWindow(fmt.Sprintf("naksu %s", version), 1, 1, false)

		mebroutines.SetMainWindow(window)
		progress.SetProgressLabel(labelStatus)

		// Run auto-update
		if doSelfUpdate() {
			mebroutines.ShowWarningMessage("naksu has been automatically updated. Please restart naksu.")
			os.Exit(0)
		}

		// Define command channel & goroutine for disabling/enabling main UI buttons
		mainUIStatus := make(chan string)
		mainUINetupdate := time.NewTicker(5 * time.Second)
		go func() {
			lastStatus := ""
			for {
				select {
				case <-mainUINetupdate.C:
					if lastStatus == "enable" {
						// Require network connection for install/update

						ui.QueueMain(func() {
							if install.TestHTTPGet(URLTest) {
								buttonGetServer.Enable()
								buttonSwitchServer.Enable()
							} else {
								buttonGetServer.Disable()
								buttonSwitchServer.Disable()
							}
						})
					}
				case newStatus := <-mainUIStatus:
					mebroutines.LogDebug(fmt.Sprintf("main_ui_status: %s", newStatus))
					// Got new status
					if newStatus == "enable" {
						mebroutines.LogDebug("enable ui")

						ui.QueueMain(func() {
							comboboxLang.Enable()

							buttonStartServer.Enable()
							buttonMebshare.Enable()

							// Require network connection for install/update
							if install.TestHTTPGet(URLTest) {
								buttonGetServer.Enable()
								buttonSwitchServer.Enable()
							} else {
								buttonGetServer.Disable()
								buttonSwitchServer.Disable()
							}
							buttonMakeBackup.Enable()
						})

						lastStatus = newStatus
					}
					if newStatus == "disable" {
						mebroutines.LogDebug("disable ui")

						ui.QueueMain(func() {
							comboboxLang.Disable()

							buttonStartServer.Disable()
							buttonMebshare.Enable()

							buttonGetServer.Disable()
							buttonSwitchServer.Disable()
							buttonMakeBackup.Disable()
						})

						lastStatus = newStatus
					}
				}
			}
		}()

		mainUIStatus <- "enable"

		window.SetMargined(true)
		window.SetChild(box)

		// Advanced group is hidden by default
		boxAdvanced.Hide()

		// Define Backup SaveAs window/dialog
		backupLabel := ui.NewLabel("Please select target path")

		backupCombobox := ui.NewCombobox()
		// Refresh media selection
		backupMediaPath := backupComboboxPopulate(backupMedia, backupCombobox)

		backupButtonSave := ui.NewButton("Save")
		backupButtonCancel := ui.NewButton("Cancel")

		backupBox := ui.NewVerticalBox()
		backupBox.SetPadded(true)
		backupBox.Append(backupLabel, false)
		backupBox.Append(backupCombobox, false)
		backupBox.Append(backupButtonSave, false)
		backupBox.Append(backupButtonCancel, false)

		backupWindow := ui.NewWindow("", 1, 1, false)

		backupWindow.SetMargined(true)
		backupWindow.SetChild(backupBox)

		// (Re)write UI labels
		rewriteUILabels := func() {
			buttonStartServer.SetText(xlate.Get("Start Stickless Exam Server"))
			buttonGetServer.SetText(xlate.Get("Install or update Abitti Stickless Exam Server"))
			buttonSwitchServer.SetText(xlate.Get("Install or update Stickless Matriculation Exam Server"))
			buttonMakeBackup.SetText(xlate.Get("Make Stickless Exam Server Backup"))
			buttonMebshare.SetText(xlate.Get("Open virtual USB stick (ktp-jako)"))

			labelBox.SetText(fmt.Sprintf(xlate.Get("Current version: %s"), mebroutines.GetVagrantBoxVersion()))

			checkboxAdvanced.SetText(xlate.Get("Show management features"))

			backupWindow.SetTitle(xlate.Get("naksu: SaveTo"))
			backupLabel.SetText(xlate.Get("Please select target path"))
			backupButtonSave.SetText(xlate.Get("Save"))
			backupButtonCancel.SetText(xlate.Get("Cancel"))
		}

		// Set UI labels with default language
		rewriteUILabels()

		// Define language selection action main window
		comboboxLang.OnSelected(func(*ui.Combobox) {
			switch comboboxLang.Selected() {
			case 0:
				xlate.SetLanguage("fi")
			case 1:
				xlate.SetLanguage("sv")
			case 2:
				xlate.SetLanguage("en")
			}

			rewriteUILabels()
		})

		// Show/hide advanced features
		checkboxAdvanced.OnToggled(func(*ui.Checkbox) {
			switch checkboxAdvanced.Checked() {
			case true:
				{
					boxAdvanced.Show()
				}
			case false:
				{
					boxAdvanced.Hide()
				}
			}
		})

		// Define actions for main window
		buttonStartServer.OnClicked(func(*ui.Button) {
			go func() {
				mainUIStatus <- "disable"
				start.StartServer()
				mainUIStatus <- "enable"
				progress.SetMessage("")
			}()
		})

		buttonGetServer.OnClicked(func(*ui.Button) {
			chFreeDisk := make(chan int)
			chDiskLowPopup := make(chan bool)

			// Check free disk
			// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
			go func() {
				freeDisk := 0
				if mebroutines.ExistsDir(mebroutines.GetVagrantDirectory()) {
					freeDisk, _ = mebroutines.GetDiskFree(mebroutines.GetVagrantDirectory())
				} else {
					freeDisk, _ = mebroutines.GetDiskFree(mebroutines.GetHomeDirectory())
				}
				chFreeDisk <- freeDisk
			}()

			go func() {
				freeDisk := <-chFreeDisk
				if freeDisk != -1 && freeDisk < lowDiskLimit {
					mebroutines.ShowWarningMessage("Your free disk size is getting low. If update/install process fails please consider freeing some disk space.")
				}

				chDiskLowPopup <- true
			}()

			go func() {
				// Wait until disk low popup has been processed
				<-chDiskLowPopup

				go func() {
					mainUIStatus <- "disable"
					install.GetServer("")
					rewriteUILabels()
					mainUIStatus <- "enable"
					progress.SetMessage("")
				}()
			}()
		})

		buttonSwitchServer.OnClicked(func(*ui.Button) {
			chFreeDisk := make(chan int)
			chDiskLowPopup := make(chan bool)
			chPathNewVagrantfile := make(chan string)

			// Check free disk
			// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
			go func() {
				freeDisk := 0
				if mebroutines.ExistsDir(mebroutines.GetVagrantDirectory()) {
					freeDisk, _ = mebroutines.GetDiskFree(mebroutines.GetVagrantDirectory())
				} else {
					freeDisk, _ = mebroutines.GetDiskFree(mebroutines.GetHomeDirectory())
				}
				chFreeDisk <- freeDisk
			}()

			go func() {
				freeDisk := <-chFreeDisk
				if freeDisk != -1 && freeDisk < lowDiskLimit {
					mebroutines.ShowWarningMessage("Your free disk size is getting low. If update/install process fails please consider freeing some disk space.")
				}

				chDiskLowPopup <- true
			}()

			go func() {
				// Wait until free disk check has been carried out
				<-chDiskLowPopup

				ui.QueueMain(func() {
					pathNewVagrantfile := ui.OpenFile(window)
					chPathNewVagrantfile <- pathNewVagrantfile
				})
			}()

			go func() {
				// Wait until you have path_new_vagrantfile
				pathNewVagrantfile := <-chPathNewVagrantfile

				// Path to ~/ktp/Vagrantfile
				pathOldVagrantfile := mebroutines.GetVagrantDirectory() + string(os.PathSeparator) + "Vagrantfile"

				if pathNewVagrantfile == "" {
					mebroutines.ShowErrorMessage(xlate.Get("Did not get a path for a new Vagrantfile"))
				} else if pathNewVagrantfile == pathOldVagrantfile {
					mebroutines.ShowErrorMessage(xlate.Get("Please place the new Exam Vagrantfile to another location (e.g. desktop or home directory)"))
				} else {
					go func() {
						mainUIStatus <- "disable"
						install.GetServer(pathNewVagrantfile)
						rewriteUILabels()
						mainUIStatus <- "enable"
						progress.SetMessage("")
					}()
				}
			}()
		})

		buttonMakeBackup.OnClicked(func(*ui.Button) {
			mainUIStatus <- "disable"
			backupWindow.Show()
		})

		buttonMebshare.OnClicked(func(*ui.Button) {
			mebroutines.LogDebug("Opening MEB share (~/ktp-jako)")
			mebroutines.OpenMebShare()
		})

		// Define actions for SaveAs window/dialog
		backupButtonSave.OnClicked(func(*ui.Button) {
			pathBackup := fmt.Sprintf("%s%s%s", backupMediaPath[backupCombobox.Selected()], string(os.PathSeparator), backup.GetBackupFilename(time.Now()))

			chFreeDisk := make(chan int)
			chDiskLowPopup := make(chan bool)

			// Check free disk
			// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
			go func() {
				freeDisk, _ := mebroutines.GetDiskFree(fmt.Sprintf("%s%s", backupMediaPath[backupCombobox.Selected()], string(os.PathSeparator)))
				chFreeDisk <- freeDisk
			}()

			go func() {
				freeDisk := <-chFreeDisk
				if freeDisk != -1 && freeDisk < lowDiskLimit {
					mebroutines.ShowWarningMessage("Your free disk size is getting low. If backup process fails please consider freeing some disk space.")
				}
				chDiskLowPopup <- true
			}()

			go func() {
				<-chDiskLowPopup

				go func() {
					backupWindow.Hide()
					backup.MakeBackup(pathBackup)
					mainUIStatus <- "enable"
				}()
			}()
		})

		backupButtonCancel.OnClicked(func(*ui.Button) {
			backupWindow.Hide()
			mainUIStatus <- "enable"
		})

		window.OnClosing(func(*ui.Window) bool {
			mebroutines.LogDebug("User exists through window exit")
			ui.Quit()
			return true
		})

		window.Show()
		mainUIStatus <- "enable"
		backupWindow.Hide()

		// Make sure we have vagrant
		if !mebroutines.IfFoundVagrant() {
			mebroutines.ShowErrorMessage(xlate.Get("Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?"))
		}

		// Make sure we have VBoxManage
		if !mebroutines.IfFoundVBoxManage() {
			mebroutines.ShowErrorMessage(xlate.Get("Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?"))
		}

		// Check if home directory contains non-american characters which may cause problems to vagrant
		if mebroutines.IfIntlCharsInPath(mebroutines.GetHomeDirectory()) {
			mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Your home directory path (%s) contains characters which may cause problems to Vagrant."), mebroutines.GetHomeDirectory()))
		}

	})

	if err != nil {
		panic(err)
	}

	mebroutines.LogDebug("Exiting GUI loop")
}

func backupComboboxPopulate(backupMedia map[string]string, combobox *ui.Combobox) []string {
	// Collect all paths to this slice
	mediaPath := make([]string, len(backupMedia))
	mediaPathN := 0

	for thisPath := range backupMedia {
		combobox.Append(fmt.Sprintf("%s [%s]", backupMedia[thisPath], thisPath))

		mediaPath[mediaPathN] = thisPath
		mediaPathN++
	}

	return mediaPath
}
