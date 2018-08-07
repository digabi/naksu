package backup

import (
  "os"
  "fmt"
  "encoding/csv"
  "strings"

  "mebroutines"
)

func Get_backup_media () map[string]string {
  media := get_backup_media_windows()

  if (os.Getenv("TEMP") != "") {
    media[os.Getenv("TEMP")] = "Temporary files"
  }
  if (os.Getenv("USERPROFILE") != "") {
    media[os.Getenv("USERPROFILE")] = "Profile directory"
  }

  return media
}

func get_backup_media_windows () map[string]string {
  var media = map[string]string{}

  run_params := []string{"wmic", "/output:stdout", "logicaldisk", "get", "DriveType,DeviceID,Description,VolumeName", "/format:csv"}

  wmic_output,wmic_err := mebroutines.Run_get_output(run_params)

  mebroutines.Message_debug("wmic says:")
  mebroutines.Message_debug(wmic_output)

  if wmic_err != nil {
    mebroutines.Message_debug("Failed to run wmic")
    // Return empty set of media
    return media
  }

  r := csv.NewReader(strings.NewReader(wmic_output))
  // Disable fields number checking
  r.FieldsPerRecord = -1

  media_records,reader_err := r.ReadAll()
  if reader_err != nil {
    mebroutines.Message_debug("Failed to process wmic csv output:")
    mebroutines.Message_debug(fmt.Sprintf("%s", reader_err))
    // Return empty set of media
    return media
  }

  for _, this_record := range media_records {
    // Select only lines with 5 columns and DriveType=2 (removable drive) or DriveType=3 (local drive)
    if len(this_record) == 5 && (this_record[3] == "2" || this_record[3] == "3") {
      mebroutines.Message_debug(fmt.Sprintf("wmic csv record: %s", strings.Join(this_record,", ")))
      this_path := fmt.Sprintf("%s\\", this_record[2])
      media[this_path] = fmt.Sprintf("%s, %s", this_record[1], this_record[4])
    } else {
      mebroutines.Message_debug(fmt.Sprintf("Skipping wmic csv record: %s", strings.Join(this_record,", ")))
    }
  }

  return media
}
