// General routines used by various MEB helper utilities
package mebroutines

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
)

var is_debug bool

func Run(command_args []string) error {
	Message_debug(fmt.Sprintf("run: %s", strings.Join(command_args, " ")))
	cmd := exec.Command(command_args[0], command_args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
  err := cmd.Run()
	if err != nil {
		Message_warning(fmt.Sprintf("command failed: %s", strings.Join(command_args, " ")))
	}

  return err
}

func Run_get_output (command_args []string) (string, error) {
  Message_debug(fmt.Sprintf("run_get_output: %s", strings.Join(command_args, " ")))
	out,err := exec.Command(command_args[0], command_args[1:]...).Output()
  if (err != nil) {
    // Executing failed, return error condition
    Message_warning(fmt.Sprintf("command failed: %s", strings.Join(command_args, " ")))
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

func Run_vagrant (args []string) {
  vagrantpath_arr := []string{get_vagrant_path()}
  run_args := append(vagrantpath_arr, args...)
  err := Run(run_args)
  if (err != nil) {
    Message_error(fmt.Sprintf("Failed to execute %s", strings.Join(run_args, " ")))
  }
}

func If_found_vagrant () bool {
  var vagrantpath = get_vagrant_path()

  run_params := []string{vagrantpath, "--version"}

  vagrant_version,err := Run_get_output (run_params)
  if (err != nil) {
    // No vagrant was found
    return false
  }

  Message_debug(fmt.Sprintf("vagrant says: %s", vagrant_version))

  return true
}

func If_found_vboxmanage () bool {
  var vboxmanagepath = get_vboxmanage_path()

  if (vboxmanagepath == "") {
    Message_debug("Could not get VBoxManage path")
    return false
  }

  run_params := []string{vboxmanagepath, "--version"}

  vboxmanage_version,err := Run_get_output (run_params)
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

func Chdir_vagrant_directory () bool {
  path_vagrant := Get_home_directory()+string(os.PathSeparator)+"ktp"
  err := os.Chdir(path_vagrant)
  if (err != nil) {
    Message_warning(fmt.Sprintf("Could not chdir to %s", path_vagrant))
    return false
  }

  return true
}

func Message_error (message string) {
  fmt.Printf("FATAL ERROR: %s\n\n", message)
  os.Exit(1)
}

func Message_warning (message string) {
  fmt.Printf("WARNING: %s\n", message)
}

func Set_debug (new_value bool) {
  is_debug = new_value
}

func Is_debug () bool {
  return is_debug
}

func Message_debug (message string) {
  if (Is_debug()) {
    fmt.Printf("DEBUG: %s\n", message)
  }
}
