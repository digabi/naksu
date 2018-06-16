package main

// required by selfupdate (needs context)
// +build go1.7

import (
  "bufio"
  "os"
  "fmt"

  "github.com/blang/semver"
  "github.com/rhysd/go-github-selfupdate/selfupdate"

  "mebroutines"
  "mebroutines/install"
  "mebroutines/start"
)

const version = "0.9.0"

func doSelfUpdate() {
  v := semver.MustParse(version)
  latest, err := selfupdate.UpdateSelf(v, "digabi/naksu")
  if err != nil {
    mebroutines.Message_warning(fmt.Sprintf("Binary update failed: %s", err))
    return
  }
  if latest.Version.Equals(v) {
    // latest version is the same as current version. It means current binary is up to date.
    mebroutines.Message_debug(fmt.Sprintf("Current binary is the latest version: %s", version))
  } else {
    mebroutines.Message_debug(fmt.Sprintf("Successfully updated to version: %s", latest.Version))
    //log.Println("Release note:\n", latest.ReleaseNotes)
  }
}

func main() {
  var selection string = ""

  doSelfUpdate()

  Askinput:

  fmt.Println("Hi! I'm Naksu "+version)
  fmt.Println("")
  fmt.Println("Choose action and press Enter:")
  fmt.Println("1) Install new or update existing Stickless Exam Server")
  fmt.Println("2) Start Stickless Exam Server")
  fmt.Println("X) Exit")
  fmt.Println("")
  fmt.Printf("Your choice (1-2 or X): ")

  reader := bufio.NewReader(os.Stdin)
  selection, _ = reader.ReadString('\n')

  selection_stripped := selection[:len(selection)-1]

  if (selection_stripped == "1") {
    mebroutines.Message_debug("Now executing install package")
    install.Do_get_server()
  } else if (selection_stripped == "2") {
    mebroutines.Message_debug("Now executing start package")
    start.Do_start_server()
  } else if (selection_stripped == "x" || selection_stripped == "X") {
    mebroutines.Message_debug("Exit")
  } else {
    goto Askinput
  }

}
