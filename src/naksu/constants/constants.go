package constants

const (
	// LowDiskLimit sets the warning level of low disk (in bytes)
	LowDiskLimit uint64 = 50 * 1024 * 1024 * 1024 // 50 Gb

	// AbittiVagrantURL is the URL for the Abittti Vagrantfile
	AbittiVagrantURL = "http://static.abitti.fi/usbimg/qa/vagrant/Vagrantfile"

	// URLTest is a testing URL for network connectivity (network.CheckIfNetworkAvailable).
	// Point this to something ultra-stable
	URLTest = "http://static.abitti.fi/usbimg/qa/latest.txt"

	// URLTestTimeout is the timeout in seconds for the test above
	URLTestTimeout = 4

	// VagrantBoxAvailVersionDetailsCacheTimeout is a timeout for vagrant box version
	// cache. See naksu/boxversion GetVagrantBoxAvailVersionDetails() for more
	// In seconds (5 minutes)
	VagrantBoxAvailVersionDetailsCacheTimeout int64 = 5 * 60

	// VBoxManageCacheTimeout is a timeout for executing VBoxManage showvminfo
	// See naksu/box getVMInfoRegexp()
	VBoxManageCacheTimeout int64 = 15
)

// AvailableSelection is a struct for a UI/configuration option
type AvailableSelection struct {
	ConfigValue string
	Legend      string
}

// AvailableLangs is an array of possible language selection values.
// The first value is the default.
var AvailableLangs = []AvailableSelection{
	AvailableSelection{
		ConfigValue: "fi",
		Legend:      "Suomeksi",
	},
	AvailableSelection{
		ConfigValue: "sv",
		Legend:      "På svenska",
	},
	AvailableSelection{
		ConfigValue: "en",
		Legend:      "In English",
	},
}

// AvailableNics is an array of possible NIC selection values.
// The first value is the default.
var AvailableNics = []AvailableSelection{
	AvailableSelection{
		ConfigValue: "virtio",
		Legend:      "virtio",
	},
	AvailableSelection{
		ConfigValue: "Am79C970A",
		Legend:      "Am79C970A",
	},
	AvailableSelection{
		ConfigValue: "Am79C973",
		Legend:      "Am79C973",
	},
	AvailableSelection{
		ConfigValue: "82540EM",
		Legend:      "82540EM",
	},
	AvailableSelection{
		ConfigValue: "82543GC",
		Legend:      "82543GC",
	},
	AvailableSelection{
		ConfigValue: "82545EM",
		Legend:      "82545EM",
	},
}

// GetAvailableSelectionID returns array id for a given ConfigValue
// in the given set of choices. Returns -1 if the configValue was not found.
func GetAvailableSelectionID(configValue string, choices []AvailableSelection) int {
	for i, thisChoice := range choices {
		if thisChoice.ConfigValue == configValue {
			return i
		}
	}

	return -1
}
