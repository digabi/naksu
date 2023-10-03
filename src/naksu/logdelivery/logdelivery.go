package logdelivery

import (
	"archive/zip"
	"encoding/base64"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"naksu/box"
	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
)

// DeleteLogCopyFiles deletes temporary files related to copying logs from the virtual machine guest
func DeleteLogCopyFiles() {
	deleteFile(filepath.Join(mebroutines.GetMebshareDirectory(), constants.LogCopyStatusFilename))
	deleteFile(filepath.Join(mebroutines.GetMebshareDirectory(), constants.LogCopyRequestFilename))
	deleteFile(filepath.Join(mebroutines.GetMebshareDirectory(), constants.LogCopyDoneFilename))
}

func deleteFile(filepath string) {
	log.Debug("Deleting file %s", filepath)
	err := os.Remove(filepath)
	if err != nil {
		log.Warning("Could not delete file %s: %s", filepath, err)
	}
}

// RequestLogsFromServer requests logs from the virtual machine and waits for them to be copied to ktp-jako
func RequestLogsFromServer() (chan bool, chan string) {
	log.Debug("Requesting logs from server")

	progressChannel := make(chan string)
	doneChannel := make(chan bool)

	statusFilepath := filepath.Join(mebroutines.GetMebshareDirectory(), constants.LogCopyStatusFilename)
	resetStatusFile(statusFilepath)

	requestFilepath := filepath.Join(mebroutines.GetMebshareDirectory(), constants.LogCopyRequestFilename)
	log.Debug("Using request file %s", requestFilepath)
	requestNumber, err := updateRequestNumber(requestFilepath)
	if err != nil {
		log.Error("Could not update request number in file %s: %s", requestFilepath, err)
		go func() {
			doneChannel <- true
		}()

		return doneChannel, progressChannel
	}

	return waitForLogs(requestNumber, constants.LogRequestTimeout, doneChannel, progressChannel, statusFilepath)
}

func resetStatusFile(statusFilepath string) {
	log.Debug("Resetting status file %s", statusFilepath)
	statusFile, err := openFileForWriting(statusFilepath)
	if err != nil {
		return
	}
	err = writeToFile(statusFilepath, statusFile, "0 %\n")
	if err != nil {
		log.Error("Error resetting status file %s", statusFilepath)
	}
}

func updateRequestNumber(requestFilepath string) (int, error) {
	_, statErr := os.Stat(requestFilepath)
	fileDidNotExists := os.IsNotExist(statErr)

	var currentNumber int
	var err error
	if fileDidNotExists {
		currentNumber = 0
	} else {
		currentNumber, err = readNumberFromFile(requestFilepath)
		if err != nil {
			log.Error("Error reading number from requestFile %s: %s", requestFilepath, err)

			return 0, err
		}
	}
	newNumber := currentNumber + 1

	requestFile, err := openFileForWriting(requestFilepath)
	if err != nil {
		return 0, err
	}
	err = writeToFile(requestFilepath, requestFile, fmt.Sprintf("%d\n", newNumber))

	return newNumber, err
}

func openFileForWriting(filepath string) (*os.File, error) {
	var file *os.File
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, constants.FilePermissionsOwnerRW)
	if err != nil {
		log.Error("Error opening file %s for writing: %s", filepath, err)
	}

	return file, err
}

func readNumberFromFile(filename string) (int, error) {
	content, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		log.Error("Error reading file %s: %s", filename, err)

		return 0, err
	}
	number, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		log.Error("Error converting value to integer %s: %s", content, err)

		// corrupted request file => restart request number sequence

		return 0, nil
	}

	return number, nil
}

func writeToFile(filename string, file *os.File, content string) error {
	_, err := file.WriteString(content)
	if err != nil {
		log.Error("Error writing to file %s: %s", filename, err)

		return err
	}
	err = file.Close()
	if err != nil {
		log.Error("Error closing file %s: %s", filename, err)
	}

	return nil
}

func waitForLogs(requestNumber int, timeout time.Duration, doneChannel chan bool, progressChannel chan string, statusFilepath string) (chan bool, chan string) {
	startTimestamp := time.Now().Local()
	log.Debug("Starting to wait for request number %d at %v with a timeout of %v", requestNumber, startTimestamp, timeout)
	endTimestamp := startTimestamp.Add(timeout)

	doneFilepath := filepath.Join(mebroutines.GetMebshareDirectory(), constants.LogCopyDoneFilename)
	go func() {
		for {
			_, err := os.Stat(doneFilepath)
			if err != nil {
				log.Debug("Done file not yet found at %s", doneFilepath)
			} else {
				var doneNumber int
				doneNumber, err = readNumberFromFile(doneFilepath)
				log.Debug("Found %d in done file %s", doneNumber, doneFilepath)
				if err == nil && doneNumber >= requestNumber {
					log.Debug("Done number %d matches request number %d", doneNumber, requestNumber)
					doneChannel <- true

					return
				}
			}

			time.Sleep(1 * time.Second)

			content, err := os.ReadFile(filepath.Clean(statusFilepath))
			if err != nil {
				log.Warning("Could not read status file %s: %s", doneFilepath, err)
			} else {
				progress := strings.TrimSpace(string(content))
				progressChannel <- progress
			}

			now := time.Now().Local()
			if now.After(endTimestamp) {
				log.Debug("Timing out copying logs at %v", now)
				doneChannel <- true

				return
			}
		}
	}()

	return doneChannel, progressChannel
}

