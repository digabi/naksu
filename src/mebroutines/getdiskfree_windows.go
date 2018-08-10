package mebroutines

import (
  "fmt"
  "regexp"
  "strconv"
  "errors"
  "os"
)

func Get_disk_free (path string) (int, error) {
  pattern_disk := regexp.MustCompile("^(\\w\\:)")
  pattern_result := pattern_disk.FindStringSubmatch(path)

  if (len(pattern_result) < 2) {
    Message_debug(fmt.Sprintf("Could not detect drive letter from path: %s", path))

    return -1,errors.New("Could not detect drive letter")
  }

  diskletter := pattern_result[1]

  // Note! At the moment this does not return anything i.e. the Windows implementation does not work
  run_params := []string{"wmic", fmt.Sprintf("/node:\"%s\"", os.Getenv("COMPUTERNAME")), "/output:stdout", "logicaldisk", "where", fmt.Sprintf("DeviceID=\"%s\"", diskletter), "get", "FreeSpace"}

  output,err := Run_get_output(run_params)

  Message_debug("wmic says:")
  Message_debug(output)

  if err != nil {
    return -1, errors.New("Get_disk_free() could not detect free disk size as it could not execute wmic")
  }

  fmt.Println(output)

  // Extract server disk image path
  pattern := regexp.MustCompile("(\\d+)")
  result := pattern.FindStringSubmatch(output)

  if (len(result)>1) {
    result_float, _ := strconv.ParseFloat(result[1], 10)
    return int(result_float), nil
  }

  return 1,errors.New("Get_disk_free() could not detect free disk size")
}
