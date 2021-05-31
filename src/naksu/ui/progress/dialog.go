package progress

import (
	"strconv"

	"github.com/andlabs/ui"

	"naksu/xlate"
)

// ProgressDialog Instance
type ProgressDialog struct {
	Window        *ui.Window
	Progress      *ui.ProgressBar
	Message       *ui.Label
	MessageString string
}

// ShowProgressDialog opens a progress dialog
func ShowProgressDialog(message string) ProgressDialog {
	progressWindow := ui.NewWindow("", 400, 1, false)
	//progressWindow.SetBorderless(true)
	progressBox := ui.NewVerticalBox()
	progressBox.SetPadded(true)
	progressBar := ui.NewProgressBar()
	status := ui.NewLabel(message)
	progressBox.Append(status, true)
	progressBox.Append(progressBar, true)
	progressWindow.SetMargined(true)
	progressWindow.SetChild(progressBox)
	ui.QueueMain(func() {
		progressWindow.Show()
	})
	return ProgressDialog{Progress: progressBar, Message: status, Window: progressWindow, MessageString: message}
}

// TranslateAndShowProgressDialog translates message and then opens the progress dialog
func TranslateAndShowProgressDialog(message string) ProgressDialog {
	return ShowProgressDialog(message)
}

// UpdateProgressDialog updates the progress bar progress
func UpdateProgressDialog(dialog ProgressDialog, progress int, message *string) {
	if dialog.Window != nil && dialog.Window.Visible() {
		dialog.Progress.SetValue(progress)
		if message != nil {
			dialog.Message.SetText(*message + " (" + strconv.Itoa(progress) + "%)")
			dialog.MessageString = *message
		} else {
			dialog.Message.SetText(dialog.MessageString + " (" + strconv.Itoa(progress) + "%)")
		}
	}
}

// TranslateAndUpdateProgressDialog translates the message, and then updates the progress bar progress
func TranslateAndUpdateProgressDialog(dialog ProgressDialog, progress int, message *string) {
	translatedMessage := xlate.Get(*message)
	UpdateProgressDialog(dialog, progress, &translatedMessage)
}

// TranslateAndUpdateProgressDialogWithMessage translates the message, and then updates the progress bar progress
func TranslateAndUpdateProgressDialogWithMessage(dialog ProgressDialog, progress int, message string) {
	translatedMessage := xlate.Get(message)
	UpdateProgressDialog(dialog, progress, &translatedMessage)
}

// CloseProgressDialog closes given progress dialog
func CloseProgressDialog(dialog ProgressDialog) {
	if dialog.Window != nil && dialog.Window.Visible() {
		dialog.Window.ControlBase.Destroy()
	}
}
