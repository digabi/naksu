package main

import (
	"fmt"
	"path/filepath"
	"time"

	"naksu/box"
	"naksu/boxversion"
	"naksu/config"
	"naksu/constants"
	"naksu/host"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/mebroutines/backup"
	"naksu/mebroutines/destroy"
	"naksu/mebroutines/install"
	"naksu/mebroutines/remove"
	"naksu/mebroutines/start"
	"naksu/network"
	"naksu/ui/progress"
	"naksu/ui/networkstatus"
	"naksu/xlate"

	"github.com/andlabs/ui"
	humanize "github.com/dustin/go-humanize"
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
var comboboxExtNic *ui.Combobox
var comboboxNic *ui.Combobox

var labelBox *ui.Label
var labelBoxAvailable *ui.Label
var labelStatus *ui.Label
var labelAdvancedExtNic *ui.Label
var labelAdvancedNic *ui.Label
var labelAdvancedUpdate *ui.Label
var labelAdvancedAnnihilate *ui.Label

var checkboxAdvanced *ui.Checkbox

var boxVersions *ui.Box
var boxBasicUpper *ui.Box
var boxBasic *ui.Box
var boxAdvancedUpdate *ui.Box
var boxAdvancedAnnihilate *ui.Box
var boxAdvanced *ui.Box
var boxStatusBar *ui.Box
var boxUI *ui.Box

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

var extInterfaces []constants.AvailableSelection

func createMainWindowElements() {
	// Define main window
	buttonStartServer = ui.NewButton("Start Exam Server")
	buttonGetServer = ui.NewButton("Abitti Exam")
	buttonSwitchServer = ui.NewButton("Matriculation Exam")
	buttonDestroyServer = ui.NewButton("Remove Exams")
	buttonRemoveServer = ui.NewButton("Remove Server")
	buttonMakeBackup = ui.NewButton("Make Exam Server Backup")
	buttonMebShare = ui.NewButton("Open virtual USB stick (ktp-jako)")

	// Define language setting combobox
	comboboxLang = ui.NewCombobox()
	for _, thisSelection := range constants.AvailableLangs {
		comboboxLang.Append(thisSelection.Legend)
	}
	comboboxLang.SetSelected(constants.GetAvailableSelectionID(config.GetLanguage(), constants.AvailableLangs, 0))

	// Define EXTNIC setting combobox
	comboboxExtNic = ui.NewCombobox()
	for _, thisSelection := range extInterfaces {
		comboboxExtNic.Append(xlate.Get(thisSelection.Legend))
	}
	comboboxExtNic.SetSelected(constants.GetAvailableSelectionID(config.GetExtNic(), extInterfaces, 0))

	// Define NIC setting combobox
	comboboxNic = ui.NewCombobox()
	for _, thisSelection := range constants.AvailableNics {
		comboboxNic.Append(thisSelection.Legend)
	}
	comboboxNic.SetSelected(constants.GetAvailableSelectionID(config.GetNic(), constants.AvailableNics, 0))

	labelBox = ui.NewLabel("")
	labelBoxAvailable = ui.NewLabel("")
	labelStatus = ui.NewLabel("")
	labelAdvancedExtNic = ui.NewLabel("")
	labelAdvancedNic = ui.NewLabel("")
	labelAdvancedUpdate = ui.NewLabel("")
	labelAdvancedAnnihilate = ui.NewLabel("")

	checkboxAdvanced = ui.NewCheckbox("")

	networkStatusArea := networkstatus.Area()

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
	boxAdvanced.Append(ui.NewHorizontalSeparator(), false)
	boxAdvanced.Append(labelAdvancedExtNic, false)
	boxAdvanced.Append(comboboxExtNic, false)
	boxAdvanced.Append(labelAdvancedNic, false)
	boxAdvanced.Append(comboboxNic, false)
	boxAdvanced.Append(ui.NewHorizontalSeparator(), false)
	boxAdvanced.Append(buttonMakeBackup, true)
	boxAdvanced.Append(ui.NewHorizontalSeparator(), false)
	boxAdvanced.Append(labelAdvancedUpdate, false)
	boxAdvanced.Append(boxAdvancedUpdate, true)
	boxAdvanced.Append(labelAdvancedAnnihilate, false)
	boxAdvanced.Append(boxAdvancedAnnihilate, true)

	// For some reason networkStatusArea doesn't get the correct height
	// when placed in its own box. A hacky fix for this issue is using a
	// HorizontalBox instead of a VerticalBox and adding an (almost) empty
	// label to the end of the horizontal layout. This label then effectively
	// sets the height of the HorizontalBox layout. Using a string with some
	// whitespace instead of a completely empty string results in the correct
	// height (reserving enough space for descenders).
	statusBarHeightSetterLabel := ui.NewLabel(" ")

	boxStatusBar = ui.NewHorizontalBox()
	boxStatusBar.SetPadded(true)
	boxStatusBar.Append(networkStatusArea, true)
	boxStatusBar.Append(statusBarHeightSetterLabel, false)

	boxUI = ui.NewVerticalBox()
	boxUI.SetPadded(true)
	boxUI.Append(boxBasic, false)
	boxUI.Append(boxAdvanced, false)
	boxUI.Append(ui.NewHorizontalSeparator(), false)
	boxUI.Append(boxStatusBar, true)

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
				networkstatus.Update()

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
				log.Debug(fmt.Sprintf("main_ui_status: %s", newStatus))
				// Got new status
				if newStatus == "enable" {
					log.Debug("enable ui")

					ui.QueueMain(func() {
						comboboxLang.Enable()
						comboboxNic.Enable()
						comboboxExtNic.Enable()

						// Require installed version to start server
						if box.GetVersion() == "" {
							buttonStartServer.Disable()
						} else {
							updateStartButtonLabel()
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
					log.Debug("disable ui")

					ui.QueueMain(func() {
						comboboxLang.Disable()
						comboboxNic.Disable()
						comboboxExtNic.Disable()

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
func checkAbittiUpdate() (bool, string) {
	abittiUpdate := false
	availAbittiVersion := ""

	currentBoxType := box.GetType()
	currentBoxVersion := box.GetVersion()

	if boxversion.GetVagrantBoxTypeIsAbitti(currentBoxType) {
		_, availBoxVersion, errAvail := boxversion.GetVagrantBoxAvailVersionDetails()
		if errAvail == nil && currentBoxVersion != availBoxVersion {
			abittiUpdate = true
			availAbittiVersion = availBoxVersion
		}
	}

	return abittiUpdate, availAbittiVersion
}

// updateStartButtonLabel updates label for start button depending on the
// installed VM style. If there is no box installed the default label is set.
func updateStartButtonLabel() {
	go func() {
		boxTypeString := boxversion.GetVagrantBoxType(box.GetType())
		ui.QueueMain(func() {
			if boxTypeString == "-" {
				buttonStartServer.SetText(xlate.Get("Start Exam Server"))
			} else {
				buttonStartServer.SetText(fmt.Sprintf(xlate.Get("Start %s"), boxTypeString))
			}
		})
	}()
}

// updateVagrantBoxAvailLabel updates UI "update available" label if the currently
// installed box is Abitti and there is new version available
func updateVagrantBoxAvailLabel() {
	go func() {
		abittiUpdate, _ := checkAbittiUpdate()
		if abittiUpdate {
			vagrantBoxAvailVersion := boxversion.GetVagrantBoxAvailVersion()
			ui.QueueMain(func() {
				labelBoxAvailable.SetText(fmt.Sprintf(xlate.Get("Update available: %s"), vagrantBoxAvailVersion))
				// Select "advanced features" checkbox
				checkboxAdvanced.SetChecked(true)
				boxAdvanced.Show()
			})
		} else {
			ui.QueueMain(func() {
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
		abittiUpdate, availAbittiVersion := checkAbittiUpdate()
		if abittiUpdate {
			ui.QueueMain(func() {
				buttonGetServer.SetText(fmt.Sprintf(xlate.Get("Abitti Exam (%s)"), availAbittiVersion))
			})
		} else {
			ui.QueueMain(func() {
				buttonGetServer.SetText(xlate.Get("Abitti Exam"))
			})
		}
	}()
}

func translateUILabels() {
	ui.QueueMain(func() {
		updateStartButtonLabel()
		updateGetServerButtonLabel()
		buttonSwitchServer.SetText(xlate.Get("Matriculation Exam"))
		buttonDestroyServer.SetText(xlate.Get("Remove Exams"))
		buttonRemoveServer.SetText(xlate.Get("Remove Server"))
		buttonMakeBackup.SetText(xlate.Get("Make Exam Server Backup"))
		buttonMebShare.SetText(xlate.Get("Open virtual USB stick (ktp-jako)"))

		labelBox.SetText(fmt.Sprintf(xlate.Get("Current version: %s"), box.GetVersion()))

		// Show available box version if we have a Abitti box
		updateVagrantBoxAvailLabel()

		// Suggest VM install if none installed
		emptyVersionProgressMessage := "Start by installing a server: Show management features"
		if (progress.GetLastMessage() == "" || progress.GetLastMessage() == emptyVersionProgressMessage) && box.GetVersion() == "" {
			progress.TranslateAndSetMessage(emptyVersionProgressMessage)
		}

		checkboxAdvanced.SetText(xlate.Get("Show management features"))
		labelAdvancedExtNic.SetText(xlate.Get("Network device:"))
		labelAdvancedNic.SetText(xlate.Get("Server networking hardware:"))
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
	networkstatus.Update()
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
		config.SetLanguage(constants.AvailableLangs[comboboxLang.Selected()].ConfigValue)

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

func bindAdvancedExtNicSwitching() {
	// Define EXTNIC selection action main window (advanced view)
	comboboxExtNic.OnSelected(func(*ui.Combobox) {
		config.SetExtNic(extInterfaces[comboboxExtNic.Selected()].ConfigValue)
	})
}

func bindAdvancedNicSwitching() {
	// Define NIC selection action main window (advanced view)
	comboboxNic.OnSelected(func(*ui.Combobox) {
		config.SetNic(constants.AvailableNics[comboboxNic.Selected()].ConfigValue)
	})
}

func bindUIDisableOnStart(mainUIStatus chan string) {
	// Define actions for main window
	buttonStartServer.OnClicked(func(*ui.Button) {
		go func() {
			// Give warnings if there is problems with configured external network device
			// and there are more than one available
			if config.GetExtNic() == "" {
				if len(extInterfaces) > 2 {
					mebroutines.ShowWarningMessage(xlate.Get("You have not set network device. Follow terminal for device selection menu."))
				}
			} else if !network.IsExtInterface(config.GetExtNic()) {
				if len(extInterfaces) > 2 {
					mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("You have selected network device '%s' which is not available. Follow terminal for device selection menu."), config.GetExtNic()))
				}
			}

			// Get defails of the current installed box and warn if we're having Matric Exam box & internet connection
			currentBoxType := box.GetType()
			if boxversion.GetVagrantBoxTypeIsMatriculationExam(currentBoxType) {
				if network.CheckIfNetworkAvailable() {
					mebroutines.ShowWarningMessage(xlate.Get("You are starting Matriculation Examination server with an Internet connection."))
				} else {
					log.Debug("Starting Matric Exam server without an internet connection - All is good!")
				}
			}

			disableUI(mainUIStatus)
			start.Server()
			enableUI(mainUIStatus)
			progress.SetMessage("")
		}()
	})

}

func checkFreeDisk(chFreeDisk chan uint64) {
	// Check free disk
	// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
	go func() {
		var freeDisk uint64
		var err error
		if mebroutines.ExistsDir(mebroutines.GetVagrantDirectory()) {
			freeDisk, err = mebroutines.GetDiskFree(mebroutines.GetVagrantDirectory())
			if err != nil {
				log.Debug("Getting free disk space from Vagrant directory failed")
				mebroutines.ShowWarningMessage("Failed to calculate free disk space of the Vagrant directory")
			}
		} else {
			freeDisk, err = mebroutines.GetDiskFree(mebroutines.GetHomeDirectory())
			if err != nil {
				log.Debug("Getting free disk space from home directory failed")
				mebroutines.ShowWarningMessage("Failed to calculate free disk space of the home directory")
			}
		}
		chFreeDisk <- freeDisk
	}()
}

func bindOnGetServer(mainUIStatus chan string) {
	buttonGetServer.OnClicked(func(*ui.Button) {
		log.Debug("Starting Abitti box update")

		chFreeDisk := make(chan uint64)
		chDiskLowPopup := make(chan bool)

		checkFreeDisk(chFreeDisk)

		go func() {
			freeDisk := <-chFreeDisk
			if freeDisk < constants.LowDiskLimit {
				mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Your free disk size is getting low (%s)."), humanize.Bytes(freeDisk)))
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

				log.Debug(fmt.Sprintf("Finished Abitti box update, version is: %s", box.GetVersion()))
			}()
		}()
	})
}

func bindOnSwitchServer(mainUIStatus chan string) {
	buttonSwitchServer.OnClicked(func(*ui.Button) {
		log.Debug("Starting Matriculation Examination box update")

		chFreeDisk := make(chan uint64)
		chDiskLowPopup := make(chan bool)
		chPathNewVagrantfile := make(chan string)

		checkFreeDisk(chFreeDisk)

		go func() {
			freeDisk := <-chFreeDisk
			if freeDisk < constants.LowDiskLimit {
				mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Your free disk size is getting low (%s)."), humanize.Bytes(freeDisk)))
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
				mebroutines.ShowWarningMessage(xlate.Get("Did not get a path for a new Vagrantfile"))
			} else if pathNewVagrantfile == pathOldVagrantfile {
				mebroutines.ShowWarningMessage(xlate.Get("Please place the new Exam Vagrantfile to another location (e.g. desktop or home directory)"))
			} else {
				go func() {
					disableUI(mainUIStatus)
					install.GetServer(pathNewVagrantfile)
					translateUILabels()
					enableUI(mainUIStatus)
					progress.SetMessage("")

					log.Debug(fmt.Sprintf("Finished Matriculation Examination box update, new version is: %s", box.GetVersion()))
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
		log.Debug("Opening MEB share (~/ktp-jako)")
		mebroutines.OpenMebShare()
	})
}

func bindOnBackup(mainUIStatus chan string) {
	// Define actions for SaveAs window/dialog
	backupButtonSave.OnClicked(func(*ui.Button) {
		pathBackup := filepath.Join(backupMediaPath[backupCombobox.Selected()], backup.GetBackupFilename(time.Now()))
		log.Debug(fmt.Sprintf("Starting backup to: %s", pathBackup))

		chFreeDisk := make(chan uint64)
		chDiskLowPopup := make(chan bool)

		checkFreeDisk(chFreeDisk)

		go func() {
			freeDisk := <-chFreeDisk
			if freeDisk < constants.LowDiskLimit {
				mebroutines.ShowWarningMessage("Your free disk size is getting low. If backup process fails please consider freeing some disk space.")
			}
			chDiskLowPopup <- true
		}()

		go func() {
			<-chDiskLowPopup

			go func() {
				backupWindow.Hide()
				err := backup.MakeBackup(pathBackup)
				if err != nil {
					mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Backup failed: %v"), err))
					progress.SetMessage(fmt.Sprintf(xlate.Get("Backup failed: %v"), err))
				} else {
					progress.SetMessage(fmt.Sprintf(xlate.Get("Backup done: %s"), pathBackup))
				}

				enableUI(mainUIStatus)

				log.Debug("Finished creating backup")
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
			log.Debug("Starting server destroy")

			destroyWindow.Hide()
			err := destroy.Server()
			if err != nil {
				mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Failed to remove exams: %v"), err))
				progress.SetMessage(fmt.Sprintf(xlate.Get("Failed to remove exams: %v"), err))
			} else {
				progress.TranslateAndSetMessage("Exams were removed successfully.")
			}

			// Update installed version label
			translateUILabels()

			enableUI(mainUIStatus)

			log.Debug("Finished server destroy")
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
			log.Debug("Starting server remove")

			removeWindow.Hide()

			err := remove.Server()
			if err != nil {
				mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Error while removing server: %v"), err))
				progress.SetMessage(fmt.Sprintf(xlate.Get("Error while removing server: %v"), err))
			} else {
				progress.TranslateAndSetMessage("Server was removed succesfully.")
			}

			// Update installed version label
			translateUILabels()

			enableUI(mainUIStatus)

			log.Debug("Finished server remove")
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

	// Same applies to Windows network interface query
	extInterfaces = network.GetExtInterfaces()

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

		networkstatus.Update()

		setupMainLoop(mainUIStatus, mainUINetupdate)

		enableUI(mainUIStatus)

		window.SetMargined(true)
		window.SetChild(boxUI)

		// Advanced group is hidden by default
		boxAdvanced.Hide()

		// Set UI labels with default language
		translateUILabels()

		bindLanguageSwitching()
		bindAdvancedToggle()
		bindAdvancedExtNicSwitching()
		bindAdvancedNicSwitching()

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
			log.Debug("User exits through window exit")
			ui.Quit()
			return true
		})

		window.Show()
		mainUIStatus <- "enable"
		backupWindow.Hide()

		// Make sure we have vagrant
		if !mebroutines.IfFoundVagrant() {
			mebroutines.ShowErrorMessage(xlate.Get("Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?"))
			log.Debug("Exiting as Vagrant is missing")
			ui.Quit()
		}

		// Make sure we have VBoxManage
		if !mebroutines.IfFoundVBoxManage() {
			mebroutines.ShowErrorMessage(xlate.Get("Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?"))
			log.Debug("Exiting as VBoxManage is missing")
			ui.Quit()
		}

		// Make sure Hyper-V is not running
		// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
		isHyperV := make(chan bool)
		go func() {
			// IsHyperV() uses Windows WMI call
			isHyperV <- host.IsHyperV()
		}()

		if <-isHyperV {
			mebroutines.ShowWarningMessage(xlate.Get("Please turn Windows Hypervisor off as it may cause problems."))
		} else {
			// Does CPU support hardware virtualisation?
			if !host.IsHWVirtualisationCPU() {
				mebroutines.ShowWarningMessage(xlate.Get("It appears your CPU does not support hardware virtualisation (VT-x or AMD-V)."))
			}

			// Make sure the hardware virtualisation is present
			if !host.IsHWVirtualisation() {
				mebroutines.ShowWarningMessage(xlate.Get("Hardware virtualisation (VT-x or AMD-V) is disabled. Please enable it before continuing."))
			}
		}

		// Check if home directory contains non-american characters which may cause problems to vagrant
		if mebroutines.IfIntlCharsInPath(mebroutines.GetHomeDirectory()) {
			mebroutines.ShowWarningMessage(fmt.Sprintf(xlate.Get("Your home directory path (%s) contains characters which may cause problems to Vagrant."), mebroutines.GetHomeDirectory()))
		}

		log.Debug(fmt.Sprintf("Currently installed box: %s %s", box.GetVersion(), box.GetType()))
	})
}