func doCollectLogsToZip(zipFilename string, progress chan uint8, errorChannel chan error) {
	progress <- 0

	zipFilepath := filepath.Join(mebroutines.GetMebshareDirectory(), zipFilename)

	zipFile, err := os.Create(zipFilepath)
	if err != nil {
		log.Error("Error creating zip file %s: %s", zipFilepath, err)
		errorChannel <- err

		return
	}

	var logFiles = []string{}

	logFiles, err = appendKtpLogs(logFiles)
	if err != nil {
		log.Warning("Error appending ktp logs: %s", err)
		// continue collecting logs after error in appending ktp logs
	}
	logFiles, err = appendVirtualBoxLogs(logFiles)
	if err != nil {
		log.Warning("Error appending VirtualBox logs: %s", err)
		// continue collecting logs after error in appending VirtualBox logs
	}
	logFiles, err = appendNaksuLastlogs(logFiles)
	if err != nil {
		log.Warning("Error appending naksu logs: %s", err)
		// continue collecting logs after error in appending naksu logs
	}

	writer := zip.NewWriter(zipFile)
	for logFileNumber, logFilepath := range logFiles {
		err = addFileToZip(logFilepath, writer)
		if err != nil {
			errorChannel <- err

			return
		}

		progress <- uint8(100 * logFileNumber / len(logFiles))
	}

	err = writer.Close()
	if err != nil {
		errorChannel <- err

		return
	}

	progress <- 127
}

// CollectLogsToZip creates a zip file of log files
func CollectLogsToZip() (string, chan uint8, chan error) {
	mebroutines.EnsureMebshareDirectory()

	log.Debug("Collecting logs")
	zipFilename := time.Now().Format("2006-01-02_15-04-05.zip")
	progress := make(chan uint8)
	errorChannel := make(chan error)

	go doCollectLogsToZip(zipFilename, progress, errorChannel)

	return zipFilename, progress, errorChannel
}

func appendKtpLogs(logFiles []string) ([]string, error) {
	ktpLogsDirectory := filepath.Join(mebroutines.GetMebshareDirectory(), "ktp_logs")
	ktpLogsDirectoryFileInfos, err := getDirectoryFileInfos(ktpLogsDirectory)
	if err != nil {
		return logFiles, err
	}
	for _, fileInfo := range ktpLogsDirectoryFileInfos {
		logFilepath := filepath.Join(ktpLogsDirectory, fileInfo.Name())
		log.Debug("Appending ktp log file %s", logFilepath)
		logFiles = append(logFiles, logFilepath)
	}

	return logFiles, nil
}

func appendVirtualBoxLogs(logFiles []string) ([]string, error) {
	virtualBoxLogsDirectory := box.GetLogDir()
	virtualBoxLogsDirectoryFileInfos, err := getDirectoryFileInfos(virtualBoxLogsDirectory)
	if err != nil {
		return logFiles, err
	}
	for _, fileInfo := range virtualBoxLogsDirectoryFileInfos {
		logFilepath := filepath.Join(virtualBoxLogsDirectory, fileInfo.Name())
		log.Debug("Appending VirtualBox log file %s", logFilepath)
		logFiles = append(logFiles, logFilepath)
	}

	return logFiles, nil
}

func appendNaksuLastlogs(logFiles []string) ([]string, error) {
	ktpDirectoryFileInfos, err := getDirectoryFileInfos(mebroutines.GetKtpDirectory())
	if err != nil {
		return logFiles, err
	}
	naksuLastlogRegexp := regexp.MustCompile(`^naksu_lastlog.*\.txt$`)
	for _, fileInfo := range ktpDirectoryFileInfos {
		matched := naksuLastlogRegexp.MatchString(fileInfo.Name())
		if matched {
			logFilepath := filepath.Join(mebroutines.GetKtpDirectory(), fileInfo.Name())
			log.Debug("Appending naksu log file %s", logFilepath)
			logFiles = append(logFiles, logFilepath)
		}
	}

	return logFiles, nil
}

