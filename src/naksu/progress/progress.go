package progress

import (
	"fmt"
	"naksu/mebroutines"
	"naksu/xlate"

	"github.com/andlabs/ui"
)

var progressLabel *ui.Label

// SetProgressLabel sets the object label object
func SetProgressLabel(newProgressLabel *ui.Label) {
	progressLabel = newProgressLabel
}

// SetMessage sets progress label text
func SetMessage(message string) {
	mebroutines.LogDebug(fmt.Sprintf("Progress message: %s", message))
	ui.QueueMain(func() {
		progressLabel.SetText(message)
	})
}

// TranslateAndSetMessage translates and sets progress label text
func TranslateAndSetMessage(message string) {
	SetMessage(xlate.Get(message))
}
