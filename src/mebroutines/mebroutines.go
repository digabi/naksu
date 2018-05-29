// General routines used by various MEB helper utilities
package mebroutines

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
)

func Run(command string) error {
	Message_debug(fmt.Sprintf("run: %s", command))
	args := strings.Split(command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
  err := cmd.Run()
	if err != nil {
		Message_warning(fmt.Sprintf("command failed: %s", command))
	}

  return err
}

func Run_get_output (command string) (string, error) {
  Message_debug(fmt.Sprintf("run_get_output: %s", command))
	args := strings.Split(command, " ")
	out,err := exec.Command(args[0], args[1:]...).Output()
  if (err != nil) {
    // Executing failed, return error condition
    Message_warning(fmt.Sprintf("command failed: %s", command))
    return "", err
  }

  return string(out), nil
}

func get_vagrant_path () string {
  var path = "vagrant"
	if (os.Getenv("VAGRANTPATH") != "") {
		path = os.Getenv("VAGRANTPATH")
	}

  return path
}

func get_vboxmanage_path () string {
  var path = "VBoxManage"
	if (os.Getenv("VBOXMANAGEPATH") != "") {
		path = os.Getenv("VBOXMANAGEPATH")
	}

  return path
}

func Run_vagrant (command string) {
  var vagrantpath = get_vagrant_path()

  err := Run(vagrantpath+" "+command)

  if (err != nil) {
    Message_error(fmt.Sprintf("Failed to execute %s %s", vagrantpath, command))
  }
}

func If_found_vagrant () bool {
  var vagrantpath = get_vagrant_path()

  vagrant_version,err := Run_get_output (vagrantpath+" --version")
  if (err != nil) {
    // No vagrant was found
    return false
  }

  Message_debug(fmt.Sprintf("vagrant says: %s", vagrant_version))

  return true
}

func If_found_vboxmanage () bool {
  var vboxmanagepath = get_vboxmanage_path()

  vboxmanage_version,err := Run_get_output (vboxmanagepath+" --version")
  if (err != nil) {
    // No VBoxManage was found
    return false
  }

  Message_debug(fmt.Sprintf("VBoxManage says: %s", vboxmanage_version))

  return true
}

func get_file_mode (path string) (os.FileMode, error) {
  fi, err := os.Lstat(path)
  if (err == nil) {
    return fi.Mode(), nil
  }

  return 0,err
}

func ExistsDir (path string) bool {
  mode,err := get_file_mode(path)

  if (err == nil && mode.IsDir()) {
    return true
  }

  return false
}

func ExistsFile (path string) bool {
  mode,err := get_file_mode(path)

  if (err == nil && mode.IsRegular()) {
    return true
  }

  return false
}

func CreateDir (path string) error {
  var err = os.Mkdir(path, os.ModePerm)
  return err
}

func Get_home_directory () string {
  home_win := os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH")
  if (home_win != "") {
    return home_win
  }

  home_linux := os.Getenv("HOME")
  if (home_linux != "") {
    return home_linux
  }

  panic("Could not get home directory")
}

func Message_error (message string) {
  fmt.Printf("FATAL ERROR: %s\n\n", message)
  os.Exit(1)
}

func Message_warning (message string) {
  fmt.Printf("WARNING: %s\n", message)
}

func Message_debug (message string) {
  fmt.Printf("DEBUG: %s\n", message)
}
