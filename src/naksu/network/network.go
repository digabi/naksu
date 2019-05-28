package network

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
)

// CheckIfNetworkAvailable tests if a pre-set utterly-reliable network setver responds to HTTP GET
func CheckIfNetworkAvailable() bool {
	return testHTTPGet(constants.URLTest, constants.URLTestTimeout)
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
		log.Debug(fmt.Sprintf("Testing HTTP GET %s and got error %v", url, err.Error()))
		return false
	}
	defer mebroutines.Close(resp.Body)

	log.Debug(fmt.Sprintf("Testing HTTP GET %s succeeded", url))

	return true
}

// DownloadFile downloads a file from the given URL and stores it to the given destFile.
// Returns error
func DownloadFile(url string, destFile string) error {
	log.Debug(fmt.Sprintf("Starting download from URL %s to file %s", url, destFile))

	out, err1 := os.Create(destFile)
	if err1 != nil {
		return errors.New("failed to create file")
	}
	defer mebroutines.Close(out)

	/* #nosec */
	resp, err2 := http.Get(url)
	if err2 != nil {
		return errors.New("failed to retrieve file")
	}
	defer mebroutines.Close(resp.Body)

	_, err3 := io.Copy(out, resp.Body)
	if err3 != nil {
		return errors.New("failed to copy body")
	}

	log.Debug(fmt.Sprintf("Finished download from URL %s to file %s", url, destFile))
	return nil
}

// DownloadString downloads a file and returns it as a string
func DownloadString(url string) (string, error) {
	fTemp, errTemp := ioutil.TempFile("", "naksu_")
	if errTemp != nil {
		log.Debug("DownloadString could not create temporary file")
		return "", errors.New("could not create temporary file")
	}

	tempname := fTemp.Name()
	errTemp = fTemp.Close()
	if errTemp != nil {
		log.Debug("DownloadString could not close temporary file")
		return "", errors.New("could not close temporary file")
	}

	errDL := DownloadFile(url, tempname)
	if errDL != nil {
		log.Debug(fmt.Sprintf("DownloadString could not download URL %s to file %s", url, tempname))
		return "", errors.New("could not download url")
	}

	// tempname originates from ioutil.TempFile() so we can turn gosec lint off here:
	/* #nosec */
	buffer, errRead := ioutil.ReadFile(tempname)
	if errRead != nil {
		log.Debug(fmt.Sprintf("DownloadString could not read file %s", tempname))
		return "", errors.New("could not read temporary file")
	}

	resultString := string(buffer)

	errRemove := os.Remove(tempname)
	if errRemove != nil {
		log.Debug(fmt.Sprintf("DownloadString could not remove temporary file %s", tempname))
		// We don't return error as the temp files will get deleted anyway by the OS
	}

	return resultString, nil
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

// IgnoreExtInterface returns true if given system-level network device should
// be ignored (i.e. not to be shown to the user).
func IgnoreExtInterface(interfaceName string) bool {
	for i := range constants.ExtNicsToIgnore {
		match, err := regexp.MatchString(constants.ExtNicsToIgnore[i], interfaceName)
		if err == nil && match {
			return true
		}
	}

	return false
}
