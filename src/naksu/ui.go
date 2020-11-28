package main

import (
	"fmt"
	"path/filepath"
	"time"

	"naksu/box"
	"naksu/box/download"
	"naksu/box/vboxmanage"
	"naksu/config"
	"naksu/constants"
	"naksu/host"
	"naksu/log"
	"naksu/logdelivery"
	"naksu/mebroutines"
	"naksu/mebroutines/backup"
	"naksu/mebroutines/destroy"
	"naksu/mebroutines/install"
	"naksu/mebroutines/remove"
	"naksu/mebroutines/start"
	"naksu/network"
	"naksu/ui/networkstatus"
	"naksu/ui/progress"
	"naksu/xlate"

	"github.com/andlabs/ui"
	"github.com/atotto/clipboard"
)

type mainUIStatusType = string

const mainUIStatusEnabled mainUIStatusType = ""
const mainUIStatusDisabled mainUIStatusType = "disable"

var window *ui.Window

var buttonSelfUpdateOn *ui.Button
var buttonStartServer *ui.Button
var buttonInstallAbittiServer *ui.Button
var buttonInstallExamServer *ui.Button
var buttonDestroyServer *ui.Button
var buttonRemoveServer *ui.Button
var buttonMakeBackup *ui.Button
var buttonDeliverLogs *ui.Button
var buttonMebShare *ui.Button

var comboboxLang *ui.Combobox
var comboboxExtNic *ui.Combobox
var comboboxNic *ui.Combobox

var labelBox *ui.Label
var labelBoxAvailable *ui.Label
var labelStatus *ui.Label
var labelExtNic *ui.Label
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

// Log Delivery Window
var logDeliveryWindow *ui.Window

var logDeliveryBox *ui.Box
var logDeliveryFilenameBox *ui.Box

var logDeliveryFilenameLabelLabel *ui.Label
var logDeliveryFilenameLabel *ui.Label
var logDeliveryFilenameCopyButton *ui.Button
var logDeliveryStatusLabel *ui.Label
var logDeliveryButtonClose *ui.Button

// Exam Install Window
var examInstallWindow *ui.Window

var examInstallBox *ui.Box
var examInstallPassphraseLabel *ui.Label
var examInstallPassphraseEntry *ui.Entry
var examInstallButtonInstall *ui.Button
var examInstallButtonCancel *ui.Button

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
	buttonSelfUpdateOn = ui.NewButton("Turn Naksu self updates back on")
	buttonStartServer = ui.NewButton("Start Exam Server")
	buttonInstallAbittiServer = ui.NewButton("Abitti Exam")
	buttonInstallExamServer = ui.NewButton("Matriculation Exam")
	buttonDestroyServer = ui.NewButton("Remove Exams")
	buttonRemoveServer = ui.NewButton("Remove Server")
	buttonMakeBackup = ui.NewButton("Make Exam Server Backup")
	buttonDeliverLogs = ui.NewButton("Send logs to Abitti support")
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
	labelExtNic = ui.NewLabel("")
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
	boxBasic.Append(buttonSelfUpdateOn, false)
	boxBasic.Append(buttonStartServer, false)
	boxBasic.Append(buttonMebShare, false)
	boxBasic.Append(labelExtNic, false)
	boxBasic.Append(comboboxExtNic, false)
	boxBasic.Append(checkboxAdvanced, true)

	boxAdvancedUpdate = ui.NewHorizontalBox()
	boxAdvancedUpdate.SetPadded(true)
	boxAdvancedUpdate.Append(buttonInstallAbittiServer, true)
	boxAdvancedUpdate.Append(buttonInstallExamServer, true)

	boxAdvancedAnnihilate = ui.NewHorizontalBox()
	boxAdvancedAnnihilate.SetPadded(true)
	boxAdvancedAnnihilate.Append(buttonDestroyServer, true)
	boxAdvancedAnnihilate.Append(buttonRemoveServer, true)

	boxAdvanced = ui.NewVerticalBox()
	boxAdvanced.SetPadded(true)
	boxAdvanced.Append(ui.NewHorizontalSeparator(), false)
	boxAdvanced.Append(labelAdvancedNic, false)
	boxAdvanced.Append(comboboxNic, false)
	boxAdvanced.Append(ui.NewHorizontalSeparator(), false)
	boxAdvanced.Append(buttonMakeBackup, true)
	boxAdvanced.Append(buttonDeliverLogs, true)
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

