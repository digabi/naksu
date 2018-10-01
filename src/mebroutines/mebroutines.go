// General routines used by various MEB helper utilities
package mebroutines

import (
  "xlate"

  "fmt"
  "os"
  "os/exec"
  "io"
  "io/ioutil"
  "strings"
  "regexp"
  "bytes"
  "errors"
  "time"
  "encoding/json"

  "github.com/andlabs/ui"
)

var is_debug bool
var main_window *ui.Window
var debug_filename string

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
    matched_timeout,err_timeout := regexp.MatchString("Timed out while waiting for the machine to boot", vagrant_output)
    matched_macaddress,err_macaddress := regexp.MatchString("error: --macaddress: RTGetOpt: Command line option needs argument", vagrant_output)
    if err_timeout == nil && matched_timeout {
      // We've obviously started the VM
      Message_debug("Running vagrant gives me timeout - things are probably ok, complete output:")
      Message_debug(vagrant_output)
    } else if err_macaddress == nil && matched_macaddress {
      // Vagrant in Windows host give this error message - just restart vagrant and you're good
      Message_info(xlate.Get("Server failed to start. This is typical in Windows after an update. Please try again to start the server."))
    } else {
      Message_debug(fmt.Sprintf("Failed to execute %s, complete output:", strings.Join(run_args, " ")))
      Message_debug(vagrant_output)
      Message_warning(fmt.Sprintf(xlate.Get("Failed to execute %s"), strings.Join(run_args, " ")))
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

func Get_vagrantbox_version () string {
  box_index_uuid := Get_vagrantbox_index_uuid()
  if box_index_uuid == "" {
    // We did not get machine index UUID so we cannot return any version string
    return ""
  }

  index_filename := Get_vagrantd_directory() + string(os.PathSeparator) + "data" + string(os.PathSeparator) + "machine-index" + string(os.PathSeparator) + "index"
  file_content, err := ioutil.ReadFile(index_filename)
  if err != nil {
    Message_debug(fmt.Sprintf("Could not read from %s", index_filename))
    return ""
  }

  Message_debug("Vagrant machine-index/index:")
  Message_debug(fmt.Sprintf("%s", file_content))

  // Get version from JSON structure
  var json_data map[string]interface{}

  json_err := json.Unmarshal([]byte(file_content), &json_data)
  if json_err != nil {
    Message_debug("Unable on decode machine-index/index response:")
    Message_debug(fmt.Sprintf("%s", json_err))
    return ""
  }

  if json_data["machines"] == nil {
    return ""
  }
  json_machines := json_data["machines"].(map[string]interface{})

  if json_machines[box_index_uuid] == nil {
    return ""
  }
  json_our_machine := json_machines[box_index_uuid].(map[string]interface{})

  if json_our_machine["extra_data"] == nil {
    return ""
  }
  json_extra_data := json_our_machine["extra_data"].(map[string]interface{})

  if json_extra_data["box"] == nil {
    return ""
  }
  json_box := json_extra_data["box"].(map[string]interface{})

  box_string := fmt.Sprintf("%s (%s %s)", Get_vagrantbox_type(json_box["name"].(string)), json_box["name"], json_box["version"] )

  return box_string
}

func Get_vagrantbox_type (name string) string {
  if name == "" {
    return "-"
  }

  if name == "digabi/ktp-qa" {
    return xlate.Get("Abitti server")
  }

  return xlate.Get("Matric Exam server")
}

func Get_vagrantbox_index_uuid () string {
  path_vagrant := Get_vagrant_directory()

  path_id := path_vagrant + string(os.PathSeparator) + ".vagrant" + string(os.PathSeparator) + "machines" + string(os.PathSeparator) + "default" + string(os.PathSeparator) + "virtualbox" + string(os.PathSeparator) + "index_uuid"

  file_content, err := ioutil.ReadFile(path_id)
  if err != nil {
    Message_debug(fmt.Sprintf(xlate.Get("Could not get vagrantbox index UUID: %d"), err))
    return ""
  }

  Message_debug(fmt.Sprintf("Vagrantbox index UUID is %s", string(file_content)))

  return string(file_content)
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

func CreateFile (path string) error {
  f, err := os.Create(path)
  defer f.Close()
  return err
}

func CopyFile (src, dst string) (err error) {
  Message_debug(fmt.Sprintf("Copying file %s to %s", src, dst))

  if ! ExistsFile(src) {
    Message_debug("Copying failed, could not find source file")
    return errors.New("Could not find source file")
  }

  in, err := os.Open(src)
  if err != nil {
    Message_debug(fmt.Sprintf("Copying failed while opening source file: %v", err))
    return
  }
  defer in.Close()
  out, err := os.Create(dst)
  if err != nil {
    Message_debug(fmt.Sprintf("Copying failed while opening destination file: %v", err))
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

func If_intl_chars_in_path (path string) bool {
  matched,err := regexp.MatchString("[^a-zA-Z0-9_\\-\\/\\:\\\\ ]", path)
  if err == nil && matched {
    return true
  }

  return false
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

func Get_vagrantd_directory () string {
  return Get_home_directory()+string(os.PathSeparator)+".vagrant.d"
}

func Get_mebshare_directory () string {
  return Get_home_directory()+string(os.PathSeparator)+"ktp-jako"
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
  append_logfile(fmt.Sprintf("FATAL ERROR: %s", message))

  // Show libui box if main window has been set with Set_main_window
  if main_window != nil {
    ui.QueueMain(func () {
      ui.MsgBoxError(main_window, xlate.Get("Error"), message)
      ui.Quit()
    })
  } else {
    os.Exit(1)
  }
}

func Message_warning (message string) {
  fmt.Printf("WARNING: %s\n", message)
  append_logfile(fmt.Sprintf("WARNING: %s", message))

  // Show libui box if main window has been set with Set_main_window
  if main_window != nil {
    ui.QueueMain(func () {
      ui.MsgBox(main_window, xlate.Get("Warning"), message)
    })
  }
}

func Message_info (message string) {
  fmt.Printf("INFO: %s\n", message)
  append_logfile(fmt.Sprintf("INFO: %s", message))

  // Show libui box if main window has been set with Set_main_window
  if main_window != nil {
    ui.QueueMain(func () {
      ui.MsgBox(main_window, xlate.Get("Info"), message)
    })
  }
}

func append_logfile (message string) {
  if (debug_filename != "") {
    // Append only if the logfile has been set

    // Current timestamp
    t := time.Now()

    f, err := os.OpenFile(debug_filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660);
    if (err != nil) {
      panic(fmt.Sprintf("Could not append to log file %s: %s", debug_filename, err))
    }

    defer f.Close()

    _, _ = f.WriteString(fmt.Sprintf("[%s] %s\n", t.Format("2006-01-02 15:04:05"), message))
    f.Sync()
    f.Close()
  }
}

func Set_debug (new_value bool) {
  is_debug = new_value
}

func Set_debug_filename (new_filename string) {
  debug_filename = new_filename

  if debug_filename != "" && ExistsFile(debug_filename) {
    // Re-create the log file
    err := os.Remove(debug_filename)
    if (err != nil) {
      panic(fmt.Sprintf("Could not open log file %s: %s", debug_filename, err))
    }
  }
}

func Is_debug () bool {
  return is_debug
}

func Message_debug (message string) {
  if (Is_debug()) {
    fmt.Printf("DEBUG: %s\n", message)
  }

  append_logfile(fmt.Sprintf("DEBUG: %s", message))
}
