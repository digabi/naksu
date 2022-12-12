package mebroutines

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"naksu/log"
)

func quoteWindowsCommandArgs(args []string) []string {
	escapedArgs := []string{}

	for _, thisArg := range args {
		escapedArgs = append(escapedArgs, fmt.Sprintf(`"%s"`, thisArg))
	}

	return escapedArgs
}

// RunAndGetOutput runs command with arguments and returns output as a string
func RunAndGetOutput(origCommandArgs []string, logAction bool) (string, error) {
	windowsComSpec := os.Getenv("ComSpec")
	if windowsComSpec == "" {
		windowsComSpec = "C:\\Windows\\system32\\cmd.exe"
		log.Warning("For some reason Windows has not set 'ComSpec' environment variable. Falling back to hard-coded default '%s'.", windowsComSpec)
	}

	unescapedCommandArgs := append([]string{windowsComSpec, "/c"}, origCommandArgs...)
	escapedCommandArgs := quoteWindowsCommandArgs(unescapedCommandArgs)

	if logAction {
		log.Debug("RunAndGetOutput: %s", strings.Join(escapedCommandArgs, " "))
	}

	cmd := exec.Command(windowsComSpec)
	cmd.SysProcAttr = &syscall.SysProcAttr{ // nolint: exhaustruct
		CmdLine:    strings.Join(escapedCommandArgs, " "),
		HideWindow: true,
	}

	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Debug("command failed: %s (%v)", strings.Join(escapedCommandArgs, " "), err)
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
