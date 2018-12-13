package network

import (
  "net/http"
  "time"
  "fmt"

  "naksu/mebroutines"
)

// urlTest is testing URL for checking network connection
const urlTest = "http://static.abitti.fi/usbimg/qa/latest.txt"
// URLTestTimeout is timeout (in seconds) for checking network connection
const urlTestTimeout = 4

// CheckIfNetworkAvailable tests if a pre-set utterly-reliable network setver responds to HTTP GET
func CheckIfNetworkAvailable () bool {
  return testHTTPGet(urlTest, urlTestTimeout)
}

// testHTTPGet tests whether HTTP get succeeds to given URL in given timeout (seconds)
func testHTTPGet(url string, timeout int) bool {
	// Set timeout for a HTTP client
	timeoutSeconds := time.Duration(timeout) * time.Second
	client := http.Client{
		Timeout: timeoutSeconds,
	}

	/* #nosec */
	resp, err := client.Get(url)
	if err != nil {
		mebroutines.LogDebug(fmt.Sprintf("Testing HTTP GET %s and got error %v", url, err.Error()))
		return false
	}
	defer mebroutines.Close(resp.Body)

	mebroutines.LogDebug(fmt.Sprintf("Testing HTTP GET %s succeeded", url))

	return true
}
