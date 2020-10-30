package vboxmanage

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	semver "github.com/blang/semver/v4"
	memory_cache "github.com/paulusrobin/go-memory-cache/memory-cache"

	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
)

var vBoxResponseCache memory_cache.Cache
var vBoxManageStarted int64

type VBoxCommand = []string

func CallRunVBoxManage(args VBoxCommand) (string, error) {
	// There is an ongoing VBoxManage call (break free after 240 loops)
	// This locking avoids executing multiple instances of VBoxManage at the same time. Calling
	// VBoxManage simulaneously tends to cause E_ACCESSDENIED errors from VBoxManage.
	tryCounter := 0
	for (vBoxManageStarted != 0) && (tryCounter < 240) {
		time.Sleep(500 * time.Millisecond)
		tryCounter++
		log.Debug(fmt.Sprintf("CallRunVBoxManage is waiting VBoxManage to exit (race condition lock count %d)", tryCounter))
	}

	vBoxManageStarted = time.Now().Unix()
	vBoxManageOutput, err := runVBoxManage(args)
	vBoxManageStarted = 0

	return vBoxManageOutput, err
}

func MultipleCallRunVBoxManage(commands []VBoxCommand) error {
	for curCommand := 0; curCommand < len(commands); curCommand++ {
		_, err := CallRunVBoxManage(commands[curCommand])
		if err != nil {
			return err
		}
	}

	return nil
}

// runVBoxManage runs vboxmanage command with given arguments
func runVBoxManage(args []string) (string, error) {
	vboxmanagepathArr := []string{getVBoxManagePath()}
	runArgs := append(vboxmanagepathArr, args...)
	vBoxManageOutput, err := mebroutines.RunAndGetOutput(runArgs)
	if err != nil {
		command := strings.Join(runArgs, " ")
		logError := func(output string, err error) {
			log.Debug(fmt.Sprintf("Failed to execute %s (%v), complete output:", command, err))
			log.Debug(output)
		}

		logError(vBoxManageOutput, err)

		fixed, fixErr := detectAndFixDuplicateHardDiskProblem(vBoxManageOutput)
		if fixErr != nil {
			log.Debug(fmt.Sprintf("Failed to fix duplicate hard disk problem with command %s: (%v)", command, fixErr))
			return "", fmt.Errorf("failed to execute %s: %v", command, err)
		}

		if fixed {
			log.Debug("Duplicate hard disk problem was fixed")
		} else {
			log.Debug("Duplicate hard disk problem was not fixed")
		}

		log.Debug(fmt.Sprintf("Retrying '%s' after fixing problem", command))
		vBoxManageOutput, err = mebroutines.RunAndGetOutput(runArgs)
		if err != nil {
			logError(vBoxManageOutput, err)
		}
	}

	return vBoxManageOutput, err
}

func ensureVBoxResponseCacheInitialised() {
	var err error

	if vBoxResponseCache == nil {
		vBoxResponseCache, err = memory_cache.New()
		if err != nil {
			log.Debug(fmt.Sprintf("Fatal error: Failed to initialise memory cache: %v", err))
			panic(err)
		}
	}
}

func ResetVBoxResponseCache() {
	vBoxResponseCache = nil
	ensureVBoxResponseCacheInitialised()
}

// GetVMInfoRegexp returns result of the given vmRegexp from the current VBoxManage showvminfo
// output. This function gets the output either from the cache or calls getVBoxManageOutput()
func GetVMInfoRegexp(boxName string, vmRegexp string) string {
	var rawVMInfo string

	ensureVBoxResponseCacheInitialised()

	rawVMInfoInterface, err := vBoxResponseCache.Get("showvminfo")
	if err != nil {
		rawVMInfo, err = CallRunVBoxManage([]string{"showvminfo", "--machinereadable", boxName})
		if err != nil {
			log.Debug(fmt.Sprintf("Could not get VM info: %v", err))
			rawVMInfo = ""
		}

		errCache := vBoxResponseCache.Set("showvminfo", rawVMInfo, constants.VBoxManageCacheTimeout)
		if errCache != nil {
			log.Debug(fmt.Sprintf("Could not store VM info to cache: %v", errCache))
		}
	} else {
		rawVMInfo = fmt.Sprintf("%v", rawVMInfoInterface)
	}

	// Extract server name
	pattern := regexp.MustCompile(vmRegexp)
	result := pattern.FindStringSubmatch(rawVMInfo)

	if len(result) > 1 {
		return result[1]
	}

	return ""
}

