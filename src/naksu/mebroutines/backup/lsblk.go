package backup

import (
	"encoding/json"
	"fmt"

	"naksu/log"
	"naksu/mebroutines"
)

// LsblkOutput represents the parsed JSON output of an lsblk -J command.
type LsblkOutput struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

// BlockDevice represents a single block device entry in the JSON output of an
// lsblk -J command.
type BlockDevice struct {
	Name       string        `json:"name"`
	FileSystem string        `json:"fstype"`
	MountPoint string        `json:"mountpoint"`
	Vendor     string        `json:"vendor"`
	Model      string        `json:"model"`
	HotPlug    interface{}   `json:"hotplug"`
	Children   []BlockDevice `json:"children"`
}

// ListBlockDevices generates a listing of the available block devices using
// lsblk. Usable only on Linux hosts.
func ListBlockDevices() (*LsblkOutput, error) {
	runParams := []string{"lsblk", "-J", "-o", "NAME,FSTYPE,MOUNTPOINT,VENDOR,MODEL,HOTPLUG"}

	lsblkJSON, lsblkErr := mebroutines.RunAndGetOutput(runParams, true)

	log.Debug("lsblk says:")
	log.Debug(lsblkJSON)

	if lsblkErr != nil {
		log.Error("Failed to run lsblk: %v", lsblkErr)

		return &LsblkOutput{BlockDevices: []BlockDevice{}}, lsblkErr
	}

	output, jsonErr := ParseLsblkJSON(lsblkJSON)
	if jsonErr != nil {
		log.Error("Unable to unmarshal lsblk response: %v", jsonErr)

		return &LsblkOutput{BlockDevices: []BlockDevice{}}, lsblkErr
	}

	return output, nil
}

// GetRemovableDisks processes lsblk output to return a listing of connected
// removable devices. The return value is a map from mount points to block
// device names.
func (blk *LsblkOutput) GetRemovableDisks() map[string]string {
	media := map[string]string{}

	if blk.BlockDevices == nil {
		return media
	}

	for blockdeviceIndex := range blk.BlockDevices {
		thisBlockdevice := blk.BlockDevices[blockdeviceIndex]
		if thisBlockdevice.IsRemovable() && thisBlockdevice.Children != nil {
			thisChildren := thisBlockdevice.Children

			for thisChildIndex := range thisChildren {
				thisChild := thisChildren[thisChildIndex]

				thisMountpoint := thisChild.MountPoint
				if thisMountpoint != "" {
					media[thisMountpoint] = fmt.Sprintf("%s, %s", thisBlockdevice.Vendor, thisBlockdevice.Model)
				}
			}
		}
	}

	return media
}

// IsRemovable returns true if the block device is a removable device (e.g. a
// USB drive) and false otherwise.
func (bd *BlockDevice) IsRemovable() bool {
	switch hotplug := bd.HotPlug.(type) {
	case bool:
		return hotplug
	case string:
		switch hotplug {
		case "0":
			return false
		case "1":
			return true
		}
	}

	log.Error("Unknown format for hotplug field in lsblk output: %v", bd.HotPlug)

	return false
}

// ParseLsblkJSON parses JSON output from an lsblk -J command.
func ParseLsblkJSON(lsblkJSON string) (*LsblkOutput, error) {
	var jsonData LsblkOutput
	jsonErr := json.Unmarshal([]byte(lsblkJSON), &jsonData)

	return &jsonData, jsonErr
}
