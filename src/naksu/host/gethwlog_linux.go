package host

import (
	"fmt"
	"os/exec"
	"strings"
)

// simpleRunAndGetOutput does almost the same as mebroutines.RunAndGetOutput, but
// it does not write anything on log and returns a string
func simpleRunAndGetOutput(commandArgs []string) string {
	/* #nosec */
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)

	out, err := cmd.CombinedOutput()

	var outString string

	if out == nil {
		outString = "n/a"
	} else {
		outString = string(out)
	}

	if err != nil {
		outString = fmt.Sprintf("command failed: %s (%v)", strings.Join(commandArgs, " "), err)
	}

	return outString
}

// GetHwLog returns a single string containing various hardware
// information to be printed to a log file
func GetHwLog() string {
	cpuinfo := simpleRunAndGetOutput([]string{"cat", "/proc/cpuinfo"})
	meminfo := simpleRunAndGetOutput([]string{"cat", "/proc/meminfo"})
	lshw := simpleRunAndGetOutput([]string{"lshw"})
	lspci := simpleRunAndGetOutput([]string{"lspci"})
	lsusb := simpleRunAndGetOutput([]string{"lsusb"})

	return fmt.Sprintf(`Output of /proc/cpuinfo
%s
Output of /proc/meminfo
%s
Output of lshw
%s
Output of lspci
%s
Output of lsusb
%s`, cpuinfo, meminfo, lshw, lspci, lsusb)
}
