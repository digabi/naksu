// Package mebroutines contains general routines used by various MEB helper utilities
package mebroutines

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"naksu/xlate"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/andlabs/ui"
)

var isDebug bool
var mainWindow *ui.Window
var debugFilename string

// Close gracefully handles closing of closable item. defer Close(item)
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Run executes command with arguments
func Run(commandArgs []string) error {
	LogDebug(fmt.Sprintf("run: %s", strings.Join(commandArgs, " ")))
	/* #nosec */
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		ShowWarningMessage(fmt.Sprintf(xlate.Get("command failed: %s"), strings.Join(commandArgs, " ")))
	}

	return err
}

// RunAndGetOutput runs command with arguments and returns output as a string
func RunAndGetOutput(commandArgs []string) (string, error) {
	LogDebug(fmt.Sprintf("RunAndGetOutput: %s", strings.Join(commandArgs, " ")))
	/* #nosec */
	out, err := exec.Command(commandArgs[0], commandArgs[1:]...).CombinedOutput()
	if err != nil {
		// Executing failed, return error condition
		ShowWarningMessage(fmt.Sprintf(xlate.Get("command failed: %s"), strings.Join(commandArgs, " ")))
		return string(out), err
	}

	if out != nil {
		LogDebug("RunAndGetOutput returns combined STDOUT and STDERR:")
		LogDebug(string(out))
	} else {
		LogDebug("RunAndGetOutput returned NIL as combined STDOUT and STDERR")
	}

	return string(out), nil
}

// RunAndGetError runs command with arguments and returns error code
func RunAndGetError(commandArgs []string) (string, error) {
	LogDebug(fmt.Sprintf("RunAndGetError: %s", strings.Join(commandArgs, " ")))

	var stderr bytes.Buffer

	/* #nosec */
	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = &stderr

	err := cmd.Run()

	LogDebug("RunAndGetError returns STDERR:")
	LogDebug(stderr.String())

	return stderr.String(), err
}

func getVagrantPath() string {
	var path = "vagrant"
	if os.Getenv("VAGRANTPATH") != "" {
		path = os.Getenv("VAGRANTPATH")
	}

	return path
}

// RunVagrant executes Vagrant with given arguments
func RunVagrant(args []string) {
	runVagrant := []string{getVagrantPath()}
	runArgs := append(runVagrant, args...)
	vagrantOutput, err := RunAndGetError(runArgs)
	if err != nil {
		matchedTimeout, errTimeout := regexp.MatchString("Timed out while waiting for the machine to boot", vagrantOutput)
		matchedMacAddress, errMacAddress := regexp.MatchString("error: --macaddress: RTGetOpt: Command line option needs argument", vagrantOutput)
		matchedConnectionRefused, errConnectionRefused := regexp.MatchString("The guest machine entered an invalid state", vagrantOutput)
		if errTimeout == nil && matchedTimeout {
			// We've obviously started the VM
			LogDebug("Running vagrant gives me timeout - things are probably ok. User was not notified. Complete output:")
			LogDebug(vagrantOutput)
		} else if errMacAddress == nil && matchedMacAddress {
			// Vagrant in Windows host give this error message - just restart vagrant and you're good
			ShowInfoMessage(xlate.Get("Server failed to start. This is typical in Windows after an update. Please try again to start the server."))
		} else if errConnectionRefused == nil && matchedConnectionRefused {
			LogDebug("Vagrant entered invalid state while booting. We expect this to occur because user has closed the VM window. User was not notified. Complete output:")
			LogDebug(vagrantOutput)
		} else {
			LogDebug(fmt.Sprintf("Failed to execute %s, complete output:", strings.Join(runArgs, " ")))
			LogDebug(vagrantOutput)
			ShowWarningMessage(fmt.Sprintf(xlate.Get("Failed to execute %s"), strings.Join(runArgs, " ")))
		}
	}
}

// RunVBoxManage runs vboxmanage command with given arguments
func RunVBoxManage(args []string) string {
	vboxmanagepathArr := []string{getVBoxManagePath()}
	runArgs := append(vboxmanagepathArr, args...)
	vBoxManageOutput, err := RunAndGetOutput(runArgs)
	if err != nil {
		LogDebug(fmt.Sprintf("Failed to execute %s, complete output:", strings.Join(runArgs, " ")))
		LogDebug(vBoxManageOutput)
		ShowErrorMessage(fmt.Sprintf(xlate.Get("Failed to execute %s"), strings.Join(runArgs, " ")))
	}

	return vBoxManageOutput
}

