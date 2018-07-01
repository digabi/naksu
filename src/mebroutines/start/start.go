package start

import (
	"mebroutines"
)

func Do_start_server() {
	// chdir ~/ktp
	if (! mebroutines.Chdir_vagrant_directory()) {
		mebroutines.Message_error("Could not change to vagrant directory ~/ktp")
	}

	// Start VM
	run_params_up := []string{"up"}
	mebroutines.Run_vagrant(run_params_up)
}
