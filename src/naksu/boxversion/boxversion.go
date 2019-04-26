package boxversion

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/network"
	"naksu/xlate"
)

// Cache for GetVagrantBoxAvailVersionDetails
type lastBoxAvail struct {
	boxString     string
	boxVersion    string
	boxTimestamp  int64
	updateStarted int64
}

// Global cache for GetVagrantBoxAvailVersionDetails()
var vagrantBoxAvailVersionDetailsCache lastBoxAvail

// GetVagrantFileVersion returns a human-readable localised version string
// for a given Vagrantfile (with "" defaults to ~/ktp/Vagrantfile)
func GetVagrantFileVersion(vagrantFilePath string) string {
	if vagrantFilePath == "" {
		vagrantFilePath = filepath.Join(mebroutines.GetVagrantDirectory(), "Vagrantfile")
	}

	boxString, boxVersion, err := GetVagrantFileVersionDetails(vagrantFilePath)
	if err != nil {
		log.LogDebug(fmt.Sprintf("Could not read from %s", vagrantFilePath))
		return ""
	}

	boxType := GetVagrantBoxType(boxString)

	versionString := fmt.Sprintf("%s (v%s)", boxType, boxVersion)
	log.LogDebug(fmt.Sprintf("GetVagrantFileVersion returns: %s", versionString))

	return versionString
}

// GetVagrantFileVersionDetails returns version string (e.g. "digabi/ktp-qa") and
// version number (e.g. "66") from the given vagrantFilePath
func GetVagrantFileVersionDetails(vagrantFilePath string) (string, string, error) {
	fileContent, err := ioutil.ReadFile(filepath.Clean(vagrantFilePath))
	if err != nil {
		log.LogDebug(fmt.Sprintf("Could not read from %s", vagrantFilePath))
		return "", "", err
	}

	boxRegexp := regexp.MustCompile(`config.vm.box = "(.+)"`)
	versionRegexp := regexp.MustCompile(`vb.name = ".+v(\d+)"`)

	boxMatches := boxRegexp.FindStringSubmatch(string(fileContent))
	versionMatches := versionRegexp.FindStringSubmatch(string(fileContent))

	if len(boxMatches) == 2 && len(versionMatches) == 2 {
		log.LogDebug(fmt.Sprintf("GetVagrantFileVersionDetails returns: [%s] [%s]", boxMatches[1], versionMatches[1]))
		return boxMatches[1], versionMatches[1], nil
	}

	return "", "", errors.New("did not find values from vagrantfile")
}

// GetVagrantBoxAvailVersion returns a human-readable localised version string
// for a vagrant box available with update
func GetVagrantBoxAvailVersion() string {
	boxString, boxVersion, err := GetVagrantBoxAvailVersionDetails()
	if err != nil {
		log.LogDebug("Could not get available version string")
		return ""
	}

	boxType := GetVagrantBoxType(boxString)

	versionString := fmt.Sprintf("%s (v%s)", boxType, boxVersion)
	log.LogDebug(fmt.Sprintf("GetVagrantBoxAvailVersion returns: %s", versionString))

	return versionString
}

// GetVagrantBoxAvailVersionDetails gets info about available vagramt box
// from ReallyGetVagrantBoxAvailVersionDetails() or global vagrantBoxAvailVersionDetailsCache
func GetVagrantBoxAvailVersionDetails() (string, string, error) {
	boxString := ""
	boxVersion := ""
	var boxError error

	// There is a avail version fetch going on (break free after 240 loops)
	tryCounter := 0
	for vagrantBoxAvailVersionDetailsCache.updateStarted != 0 && tryCounter < 240 {
		time.Sleep(500)
		tryCounter++
	}

	if vagrantBoxAvailVersionDetailsCache.boxTimestamp < (time.Now().Unix() - constants.VagrantBoxAvailVersionDetailsCacheTimeout) {
		// We need to update the cache
		vagrantBoxAvailVersionDetailsCache.updateStarted = time.Now().Unix()

		boxString, boxVersion, boxError = reallyGetVagrantBoxAvailVersionDetails()
		if boxError == nil {
			vagrantBoxAvailVersionDetailsCache.boxString = boxString
			vagrantBoxAvailVersionDetailsCache.boxVersion = boxVersion
			vagrantBoxAvailVersionDetailsCache.boxTimestamp = time.Now().Unix()
		}

		vagrantBoxAvailVersionDetailsCache.updateStarted = 0
	} else {
		// Return data from the cache

		boxString = vagrantBoxAvailVersionDetailsCache.boxString
		boxVersion = vagrantBoxAvailVersionDetailsCache.boxVersion
	}

	log.LogDebug(fmt.Sprintf("GetVagrantBoxAvailVersionDetails returns: [%s] [%s]", boxString, boxVersion))
	return boxString, boxVersion, boxError
}

