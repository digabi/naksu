package network

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"naksu/constants"
	"naksu/log"
)

// CheckIfNetworkAvailable tests if a pre-set utterly-reliable network setver responds to HTTP GET
func CheckIfNetworkAvailable() bool {
	return testHTTPGet(constants.URLTest, constants.URLTestTimeout)
}

// testHTTPGet tests whether HTTP get succeeds to given URL in given timeout (seconds)
func testHTTPGet(url string, timeout int) bool {
	// Set timeout for a HTTP client
	timeoutDuration := time.Duration(timeout) * time.Second
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       timeoutDuration,
	}

	ctx := context.Background()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error("Testing HTTP GET %s and got error %v when creating request", url, err)

		return false
	}

	response, err := client.Do(request)
	if err != nil {
		return false
	}

	defer response.Body.Close()

	const httpStatusOK = 200
	if response.StatusCode != httpStatusOK {
		log.Error("Testing HTTP GET %s and got HTTP status code %d which is not OK", url, response.StatusCode)

		return false
	}

	return true
}

// StartEnvironmentStatusUpdate starts periodically updating given
// environmentStatus.NetAvailable value
func StartEnvironmentStatusUpdate(environmentStatus *constants.EnvironmentStatus, tickerDuration time.Duration) {
	ticker := time.NewTicker(tickerDuration)

	go func() {
		for {
			<-ticker.C
			environmentStatus.NetAvailable = CheckIfNetworkAvailable()
		}
	}()
}

// IsExtInterface returns true if given interfaceName is a name for a valid interface.
func IsExtInterface(interfaceName string) bool {
	interfaces := GetExtInterfaces()

	result := false

	for _, thisInterface := range interfaces {
		if thisInterface.ConfigValue == interfaceName {
			result = true
		}
	}

	return result
}

// isIgnoredExtInterface returns true if given system-level network device should
// be ignored (i.e. not to be shown to the user).
func isIgnoredExtInterface(interfaceName string, ignoredExtNics []string) bool {
	for _, ignoredExtNic := range ignoredExtNics {
		match, err := regexp.MatchString(ignoredExtNic, interfaceName)
		if err == nil && match {
			return true
		}
	}

	return false
}

func isIgnoredExtInterfaceLinux(interfaceName string) bool {
	return isIgnoredExtInterface(interfaceName, []string{
		"^lo$",
		"^tap\\d",
		"^vboxnet\\d",
		"^virbr",
	})
}

func isIgnoredExtInterfaceWindows(interfaceName string) bool {
	return isIgnoredExtInterface(interfaceName, []string{
		"^VirtualBox Host-Only Ethernet Adapter$",
	})
}

func isIgnoredExtInterfaceDarwin(interfaceName string) bool {
	return isIgnoredExtInterface(interfaceName, []string{
		"^lo\\d*$",
		"^gif\\d+$",
		"^XHC\\d+$",
		"^awdl\\d+$",
		"^utun\\d+$",
		"^bridge\\d+$",
		"^stf\\d+$",
		"^vboxnet\\d*",
	})
}

func bpsToMbps(bps uint64) uint64 {
	return bps / 1000000 // nolint:gomnd
}
