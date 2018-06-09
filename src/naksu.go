package main

import (
  "bufio"
  "os"
  "fmt"

  "mebroutines"
  "mebroutines/install"
  "mebroutines/start"
)

const version = "0.1.0"

func main() {
  var selection string = ""

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
