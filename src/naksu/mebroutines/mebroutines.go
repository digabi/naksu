// Package mebroutines contains general routines used by various MEB helper utilities
package mebroutines

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	golog "log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"naksu/config"
	"naksu/log"
	"naksu/xlate"

	"github.com/andlabs/ui"
	"github.com/mitchellh/go-homedir"
)

var mainWindow *ui.Window

// Close gracefully handles closing of closable item. defer Close(item)
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		golog.Fatal(err)
	}
}

// getRunEnvironment returns array of strings containing environment strings
func getRunEnvironment() []string {
	runEnv := os.Environ()

	config.Load()

	envs := []struct {
		envName  string
		envValue string
	}{
		{"NIC", config.GetNic()},
		{"EXTNIC", config.GetExtNic()},
	}

	for _, thisEnv := range envs {
		if thisEnv.envValue != "" {
			runEnv = append(runEnv, fmt.Sprintf("%s=%s", thisEnv.envName, thisEnv.envValue))
			log.Debug(fmt.Sprintf("Adding environment value %s=%s", thisEnv.envName, thisEnv.envValue))
		}
	}

	return runEnv
}

// RunAndGetOutput runs command with arguments and returns output as a string
func RunAndGetOutput(commandArgs []string) (string, error) {
	log.Debug(fmt.Sprintf("RunAndGetOutput: %s", strings.Join(commandArgs, " ")))
	/* #nosec */
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	cmd.Env = getRunEnvironment()

	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Debug(fmt.Sprintf(xlate.Get("command failed: %s (%v)"), strings.Join(commandArgs, " "), err))
	}

	if out != nil {
		log.Debug("RunAndGetOutput returns combined STDOUT and STDERR:")
		log.Debug(string(out))
	} else {
		log.Debug("RunAndGetOutput returned NIL as combined STDOUT and STDERR")
	}

	return string(out), err
}

// RunAndGetError runs command with arguments and returns error code
func RunAndGetError(commandArgs []string) (string, error) {
	log.Debug(fmt.Sprintf("RunAndGetError: %s", strings.Join(commandArgs, " ")))

	var stderr bytes.Buffer

	/* #nosec */
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = &stderr
	cmd.Env = getRunEnvironment()

	err := cmd.Run()

	log.Debug("RunAndGetError returns STDERR:")
	log.Debug(stderr.String())

	return stderr.String(), err
}

// GetVagrantPath returns path for Vagrant binary
func GetVagrantPath() string {
	var path = "vagrant"
	if os.Getenv("VAGRANTPATH") != "" {
		path = os.Getenv("VAGRANTPATH")
	}

	return path
}

// RunVagrant executes Vagrant with given arguments
func RunVagrant(args []string) error {
	runVagrant := []string{GetVagrantPath()}
	runArgs := append(runVagrant, args...)
	vagrantOutput, err := RunAndGetError(runArgs)
	if err != nil {
		matchedTimeout, errTimeout := regexp.MatchString("Timed out while waiting for the machine to boot", vagrantOutput)
		matchedMacAddress, errMacAddress := regexp.MatchString("error: --macaddress: RTGetOpt: Command line option needs argument", vagrantOutput)
		matchedConnectionRefused, errConnectionRefused := regexp.MatchString("The guest machine entered an invalid state", vagrantOutput)
		if errTimeout == nil && matchedTimeout {
			// We've obviously started the VM
			log.Debug("Running vagrant gives me timeout - things are probably ok. User was not notified. Complete output:")
			log.Debug(vagrantOutput)
			// This case should not be considered as an error
			return nil
		} else if errMacAddress == nil && matchedMacAddress {
			// Vagrant in Windows host give this error message - just restart vagrant and you're good
			log.Debug("Running vagrant gives me a known --macaddress/RTGetOpt error (generated every time by vagrant after a new box in Windows). Complete output:")
			log.Debug(vagrantOutput)
			return errors.New("macaddress/rtgetopt")
		} else if errConnectionRefused == nil && matchedConnectionRefused {
			log.Debug("Vagrant entered invalid state while booting. We expect this to occur because user has closed the VM window. User was not notified. Complete output:")
			log.Debug(vagrantOutput)
			// This case should not be considered as an error
			return nil
		}

		log.Debug(fmt.Sprintf("Failed to execute %s (%v), complete output:", strings.Join(runArgs, " "), err))
		log.Debug(vagrantOutput)
	}

	return err
}

