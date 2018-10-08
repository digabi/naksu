package progress

import (
  "fmt"

  "naksu/xlate"
  "naksu/mebroutines"

  "github.com/andlabs/ui"
)

var progress_label *ui.Label

func Set_label_object (new_progress_label *ui.Label) {
  progress_label = new_progress_label
}

func Set_message (message string) {
  mebroutines.Message_debug(fmt.Sprintf("Progress message: %s", message))
  ui.QueueMain(func () {
    progress_label.SetText(message)
  })
}

func Set_message_xlate (message string) {
  Set_message(xlate.Get(message))
}
