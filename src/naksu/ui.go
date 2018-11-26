package main

import (
	"fmt"
	"naksu/mebroutines"
	"naksu/mebroutines/backup"
	"naksu/mebroutines/install"
	"naksu/mebroutines/start"
	"naksu/mebroutines/destroy"
	"naksu/mebroutines/remove"
	"naksu/progress"
	"naksu/xlate"
	"os"
	"time"

	"github.com/andlabs/ui"
	"github.com/dustin/go-humanize"
)

var window *ui.Window

var buttonStartServer *ui.Button
var buttonGetServer *ui.Button
var buttonSwitchServer *ui.Button
var buttonDestroyServer *ui.Button
var buttonRemoveServer *ui.Button
var buttonMakeBackup *ui.Button
var buttonMebShare *ui.Button

var comboboxLang *ui.Combobox

var labelBox *ui.Label
var labelStatus *ui.Label
var labelAdvancedUpdate *ui.Label
var labelAdvancedAnnihilate *ui.Label

var checkboxAdvanced *ui.Checkbox

var boxBasicUpper *ui.Box
var boxBasic *ui.Box
var boxAdvancedUpdate *ui.Box
var boxAdvancedAnnihilate *ui.Box
var boxAdvanced *ui.Box
var box *ui.Box

// Backup Dialog Window
var backupWindow *ui.Window

var backupCombobox *ui.Combobox
var backupButtonSave *ui.Button
var backupButtonCancel *ui.Button

var backupBox *ui.Box

var backupLabel *ui.Label

var backupMediaPath []string

// Destroy Confirmation Window
var destroyWindow *ui.Window

var destroyButtonDestroy *ui.Button
var destroyButtonCancel *ui.Button

var destroyBox *ui.Box

var destroyInfoLabel *ui.Label
var destroyQuestionLabel *ui.Label

// Remove Confirmation Window
var removeWindow *ui.Window

var removeButtonRemove *ui.Button
var removeButtonCancel *ui.Button

var removeBox *ui.Box

var removeInfoLabel *ui.Label
var removeQuestionLabel *ui.Label


func createMainWindowElements() {
	// Define main window
	buttonStartServer = ui.NewButton(xlate.Get("Start Exam Server"))
	buttonGetServer = ui.NewButton("Abitti Exam")
	buttonSwitchServer = ui.NewButton("Matriculation Exam")
	buttonDestroyServer = ui.NewButton("Remove Exams")
	buttonRemoveServer = ui.NewButton("Remove Server")
	buttonMakeBackup = ui.NewButton("Make Exam Server Backup")
	buttonMebShare = ui.NewButton("Open virtual USB stick (ktp-jako)")

	comboboxLang = ui.NewCombobox()
	comboboxLang.Append("Suomeksi")
	comboboxLang.Append("På svenska")
	comboboxLang.Append("In English")
	comboboxLang.SetSelected(0)

	labelBox = ui.NewLabel("")
	labelStatus = ui.NewLabel("")
	labelAdvancedUpdate = ui.NewLabel("")
	labelAdvancedAnnihilate = ui.NewLabel("")

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

	boxAdvancedUpdate = ui.NewHorizontalBox()
	boxAdvancedUpdate.SetPadded(true)
	boxAdvancedUpdate.Append(buttonGetServer, true)
	boxAdvancedUpdate.Append(buttonSwitchServer, true)

	boxAdvancedAnnihilate = ui.NewHorizontalBox()
	boxAdvancedAnnihilate.SetPadded(true)
	boxAdvancedAnnihilate.Append(buttonDestroyServer, true)
	boxAdvancedAnnihilate.Append(buttonRemoveServer, true)

	boxAdvanced = ui.NewVerticalBox()
	boxAdvanced.SetPadded(true)
	boxAdvanced.Append(buttonMakeBackup, true)
	boxAdvanced.Append(labelAdvancedUpdate, false)
	boxAdvanced.Append(boxAdvancedUpdate, true)
	boxAdvanced.Append(labelAdvancedAnnihilate, false)
	boxAdvanced.Append(boxAdvancedAnnihilate, true)

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

func createDestroyElements() {
	// Define Destroy Confirmation window/dialog
	destroyInfoLabel = ui.NewLabel("destroyInfoLabel")
	destroyQuestionLabel = ui.NewLabel("destroyQuestionLabel")

	destroyButtonDestroy = ui.NewButton("Yes, Remove")
	destroyButtonCancel = ui.NewButton("Cancel")

	destroyBox = ui.NewVerticalBox()
	destroyBox.SetPadded(true)
	destroyBox.Append(destroyInfoLabel,true)
	destroyBox.Append(destroyQuestionLabel, true)
	destroyBox.Append(destroyButtonDestroy, false)
	destroyBox.Append(destroyButtonCancel, false)

	destroyWindow = ui.NewWindow("", 1, 1, false)

	destroyWindow.SetMargined(true)
	destroyWindow.SetChild(destroyBox)
}

func createRemoveElements() {
	// Define Destroy Confirmation window/dialog
	removeInfoLabel = ui.NewLabel("removeInfoLabel")
	removeQuestionLabel = ui.NewLabel("removeQuestionLabel")

	removeButtonRemove = ui.NewButton("Yes, Remove")
	removeButtonCancel = ui.NewButton("Cancel")

	removeBox = ui.NewVerticalBox()
	removeBox.SetPadded(true)
	removeBox.Append(removeInfoLabel,true)
	removeBox.Append(removeQuestionLabel, true)
	removeBox.Append(removeButtonRemove, false)
	removeBox.Append(removeButtonCancel, false)

	removeWindow = ui.NewWindow("", 1, 1, false)

	removeWindow.SetMargined(true)
	removeWindow.SetChild(removeBox)
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
						buttonDestroyServer.Enable()
						buttonRemoveServer.Enable()
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
						buttonDestroyServer.Disable()
						buttonRemoveServer.Disable()
					})

					lastStatus = newStatus
				}
			}
		}
	}()
}