func getVBoxManageVersionSemanticPart() (string, error) {
	output, errVBM := CallRunVBoxManage([]string{"--version"})
	if errVBM != nil {
		log.Debug(fmt.Sprintf("GetVBoxManageVersion() failed to get VBoxManage version: %v", errVBM))
		return "", fmt.Errorf("failed to get vboxmanage version: %v", errVBM)
	}

	re := regexp.MustCompile(`^(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) == 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("could not find semantic version string from vboxmanage version '%s'", output)
}

func GetVBoxManageVersion() (semver.Version, error) {
	ensureVBoxResponseCacheInitialised()

	errorVersion, _ := semver.Make("0.0.0")

	cachedVBoxManageVersion, errCache := vBoxResponseCache.Get("vboxmanageversion")
	if errCache != nil {
		vBoxManageVersionString, errVersionString := getVBoxManageVersionSemanticPart()
		if errVersionString != nil {
			log.Debug(fmt.Sprintf("GetVBoxManageVersion() could not get VBoxManage version: %v", errVersionString))
			return errorVersion, errVersionString
		}

		vBoxManageVersion, errSemVer := semver.Make(vBoxManageVersionString)
		if errSemVer != nil {
			log.Debug(fmt.Sprintf("GetVBoxManageVersion() got VBoxManage version code '%s' but it is not semantic version number: %v", vBoxManageVersionString, errSemVer))
			return errorVersion, fmt.Errorf("vboxmanage version %s is not a semantic version number: %v", vBoxManageVersionString, errSemVer)
		}

		errCache = vBoxResponseCache.Set("vboxmanageversion", vBoxManageVersion.String(), constants.VBoxManageCacheTimeout)
		if errCache != nil {
			log.Debug(fmt.Sprintf("GetVBoxManageVersion() could not store version to cache: %v", errCache))
		}

		return vBoxManageVersion, nil
	}

	cachedVBoxManageVersionSemVer, _ := semver.Make(fmt.Sprintf("%v", cachedVBoxManageVersion))

	return cachedVBoxManageVersionSemVer, nil
}

func GetBoxProperty(boxName string, property string) string {
	ensureVBoxResponseCacheInitialised()

	propertyValue := ""

	propertyValueInterface, errCache := vBoxResponseCache.Get(property)
	if errCache != nil {
		output, errVBoxManage := CallRunVBoxManage([]string{"guestproperty", "get", boxName, property})
		if errVBoxManage != nil {
			log.Debug(fmt.Sprintf("Could not get VM guest property: %v", errVBoxManage))
			output = ""
		}

		propRegexp := regexp.MustCompile(`Value: (.+)`)
		propMatches := propRegexp.FindStringSubmatch(output)
		if len(propMatches) == 2 {
			propertyValue = propMatches[1]
		}

		errCacheSet := vBoxResponseCache.Set(property, propertyValue, constants.VBoxManageCacheTimeout)
		if errCacheSet == nil {
			log.Debug(fmt.Sprintf("Stored VM guest property %s to cache: %s", property, propertyValue))
		} else {
			log.Debug(fmt.Sprintf("Could not store VM guest property %s, value %s to cache: %v", property, propertyValue, errCacheSet))
		}
	} else {
		propertyValue = fmt.Sprintf("%v", propertyValueInterface)
		log.Debug(fmt.Sprintf("Got VM guest property %s from cache: %s", property, propertyValue))
	}

	return propertyValue
}

func checkOutputIfNoVMInstalled(output string) bool {
	re := regexp.MustCompile(`Could not find a registered machine named`)
	return re.MatchString(output)
}

func checkOutputGetVMState(output string) string {
	re := regexp.MustCompile(`VMState="(.+)"`)
	result := re.FindStringSubmatch(output)

	if len(result) > 1 {
		return result[1]
	}

	return ""
}

func getVMState(boxName string) (string, error) {
	ensureVBoxResponseCacheInitialised()

	vmState, err := vBoxResponseCache.Get("vmstate")
	if err != nil {
		rawVMInfo, err := CallRunVBoxManage([]string{"showvminfo", "--machinereadable", boxName})

		// Check whether VM is installed
		if checkOutputIfNoVMInstalled(rawVMInfo) {
			log.Debug("When trying to get VM state, VM is not installed")
			return "", nil
		}

		// Process other VBoxManage errors
		if err != nil {
			log.Debug(fmt.Sprintf("When trying to get VM state, could not get VM info: %v", err))
			return "", err
		}

		// Extract state string
		vmState := checkOutputGetVMState(rawVMInfo)
		if vmState == "" {
			log.Debug("Could not find VM state from the VM info")
			return "", errors.New("could not find vm state from the vm info")
		}

		errCache := vBoxResponseCache.Set("vmstate", vmState, constants.VBoxRunningCacheTimeout)
		if errCache != nil {
			log.Debug(fmt.Sprintf("Could not store VM state to cache: %v", errCache))
		}
	}

	return fmt.Sprintf("%v", vmState), nil
}

func Running(boxName string) (bool, error) {
	vmState, err := getVMState(boxName)

	log.Debug(fmt.Sprintf("vboxmanage.Running() got following state string: '%s'", vmState))
	if vmState == "running" {
		return true, err
	}

	return false, err
}

func Installed(boxName string) (bool, error) {
	rawVMInfo, err := CallRunVBoxManage([]string{"showvminfo", "--machinereadable", boxName})

	if err != nil {
		if checkOutputIfNoVMInstalled(rawVMInfo) {
			log.Debug("box.Installed: Box is not installed")
			return false, nil
		}

		// Other error, return it to the caller
		return false, err
	}

	// We got the showvminfo all right, so the machine is installed
	return true, nil
}

// InstalledVBoxManage returns true if VBoxManage has been installed
func InstalledVBoxManage() bool {
	var vboxmanagepath = getVBoxManagePath()

	if vboxmanagepath == "" {
		log.Debug("Could not get VBoxManage path")
		return false
	}

	vBoxManageVersion, err := CallRunVBoxManage([]string{"--version"})
	if err != nil {
		// No VBoxManage was found
		log.Debug("VBoxManage was not found")
		return false
	}

	log.Debug(fmt.Sprintf("VBoxManage version: %s", vBoxManageVersion))
	return true
}
