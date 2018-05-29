package main

import (
	"mebroutines"
	"fmt"
)

func main() {
	// Make sure we have vagrant
	if (! mebroutines.If_found_vagrant()) {
		mebroutines.Message_error(fmt.Sprintf("Could not execute vagrant. Are you sure you have installed HashiCorp Vagrant?"))
	}

	// Make sure we have VBoxManage
	if (! mebroutines.If_found_vboxmanage()) {
		mebroutines.Message_error(fmt.Sprintf("Could not execute VBoxManage. Are you sure you have installed Oracle VirtualBox?"))
	}

	// Start VM
	mebroutines.Run_vagrant("up")
}
