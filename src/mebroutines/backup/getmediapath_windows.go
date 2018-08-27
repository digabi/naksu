package backup

import (
  "os"
  "fmt"

  "github.com/StackExchange/wmi"

  "mebroutines"
  "xlate"
)

func Get_backup_media () map[string]string {
  media := get_backup_media_windows()

  if (os.Getenv("TEMP") != "") {
    media[os.Getenv("TEMP")] = xlate.Get("Temporary files")
  }
  if (os.Getenv("USERPROFILE") != "") {
    media[os.Getenv("USERPROFILE")] = xlate.Get("Profile directory")
  }

  return media
}

func get_backup_media_windows () map[string]string {
  type Win32_LogicalDisk struct {
    DeviceID string
    DriveType int
    Description string
    VolumeName string
  }

  var media = map[string]string{}

  var dst []Win32_LogicalDisk
  query := wmi.CreateQuery(&dst, "")
  err := wmi.Query(query, &dst);
  if err != nil {
    mebroutines.Message_debug("get_backup_media_windows() could not detect removable/hard drives as it could not query WMI")
    mebroutines.Message_debug(fmt.Sprint(err))
    return media
  }

  for this_drive := range dst {
    if (dst[this_drive].DriveType == 2 || dst[this_drive].DriveType == 3) {
      // We have either hard or removable drive
      this_path := fmt.Sprintf("%s\\", dst[this_drive].DeviceID)
      media[this_path] = fmt.Sprintf("%s, %s", dst[this_drive].VolumeName, dst[this_drive].Description)
    }
  }

  return media
}
