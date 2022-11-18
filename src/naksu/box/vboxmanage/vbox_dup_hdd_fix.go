package vboxmanage

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"

	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
)

var duplicateHardDiskRegexp = regexp.MustCompile("because a hard disk '[^']*' with UUID {([0-9a-fA-F-]+)} already exists")

// detectAndFixDuplicateHardDiskProblem fixes VBoxManage crashes in VirtualBoxes before 6.0.8 or 5.2.30
// if VirtualBox.xml has multiple HardDisk with the same location but different UUIDs
func detectAndFixDuplicateHardDiskProblem(vBoxManageOutput string) (wasFixed bool, err error) {
	duplicateHardDiskRegexpMatch := duplicateHardDiskRegexp.FindStringSubmatch(vBoxManageOutput)
	if len(duplicateHardDiskRegexpMatch) != 0 {
		orphanedHardDiskUUID := duplicateHardDiskRegexpMatch[1]
		log.Debug("Detected duplicate VirtualBox disk %s, fixing...", orphanedHardDiskUUID)

		fixedVirtualBoxConfigPath, err := writeFixedVirtualBoxConfig(orphanedHardDiskUUID)
		if err != nil {
			return false, err
		}

		virtualBoxConfigBackupPath, err := backupVirtualBoxConfig()
		if err != nil {
			return false, err
		}

		err = replaceVirtualBoxConfigWithFixedOne(fixedVirtualBoxConfigPath, virtualBoxConfigBackupPath)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func writeFixedVirtualBoxConfig(orphanedHardDiskUUID string) (string, error) {
	virtualBoxConfigFile, err := os.Open(getVirtualBoxConfigPath())
	if err != nil {
		log.Error("Failed to open virtualbox configuration file %s", getVirtualBoxConfigPath())

		return "", err
	}
	defer func() {
		_ = virtualBoxConfigFile.Close() // #nosec
	}()

	fixedVirtualBoxConfigPath := getVirtualBoxConfigPath() + ".new"
	fixedVirtualBoxConfigFile, err := os.OpenFile(fixedVirtualBoxConfigPath, os.O_RDWR|os.O_CREATE, constants.FilePermissionsOwnerRW)
	if err != nil {
		log.Error("Failed to open replacement virtual box configuration file %s", fixedVirtualBoxConfigPath)

		return "", err
	}
	defer func() {
		_ = fixedVirtualBoxConfigFile.Close() // #nosec
	}()

	scanner := bufio.NewScanner(virtualBoxConfigFile)
	fixedVirtualBoxConfigWriter := bufio.NewWriter(fixedVirtualBoxConfigFile)
	for scanner.Scan() {
		line := scanner.Text()
		isDuplicateHardDiskLine, err := regexp.MatchString("<HardDisk uuid=\"{"+orphanedHardDiskUUID+"}\"", line)

		if err != nil {
			log.Error("Error trying to regexp, UUID=%s", orphanedHardDiskUUID)

			return "", err
		}

		if !isDuplicateHardDiskLine {
			_, err = fixedVirtualBoxConfigWriter.WriteString(line + "\r\n")
			if err != nil {
				log.Error("Failed to write to %s", fixedVirtualBoxConfigPath)

				return "", err
			}
		}
	}

	if err := fixedVirtualBoxConfigWriter.Flush(); err != nil {
		log.Error("Failed to flush writer of %s", fixedVirtualBoxConfigPath)

		return "", err
	}
	if err := fixedVirtualBoxConfigFile.Close(); err != nil {
		log.Error("Failed to close %s", fixedVirtualBoxConfigPath)

		return "", err
	}

	return fixedVirtualBoxConfigPath, nil
}

func backupVirtualBoxConfig() (string, error) {
	virtualBoxConfigBackupPath := getVirtualBoxConfigPath() + ".naksubackup"
	if err := os.Rename(getVirtualBoxConfigPath(), virtualBoxConfigBackupPath); err != nil {
		log.Error("Failed to backup %s to %s", getVirtualBoxConfigPath(), virtualBoxConfigBackupPath)

		return "", err
	}

	return virtualBoxConfigBackupPath, nil
}

func replaceVirtualBoxConfigWithFixedOne(fixedVirtualBoxConfigPath string, virtualBoxConfigBackupPath string) error {
	if err := os.Rename(fixedVirtualBoxConfigPath, getVirtualBoxConfigPath()); err != nil {
		log.Error("Failed to move %s to %s, trying to restore %s from %s", fixedVirtualBoxConfigPath, getVirtualBoxConfigPath(), getVirtualBoxConfigPath(), virtualBoxConfigBackupPath)

		if recoveryErr := os.Rename(virtualBoxConfigBackupPath, getVirtualBoxConfigPath()); recoveryErr != nil {
			mebroutines.ShowTranslatedErrorMessage("Naksu encountered an error while trying to fix a problem with VirtualBox. VirtualBox configuration file %s has been moved to %s. Manually rename it to %s to fix this.", getVirtualBoxConfigPath(), virtualBoxConfigBackupPath, getVirtualBoxConfigPath())
		}

		return err
	}

	fmt.Print("Please wait a few seconds while Naksu is fixing a duplicate hard disk problem that it detected with VirtualBox.\n")
	// VBoxManage takes 5 seconds from last run to notice that the config file is now on a different inode
	const timeToWaitForVirtualbox = 5500 * time.Millisecond
	time.Sleep(timeToWaitForVirtualbox)

	return nil
}
