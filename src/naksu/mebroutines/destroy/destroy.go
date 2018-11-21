package destroy

import (
	"naksu/mebroutines"
	"naksu/progress"
)

// Server destroys existing exam server
func Server() {
	// chdir ~/ktp
	if !mebroutines.ChdirVagrantDirectory() {
		mebroutines.ShowErrorMessage("Could not change to vagrant directory ~/ktp")
	}

	// Start VM
	progress.TranslateAndSetMessage("Destroying Exam server. This takes a while.")
	destroyRunParams := []string{"destroy", "-f"}
	mebroutines.RunVagrant(destroyRunParams)
}
