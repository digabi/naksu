package progress

import (
	"naksu/log"
	"naksu/xlate"

	"github.com/andlabs/ui"
)

var progressLabel *ui.Label
var lastMessage string

// SetProgressLabel sets the object label object
func SetProgressLabel(newProgressLabel *ui.Label) {
	progressLabel = newProgressLabel
	lastMessage = ""
}

// setMessage does the actual message label updating
func setMessage(message string) {
	log.Debug("Progress message: %s", message)
	ui.QueueMain(func() {
		progressLabel.SetText(message)
	})
}

// SetMessage sets progress label text
func SetMessage(message string) {
	lastMessage = message
	setMessage(message)
}

// TranslateAndSetMessage translates and sets progress label text
func TranslateAndSetMessage(str string, vars ...interface{}) {
	message := xlate.Get(str, vars...)

	lastMessage = message
	setMessage(xlate.Get(message))
}

// GetLastMessage returns last message string
func GetLastMessage() string {
	return lastMessage
}
