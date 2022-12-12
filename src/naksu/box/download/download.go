package download

// Package "download" can be used to query info about box versions available
// at the cloud (Abitti, yo?) and retrieve these box images. The images
// can be installed using package "box".

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	humanize "github.com/dustin/go-humanize"
	memory_cache "github.com/paulusrobin/go-memory-cache/memory-cache"

	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/xlate"
)

const (
	// Suppress progress messages if there has been less than 2 seconds from a message
	progressLastMessageTimeout = 2 * time.Second

	downloadProgressPercentageContactingServer = 0
	downloadProgressPercentageOpeningFile      = 1
	downloadProgressPercentageDownloading      = 2
	downloadProgressPercentageFinished         = 100

	unzipProgressPercentageStarting = 1
	unzipProgressPercentageFinished = 100

	httpStatusOK = 200
)

var ErrDownloadedDiskImageCorrupted = errors.New("downloaded image is corrupted")

// writeCounter implements io.Writer interface, see
//   * downloadServerImage
//   * unZipServerImageFile
//   * GetSHA256ChecksumFromFile

type writeCounter struct {
	Total              uint64
	FileSize           uint64
	ProgressCallbackFn func(string, int)
	ProgressString     string
}

var progressLastMessageTime = time.Now()

var cloudStatusCache memory_cache.Cache

func (wc *writeCounter) Write(buffer []byte) (int, error) {
	bufferLength := len(buffer)
	wc.Total += uint64(bufferLength)

	if time.Now().After(progressLastMessageTime.Add(progressLastMessageTimeout)) {
		wc.ProgressCallbackFn(wc.ProgressString, int((100*wc.Total)/wc.FileSize)) // nolint:gomnd
		progressLastMessageTime = time.Now()
	}

	return bufferLength, nil
}

// GetServerImagePath returns path to cached server image (~/ktp/naksu_last_image.zip)
func GetServerImagePath() string {
	return filepath.Join(mebroutines.GetKtpDirectory(), "naksu_last_image.zip")
}

func makeHTTPGet(url string) (http.Response, error) {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	ctx := context.Background()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error("Creating HTTP GET request to '%s' resulted an error: %v", url, err)

		var emptyResponse http.Response

		return emptyResponse, err
	}

	// response should be called by the caller
	response, err := client.Do(request) // nolint:bodyclose
	if err != nil {
		log.Error("Making HTTP GET request to '%s' resulted an error: %v", url, err)

		var emptyResponse http.Response

		return emptyResponse, err
	}

	return *response, nil
}

func downloadServerImage(url string, progressCallbackFn func(string, int)) error {
	if mebroutines.ExistsFile(mebroutines.GetZipImagePath()) {
		err := os.Remove(mebroutines.GetZipImagePath())
		if err != nil {
			return fmt.Errorf("could not remove old image file: %w", err)
		}
	}

	progressCallbackFn(xlate.Get("Contacting server"), downloadProgressPercentageContactingServer)
	log.Debug("Starting to download image from '%s'", url)

	response, err := makeHTTPGet(url)
	if err != nil {
		log.Error("Getting available version from '%s' resulted an error: %v", url, err)

		return err
	}

	defer response.Body.Close()

	if response.StatusCode != httpStatusOK {
		log.Error("HTTP GET from url '%s' gives a status code %d", url, response.StatusCode)

		return fmt.Errorf("%d", response.StatusCode)
	}

	fileSize := uint64(response.ContentLength)

	progressCallbackFn(xlate.Get("Opening file"), downloadProgressPercentageOpeningFile)
	zipFile, errFile := os.OpenFile(mebroutines.GetZipImagePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, constants.FilePermissionsOwnerRW)
	if errFile != nil {
		log.Error("Could not open file '%s' for server image zip: %v", mebroutines.GetZipImagePath(), errFile)

		return errFile
	}
	defer zipFile.Close()

	progressCallbackFn(xlate.Get("Downloading server image"), downloadProgressPercentageDownloading)

	serverImageWriteCounter := writeCounter{
		ProgressCallbackFn: progressCallbackFn,
		FileSize:           fileSize,
		ProgressString:     xlate.GetRaw("Downloading server image"),
		Total:              0,
	}
	counter := &serverImageWriteCounter

	var errCopy error
	if _, errCopy = io.Copy(zipFile, io.TeeReader(response.Body, counter)); errCopy != nil {
		return errCopy
	}

	progressCallbackFn(xlate.Get("Server image downloaded"), downloadProgressPercentageFinished)

	return nil
}

func unZipServerImageChecksum(file *zip.File) (string, error) {
	fZipped, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("could not open file inside the zip: %w", err)
	}

	defer fZipped.Close()

	definedChecksumFileContent, err := io.ReadAll(fZipped)
	if err != nil {
		return "", fmt.Errorf("could not read image checksum file inside the zip: %w", err)
	}

	fZipped.Close()

	return CleanSHA256ChecksumString(string(definedChecksumFileContent)), nil
}

