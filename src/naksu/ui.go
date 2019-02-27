package main

import (
	"fmt"
	"path/filepath"
	"naksu/constants"
	"naksu/config"
	"naksu/mebroutines"
	"naksu/mebroutines/backup"
	"naksu/mebroutines/destroy"
	"naksu/mebroutines/install"
	"naksu/mebroutines/remove"
	"naksu/mebroutines/start"
	"naksu/network"
	"naksu/boxversion"
	"naksu/progress"
	"naksu/xlate"
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
var labelBoxAvailable *ui.Label
var labelStatus *ui.Label
var labelAdvancedUpdate *ui.Label
var labelAdvancedAnnihilate *ui.Label

var checkboxAdvanced *ui.Checkbox

var boxVersions *ui.Box
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

var destroyInfoLabel [5]*ui.Label

// Remove Confirmation Window
var removeWindow *ui.Window

var removeButtonRemove *ui.Button
var removeButtonCancel *ui.Button

var removeBox *ui.Box

var removeInfoLabel [5]*ui.Label

func createMainWindowElements() {
	// Define main window
	buttonStartServer = ui.NewButton("Start Exam Server")
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
	switch config.GetLanguage() {
	case "fi":
		comboboxLang.SetSelected(0)
	case "sv":
		comboboxLang.SetSelected(1)
	case "en":
		comboboxLang.SetSelected(2)
	default:
		comboboxLang.SetSelected(0)
	}

	labelBox = ui.NewLabel("")
	labelBoxAvailable = ui.NewLabel("")
	labelStatus = ui.NewLabel("")
	labelAdvancedUpdate = ui.NewLabel("")
	labelAdvancedAnnihilate = ui.NewLabel("")

	checkboxAdvanced = ui.NewCheckbox("")

	// Box versions
	boxVersions = ui.NewVerticalBox()
	boxVersions.SetPadded(true)
	boxVersions.Append(labelBox, true)
	boxVersions.Append(labelBoxAvailable, true)

	// Box version and language selection dropdown
	boxBasicUpper = ui.NewHorizontalBox()
	boxBasicUpper.SetPadded(true)
	boxBasicUpper.Append(boxVersions, true)
	boxBasicUpper.Append(comboboxLang, false)

	boxBasic = ui.NewVerticalBox()
	boxBasic.SetPadded(true)
	boxBasic.Append(boxBasicUpper, false)
	boxBasic.Append(labelStatus, true)
	boxBasic.Append(buttonStartServer, false)
	boxBasic.Append(buttonMebShare, false)
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
	for i := 0; i <= 4; i++ {
		destroyInfoLabel[i] = ui.NewLabel("destroyInfoLabel")
	}

	destroyButtonDestroy = ui.NewButton("Yes, Remove")
	destroyButtonCancel = ui.NewButton("Cancel")

	destroyBox = ui.NewVerticalBox()
	destroyBox.SetPadded(true)
	for i := 0; i <= 4; i++ {
		destroyBox.Append(destroyInfoLabel[i], false)
	}
	destroyBox.Append(destroyButtonDestroy, false)
	destroyBox.Append(destroyButtonCancel, false)

	destroyWindow = ui.NewWindow("", 1, 1, false)

	destroyWindow.SetMargined(true)
	destroyWindow.SetChild(destroyBox)
}

func createRemoveElements() {
	// Define Destroy Confirmation window/dialog
	for i := 0; i <= 4; i++ {
		removeInfoLabel[i] = ui.NewLabel("removeInfoLabel")
	}

	removeButtonRemove = ui.NewButton("Yes, Remove")
	removeButtonCancel = ui.NewButton("Cancel")

	removeBox = ui.NewVerticalBox()
	removeBox.SetPadded(true)
	for i := 0; i <= 4; i++ {
		removeBox.Append(removeInfoLabel[i], false)
	}
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
						if network.CheckIfNetworkAvailable() {
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

						// Require installed version to start server
						if boxversion.GetVagrantFileVersion("") == "" {
							buttonStartServer.Disable()
						} else {
							buttonStartServer.Enable()
						}

						buttonMebShare.Enable()

						// Require network connection for install/update
						if network.CheckIfNetworkAvailable() {
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

// checkAbittiUpdate checks
// 1) if currently installed box is Abitti
// 2) and there is a new version available
func checkAbittiUpdate() (bool, string, string) {
	abittiUpdate := false
	currentAbittiVersion := ""
	availAbittiVersion := ""

	currentBoxType, currentBoxVersion, errCurrent := boxversion.GetVagrantFileVersionDetails(filepath.Join(mebroutines.GetVagrantDirectory(), "Vagrantfile"))
	if (errCurrent == nil && boxversion.GetVagrantBoxTypeIsAbitti(currentBoxType)) {
		currentAbittiVersion = currentBoxVersion
		_, availBoxVersion, errAvail := boxversion.GetVagrantBoxAvailVersionDetails()
		if (errAvail == nil && currentBoxVersion != availBoxVersion) {
			abittiUpdate = true
			availAbittiVersion = availBoxVersion
		}
	}

	return abittiUpdate, currentAbittiVersion, availAbittiVersion
}

// updateVagrantBoxAvailLabel updates UI "update available" label if the currently
// installed box is Abitti and there is new version available
func updateVagrantBoxAvailLabel () {
	go func() {
		abittiUpdate, _, _ := checkAbittiUpdate()
		if abittiUpdate {
			vagrantBoxAvailVersion := boxversion.GetVagrantBoxAvailVersion()
			ui.QueueMain(func () {
				labelBoxAvailable.SetText(fmt.Sprintf(xlate.Get("Update available: %s"), vagrantBoxAvailVersion))
				// Select "advanced features" checkbox
				checkboxAdvanced.SetChecked(true)
				boxAdvanced.Show()
			})
		} else {
			ui.QueueMain(func () {
				labelBoxAvailable.SetText("")
			})
		}
	}()
}

// updateGetServerButtonLabel updates UI "Abitti update" button label
// If there is new version available it shows current and new version numbers
// Make sure you call this inside ui.Queuemain() only
func updateGetServerButtonLabel() {
	// Set default text for an empty button before doing time-consuming Abitti update check
	if buttonGetServer.Text() == "" {
		buttonGetServer.SetText(xlate.Get("Abitti Exam"))
	}

	go func() {
		abittiUpdate, currentAbittiVersion, availAbittiVersion := checkAbittiUpdate()
		if abittiUpdate {
			ui.QueueMain(func () {
				buttonGetServer.SetText(fmt.Sprintf(xlate.Get("Abitti Exam (v%s > v%s)"), currentAbittiVersion, availAbittiVersion))
			})
		} else {
			ui.QueueMain(func () {
				buttonGetServer.SetText(xlate.Get("Abitti Exam"))
			})
		}
	}()
}

func translateUILabels() {
	ui.QueueMain(func() {
		buttonStartServer.SetText(xlate.Get("Start Exam Server"))
		updateGetServerButtonLabel()
		buttonSwitchServer.SetText(xlate.Get("Matriculation Exam"))
		buttonDestroyServer.SetText(xlate.Get("Remove Exams"))
		buttonRemoveServer.SetText(xlate.Get("Remove Server"))
		buttonMakeBackup.SetText(xlate.Get("Make Exam Server Backup"))
		buttonMebShare.SetText(xlate.Get("Open virtual USB stick (ktp-jako)"))

		labelBox.SetText(fmt.Sprintf(xlate.Get("Current version: %s"), boxversion.GetVagrantFileVersion("")))

		// Show available box version if we have a Abitti box
		updateVagrantBoxAvailLabel()

		// Suggest VM install if none installed
		emptyVersionProgressMessage := "Start by installing a server: Show management features"
		if (progress.GetLastMessage() == "" || progress.GetLastMessage() == emptyVersionProgressMessage) && boxversion.GetVagrantFileVersion("") == "" {
			progress.TranslateAndSetMessage(emptyVersionProgressMessage)
		}

		checkboxAdvanced.SetText(xlate.Get("Show management features"))
		labelAdvancedUpdate.SetText(xlate.Get("Install/update server for:"))
		labelAdvancedAnnihilate.SetText(xlate.Get("DANGER! Annihilate your server:"))

		backupWindow.SetTitle(xlate.Get("naksu: SaveTo"))
		backupLabel.SetText(xlate.Get("Please select target path"))
		backupButtonSave.SetText(xlate.Get("Save"))
		backupButtonCancel.SetText(xlate.Get("Cancel"))

		destroyWindow.SetTitle(xlate.Get("naksu: Remove Exams"))
		destroyInfoLabel[0].SetText(xlate.Get("Remove Exams restores server to its initial status."))
		destroyInfoLabel[1].SetText(xlate.Get("Exams, responses and logs in the server will be irreversibly deleted."))
		destroyInfoLabel[2].SetText(xlate.Get("It is recommended to back up your server before removing exams."))
		destroyInfoLabel[3].SetText(xlate.Get(""))
		destroyInfoLabel[4].SetText(xlate.Get("Do you wish to remove all exams?"))
		destroyButtonDestroy.SetText(xlate.Get("Yes, Remove"))
		destroyButtonCancel.SetText(xlate.Get("Cancel"))

		removeWindow.SetTitle(xlate.Get("naksu: Remove Server"))
		removeInfoLabel[0].SetText(xlate.Get("Removing server destroys it and all downloaded disk images."))
		removeInfoLabel[1].SetText(xlate.Get("Exams, responses and logs in the server will be irreversibly deleted."))
		removeInfoLabel[2].SetText(xlate.Get("It is recommended to back up your server before removing server."))
		removeInfoLabel[3].SetText(xlate.Get(""))
		removeInfoLabel[4].SetText(xlate.Get("Do you wish to remove the server?"))
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
			config.SetLanguage("fi")
		case 1:
			config.SetLanguage("sv")
		case 2:
			config.SetLanguage("en")
		default:
			config.SetLanguage("fi")
		}

		xlate.SetLanguage(config.GetLanguage())
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
			// Get defails of the current installed box and warn if we're having Matric Exam box & internet connection
			boxVersionString, _, boxErr := boxversion.GetVagrantFileVersionDetails(filepath.Join(mebroutines.GetVagrantDirectory(), "Vagrantfile"))
			if (boxErr == nil && boxversion.GetVagrantBoxTypeIsMatriculationExam(boxVersionString)) {
				if network.CheckIfNetworkAvailable() {
					mebroutines.ShowWarningMessage(xlate.Get("You are starting Matriculation Examination server with an Internet connection."))
				} else {
					mebroutines.LogDebug("Starting Matric Exam server without an internet connection - All is good!")
				}
			}

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
		mebroutines.LogDebug("Starting Abitti box update")

		chFreeDisk := make(chan int)
		chDiskLowPopup := make(chan bool)

		checkFreeDisk(chFreeDisk)

		go func() {
			freeDisk := <-chFreeDisk
			if freeDisk != -1 && freeDisk < constants.LowDiskLimit {
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

				mebroutines.LogDebug(fmt.Sprintf("Finished Abitti box update, version is: %s", boxversion.GetVagrantFileVersion("")))
			}()
		}()
	})
}

func bindOnSwitchServer(mainUIStatus chan string) {
	buttonSwitchServer.OnClicked(func(*ui.Button) {
		mebroutines.LogDebug("Starting Matriculation Examination box update")

		chFreeDisk := make(chan int)
		chDiskLowPopup := make(chan bool)
		chPathNewVagrantfile := make(chan string)

		checkFreeDisk(chFreeDisk)

		go func() {
			freeDisk := <-chFreeDisk
			if freeDisk != -1 && freeDisk < constants.LowDiskLimit {
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
			pathOldVagrantfile := filepath.Join(mebroutines.GetVagrantDirectory(), "Vagrantfile")

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

					mebroutines.LogDebug(fmt.Sprintf("Finished Matriculation Examination box update, new version is: %s", boxversion.GetVagrantFileVersion("")))
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
		pathBackup := filepath.Join(backupMediaPath[backupCombobox.Selected()], backup.GetBackupFilename(time.Now()))
		mebroutines.LogDebug(fmt.Sprintf("Starting backup to: %s", pathBackup))

		chFreeDisk := make(chan int)
		chDiskLowPopup := make(chan bool)

		checkFreeDisk(chFreeDisk)

		go func() {
			freeDisk := <-chFreeDisk
			if freeDisk != -1 && freeDisk < constants.LowDiskLimit {
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

				mebroutines.LogDebug("Finished creating backup")
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
		go func() {
			mebroutines.LogDebug("Starting server destroy")

			destroyWindow.Hide()
			destroy.Server()
			progress.SetMessage("")

			// Update installed version label
			translateUILabels()

			enableUI(mainUIStatus)

			mebroutines.LogDebug("Finished server destroy")
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
		go func() {
			mebroutines.LogDebug("Starting server remove")

			removeWindow.Hide()

			err := remove.Server()
			if err != nil {
				mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("Error while removing server: %v"), err))
				progress.SetMessage(fmt.Sprintf(xlate.Get("Error while removing server: %v"), err))
			} else {
				mebroutines.ShowInfoMessage(xlate.Get("Server was removed succesfully."))
				progress.TranslateAndSetMessage("Server was removed succesfully.")
			}

			// Update installed version label
			translateUILabels()

			enableUI(mainUIStatus)

			mebroutines.LogDebug("Finished server remove")
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