func getDirectoryFileInfos(directory string) ([]fs.DirEntry, error) {
	fileInfos, err := os.ReadDir(directory)
	if err != nil {
		log.Error("Error listing directory %s: %s", directory, err)

		return nil, err
	}

	return fileInfos, nil
}

func addFileToZip(logFilepath string, zipWriter *zip.Writer) error {
	fileInfo, err := os.Stat(logFilepath)
	if err != nil {
		log.Warning("Could not stat %s: %s", logFilepath, err)

		return nil
	}
	fileInfoHeader, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		log.Warning("Converting FileInfo to FileInfoHeader %s: %s", logFilepath, err)

		return nil
	}
	fileInfoHeader.Method = zip.Deflate
	outFile, err := zipWriter.CreateHeader(fileInfoHeader)
	if err != nil {
		log.Error("Error creating zip entry for %s: %s", logFilepath, err)

		return err
	}
	inFile, err := os.Open(filepath.Clean(logFilepath))
	if err != nil {
		log.Warning("Could not open %s: %s", logFilepath, err)

		return nil
	}
	_, err = io.Copy(outFile, inFile)
	if err != nil {
		log.Warning("Could not add %s to zip: %s", logFilepath, err)

		return nil
	}

	return nil
}

// ProgressReadSeeker is a read seeker that can report progress via a callback
type ProgressReadSeeker struct {
	fp               *os.File
	size             int64
	progressCallback func(uint8)
	read             int64
}

// Read reads from file
func (progressReadSeeker *ProgressReadSeeker) Read(buffer []byte) (int, error) {
	return progressReadSeeker.fp.Read(buffer)
}

// ReadAt reads from file at specified offset
func (progressReadSeeker *ProgressReadSeeker) ReadAt(buffer []byte, off int64) (int, error) {
	bytesRead, err := progressReadSeeker.fp.ReadAt(buffer, off)
	if err != nil {
		return bytesRead, err
	}

	atomic.AddInt64(&progressReadSeeker.read, int64(bytesRead))

	// read length is divided by two because s3manager reads the file twice
	progressReadSeeker.progressCallback(uint8(float32(progressReadSeeker.read*100/2) / float32(progressReadSeeker.size))) // nolint:gomnd

	return bytesRead, err
}

// Seek seeks to a specified offset
func (progressReadSeeker *ProgressReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return progressReadSeeker.fp.Seek(offset, whence)
}

// SendLogs sends logs to S3
func SendLogs(filename string, progressCallback func(uint8)) error {
	log.Debug("Sending log file %s", filename)

	setAwsCredentials()

	sess, err := session.NewSession(&aws.Config{ // nolint:exhaustruct
		Region: aws.String("eu-north-1")},
	)
	if err != nil {
		log.Debug("Could not create AWS session")

		return err
	}

	logZipFilePath := filepath.Join(mebroutines.GetMebshareDirectory(), filename)
	logZipFile, err := os.Open(filepath.Clean(logZipFilePath))
	if err != nil {
		log.Error("Could not open %s", logZipFilePath)

		return err
	}

	defer logZipFile.Close()

	fileInfo, err := logZipFile.Stat()
	if err != nil {
		log.Error("Could not stat %s", logZipFilePath)

		return err
	}

	reader := &ProgressReadSeeker{
		fp:               logZipFile,
		size:             fileInfo.Size(),
		progressCallback: progressCallback,
		read:             0,
	}

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		const megabyteInBytes = 1024 * 1024
		u.PartSize = 5 * megabyteInBytes // nolint:gomnd
	})

	output, err := uploader.Upload(&s3manager.UploadInput{ // nolint:exhaustruct
		Bucket: aws.String("naksulogs.yo-prod"),
		Key:    aws.String(filename),
		Body:   reader,
	})
	if err != nil {
		log.Error("Uploading %s failed: %s", filename, err)

		return err
	}

	log.Debug("Log file %s sent to %s", filename, output.Location)

	return nil
}

func decodeBase64(base64String string) string {
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		log.Error("Error decoding base64 '%s'", base64String)
	}

	return string(decoded)
}

func setAwsCredentials() {
	keyIDKey := decodeBase64("QVdTX0FDQ0VTU19LRVlfSUQ=")
	keyID := decodeBase64("QUtJQVQyU1JDRTJVRTNRVllMRlM=")
	secretKeyKey := decodeBase64("QVdTX1NFQ1JFVF9BQ0NFU1NfS0VZ")
	secretKey := decodeBase64("U2I0WEhYcWhYeG5LMnFUQVF6TFd6OFpZVFN1b0w3ZlpZQUVZcjBRaA==")

	setEnvAndLogError(keyIDKey, keyID)
	setEnvAndLogError(secretKeyKey, secretKey)
}

func setEnvAndLogError(key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		log.Error("Could not set environment variable %s", key)
	}
}