func unZipServerImageFile(file *zip.File, progressCallbackFn func(string, int)) error {
	fImage, err := os.OpenFile(mebroutines.GetImagePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, constants.FilePermissionsOwnerRW)
	if err != nil {
		return fmt.Errorf("could not create image file %s: %w", mebroutines.GetImagePath(), err)
	}

	defer fImage.Close()

	fZipped, err := file.Open()
	if err != nil {
		return fmt.Errorf("could not open file inside the zip: %w", err)
	}

	defer fZipped.Close()

	progressCallbackFn(xlate.Get("Starting to uncompress raw image"), unzipProgressPercentageStarting)

	serverImageUnzipCounter := writeCounter{
		ProgressCallbackFn: progressCallbackFn,
		FileSize:           file.UncompressedSize64,
		ProgressString:     xlate.GetRaw("Uncompressing image..."),
		Total:              0,
	}
	counter := &serverImageUnzipCounter

	if _, err = io.Copy(fImage, io.TeeReader(fZipped, counter)); err != nil {
		return err
	}

	progressCallbackFn(xlate.Get("Uncompressing finished"), unzipProgressPercentageFinished)

	return nil
}

func unZipServerImage(progressCallbackFn func(string, int)) error {
	definedChecksum := ""

	zipReader, err := zip.OpenReader(mebroutines.GetZipImagePath())
	if err != nil {
		return fmt.Errorf("could not open zip %s: %w", mebroutines.GetZipImagePath(), err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		log.Debug("Etcher zip contains file %s, size %s", file.Name, humanize.Bytes(file.UncompressedSize64))

		if file.Name == "ytl/ktp.img.sha256" {
			definedChecksum, err = unZipServerImageChecksum(file)
			if err != nil {
				return err
			}
		}

		if file.Name == "ytl/ktp.img" {
			err = unZipServerImageFile(file, progressCallbackFn)
			if err != nil {
				return err
			}
		}
	}
	if definedChecksum != "" {
		log.Debug("Checking that uncompressed image meets defined checksum '%s'", definedChecksum)

		calculatedChecksum, err := GetSHA256ChecksumFromFile(mebroutines.GetImagePath(), progressCallbackFn)
		if err != nil {
			return fmt.Errorf("could not calculate sha256: %w", err)
		}

		if definedChecksum != calculatedChecksum {
			log.Error("Image checksums differ, defined: %s, calculated: %s", definedChecksum, calculatedChecksum)

			return ErrDownloadedDiskImageCorrupted
		}

		log.Debug("Image checksum verified without errors")
	}

	return nil
}

func GetServerImage(url string, progressCallbackFn func(string, int)) error {
	err := downloadServerImage(url, progressCallbackFn)
	if err != nil {
		log.Error("Failed to download server image from '%s': %v", url, err)

		return err
	}

	err = unZipServerImage(progressCallbackFn)
	if err != nil {
		log.Error("Failed to unZipServerImage: %v", err)

		return err
	}

	return nil
}

func GetAvailableVersion(versionURL string) (string, error) {
	ensureCloudStatusCacheInitialised()

	var version string

	cachedVersion, err := cloudStatusCache.Get(versionURL)

	if err == nil {
		version = fmt.Sprintf("%v", cachedVersion)

		return version, nil
	}

	response, err := makeHTTPGet(versionURL)
	if err != nil {
		log.Error("Getting available version from '%s' resulted an error: %v", versionURL, err)

		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != httpStatusOK {
		log.Error(fmt.Sprintf("Getting available version from '%s' gives a status code %d", versionURL, response.StatusCode))

		return "", fmt.Errorf("%d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Reading available version from '%s' resulted and error: %v", versionURL, err)

		return "", err
	}

	version = sanitizeBoxVersionString(string(body))

	err = cloudStatusCache.Set(versionURL, version, constants.CloudStatusTimeout)
	if err != nil {
		log.Warning("Could not set cloud status cache: %v", err)
	}

	log.Debug("Box version from '%s' is '%s'", versionURL, version)

	return version, nil
}

// sanitizeBoxVersionString removes all unallowed characters from box version string
func sanitizeBoxVersionString(str string) string {
	re := regexp.MustCompile(`\W`)

	return re.ReplaceAllString(str, "")
}

func ensureCloudStatusCacheInitialised() {
	var err error

	if cloudStatusCache == nil {
		cloudStatusCache, err = memory_cache.New()
		if err != nil {
			log.Error("Fatal error: Failed to initialise memory cache: %v", err)
			panic(err)
		}
	}
}
