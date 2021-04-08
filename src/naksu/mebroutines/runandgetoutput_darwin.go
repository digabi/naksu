package mebroutines

import (
	"os/exec"
	"strings"

	"naksu/log"
)

// RunAndGetOutput runs command with arguments and returns output as a string
func RunAndGetOutput(commandArgs []string, logAction bool) (string, error) {
	if logAction {
		log.Debug("RunAndGetOutput: %s", strings.Join(commandArgs, " "))
	}

	/* #nosec */
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)

	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Debug("command failed: %s (%v)", strings.Join(commandArgs, " "), err)
	}

	if out != nil {
		if logAction {
			log.Debug("RunAndGetOutput returns combined STDOUT and STDERR:")
			log.Debug(string(out))
		}
	} else {
		log.Debug("RunAndGetOutput returned NIL as combined STDOUT and STDERR")
	}

	return string(out), err
}
