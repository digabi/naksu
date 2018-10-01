package mebroutines

import (
  "fmt"
)

func Open_meb_share () {
  meb_share_path := Get_mebshare_directory()

  Message_debug(fmt.Sprintf("MEB share directory: %s", meb_share_path))

  if ! ExistsDir(meb_share_path) {
    Message_warning("Cannot open MEB share directory since it does not exist")
    return
  }

  run_params := []string{"explorer", meb_share_path}

  // For some not-obvious reason Run_get_output() results err
  output,err := Run_get_error(run_params)

  if err != nil {
    err_str := fmt.Sprintf("%v", err)
    // Opening explorer results exit code 1
    if err_str != "exit status 1" {
      Message_warning("Could not open MEB share directory")
      Message_debug(fmt.Sprintf("Could not open MEB share directory: %v", err))
    }
  }

  Message_debug("MEB share directory open output:")
  Message_debug(output)
}
