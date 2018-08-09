package mebroutines

import (
  "os"
  "fmt"
)

func ExistsStdin () bool {
  fi, err := os.Stdin.Stat()
  if err != nil {
    Message_debug(fmt.Sprintf("Checking status for Stdin results error: %s", fmt.Sprint(err)))
    return false
  }

  stat := fmt.Sprintf("%v", fi.Mode())

  Message_debug(fmt.Sprintf("STDIN stat is: %s", stat))

  if (stat == "Dcrw-rw-rw-") {
    // No stdin
    Message_debug("No STDIN detected - naksu is executed without a terminal")
    return false
  }

  Message_debug("STDIN detected - naksu is executed from a terminal")

  return true
}
