package constants

import "time"

const (
	// LowDiskLimit sets the warning level of low disk (in bytes)
	LowDiskLimit uint64 = 50 * 1024 * 1024 * 1024 // 50 Gb

	// AbittiEtcherURL is the URL for the latest Abitti Etcher zip
	AbittiEtcherURL  = "http://static.abitti.fi/etcher-usb/ktp-etcher.zip"
	AbittiVersionURL = "http://static.abitti.fi/etcher-usb/ktp-etcher.ver"
	AbittiBoxType    = "abitti"

	// ExamEtcherURL is the URL for an Exam Etcher zip
	ExamEtcherURL  = "http://static.abitti.fi/etcher-usb/releases/###PASSPHRASEHASH###/ktp-etcher.zip"
	ExamVersionURL = "http://static.abitti.fi/etcher-usb/releases/###PASSPHRASEHASH###/ktp-etcher.ver"
	ExamBoxType    = "exam"

	// URLTest is a testing URL for network connectivity (network.CheckIfNetworkAvailable).
	// Point this to something ultra-stable
	URLTest = "http://static.abitti.fi/usbimg/qa/latest.txt"

	// URLTestTimeout is the timeout in seconds for the test above
	URLTestTimeout = 4

	// VBoxManageCacheTimeout is a timeout for VBoxManage cache
	// See naksu/box
	VBoxManageCacheTimeout = 30 * time.Second

	// VBoxRunningCacheTimeout is a timeout for VM state cache
	// See Running() at naksu/box
	VBoxRunningCacheTimeout = 2 * time.Second

	// CloudStatusTimeout is a timeout for cloud status cache
	// See naksu/cloud
	CloudStatusTimeout = 10 * time.Minute

	// LogCopyRequestFilename is for requesting logs from ktp
	LogCopyRequestFilename = "_log_copy_requested"
	// LogCopyDoneFilename is for detecting when log request is done
	LogCopyDoneFilename = "_log_copy_done"
	// LogCopyStatusFilename is for progress info on log copying
	LogCopyStatusFilename = "_log_copy_status"

	// LogRequestTimeout is the timeout for log request from ktp
	LogRequestTimeout = 1 * time.Minute
)

// AvailableSelection is a struct for a UI/configuration option
type AvailableSelection struct {
	ConfigValue string
	Legend      string
}

// AvailableLangs is an array of possible language selection values.
// The first value is the default.
var AvailableLangs = []AvailableSelection{
	{
		ConfigValue: "fi",
		Legend:      "Suomeksi",
	},
	{
		ConfigValue: "sv",
		Legend:      "På svenska",
	},
	{
		ConfigValue: "en",
		Legend:      "In English",
	},
}

// AvailableNics is an array of possible NIC selection values.
// The first value is the default.
var AvailableNics = []AvailableSelection{
	{
		ConfigValue: "virtio",
		Legend:      "virtio",
	},
	{
		ConfigValue: "Am79C970A",
		Legend:      "Am79C970A",
	},
	{
		ConfigValue: "Am79C973",
		Legend:      "Am79C973",
	},
	{
		ConfigValue: "82540EM",
		Legend:      "82540EM",
	},
	{
		ConfigValue: "82543GC",
		Legend:      "82543GC",
	},
	{
		ConfigValue: "82545EM",
		Legend:      "82545EM",
	},
}

// DefaultExtNicArray is an array holding the default EXTNIC value
var DefaultExtNicArray = []AvailableSelection{
	{
		ConfigValue: "",
		Legend:      "Select network device",
	},
}

// GetAvailableSelectionID returns array id for a given ConfigValue
// in the given set of choices. Returns -1 if the configValue was not found.
func GetAvailableSelectionID(configValue string, choices []AvailableSelection, valueIfNotFound int) int {
	for i, thisChoice := range choices {
		if thisChoice.ConfigValue == configValue {
			return i
		}
	}

	return valueIfNotFound
}
