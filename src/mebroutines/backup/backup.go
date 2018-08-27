package backup

import (
  "xlate"

  "os"
  "fmt"
  "io/ioutil"
  "regexp"
  "time"

  "mebroutines"
)

func Do_make_backup (path_backup string) {
  // Get box
  box_id := get_vagrantbox_id()
  mebroutines.Message_debug(fmt.Sprintf("Vagrantbox ID: %s", box_id))

  // Get disk UUID
  disk_uuid := get_disk_uuid(box_id)
  mebroutines.Message_debug(fmt.Sprintf("Disk UUID: %s", disk_uuid))

  // Make clone to path_backup
  if (mebroutines.ExistsFile(path_backup)) {
    mebroutines.Message_error(fmt.Sprintf(xlate.Get("File %s already exists"), path_backup))
  }
  make_clone(disk_uuid, path_backup)

  // Close backup media (detach it from VirtualBox disk management)
  delete_clone(path_backup)

  mebroutines.Message_info(fmt.Sprintf(xlate.Get("Backup has been made to %s"), path_backup))
}

func get_vagrantbox_id () string {
  path_vagrant := mebroutines.Get_vagrant_directory()

  path_id := path_vagrant + string(os.PathSeparator) + ".vagrant" + string(os.PathSeparator) + "machines" + string(os.PathSeparator) + "default" + string(os.PathSeparator) + "virtualbox" + string(os.PathSeparator) + "id"

  file_content, err := ioutil.ReadFile(path_id)
  if err != nil {
    mebroutines.Message_error(fmt.Sprintf(xlate.Get("Could not get vagrantbox ID: %d"), err))
  }

  return string(file_content)
}

func get_disk_uuid(box_id string) string {
  vboxmanage_output := mebroutines.Run_vboxmanage([]string{"showvminfo", "-machinereadable", box_id})

  // Extract server disk image path
  pattern := regexp.MustCompile("\"SATA Controller-ImageUUID-0-0\"=\"(.*?)\"")
  result := pattern.FindStringSubmatch(vboxmanage_output)

  if (len(result)>1) {
    return result[1]
  }

  // No match
  mebroutines.Message_debug(vboxmanage_output)
  mebroutines.Message_error(xlate.Get("Could not make backup: failed to get disk UUID"))

  return ""
}

func make_clone(disk_uuid string, path_backup string) {
  vboxmanage_output := mebroutines.Run_vboxmanage([]string{"clonemedium", disk_uuid, path_backup})

  // Check whether clone was successful or not
  matched,err_re := regexp.MatchString("Clone medium created in format 'VMDK'", vboxmanage_output)
  if err_re != nil || !matched {
    // Failure
    mebroutines.Message_error(fmt.Sprintf(xlate.Get("Could not back up disk %s to %s"), disk_uuid, path_backup))
  }
}

func delete_clone(path_backup string) {
  _ = mebroutines.Run_vboxmanage([]string{"closemedium", path_backup})
}

func Get_backup_filename (timestamp time.Time) string {
  return timestamp.Format("2006-01-02_15-04-05.vmdk")
}
