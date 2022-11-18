package progress

import (
	"strconv"

	"github.com/andlabs/ui"

	"naksu/xlate"
)

// Dialog ProgressDialog Instance
type Dialog struct {
	Window        *ui.Window
	Progress      *ui.ProgressBar
	Message       *ui.Label
	MessageString string
}

// ShowProgressDialog opens a progress dialog
func ShowProgressDialog(message string) Dialog {
	const progressDialogDefaultWidth = 400

	dialogChannel := make(chan Dialog)

	ui.QueueMain(func() {
		progressWindow := ui.NewWindow("", progressDialogDefaultWidth, 1, false)
		progressBox := ui.NewVerticalBox()
		progressBox.SetPadded(true)
		progressBar := ui.NewProgressBar()
		status := ui.NewLabel(message)
		progressBox.Append(status, true)
		progressBox.Append(progressBar, true)
		progressWindow.SetMargined(true)
		progressWindow.SetChild(progressBox)
		progressWindow.Show()

		dialogChannel <- Dialog{Progress: progressBar, Message: status, Window: progressWindow, MessageString: message}
	})

	return <-dialogChannel
}

// TranslateAndShowProgressDialog translates message and then opens the progress dialog
func TranslateAndShowProgressDialog(message string) Dialog {
	return ShowProgressDialog(xlate.Get(message))
}

// UpdateProgressDialog updates the progress bar progress
func UpdateProgressDialog(dialog Dialog, progress int, message *string) {
	if dialog.Window != nil && dialog.Window.Visible() {
		ui.QueueMain(func() {
			dialog.Progress.SetValue(progress)
			if message != nil {
				dialog.Message.SetText(*message + " (" + strconv.Itoa(progress) + "%)")
				dialog.MessageString = *message
			} else {
				dialog.Message.SetText(dialog.MessageString + " (" + strconv.Itoa(progress) + "%)")
			}
		})
	}
}

// TranslateAndUpdateProgressDialog translates the message, and then updates the progress bar progress
func TranslateAndUpdateProgressDialog(dialog Dialog, progress int, message *string) {
	translatedMessage := xlate.Get(*message)
	UpdateProgressDialog(dialog, progress, &translatedMessage)
}

// TranslateAndUpdateProgressDialogWithMessage translates the message, and then updates the progress bar progress
func TranslateAndUpdateProgressDialogWithMessage(dialog Dialog, progress int, message string) {
	translatedMessage := xlate.Get(message)
	UpdateProgressDialog(dialog, progress, &translatedMessage)
}

// CloseProgressDialog closes given progress dialog
func CloseProgressDialog(dialog Dialog) {
	if dialog.Window != nil && dialog.Window.Visible() {
		ui.QueueMain(func() {
			dialog.Window.ControlBase.Destroy()
		})
	}
}