func translateUILabels() {
	ui.QueueMain(func() {
		buttonStartServer.SetText(xlate.Get("Start Exam Server"))
		buttonGetServer.SetText(xlate.Get("Abitti Exam"))
		buttonSwitchServer.SetText(xlate.Get("Matriculation Exam"))
		buttonDestroyServer.SetText(xlate.Get("Remove Exams"))
		buttonRemoveServer.SetText(xlate.Get("Remove Server"))
		buttonMakeBackup.SetText(xlate.Get("Make Exam Server Backup"))
		buttonMebShare.SetText(xlate.Get("Open virtual USB stick (ktp-jako)"))

		labelBox.SetText(fmt.Sprintf(xlate.Get("Current version: %s"), mebroutines.GetVagrantBoxVersion()))

		checkboxAdvanced.SetText(xlate.Get("Show management features"))
		labelAdvancedUpdate.SetText(xlate.Get("Install/update server for:"))
		labelAdvancedAnnihilate.SetText(xlate.Get("DANGER! Annihilate your server:"))

		backupWindow.SetTitle(xlate.Get("naksu: SaveTo"))
		backupLabel.SetText(xlate.Get("Please select target path"))
		backupButtonSave.SetText(xlate.Get("Save"))
		backupButtonCancel.SetText(xlate.Get("Cancel"))

		destroyWindow.SetTitle(xlate.Get("naksu: Remove Exams"))
		destroyInfoLabel.SetText(xlate.Get("Remove Exams restores server to its initial status.\nExams, responses and logs in the server will be irreversibly deleted.\nIt is recommended to back up your server before doing this."))
		destroyQuestionLabel.SetText(xlate.Get("Do you wish to remove all exams?"))
		destroyButtonDestroy.SetText(xlate.Get("Yes, Remove"))
		destroyButtonCancel.SetText(xlate.Get("Cancel"))

		removeWindow.SetTitle(xlate.Get("naksu: Remove Server"))
		removeInfoLabel.SetText(xlate.Get("Removing server destroys it and all downloaded disk images.\nExams, responses and logs in the server will be irreversibly deleted.\nIt is recommended to back up your server before doing this."))
		removeQuestionLabel.SetText(xlate.Get("Do you wish to remove the server?"))
		removeButtonRemove.SetText(xlate.Get("Yes, Remove"))
		removeButtonCancel.SetText(xlate.Get("Cancel"))
	})
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
			start.Server()
			enableUI(mainUIStatus)
			progress.SetMessage("")
		}()
	})

}

