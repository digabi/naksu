package mebroutines

import (
  "fmt"
  "regexp"
  "errors"

  "github.com/StackExchange/wmi"
)

type Win32_LogicalDisk struct {
  Size int
  FreeSpace int
  DeviceID string
}

func Get_disk_free (path string) (int, error) {
  pattern_disk := regexp.MustCompile("^(\\w\\:)")
  pattern_result := pattern_disk.FindStringSubmatch(path)

  if (len(pattern_result) < 2) {
    Message_debug(fmt.Sprintf("Could not detect drive letter from path: %s", path))

    return -1,errors.New("Could not detect drive letter")
  }

  diskletter := pattern_result[1]

  var dst []Win32_LogicalDisk
  query := wmi.CreateQuery(&dst, fmt.Sprintf("WHERE DeviceID=\"%s\"", diskletter))
  err := wmi.Query(query, &dst);
  if err != nil {
    Message_debug(fmt.Sprintf("Get_disk_free() could not make WMI query: %s", fmt.Sprint(err)))
    return -1, errors.New("Get_disk_free() could not detect free disk size as it could not query WMI")
  }

  if (len(dst) > 0) {
    free_space := dst[0].FreeSpace / 1000
    Message_debug(fmt.Sprintf("Disk free for path %s: %d", path, free_space))
    return free_space, nil
  }

  return -1,errors.New("Get_disk_free() could not detect free disk size")
}
