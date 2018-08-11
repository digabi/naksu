package backup

import (
  "os"
  "fmt"
  "encoding/json"

  "mebroutines"
)

func Get_backup_media () map[string]string {
  media := get_backup_media_linux()

  // Add some entries from environment variables
  if (os.Getenv("HOME") != "") {
    media[os.Getenv("HOME")] = "Home directory"
  }
  if (os.Getenv("TEMP") != "") {
    media[os.Getenv("TEMP")] = "Temporary files"
  }

  return media
}

func get_backup_media_linux () map[string]string {
  var media = map[string]string{}

  run_params := []string{"lsblk", "-J", "-o", "NAME,FSTYPE,MOUNTPOINT,VENDOR,MODEL,HOTPLUG"}

  lsblk_json,lsblk_err := mebroutines.Run_get_output(run_params)

  mebroutines.Message_debug("lsblk says:")
  mebroutines.Message_debug(lsblk_json)

  if lsblk_err != nil {
    mebroutines.Message_debug("Failed to run lsblk")
    // Return empty set of media
    return media
  }

  var json_data map[string]interface{}

  json_err := json.Unmarshal([]byte(lsblk_json), &json_data)
  if json_err != nil {
    mebroutines.Message_debug("Unable on decode lsblk response:")
    mebroutines.Message_debug(fmt.Sprintf("%s", json_err))
    // Return empty set of media
    return media
  }

  blockdevices := json_data["blockdevices"].([]interface{})

  media = get_removable_disks(blockdevices)

  return media
}


func get_removable_disks (blockdevices []interface{}) map[string]string {
  var media = map[string]string{}
  //media_n := 0

  if blockdevices == nil {
    return media
  }

  for blockdevice_n := range blockdevices {
    //fmt.Println(blockdevices[blockdevice_n])
    this_blockdevice := blockdevices[blockdevice_n].(map[string]interface{})
    if device_field_string(this_blockdevice["hotplug"]) == "1" && this_blockdevice["children"] != nil {
      this_children := this_blockdevice["children"].([]interface{})

      for this_child_n := range this_children {
        this_child := this_children[this_child_n].(map[string]interface{})

        this_mountpoint := device_field_string(this_child["mountpoint"])
        if this_mountpoint != "" {
          media[this_mountpoint] = fmt.Sprintf("%s, %s", device_field_string(this_blockdevice["vendor"]), device_field_string(this_blockdevice["model"]))
        }
      }
    }
  }

  return media
}

func device_field_string (this_field interface{}) string {
  if this_field == nil {
    return ""
  }

  return this_field.(string)
}
