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

  run_params := []string{"explorer.exe", meb_share_path}

  output,err := Run_get_output(run_params)

  if err != nil {
    Message_warning("Could not open MEB share directory")
  }

  Message_debug("MEB share directory open output:")
  Message_debug(output)
}
