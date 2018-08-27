// General routines used by various MEB helper utilities
package mebroutines

import (
  "xlate"

  "fmt"
  "os"
  "os/exec"
  "io"
  "strings"
  "regexp"
  "bytes"
  "errors"

  "github.com/andlabs/ui"
)

var is_debug bool
var main_window *ui.Window

func Run(command_args []string) error {
	Message_debug(fmt.Sprintf("run: %s", strings.Join(command_args, " ")))
	cmd := exec.Command(command_args[0], command_args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
  err := cmd.Run()
	if err != nil {
		Message_warning(fmt.Sprintf(xlate.Get("command failed: %s"), strings.Join(command_args, " ")))
	}

  return err
}

func Run_get_output (command_args []string) (string, error) {
  Message_debug(fmt.Sprintf("Run_get_output: %s", strings.Join(command_args, " ")))
	out,err := exec.Command(command_args[0], command_args[1:]...).CombinedOutput()
  if (err != nil) {
    // Executing failed, return error condition
    Message_warning(fmt.Sprintf(xlate.Get("command failed: %s"), strings.Join(command_args, " ")))
    return string(out), err
  }

  if (out != nil) {
    Message_debug("Run_get_output returns combined STDOUT and STDERR:")
    Message_debug(string(out))
  } else {
    Message_debug("Run_get_output returned NIL as combined STDOUT and STDERR")
  }

  return string(out), nil
}

func Run_get_error (command_args []string) (string, error) {
  Message_debug(fmt.Sprintf("Run_get_error: %s", strings.Join(command_args, " ")))

  var stderr bytes.Buffer

  cmd := exec.Command(command_args[0], command_args[1:]...)
  cmd.Stdout = os.Stdout
  cmd.Stdin = os.Stdin
  cmd.Stderr = &stderr

  err := cmd.Run()

  Message_debug("Run_get_error returns STDERR:")
  Message_debug(stderr.String())

  return stderr.String(), err
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
  vagrant_output,err := Run_get_error(run_args)
  if (err != nil) {
    matched,err_re := regexp.MatchString("Timed out while waiting for the machine to boot", vagrant_output)
    if err_re == nil && matched {
      // We've obviously started the VM
      Message_debug("Running vagrant gives me timeout - things are probably ok, complete output:")
      Message_debug(vagrant_output)
    } else {
      Message_debug(fmt.Sprintf("Failed to execute %s, complete output:", strings.Join(run_args, " ")))
      Message_debug(vagrant_output)
      Message_error(fmt.Sprintf(xlate.Get("Failed to execute %s"), strings.Join(run_args, " ")))
    }
  }
}

func Run_vboxmanage (args []string) string {
  vboxmanagepath_arr := []string{get_vboxmanage_path()}
  run_args := append(vboxmanagepath_arr, args...)
  vboxmanage_output,err := Run_get_output(run_args)
  if (err != nil) {
    Message_debug(fmt.Sprintf("Failed to execute %s, complete output:", strings.Join(run_args, " ")))
    Message_debug(vboxmanage_output)
    Message_error(fmt.Sprintf(xlate.Get("Failed to execute %s"), strings.Join(run_args, " ")))
  }

  return vboxmanage_output
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

func CopyFile (src, dst string) (err error) {
  if ! ExistsFile(src) {
    return errors.New("Could not find source file")
  }

  in, err := os.Open(src)
  if err != nil {
    return
  }
  defer in.Close()
  out, err := os.Create(dst)
  if err != nil {
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

func Get_vagrant_directory () string {
  return Get_home_directory()+string(os.PathSeparator)+"ktp"
}

func Chdir_vagrant_directory () bool {
  path_vagrant := Get_vagrant_directory()
  Message_debug(fmt.Sprintf("chdir %s", path_vagrant))
  err := os.Chdir(path_vagrant)
  if (err != nil) {
    Message_warning(fmt.Sprintf(xlate.Get("Could not chdir to %s"), path_vagrant))
    return false
  }

  return true
}

func Set_main_window (win *ui.Window) {
  // Set libui main window pointer used by Message_error and Message_warning
  main_window = win
}

func Message_error (message string) {
  fmt.Printf("FATAL ERROR: %s\n\n", message)

  // Show libui box if main window has been set with Set_main_window
  if main_window != nil {
    ui.MsgBoxError(main_window, xlate.Get("Error"), message)
    ui.Quit()
  } else {
    os.Exit(1)
  }
}

func Message_warning (message string) {
  fmt.Printf("WARNING: %s\n", message)

  // Show libui box if main window has been set with Set_main_window
  if main_window != nil {
    ui.MsgBox(main_window, xlate.Get("Warning"), message)
  }
}

func Message_info (message string) {
  fmt.Printf("INFO: %s\n", message)

  // Show libui box if main window has been set with Set_main_window
  if main_window != nil {
    ui.MsgBox(main_window, xlate.Get("Info"), message)
  }
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