// IfFoundVagrant returns true if Vagrant can be found in path
func IfFoundVagrant() bool {
	var vagrantpath = getVagrantPath()

	runParams := []string{vagrantpath, "--version"}

	vagrantVersion, err := RunAndGetOutput(runParams)
	if err != nil {
		// No vagrant was found
		return false
	}

	LogDebug(fmt.Sprintf("vagrant says: %s", vagrantVersion))

	return true
}

// IfFoundVBoxManage returns true if vboxmanage can be found in path
func IfFoundVBoxManage() bool {
	var vboxmanagepath = getVBoxManagePath()

	if vboxmanagepath == "" {
		LogDebug("Could not get VBoxManage path")
		return false
	}

	runParams := []string{vboxmanagepath, "--version"}

	vBoxManageVersion, err := RunAndGetOutput(runParams)
	if err != nil {
		// No VBoxManage was found
		return false
	}

	LogDebug(fmt.Sprintf("VBoxManage says: %s", vBoxManageVersion))

	return true
}

// GetVagrantBoxVersion returns version string for vagrant box
func GetVagrantBoxVersion() string {
	boxIndexUUID := GetVagrantBoxIndexUUID()
	if boxIndexUUID == "" {
		// We did not get machine index UUID so we cannot return any version string
		return ""
	}

	indexFilename := GetVagrantdDirectory() + string(os.PathSeparator) + "data" + string(os.PathSeparator) + "machine-index" + string(os.PathSeparator) + "index"
	/* #nosec */
	fileContent, err := ioutil.ReadFile(indexFilename)
	if err != nil {
		LogDebug(fmt.Sprintf("Could not read from %s", indexFilename))
		return ""
	}

	LogDebug("Vagrant machine-index/index:")
	LogDebug(fmt.Sprintf("%s", fileContent))

	// Get version from JSON structure
	var jsonData map[string]interface{}

	jsonErr := json.Unmarshal(fileContent, &jsonData)
	if jsonErr != nil {
		LogDebug("Unable on decode machine-index/index response:")
		LogDebug(fmt.Sprintf("%s", jsonErr))
		return ""
	}

	if jsonData["machines"] == nil {
		return ""
	}
	jsonMachines := jsonData["machines"].(map[string]interface{})

	if jsonMachines[boxIndexUUID] == nil {
		return ""
	}
	jsonOurMachine := jsonMachines[boxIndexUUID].(map[string]interface{})

	if jsonOurMachine["extra_data"] == nil {
		return ""
	}
	jsonExtraData := jsonOurMachine["extra_data"].(map[string]interface{})

	if jsonExtraData["box"] == nil {
		return ""
	}
	jsonBox := jsonExtraData["box"].(map[string]interface{})

	boxString := fmt.Sprintf("%s (%s %s)", GetVagrantBoxType(jsonBox["name"].(string)), jsonBox["name"], jsonBox["version"])

	return boxString
}

// GetVagrantBoxType returns the type string (Abitti server or Matric Exam server) for vagrant box name
func GetVagrantBoxType(name string) string {
	if name == "" {
		return "-"
	}

	if name == "digabi/ktp-qa" {
		return xlate.Get("Abitti server")
	}

	return xlate.Get("Matric Exam server")
}

