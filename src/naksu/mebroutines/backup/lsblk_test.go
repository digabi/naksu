package backup_test

import (
	"testing"

	"naksu/mebroutines/backup"
)

// The old style JSON output from lsblk represents boolean values as string
// values that can be either "1" or "0". This variant is used by the lsblk
// version in Ubuntu 18.04.
var oldStyleJSON = `
{
   "blockdevices": [
      {"name": "loop0", "fstype": "squashfs", "mountpoint": "/snap/core/4917", "vendor": null, "model": null, "hotplug": "0"},
      {"name": "sda", "fstype": null, "mountpoint": null, "vendor": "ATA     ", "model": "TS128GSSD370S   ", "hotplug": "0",
         "children": [
            {"name": "sda1", "fstype": "ext4", "mountpoint": "/", "vendor": null, "model": null, "hotplug": "0"}
         ]
      },
      {"name": "sdb", "fstype": null, "mountpoint": null, "vendor": "YTL", "model": "SSD", "hotplug": "1",
         "children": [
            {"name": "sdb1", "fstype": "ntfs", "mountpoint": "/media/abitti/BACKUP", "vendor": null, "model": null, "hotplug": "1"}
         ]
      },
      {"name": "sr0", "fstype": null, "mountpoint": null, "vendor": "DVD Drive Vendor", "model": "12345F", "hotplug": "1"}
   ]
}
`

// The new style JSON output from lsblk represents boolean values as actual
// JSON booleans. This is used by later versions of lsblk.
var newStyleJSON = `
{
   "blockdevices": [
      {"name":"loop0", "fstype":"squashfs", "mountpoint":"/snap/core/4917", "vendor":null, "model":null, "hotplug":false},
      {"name":"sda", "fstype":null, "mountpoint":null, "vendor":"ATA     ", "model":"TS128GSSD370S   ", "hotplug":false,
         "children": [
            {"name":"sda1", "fstype":"ext4", "mountpoint":"/", "vendor":null, "model":null, "hotplug":false}
         ]
      },
      {"name":"sdb", "fstype":null, "mountpoint":null, "vendor":"YTL", "model":"SSD", "hotplug":true,
         "children": [
            {"name":"sdb1", "fstype":"ntfs", "mountpoint":"/media/abitti/BACKUP", "vendor":null, "model":null, "hotplug":true}
         ]
      },
      {"name":"sr0", "fstype":null, "mountpoint":null, "vendor":"DVD Drive Vendor", "model":"12345F", "hotplug":true}
   ]
}
`

func TestOldStyleLsblkParsing(t *testing.T) {
	lsblk, err := backup.ParseLsblkJSON(oldStyleJSON)
	checkOutput(t, lsblk, err)
}

func TestNewStyleLsblkParsing(t *testing.T) {
	lsblk, err := backup.ParseLsblkJSON(newStyleJSON)
	checkOutput(t, lsblk, err)
}

func checkOutput(t *testing.T, lsblk *backup.LsblkOutput, err error) {
	if err != nil {
		t.Errorf("Error parsing lsblk JSON: %s", err.Error())
	}

	disks := lsblk.GetRemovableDisks()
	if diskName, ok := disks["/media/abitti/BACKUP"]; ok {
		if diskName != "YTL, SSD" {
			t.Errorf("Expected backup disk name to be 'YTL, SSD' but was '%s'", diskName)
		}
	} else {
		t.Errorf("/media/abitti/BACKUP was not found in the removable disks listing")
	}

	if _, ok := disks["/"]; ok {
		t.Errorf("/ was listed as a removable device when it shouldn't have been")
	}

	if len(disks) > 1 {
		t.Errorf("Expected only one removable device to be listed, but there were %d", len(disks))
	}
}
