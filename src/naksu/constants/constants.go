package constants

const (
  // LowDiskLimit sets the warning level of low disk
	LowDiskLimit = 50 * 1024 * 1024 // 50 Gb

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
	VagrantBoxAvailVersionDetailsCacheTimeout int64 = 5*60
)
