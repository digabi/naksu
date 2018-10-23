package main

import (
	"fmt"
	"naksu/mebroutines"
	"naksu/mebroutines/backup"
	"naksu/mebroutines/install"
	"naksu/mebroutines/start"
	"naksu/progress"
	"naksu/xlate"
	"os"
	"time"

	"github.com/andlabs/ui"
)

var window *ui.Window

var buttonStartServer *ui.Button
var buttonGetServer *ui.Button
var buttonSwitchServer *ui.Button
var buttonMakeBackup *ui.Button
var buttonMebShare *ui.Button

var comboboxLang *ui.Combobox

var labelBox *ui.Label
var labelStatus *ui.Label

var checkboxAdvanced *ui.Checkbox

var boxBasicUpper *ui.Box
var boxBasic *ui.Box
var boxAdvanced *ui.Box
var box *ui.Box

var backupWindow *ui.Window

var backupCombobox *ui.Combobox
var backupButtonSave *ui.Button
var backupButtonCancel *ui.Button

var backupBox *ui.Box

var backupLabel *ui.Label

var backupMediaPath []string

func createMainWindowElements() {
	// Define main window
	buttonStartServer = ui.NewButton(xlate.Get("Start Stickless Exam Server"))
	buttonGetServer = ui.NewButton("Install or update Abitti Stickless Exam Server")
	buttonSwitchServer = ui.NewButton("Install or update Stickless Matriculation Exam Server")
	buttonMakeBackup = ui.NewButton("Make Stickless Exam Server Backup")
	buttonMebShare = ui.NewButton("Open virtual USB stick (ktp-jako)")

	comboboxLang = ui.NewCombobox()
	comboboxLang.Append("Suomeksi")
	comboboxLang.Append("På svenska")
	comboboxLang.Append("In English")
	comboboxLang.SetSelected(0)

	labelBox = ui.NewLabel("")
	labelStatus = ui.NewLabel("")

	checkboxAdvanced = ui.NewCheckbox("")

	// Box version and language selection dropdown
	boxBasicUpper = ui.NewHorizontalBox()
	boxBasicUpper.SetPadded(true)
	boxBasicUpper.Append(labelBox, true)
	boxBasicUpper.Append(comboboxLang, false)

	boxBasic = ui.NewVerticalBox()
	boxBasic.SetPadded(true)
	boxBasic.Append(boxBasicUpper, true)
	boxBasic.Append(labelStatus, true)
	boxBasic.Append(buttonStartServer, true)
	boxBasic.Append(buttonMebShare, true)
	boxBasic.Append(checkboxAdvanced, true)

	boxAdvanced = ui.NewVerticalBox()
	boxAdvanced.SetPadded(true)
	boxAdvanced.Append(buttonGetServer, true)
	boxAdvanced.Append(buttonSwitchServer, true)
	boxAdvanced.Append(buttonMakeBackup, true)

	box = ui.NewVerticalBox()
	box.Append(boxBasic, false)
	box.Append(boxAdvanced, false)

	window = ui.NewWindow(fmt.Sprintf("naksu %s", version), 1, 1, false)
}

func createBackupElements(backupMedia map[string]string) {
	// Define Backup SaveAs window/dialog
	backupLabel = ui.NewLabel("Please select target path")

	backupCombobox = ui.NewCombobox()
	// Refresh media selection
	backupMediaPath = populateBackupCombobox(backupMedia, backupCombobox)

	backupButtonSave = ui.NewButton("Save")
	backupButtonCancel = ui.NewButton("Cancel")

	backupBox = ui.NewVerticalBox()
	backupBox.SetPadded(true)
	backupBox.Append(backupLabel, false)
	backupBox.Append(backupCombobox, false)
	backupBox.Append(backupButtonSave, false)
	backupBox.Append(backupButtonCancel, false)

	backupWindow = ui.NewWindow("", 1, 1, false)

	backupWindow.SetMargined(true)
	backupWindow.SetChild(backupBox)
}

func populateBackupCombobox(backupMedia map[string]string, combobox *ui.Combobox) []string {
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

func setupMainLoop(mainUIStatus chan string, mainUINetupdate *time.Ticker) {
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
						buttonMebShare.Enable()

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
						buttonMebShare.Enable()

						buttonGetServer.Disable()
						buttonSwitchServer.Disable()
						buttonMakeBackup.Disable()
					})

					lastStatus = newStatus
				}
			}
		}
	}()
}

