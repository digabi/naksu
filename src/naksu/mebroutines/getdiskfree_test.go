package mebroutines_test

import (
	"naksu/mebroutines"
	"testing"
)

func TestExtractDiskFreeDarwin(t *testing.T) {
	tables := []struct {
		dfOutput string
		freeSize uint64
		isError  bool
	}{
		{`Filesystem 512-blocks       Used Available Capacity iused      ifree %iused  Mounted on
  /dev/disk1 1951662080 1365304512 585845568    70% 7050261 4287917018    0%   /`, 585845568 * 512, false},
	}

	for _, table := range tables {
		freeSize, err := mebroutines.ExtractDiskFreeDarwin(table.dfOutput)
		if freeSize != table.freeSize {
			t.Errorf("ExtractDiskFreeDarwin gives %d instead of %d for df output '%s'", freeSize, table.freeSize, table.dfOutput)
		}

		hasError := err != nil

		if hasError != table.isError {
			t.Errorf("ExtractDiskFreeDarwin error mismatch for df output '%s'", table.dfOutput)
		}
	}
}

func TestExtractDiskFreeLinux(t *testing.T) {
	tables := []struct {
		dfOutput string
		freeSize uint64
		isError  bool
	}{
		{`  Avail
4981668`, 4981668 * 1024, false},
		// This should fail as there is only the header
		{`Avail`, 0, true},
	}

	for _, table := range tables {
		freeSize, err := mebroutines.ExtractDiskFreeLinux(table.dfOutput)
		if freeSize != table.freeSize {
			t.Errorf("ExtractDiskFreeLinux gives %d instead of %d for df output '%s'", freeSize, table.freeSize, table.dfOutput)
		}

		hasError := err != nil

		if hasError != table.isError {
			t.Errorf("ExtractDiskFreeDarwin error mismatch for df output '%s'", table.dfOutput)
		}
	}
}

func TestExtractDiskFreeWindows(t *testing.T) {
	tables := []struct {
		wmiData  mebroutines.Win32_LogicalDisk
		freeSize uint64
		isError  bool
	}{
		{mebroutines.Win32_LogicalDisk{Size: 1234, FreeSpace: 1024, DeviceID: "C:"}, 1024, false},
		// 1024 Gb drive with 763something Gb free
		{mebroutines.Win32_LogicalDisk{Size: 1024000000000, FreeSpace: 763123736200, DeviceID: "D:"}, 763123736200, false},
	}

	for _, table := range tables {
		var wmiDataArray []mebroutines.Win32_LogicalDisk
		wmiDataArray = append(wmiDataArray, table.wmiData)
		freeSize, err := mebroutines.ExtractDiskFreeWindows(wmiDataArray)
		if freeSize != table.freeSize {
			t.Errorf("ExtractDiskFreeWindows gives %d instead of %d for DeviceID '%s'", freeSize, table.freeSize, table.wmiData.DeviceID)
		}

		hasError := err != nil

		if hasError != table.isError {
			t.Errorf("ExtractDiskFreeDarwin error mismatch for DeviceID '%s'", table.wmiData.DeviceID)
		}
	}
}
