package mebroutines

import (
  "os"
)

func get_vboxmanage_path () string {
  var path = "VBoxManage"
	if (os.Getenv("VBOXMANAGEPATH") != "") {
		path = os.Getenv("VBOXMANAGEPATH")
	} else {
    var path_virtualbox = os.Getenv("VBOX_MSI_INSTALL_PATH")
    if (path_virtualbox != "") {
      path = path_virtualbox+string(os.PathSeparator)+path
    }
  }

  return path
}