// RunVBoxManage runs vboxmanage command with given arguments
func RunVBoxManage(args []string) (string, error) {
	vboxmanagepathArr := []string{getVBoxManagePath()}
	runArgs := append(vboxmanagepathArr, args...)
	vBoxManageOutput, err := RunAndGetOutput(runArgs)
	if err != nil {
		logError := func(output string, err error) {
			log.Debug(fmt.Sprintf("Failed to execute %s (%v), complete output:", strings.Join(runArgs, " "), err))
			log.Debug(output)
		}

		logError(vBoxManageOutput, err)

		fixed, fixErr := detectAndFixDuplicateHardDiskProblem(vBoxManageOutput)
		if !fixed || fixErr != nil {
			log.Debug(fmt.Sprintf("Failed to detect & fix duplicate hard disk problem: %v", fixErr))
			return "",errors.New("failed to fix duplicate hard disk problem")
		}

		log.Debug(fmt.Sprintf("Retrying '%s' after fixing problem", strings.Join(runArgs, " ")))
		vBoxManageOutput, err = RunAndGetOutput(runArgs)
		if err != nil {
			logError(vBoxManageOutput, err)
		}
	}

	return vBoxManageOutput, err
}

// IfFoundVagrant returns true if Vagrant can be found in path
func IfFoundVagrant() bool {
	var vagrantpath = GetVagrantPath()

	runParams := []string{vagrantpath, "--version"}

	vagrantVersion, err := RunAndGetOutput(runParams)
	if err != nil {
		// No vagrant was found
		return false
	}

	log.Debug(fmt.Sprintf("vagrant says: %s", vagrantVersion))

	return true
}

// IfFoundVBoxManage returns true if vboxmanage can be found in path
func IfFoundVBoxManage() bool {
	var vboxmanagepath = getVBoxManagePath()

	if vboxmanagepath == "" {
		log.Debug("Could not get VBoxManage path")
		return false
	}

	runParams := []string{vboxmanagepath, "--version"}

	vBoxManageVersion, err := RunAndGetOutput(runParams)
	if err != nil {
		// No VBoxManage was found
		return false
	}

	log.Debug(fmt.Sprintf("VBoxManage says: %s", vBoxManageVersion))

	return true
}

func getFileMode(path string) (os.FileMode, error) {
	fi, err := os.Lstat(path)
	if err == nil {
		return fi.Mode(), nil
	}

	return 0, err
}

// ExistsDir returns true if given path exists
func ExistsDir(path string) bool {
	mode, err := getFileMode(path)

	if err == nil && mode.IsDir() {
		return true
	}

	return false
}

// ExistsFile returns true if given file exists
func ExistsFile(path string) bool {
	mode, err := getFileMode(path)

	if err == nil && mode.IsRegular() {
		return true
	}

	return false
}

// ExistsCharDevice returns true if given file is a Linux device file
func ExistsCharDevice(path string) bool {
	mode, err := getFileMode(path)

	return err == nil && mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0
}

// CreateDir creates new directory
func CreateDir(path string) error {
	var err = os.Mkdir(path, os.ModePerm)
	return err
}

// CreateFile creates empty new file
func CreateFile(path string) error {
	f, err := os.Create(path)
	if err == nil {
		defer Close(f)
	}
	return err
}

// RemoveDir removes directory and all its contents
func RemoveDir(path string) error {
	err := os.RemoveAll(path)
	return err
}

