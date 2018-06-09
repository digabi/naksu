package start

import (
	"mebroutines"
)

func Do_start_server() {
	// Make sure we have vagrant
	if (! mebroutines.If_found_vagrant()) {
		mebroutines.Message_error("Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?")
	}

	// Make sure we have VBoxManage
	if (! mebroutines.If_found_vboxmanage()) {
		mebroutines.Message_error("Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?")
	}

	// chdir ~/ktp
	if (! mebroutines.Chdir_vagrant_directory()) {
		mebroutines.Message_error("Could not change to vagrant directory ~/ktp")
	}

	// Start VM
	run_params_up := []string{"up"}
	mebroutines.Run_vagrant(run_params_up)
}