func createLogDeliveryElements() {
	// Define Backup SaveAs window/dialog
	logDeliveryFilenameLabelLabel = ui.NewLabel(xlate.Get("Filename for Abitti support:"))
	logDeliveryFilenameLabel = ui.NewLabel(xlate.Get("Wait..."))
	logDeliveryFilenameCopyButton = ui.NewButton(xlate.Get("Copy to clipboard"))
	logDeliveryStatusLabel = ui.NewLabel(fmt.Sprintf(xlate.Get("Copying logs: %s"), xlate.Get("0 % (this can take a while...)")))
	logDeliveryButtonClose = ui.NewButton(xlate.Get("Close"))

	logDeliveryBox = ui.NewVerticalBox()
	logDeliveryBox.SetPadded(true)

	logDeliveryBox.Append(logDeliveryFilenameLabelLabel, false)

	logDeliveryFilenameBox = ui.NewHorizontalBox()
	logDeliveryFilenameBox.Append(logDeliveryFilenameLabel, true)

	if clipboard.Unsupported {
		log.Debug("Not adding logDeliveryFilenameCopyButton as clipboard is unsupported. If you're on Linux, install xsel or xclip.")
	} else {
		logDeliveryFilenameBox.Append(logDeliveryFilenameCopyButton, true)
	}

	logDeliveryFilenameBox.SetPadded(true)
	logDeliveryBox.Append(logDeliveryFilenameBox, true)

	logDeliveryBox.Append(ui.NewHorizontalSeparator(), false)
	logDeliveryBox.Append(logDeliveryStatusLabel, false)
	logDeliveryBox.Append(logDeliveryButtonClose, false)

	logDeliveryWindow = ui.NewWindow("", 400, 1, false)
	logDeliveryWindow.SetMargined(true)
	logDeliveryWindow.SetChild(logDeliveryBox)
}

