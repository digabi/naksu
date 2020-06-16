package host

import "fmt"

import "naksu/mebroutines"

// getOutput executes `cat filename` and returns output
func getOutput(commandArgs []string) string {
	output, err := mebroutines.RunAndGetOutput(commandArgs, false)

	if err != nil {
		output = "n/a"
	}

	return output
}

// GetHwLog returns a single string containing various hardware
// information to be printed to a log file
func GetHwLog() string {
	cpuinfo := getOutput([]string{"cat", "/proc/cpuinfo"})
	meminfo := getOutput([]string{"cat", "/proc/meminfo"})
	lshw := getOutput([]string{"lshw"})
	lspci := getOutput([]string{"lspci"})
	lsusb := getOutput([]string{"lsusb"})

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