func translateUILabels() {
	buttonStartServer.SetText(xlate.Get("Start Stickless Exam Server"))
	buttonGetServer.SetText(xlate.Get("Install or update Abitti Stickless Exam Server"))
	buttonSwitchServer.SetText(xlate.Get("Install or update Stickless Matriculation Exam Server"))
	buttonMakeBackup.SetText(xlate.Get("Make Stickless Exam Server Backup"))
	buttonMebShare.SetText(xlate.Get("Open virtual USB stick (ktp-jako)"))

	labelBox.SetText(fmt.Sprintf(xlate.Get("Current version: %s"), mebroutines.GetVagrantBoxVersion()))

	checkboxAdvanced.SetText(xlate.Get("Show management features"))

	backupWindow.SetTitle(xlate.Get("naksu: SaveTo"))
	backupLabel.SetText(xlate.Get("Please select target path"))
	backupButtonSave.SetText(xlate.Get("Save"))
	backupButtonCancel.SetText(xlate.Get("Cancel"))

}

func disableUI(mainUIStatus chan string) {
	mainUIStatus <- "disable"
}

func enableUI(mainUIStatus chan string) {
	mainUIStatus <- "enable"
}

func bindLanguageSwitching() {
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

		translateUILabels()
	})
}

func bindAdvancedToggle() {
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
}

func bindUIDisableOnStart(mainUIStatus chan string) {
	// Define actions for main window
	buttonStartServer.OnClicked(func(*ui.Button) {
		go func() {
			disableUI(mainUIStatus)
			start.StartServer()
			enableUI(mainUIStatus)
			progress.SetMessage("")
		}()
	})

}

func bindOnGetServer(mainUIStatus chan string) {
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
				disableUI(mainUIStatus)
				install.GetServer("")
				translateUILabels()
				enableUI(mainUIStatus)
				progress.SetMessage("")
			}()
		}()
	})
}

func bindOnSwitchServer(mainUIStatus chan string) {
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
					disableUI(mainUIStatus)
					install.GetServer(pathNewVagrantfile)
					translateUILabels()
					enableUI(mainUIStatus)
					progress.SetMessage("")
				}()
			}
		}()
	})

}

func bindOnMakeBackup(mainUIStatus chan string) {
	buttonMakeBackup.OnClicked(func(*ui.Button) {
		disableUI(mainUIStatus)
		backupWindow.Show()
	})
}

func bindOnMebShare() {
	buttonMebShare.OnClicked(func(*ui.Button) {
		mebroutines.LogDebug("Opening MEB share (~/ktp-jako)")
		mebroutines.OpenMebShare()
	})
}

func bindOnBackup(mainUIStatus chan string) {
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
				enableUI(mainUIStatus)
			}()
		}()
	})

	backupButtonCancel.OnClicked(func(*ui.Button) {
		backupWindow.Hide()
		enableUI(mainUIStatus)
	})
}

// RunUI sets up user interface and starts running it. function exists when application exits
func RunUI() error {

	// Get list of backup locations (as there is not SaveAs/directory dialog in libui)
	// We do this before starting GUI to avoid "cannot change thread mode" in Windows WMI call
	backupMedia := backup.GetBackupMedia()

	// UI (main menu)
	return ui.Main(func() {

		createMainWindowElements()

		createBackupElements(backupMedia)

		mebroutines.SetMainWindow(window)
		progress.SetProgressLabel(labelStatus)

		// Define command channel & goroutine for disabling/enabling main UI buttons
		mainUIStatus := make(chan string)
		mainUINetupdate := time.NewTicker(5 * time.Second)

		setupMainLoop(mainUIStatus, mainUINetupdate)

		enableUI(mainUIStatus)

		window.SetMargined(true)
		window.SetChild(box)

		// Advanced group is hidden by default
		boxAdvanced.Hide()

		// Set UI labels with default language
		translateUILabels()

		bindLanguageSwitching()
		bindAdvancedToggle()

		bindUIDisableOnStart(mainUIStatus)

		// Bind buttons
		bindOnGetServer(mainUIStatus)
		bindOnSwitchServer(mainUIStatus)
		bindOnMakeBackup(mainUIStatus)
		bindOnMebShare()

		bindOnBackup(mainUIStatus)

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
}