// GetVagrantBoxIndexUUID return vagrant box index UUID for current vagrant box
func GetVagrantBoxIndexUUID() string {
	pathVagrant := GetVagrantDirectory()

	pathID := pathVagrant + string(os.PathSeparator) + ".vagrant" + string(os.PathSeparator) + "machines" + string(os.PathSeparator) + "default" + string(os.PathSeparator) + "virtualbox" + string(os.PathSeparator) + "index_uuid"

	/* #nosec */
	fileContent, err := ioutil.ReadFile(pathID)
	if err != nil {
		LogDebug(fmt.Sprintf(xlate.Get("Could not get vagrantbox index UUID: %d"), err))
		return ""
	}

	LogDebug(fmt.Sprintf("Vagrantbox index UUID is %s", string(fileContent)))

	return string(fileContent)
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

// CopyFile copies existing file
func CopyFile(src, dst string) (err error) {
	LogDebug(fmt.Sprintf("Copying file %s to %s", src, dst))

	if !ExistsFile(src) {
		LogDebug("Copying failed, could not find source file")
		return errors.New("Could not find source file")
	}

	/* #nosec */
	in, err := os.Open(src)
	if err != nil {
		LogDebug(fmt.Sprintf("Copying failed while opening source file: %v", err))
		return
	}
	defer Close(in)

	out, err := os.Create(dst)
	if err != nil {
		LogDebug(fmt.Sprintf("Copying failed while opening destination file: %v", err))
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
	homeWin := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if homeWin != "" {
		return homeWin
	}

	homeLinux := os.Getenv("HOME")
	if homeLinux != "" {
		return homeLinux
	}

	panic("Could not get home directory")
}

// GetVagrantDirectory returns ktp-directory path from under home directory
func GetVagrantDirectory() string {
	return GetHomeDirectory() + string(os.PathSeparator) + "ktp"
}

// GetVagrantdDirectory returns .vagrantd-directory path from under home directory
func GetVagrantdDirectory() string {
	return GetHomeDirectory() + string(os.PathSeparator) + ".vagrant.d"
}

// GetMebshareDirectory returns ktp-jako path from under home directory
func GetMebshareDirectory() string {
	return GetHomeDirectory() + string(os.PathSeparator) + "ktp-jako"
}

// ChdirVagrantDirectory changes current working directory to vagrant path (ktp)
func ChdirVagrantDirectory() bool {
	pathVagrant := GetVagrantDirectory()
	LogDebug(fmt.Sprintf("chdir %s", pathVagrant))
	err := os.Chdir(pathVagrant)
	if err != nil {
		ShowWarningMessage(fmt.Sprintf(xlate.Get("Could not chdir to %s"), pathVagrant))
		return false
	}

	return true
}

// SetMainWindow sets libui main window pointer used by ShowErrorMessage and ShowWarningMessage
func SetMainWindow(win *ui.Window) {
	mainWindow = win
}

// ShowErrorMessage shows error message popup to user
func ShowErrorMessage(message string) {
	fmt.Printf("FATAL ERROR: %s\n\n", message)
	appendLogFile(fmt.Sprintf("FATAL ERROR: %s", message))

	// Show libui box if main window has been set with Set_main_window
	if mainWindow != nil {
		ui.QueueMain(func() {
			ui.MsgBoxError(mainWindow, xlate.Get("Error"), message)
			ui.Quit()
		})
	} else {
		os.Exit(1)
	}
}

// ShowWarningMessage shows warning message popup to user
func ShowWarningMessage(message string) {
	fmt.Printf("WARNING: %s\n", message)
	appendLogFile(fmt.Sprintf("WARNING: %s", message))

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
	appendLogFile(fmt.Sprintf("INFO: %s", message))

	// Show libui box if main window has been set with Set_main_window
	if mainWindow != nil {
		ui.QueueMain(func() {
			ui.MsgBox(mainWindow, xlate.Get("Info"), message)
		})
	}
}

func appendLogFile(message string) {
	if debugFilename != "" {
		// Append only if the logfile has been set

		// Current timestamp
		t := time.Now()

		/* #nosec */
		f, err := os.OpenFile(debugFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			panic(fmt.Sprintf("Could not append to log file %s: %s", debugFilename, err))
		}
		defer Close(f)

		_, err = f.WriteString(fmt.Sprintf("[%s] %s\n", t.Format("2006-01-02 15:04:05"), message))
		if err != nil {
			if f.Sync() != nil {
				defer Close(f)
			}
		}
	}
}

// SetDebug enables debug printing if set to true
func SetDebug(newValue bool) {
	isDebug = newValue
}

// SetDebugFilename sets debug log path
func SetDebugFilename(newFilename string) {
	debugFilename = newFilename

	if debugFilename != "" && ExistsFile(debugFilename) {
		// Re-create the log file
		err := os.Remove(debugFilename)
		if err != nil {
			panic(fmt.Sprintf("Could not open log file %s: %s", debugFilename, err))
		}
	}
}

// IsDebug returns true if we need to log debug information
func IsDebug() bool {
	return isDebug
}

// LogDebug logs debug information to log file
func LogDebug(message string) {
	if IsDebug() {
		fmt.Printf("DEBUG: %s\n", message)
	}

	appendLogFile(fmt.Sprintf("DEBUG: %s", message))
}
