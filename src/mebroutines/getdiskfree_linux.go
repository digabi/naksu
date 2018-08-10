package mebroutines

import (
  "regexp"
  "strconv"
)

func Get_disk_free (path string) (int, error) {
  run_params := []string{"df", "--output=avail", path}

  output,err := Run_get_output(run_params)

  if err != nil {
    return -1, err
  }

  // Extract server disk image path
  pattern := regexp.MustCompile("(\\d+)")
  result := pattern.FindStringSubmatch(output)

  if (len(result)>1) {
    result_float, _ := strconv.ParseFloat(result[1], 10)
    return int(result_float), nil
  }

  return -1,nil
}
