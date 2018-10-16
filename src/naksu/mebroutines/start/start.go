package start

import (
	"naksu/mebroutines"
	"naksu/progress"
)

// StartServer starts exam server by running vagrant
func StartServer() {
	// chdir ~/ktp
	if !mebroutines.ChdirVagrantDirectory() {
		mebroutines.ShowErrorMessage("Could not change to vagrant directory ~/ktp")
	}

	// Start VM
	progress.TranslateAndSetMessage("Starting Exam server. This takes a while.")
	upRunParams := []string{"up"}
	mebroutines.RunVagrant(upRunParams)
}