func createExamInstallElements() {
	// Define exam install passhrase dialog window
	examInstallPassphraseLabel = ui.NewLabel(xlate.Get("Give passphrase for exam server:"))
	examInstallPassphraseEntry = ui.NewEntry()
	examInstallButtonCancel = ui.NewButton(xlate.Get("Cancel"))
	examInstallButtonInstall = ui.NewButton(xlate.Get("Install"))

	examInstallBox = ui.NewVerticalBox()
	examInstallBox.SetPadded(true)

	examInstallBox.Append(examInstallPassphraseLabel, false)
	examInstallBox.Append(examInstallPassphraseEntry, false)
	examInstallBox.Append(examInstallButtonInstall, false)
	examInstallBox.Append(examInstallButtonCancel, false)

	examInstallWindow = ui.NewWindow("", 400, 1, false)
	examInstallWindow.SetMargined(true)
	examInstallWindow.SetChild(examInstallBox)
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

func setupMainLoop(mainUIStatus chan string, mainUIUpdate *time.Ticker) {
	go func() {

		// Start with all functions disabled
		currentMainUIStatus := mainUIStatusDisabled

		for {
			select {
			case <-mainUIUpdate.C:
				mainUIStatusHandler(currentMainUIStatus)
			case newStatus := <-mainUIStatus:
				currentMainUIStatus = newStatus
				mainUIStatusHandler(currentMainUIStatus)
			}
		}
	}()
}

// Define UI button status handler
func mainUIStatusHandler(currentMainUIStatus mainUIStatusType) { //nolint:gocyclo
	// The cyclomatic complexity calculated by gocyclo is 20

	networkstatus.Update()

	// Check general UI status
	mainUIEnabled := (currentMainUIStatus == mainUIStatusEnabled)

	boxInstalled := false
	boxRunning := false
	netAvailable := false

	// Make these checks only if main UI is enable to save time when disabling the UI
	if mainUIEnabled {
		boxInstalled, boxRunning, netAvailable = getStatusForMainUIStatusHandler()
	}

	// Create rule for "enabled" for each button
	buttonRules := []struct {
		element *ui.Button
		enable  bool
	}{
		{buttonSelfUpdateOn, config.IsSelfUpdateDisabled()},
		{buttonStartServer, mainUIEnabled && boxInstalled && !boxRunning},
		{buttonMebShare, true},
		{buttonMakeBackup, mainUIEnabled && boxInstalled && !boxRunning},
		{buttonDeliverLogs, mainUIEnabled && true},
		{buttonInstallAbittiServer, mainUIEnabled && !boxRunning && netAvailable},
		{buttonInstallExamServer, mainUIEnabled && !boxRunning && netAvailable},
		{buttonDestroyServer, mainUIEnabled && boxInstalled && !boxRunning},
		{buttonRemoveServer, true},
	}

	comboboxRules := []struct {
		element *ui.Combobox
		enable  bool
	}{
		{comboboxLang, mainUIEnabled && !boxRunning},
		{comboboxNic, mainUIEnabled && !boxRunning},
		{comboboxExtNic, mainUIEnabled && !boxRunning},
	}

	ui.QueueMain(func() {
		for _, buttonRule := range buttonRules {
			if buttonRule.enable {
				buttonRule.element.Enable()
			} else {
				buttonRule.element.Disable()
			}
		}

		for _, comboboxRule := range comboboxRules {
			if comboboxRule.enable {
				comboboxRule.element.Enable()
			} else {
				comboboxRule.element.Disable()
			}
		}
	})
}

func getStatusForMainUIStatusHandler() (bool, bool, bool) {
	// Query box installation status
	boxInstalled, boxInstalledErr := box.Installed()
	if boxInstalledErr != nil {
		log.Debug(fmt.Sprintf("Could not query whether VM is installed: %v", boxInstalledErr))
	}

	boxInstalled = (boxInstalledErr == nil) && boxInstalled

	// Query box running status
	boxRunning, boxRunningErr := box.Running()
	if boxRunningErr != nil {
		log.Debug(fmt.Sprintf("Could not query whether VM is running: %v", boxRunningErr))
	}

	boxRunning = (boxRunningErr == nil) && boxRunning

	// Query network (internet) connection status
	netAvailable := network.CheckIfNetworkAvailable()

	return boxInstalled, boxRunning, netAvailable
}

// checkAbittiUpdate checks
// 1) if currently installed box is Abitti or if no box is installed
// 2) and there is a new version available
func checkAbittiUpdate() (bool, string) {
	availAbittiVersion := ""

	currentBoxVersion := box.GetVersion()

	boxInstalled, err := box.Installed()

	if err != nil {
		log.Debug(fmt.Sprintf("Could not detect whether VM is installed: %v", err))
	}

	if (err == nil && !boxInstalled) || box.TypeIsAbitti() {
		availAbittiVersion, err = download.GetAvailableVersion(constants.AbittiVersionURL)
		if err == nil && currentBoxVersion != availAbittiVersion {
			return true, availAbittiVersion
		}
	}

	return false, availAbittiVersion
}

// updateStartButtonLabel updates label for start button depending on the
// installed VM style. If there is no box installed the default label is set.
func updateStartButtonLabel() {
	go func() {
		boxTypeString := box.GetTypeLegend()
		ui.QueueMain(func() {
			if boxTypeString == "-" {
				buttonStartServer.SetText(xlate.Get("Start Exam Server"))
			} else {
				buttonStartServer.SetText(fmt.Sprintf(xlate.Get("Start %s"), boxTypeString))
			}
		})
	}()
}

// updateBoxAvailabilityLabel updates UI "update available" label if the currently
// installed box is Abitti and there is new version available
func updateBoxAvailabilityLabel() {
	go func() {
		abittiUpdate, availAbittiVersion := checkAbittiUpdate()
		if abittiUpdate {
			ui.QueueMain(func() {
				labelBoxAvailable.SetText(fmt.Sprintf(xlate.Get("Update available: %s"), availAbittiVersion))
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
	if buttonInstallAbittiServer.Text() == "" {
		buttonInstallAbittiServer.SetText(xlate.Get("Abitti Exam"))
	}

	go func() {
		abittiUpdate, availAbittiVersion := checkAbittiUpdate()
		if abittiUpdate {
			ui.QueueMain(func() {
				buttonInstallAbittiServer.SetText(fmt.Sprintf(xlate.Get("Abitti Exam (%s)"), availAbittiVersion))
			})
		} else {
			ui.QueueMain(func() {
				buttonInstallAbittiServer.SetText(xlate.Get("Abitti Exam"))
			})
		}
	}()
}

func translateUILabels() {
	ui.QueueMain(func() {
		updateStartButtonLabel()
		updateGetServerButtonLabel()
		buttonInstallExamServer.SetText(xlate.Get("Matriculation Exam"))
		buttonDestroyServer.SetText(xlate.Get("Remove Exams"))
		buttonRemoveServer.SetText(xlate.Get("Remove Server"))
		buttonMakeBackup.SetText(xlate.Get("Make Exam Server Backup"))
		buttonDeliverLogs.SetText(xlate.Get("Send logs to Abitti support"))
		buttonMebShare.SetText(xlate.Get("Open virtual USB stick (ktp-jako)"))
		labelExtNic.SetText(xlate.Get("Network device:"))

		labelBox.SetText(fmt.Sprintf(xlate.Get("Current version: %s"), box.GetVersion()))

		// Show available box version if we have a Abitti box
		updateBoxAvailabilityLabel()

		// Suggest VM install if none installed
		if progress.GetLastMessage() == "" && box.GetVersion() == "" {
			progress.TranslateAndSetMessage("Start by installing a server: Show management features")
		}

		checkboxAdvanced.SetText(xlate.Get("Show management features"))
		labelAdvancedNic.SetText(xlate.Get("Server networking hardware:"))
		labelAdvancedUpdate.SetText(xlate.Get("Install/update server for:"))
		labelAdvancedAnnihilate.SetText(xlate.Get("DANGER! Annihilate your server:"))

		backupWindow.SetTitle(xlate.Get("naksu: SaveTo"))
		backupLabel.SetText(xlate.Get("Please select target path"))
		backupButtonSave.SetText(xlate.Get("Save"))
		backupButtonCancel.SetText(xlate.Get("Cancel"))

		logDeliveryWindow.SetTitle(xlate.Get("naksu: Send Logs"))
		logDeliveryFilenameLabelLabel.SetText(xlate.Get("Filename for Abitti support:"))
		logDeliveryFilenameCopyButton.SetText(xlate.Get("Copy to clipboard"))
		logDeliveryButtonClose.SetText(xlate.Get("Close"))

		examInstallWindow.SetTitle(xlate.Get("naksu: Install Exam Server"))
		examInstallPassphraseLabel.SetText(xlate.Get("Enter Exam Server passphrase:"))
		examInstallButtonInstall.SetText(xlate.Get("Install"))
		examInstallButtonCancel.SetText(xlate.Get("Cancel"))

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

// disableUI sends
func disableUI(mainUIStatus chan string) {
	mainUIStatus <- mainUIStatusDisabled
}

func enableUI(mainUIStatus chan string) {
	mainUIStatus <- mainUIStatusEnabled
}

func bindLanguageSwitching() {
	// Define language selection action main window
	comboboxLang.OnSelected(func(*ui.Combobox) {
		config.SetLanguage(constants.AvailableLangs[comboboxLang.Selected()].ConfigValue)

		xlate.SetLanguage(config.GetLanguage())
		progress.SetMessage("")
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

	if config.IsSelfUpdateDisabled() {
		buttonSelfUpdateOn.Show()
	} else {
		buttonSelfUpdateOn.Hide()
	}

	buttonSelfUpdateOn.OnClicked(func(*ui.Button) {
		config.SetSelfUpdateDisabled(false)
	})

	buttonStartServer.OnClicked(func(*ui.Button) {
		go func() {
			// Give warnings if there is problems with configured external network device
			// and there are more than one available
			if config.GetExtNic() == "" {
				mebroutines.ShowErrorMessage(xlate.Get("Please select the network device which is connected to your exam network."))
				return
			}

			if !network.IsExtInterface(config.GetExtNic()) {
				mebroutines.ShowErrorMessage(fmt.Sprintf(xlate.Get("You have selected network device '%s' which is not available."), config.GetExtNic()))
				return
			}

			// Get defails of the current installed box and warn if we're having Matric Exam box & internet connection
			if box.TypeIsMatriculationExam() {
				if network.CheckIfNetworkAvailable() {
					mebroutines.ShowWarningMessage(xlate.Get("You are starting Matriculation Examination server with an Internet connection."))
				} else {
					log.Debug("Starting Matric Exam server without an internet connection - All is good!")
				}
			}

			// Disable UI to prevent multiple simultaneous server starts
			disableUI(mainUIStatus)
			start.Server()
			// Wait over one UI loop
			time.Sleep(5000)
			enableUI(mainUIStatus)

			progress.SetMessage("")
		}()
	})

}

func bindOnInstallAbittiServer(mainUIStatus chan string) {
	buttonInstallAbittiServer.OnClicked(func(*ui.Button) {
		go func() {
			log.Debug("Starting Abitti box update")

			disableUI(mainUIStatus)
			install.NewAbittiServer()
			translateUILabels()
			enableUI(mainUIStatus)
			progress.SetMessage("")

			log.Debug(fmt.Sprintf("Finished Abitti box update, version is: %s", box.GetVersion()))
		}()
	})
}

func bindOnInstallExamServer(mainUIStatus chan string) {
	buttonInstallExamServer.OnClicked(func(*ui.Button) {
		disableUI(mainUIStatus)
		examInstallWindow.Show()
	})

	examInstallButtonInstall.OnClicked(func(*ui.Button) {
		go func() {
			passphrase := examInstallPassphraseEntry.Text()
			examInstallPassphraseEntry.SetText("")
			disableUI(mainUIStatus)
			if passphrase != "" {
				examInstallWindow.Hide()
				install.NewExamServer(passphrase)
				translateUILabels()
				log.Debug(fmt.Sprintf("Finished Exam box update, version is: %s", box.GetVersion()))
			} else {
				mebroutines.ShowWarningMessage(xlate.Get("Please enter passphrase to install the exam server"))
			}
			enableUI(mainUIStatus)
		}()
	})

	examInstallButtonCancel.OnClicked(func(*ui.Button) {
		examInstallWindow.Hide()
		examInstallPassphraseEntry.SetText("")
		enableUI(mainUIStatus)
	})

	examInstallWindow.OnClosing(func(*ui.Window) bool {
		examInstallWindow.Hide()
		enableUI(mainUIStatus)
		examInstallPassphraseEntry.SetText("")
		return false
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

func setLogDeliveryLabelTextInGoroutine(text string) {
	log.Debug(fmt.Sprintf("Log delivery status: %s", text))
	ui.QueueMain(func() {
		logDeliveryStatusLabel.SetText(text)
	})
}

func bindOnDeliverLogs(mainUIStatus chan string) {
	buttonDeliverLogs.OnClicked(func(*ui.Button) {
		disableUI(mainUIStatus)
		buttonDeliverLogs.Disable()

		logDeliveryFilenameLabel.SetText(xlate.Get("Wait..."))
		logDeliveryStatusLabel.SetText(fmt.Sprintf(xlate.Get("Copying logs: %s"), xlate.Get("0 % (this can take a while...)")))
		logDeliveryWindow.Show()

		go func() {
			copyDoneChannel, copyProgressChannel := logdelivery.RequestLogsFromServer()
			followLogCopyProgress(copyDoneChannel, copyProgressChannel)

			logFilename, zipProgressChannel, zipErrorChannel := logdelivery.CollectLogsToZip()

			ui.QueueMain(func() {
				logDeliveryFilenameLabel.SetText(logFilename)
			})

			if followLogDeliveryZippingProgress(zipProgressChannel, zipErrorChannel) != nil {
				return
			}

			if network.CheckIfNetworkAvailable() {
				setLogDeliveryLabelTextInGoroutine(xlate.Get("Sending logs"))
				err := logdelivery.SendLogs(logFilename, func(progress uint8) {
					setLogDeliveryLabelTextInGoroutine(fmt.Sprintf(xlate.Get("Sending logs: %d %%"), progress))
				})
				if err != nil {
					setLogDeliveryLabelTextInGoroutine(fmt.Sprintf(xlate.Get("Error sending logs: %s"), err))
				} else {
					setLogDeliveryLabelTextInGoroutine(xlate.Get("Logs sent!"))
				}
			} else {
				setLogDeliveryLabelTextInGoroutine(xlate.Get("Cannot send logs because there is no Internet connection. Logs are in a zip archive in the ktp-jako folder."))
			}
		}()
	})
}

func followLogCopyProgress(copyDoneChannel chan bool, copyProgressChannel chan string) {
	for {
		select {
		case copyDone := <-copyDoneChannel:
			if copyDone {
				setLogDeliveryLabelTextInGoroutine(xlate.Get("Done copying"))
				return
			}
		case copyProgress := <-copyProgressChannel:
			if copyProgress != "0 %" {
				setLogDeliveryLabelTextInGoroutine(fmt.Sprintf(xlate.Get("Copying logs: %s"), copyProgress))
			}
		}
	}
}

func followLogDeliveryZippingProgress(zipProgressChannel chan uint8, zipErrorChannel chan error) error {
	for {
		select {
		case zipProgress := <-zipProgressChannel:
			if zipProgress <= 100 {
				setLogDeliveryLabelTextInGoroutine(fmt.Sprintf(xlate.Get("Zipping logs: %d %%"), zipProgress))
			} else {
				setLogDeliveryLabelTextInGoroutine(xlate.Get("Done zipping"))
				return nil
			}
		case zipError := <-zipErrorChannel:
			setLogDeliveryLabelTextInGoroutine(fmt.Sprintf(xlate.Get("Error zipping logs: %s"), zipError))
			return zipError
		}
	}
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
		go func() {
			pathBackup := filepath.Join(backupMediaPath[backupCombobox.Selected()], backup.GetBackupFilename(time.Now()))
			log.Debug(fmt.Sprintf("Starting backup to: %s", pathBackup))

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
	})

	backupButtonCancel.OnClicked(func(*ui.Button) {
		backupWindow.Hide()
		enableUI(mainUIStatus)
	})

	backupWindow.OnClosing(func(*ui.Window) bool {
		backupWindow.Hide()
		enableUI(mainUIStatus)
		return false
	})
}

func bindOnLogDelivery(mainUIStatus chan string) {
	logDeliveryFilenameCopyButton.OnClicked(func(*ui.Button) {
		err := clipboard.WriteAll(logDeliveryFilenameLabel.Text())
		if err != nil {
			log.Debug(fmt.Sprintf("Could not write to clipboard: %v", err))
		}
	})

	logDeliveryButtonClose.OnClicked(func(*ui.Button) {
		logDeliveryWindow.Hide()
		buttonDeliverLogs.Enable()
		enableUI(mainUIStatus)
	})

	logDeliveryWindow.OnClosing(func(*ui.Window) bool {
		logDeliveryWindow.Hide()
		buttonDeliverLogs.Enable()
		enableUI(mainUIStatus)
		return false
	})
}

// dupl linter finds this too similar with bindOnRemove()
// nolint: dupl
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

// dupl linter finds this too similar with bindOnDestroy()
// nolint: dupl
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
				progress.TranslateAndSetMessage("Server was removed successfully.")
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
		createLogDeliveryElements()
		createExamInstallElements()
		createDestroyElements()
		createRemoveElements()

		mebroutines.SetMainWindow(window)
		progress.SetProgressLabel(labelStatus)

		// Define command channel & goroutine for disabling/enabling main UI buttons
		mainUIStatus := make(chan string)
		mainUIUpdate := time.NewTicker(5 * time.Second)

		networkstatus.Update()

		setupMainLoop(mainUIStatus, mainUIUpdate)

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
		bindOnInstallAbittiServer(mainUIStatus)
		bindOnInstallExamServer(mainUIStatus)
		bindOnMakeBackup(mainUIStatus)
		bindOnDeliverLogs(mainUIStatus)
		bindOnDestroyServer(mainUIStatus)
		bindOnRemoveServer(mainUIStatus)
		bindOnMebShare()

		bindOnBackup(mainUIStatus)
		bindOnLogDelivery(mainUIStatus)
		bindOnDestroy(mainUIStatus)
		bindOnRemove(mainUIStatus)

		window.OnClosing(func(*ui.Window) bool {
			log.Debug("User exits through window exit")
			ui.Quit()
			return false
		})

		window.Show()
		enableUI(mainUIStatus)
		backupWindow.Hide()

		// Make sure we have VBoxManage
		if !vboxmanage.IsInstalled() {
			mebroutines.ShowTranslatedErrorMessage("Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?")
			log.Debug("VBoxManage is missing, disabling UI")
			disableUI(mainUIStatus)
		}

		// Make sure Hyper-V is not running
		// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
		isHyperV := make(chan bool)
		go func() {
			// IsHyperV() uses Windows WMI call
			isHyperV <- host.IsHyperV()
		}()

		if <-isHyperV {
			mebroutines.ShowWarningMessage(xlate.Get("Please turn Windows Hypervisor and Windows Virtualization Platform off as those may cause problems."))
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

		log.Debug(fmt.Sprintf("Currently installed box: %s %s", box.GetVersion(), box.GetType()))

		logdelivery.DeleteLogCopyFiles()
	})
}