// reallyGetVagrantBoxAvailVersionDetails returns version string (e.g. "digabi/ktp-qa") and
// version number (e.g. "69") by getting Vagrantfile -> metadata.json
func reallyGetVagrantBoxAvailVersionDetails() (string, string, error) {
	// Phase 1: Get Abitti Vagrantfile
	strVagrantfile, errVagrantfile := network.DownloadString(constants.AbittiVagrantURL)
	if errVagrantfile != nil {
		log.LogDebug(fmt.Sprintf("Could not download Abitti Vagrantfile from %s", constants.AbittiVagrantURL))
		return "", "", errors.New("could not download abitti vagrantfile")
	}

	reBoxString := regexp.MustCompile(`config.vm.box = "(.+)"`)
	reMetadata := regexp.MustCompile(`config.vm.box_url = "(.+)"`)

	boxStringMatches := reBoxString.FindStringSubmatch(strVagrantfile)
	boxMetadataMatches := reMetadata.FindStringSubmatch(strVagrantfile)

	if len(boxStringMatches) != 2 {
		log.LogDebug("Could not find config.vm.box from Abitti Vagrantfile:")
		log.LogDebug(strVagrantfile)
		return "", "", errors.New("could not find config.vm.box")
	}

	if len(boxMetadataMatches) != 2 {
		log.LogDebug("Could not find config.vm.box_url from Abitti Vagrantfile:")
		log.LogDebug(strVagrantfile)
		return "", "", errors.New("could not find config.vm.box_url")
	}

	// Phase 2: Get vagrant metadata.json
	strMetadata, errMetadata := network.DownloadString(boxMetadataMatches[1])
	if errMetadata != nil {
		log.LogDebug(fmt.Sprintf("Could not download Abitti metadata from %s", boxMetadataMatches[1]))
		return "", "", errors.New("could not download abitti metadata")
	}

	reVersion := regexp.MustCompile(`"version": "(\d+)"`)

	versionMatches := reVersion.FindStringSubmatch(strMetadata)

	if len(versionMatches) != 2 {
		log.LogDebug("Could not find version number from Vagrant metadata:")
		log.LogDebug(strMetadata)
		return "", "", errors.New("could not find version number from vagrant metadata")
	}

	return boxStringMatches[1], versionMatches[1], nil
}

// GetVagrantBoxType returns the type string (Abitti server or Matric Exam server) for vagrant box name
func GetVagrantBoxType(name string) string {
	if GetVagrantBoxTypeIsAbitti(name) {
		return xlate.Get("Abitti server")
	}

	if GetVagrantBoxTypeIsMatriculationExam(name) {
		return xlate.Get("Matric Exam server")
	}

	// Unknown box type
	log.LogDebug(fmt.Sprintf("Warning: We have a vagrant box type string '%s' which does not resolve to Abitti/Matriculation box type (GetVagrantBoxType)", name))
	return "-"
}

// GetVagrantBoxTypeIsAbitti returns true if given box name string
// belongs to an Abitti vagrant box
func GetVagrantBoxTypeIsAbitti(name string) bool {
	return (name == "digabi/ktp-qa")
}

// GetVagrantBoxTypeIsMatriculationExam returns true if given box name string
// belongs to a Matriculation Examination vagrant box
func GetVagrantBoxTypeIsMatriculationExam(name string) bool {
	re := regexp.MustCompile(`[ksKS]*\d\d\d\d[ksKS]*-\d+`)
	return re.MatchString(name)
}