func checkFreeDisk(chFreeDisk chan int) {
	// Check free disk
	// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
	go func() {
		freeDisk := 0
		var err error
		if mebroutines.ExistsDir(mebroutines.GetVagrantDirectory()) {
			freeDisk, err = mebroutines.GetDiskFree(mebroutines.GetVagrantDirectory())
			if err != nil {
				mebroutines.LogDebug("Getting free disk space from Vagrant directory failed")
			}
		} else {
			freeDisk, err = mebroutines.GetDiskFree(mebroutines.GetHomeDirectory())
			if err != nil {
				mebroutines.LogDebug("Getting free disk space from home directory failed")
			}
		}
		chFreeDisk <- freeDisk
	}()
}

func bindOnGetServer(mainUIStatus chan string) {
	buttonGetServer.OnClicked(func(*ui.Button) {
		chFreeDisk := make(chan int)
		chDiskLowPopup := make(chan bool)

		checkFreeDisk(chFreeDisk)

		go func() {
			freeDisk := <-chFreeDisk
			if freeDisk != -1 && freeDisk < lowDiskLimit {
				mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Your free disk size is getting low (%s)."), humanize.Bytes(uint64(freeDisk))))
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

		checkFreeDisk(chFreeDisk)

		go func() {
			freeDisk := <-chFreeDisk
			if freeDisk != -1 && freeDisk < lowDiskLimit {
				mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Your free disk size is getting low (%s)."), humanize.Bytes(uint64(freeDisk))))
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

func bindOnDestroyServer(mainUIStatus chan string) {
	// Define actions for Destroy popup/window
	buttonDestroyServer.OnClicked(func(*ui.Button) {
		disableUI(mainUIStatus)
		destroyWindow.Show()
	})
}

func bindOnRemoveServer(mainUIStatus chan string) {
	// Define actions for Remove popup/window
	buttonRemoveServer.OnClicked(func(*ui.Button) {
		disableUI(mainUIStatus)
		removeWindow.Show()
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

		checkFreeDisk(chFreeDisk)

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

	backupWindow.OnClosing(func(*ui.Window) bool {
		backupWindow.Hide()
		enableUI(mainUIStatus)
		return true
	})
}

func bindOnDestroy(mainUIStatus chan string) {
	// Define actions for Destroy window/dialog

	destroyButtonDestroy.OnClicked(func(*ui.Button) {
		go func () {
			destroyWindow.Hide()
			destroy.Server()
			enableUI(mainUIStatus)
		}()
	})

	destroyButtonCancel.OnClicked(func(*ui.Button) {
		destroyWindow.Hide()
		enableUI(mainUIStatus)
	})

	destroyWindow.OnClosing(func(*ui.Window) bool {
		destroyWindow.Hide()
		enableUI(mainUIStatus)
		return true
	})
}

func bindOnRemove(mainUIStatus chan string) {
	// Define actions for Remove window/dialog

	removeButtonRemove.OnClicked(func(*ui.Button) {
		go func () {
			removeWindow.Hide()

			err := remove.Server()
			if err != nil {
				mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Error while removing server: %v"), err))
				progress.SetMessage(fmt.Sprintf(xlate.Get("Error while removing server: %v"), err))
			} else {
				mebroutines.ShowInfoMessage(xlate.Get("Server was removed succesfully."))
				progress.TranslateAndSetMessage("Server was removed succesfully.")
			}

			enableUI(mainUIStatus)
		}()
	})

	removeButtonCancel.OnClicked(func(*ui.Button) {
		removeWindow.Hide()
		enableUI(mainUIStatus)
	})

	removeWindow.OnClosing(func(*ui.Window) bool {
		removeWindow.Hide()
		enableUI(mainUIStatus)
		return true
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
		createDestroyElements()
		createRemoveElements()

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
		bindOnDestroyServer(mainUIStatus)
		bindOnRemoveServer(mainUIStatus)
		bindOnMebShare()

		bindOnBackup(mainUIStatus)
		bindOnDestroy(mainUIStatus)
		bindOnRemove(mainUIStatus)

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
