package mebroutines

import (
  "os"
)

func get_vboxmanage_path () string {
  var path = "VBoxManage"
	if (os.Getenv("VBOXMANAGEPATH") != "") {
		path = os.Getenv("VBOXMANAGEPATH")
	}

  return path
}