// CopyFile copies existing file
func CopyFile(src, dst string) (err error) {
	log.Debug(fmt.Sprintf("Copying file %s to %s", src, dst))

	if !ExistsFile(src) {
		log.Debug("Copying failed, could not find source file")
		return errors.New("could not find source file")
	}

	/* #nosec */
	in, err := os.Open(src)
	if err != nil {
		log.Debug(fmt.Sprintf("Copying failed while opening source file: %v", err))
		return
	}
	defer Close(in)

	out, err := os.Create(dst)
	if err != nil {
		log.Debug(fmt.Sprintf("Copying failed while opening destination file: %v", err))
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// IfIntlCharsInPath returns true if path contains non-ASCII characters
func IfIntlCharsInPath(path string) bool {
	matched, err := regexp.MatchString(`[^a-zA-Z0-9_\-\/\:\\ \.]`, path)
	if err == nil && matched {
		return true
	}

	return false
}

// GetHomeDirectory returns home directory path
func GetHomeDirectory() string {
	homeDir, err := homedir.Dir()

	if err != nil {
		panic("Could not get home directory")
	}

	return homeDir
}

// GetVagrantDirectory returns ktp-directory path from under home directory
func GetVagrantDirectory() string {
	return filepath.Join(GetHomeDirectory(), "ktp")
}

// GetVagrantdDirectory returns .vagrantd-directory path from under home directory
func GetVagrantdDirectory() string {
	return filepath.Join(GetHomeDirectory(), ".vagrant.d")
}

// GetMebshareDirectory returns ktp-jako path from under home directory
func GetMebshareDirectory() string {
	return filepath.Join(GetHomeDirectory(), "ktp-jako")
}

// GetVirtualBoxHiddenDirectory returns ".VirtualBox" path from under home directory
func GetVirtualBoxHiddenDirectory() string {
	return filepath.Join(GetHomeDirectory(), ".VirtualBox")
}

// GetVirtualBoxVMsDirectory returns "VirtualBox VMs" path from under home directory
func GetVirtualBoxVMsDirectory() string {
	return filepath.Join(GetHomeDirectory(), "VirtualBox VMs")
}

// chdir changes current working directory to the given directory
func chdir(chdirTo string) bool {
	log.Debug(fmt.Sprintf("chdir %s", chdirTo))
	err := os.Chdir(chdirTo)
	if err != nil {
		log.Debug(fmt.Sprintf("Could not chdir to %s: %v", chdirTo, err))
		return false
	}

	return true
}

// ChdirVagrantDirectory changes current working directory to vagrant path (ktp)
func ChdirVagrantDirectory() bool {
	return chdir(GetVagrantDirectory())
}

// ChdirHomeDirectory changes current working directory to home directory
func ChdirHomeDirectory() bool {
	return chdir(GetHomeDirectory())
}

// SetMainWindow sets libui main window pointer used by ShowErrorMessage and ShowWarningMessage
func SetMainWindow(win *ui.Window) {
	mainWindow = win
}

// ShowErrorMessage shows error message popup to user
func ShowErrorMessage(message string) {
	fmt.Printf("ERROR: %s\n\n", message)
	log.Debug(fmt.Sprintf("ERROR: %s", message))

	// Show libui box if main window has been set with Set_main_window
	if mainWindow != nil {
		ui.QueueMain(func() {
			ui.MsgBoxError(mainWindow, xlate.Get("Error"), message)
		})
	}
}

// ShowWarningMessage shows warning message popup to user
func ShowWarningMessage(message string) {
	fmt.Printf("WARNING: %s\n", message)
	log.Debug(fmt.Sprintf("WARNING: %s", message))

	// Show libui box if main window has been set with Set_main_window
	if mainWindow != nil {
		ui.QueueMain(func() {
			ui.MsgBox(mainWindow, xlate.Get("Warning"), message)
		})
	}
}

// ShowInfoMessage shows warning message popup to user
func ShowInfoMessage(message string) {
	fmt.Printf("INFO: %s\n", message)
	log.Debug(fmt.Sprintf("INFO: %s", message))

	// Show libui box if main window has been set with Set_main_window
	if mainWindow != nil {
		ui.QueueMain(func() {
			ui.MsgBox(mainWindow, xlate.Get("Info"), message)
		})
	}
}
