package start

import (
	"naksu/mebroutines"
	"naksu/progress"
)

func Do_start_server() {
	// chdir ~/ktp
	if (! mebroutines.Chdir_vagrant_directory()) {
		mebroutines.Message_error("Could not change to vagrant directory ~/ktp")
	}

	// Start VM
	progress.Set_message_xlate("Starting Exam server. This takes a while.")
	run_params_up := []string{"up"}
	mebroutines.Run_vagrant(run_params_up)
}
