package download

// Package "download" can be used to query info about box versions available
// at the cloud (Abitti, yo?) and retrieve these box images. The images
// can be installed using package "box".

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
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

// Suppress progress messages if there has been less than 2 seconds from a message
const progressLastMessageTimeout = 2 * time.Second

// writeCounter implements io.Writer interface (see downloadServerImage, unZipServerImage)
type writeCounter struct {
	Total              uint64
	FileSize           uint64
	ProgressCallbackFn func(string, int)
	ProgressString     string
}

var progressLastMessageTime = time.Now()

var cloudStatusCache memory_cache.Cache

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)

	if time.Now().After(progressLastMessageTime.Add(progressLastMessageTimeout)) {
		wc.ProgressCallbackFn(wc.ProgressString, int((100*wc.Total)/wc.FileSize))
		progressLastMessageTime = time.Now()
	}

	return n, nil
}

// GetServerImagePath returns path to cached server image (~/ktp/naksu_last_image.zip)
func GetServerImagePath() string {
	return filepath.Join(mebroutines.GetKtpDirectory(), "naksu_last_image.zip")
}

func downloadServerImage(url string, progressCallbackFn func(string, int)) error {
	if mebroutines.ExistsFile(mebroutines.GetZipImagePath()) {
		err := os.Remove(mebroutines.GetZipImagePath())
		if err != nil {
			return fmt.Errorf("could not remove old image file: %v", err)
		}
	}

	progressCallbackFn(xlate.Get("Contacting server"), 0)
	log.Debug(fmt.Sprintf("Starting to download image from '%s'", url))
	response, errHTTPGet := http.Get(url) // #nosec

	if errHTTPGet != nil {
		log.Debug(fmt.Sprintf("HTTP GET from url '%s' gives an error: %v", url, errHTTPGet))
		return errHTTPGet
	}

	if response.StatusCode != 200 {
		log.Debug(fmt.Sprintf("HTTP GET from url '%s' gives a status code %d", url, response.StatusCode))
		return fmt.Errorf("%d", response.StatusCode)
	}

	defer response.Body.Close()

	fileSize := uint64(response.ContentLength)

	progressCallbackFn(xlate.Get("Opening file"), 1)
	zipFile, errFile := os.OpenFile(mebroutines.GetZipImagePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if errFile != nil {
		log.Debug(fmt.Sprintf("Could not open file '%s' for server image zip: %v", mebroutines.GetZipImagePath(), errFile))
		return errFile
	}
	defer zipFile.Close()

	progressCallbackFn(xlate.Get("Downloading server image"), 2)

	counter := &writeCounter{}
	counter.ProgressCallbackFn = progressCallbackFn
	counter.FileSize = fileSize
	counter.ProgressString = xlate.GetRaw("Downloading server image")

	var errCopy error
	if _, errCopy = io.Copy(zipFile, io.TeeReader(response.Body, counter)); errCopy != nil {
		return errCopy
	}

	progressCallbackFn(xlate.Get("Server image downloaded"), 100)

	return nil
}

func unZipServerImage(progressCallbackFn func(string, int)) error {
	r, err := zip.OpenReader(mebroutines.GetZipImagePath())
	if err != nil {
		return fmt.Errorf("could not open zip %s: %v", mebroutines.GetZipImagePath(), err)
	}
	defer r.Close()

	for _, file := range r.File {
		log.Debug(fmt.Sprintf("Etcher zip contains file %s, size %s", file.Name, humanize.Bytes(file.UncompressedSize64)))

		if file.Name == "ytl/ktp.img" {

			fImage, err := os.OpenFile(mebroutines.GetImagePath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				return fmt.Errorf("could not create image file %s: %v", mebroutines.GetImagePath(), err)
			}

			defer fImage.Close()

			fZipped, err := file.Open()
			if err != nil {
				return fmt.Errorf("could not open file inside the zip: %v", err)
			}

			defer fZipped.Close()

			progressCallbackFn(xlate.Get("Starting to uncompress raw image"), 1)

			counter := &writeCounter{}
			counter.ProgressCallbackFn = progressCallbackFn
			counter.FileSize = file.UncompressedSize64
			counter.ProgressString = xlate.GetRaw("Uncompressing image...")

			if _, err = io.Copy(fImage, io.TeeReader(fZipped, counter)); err != nil {
				return err
			}

			progressCallbackFn(xlate.Get("Uncompressing finished"), 100)
		}
	}

	return nil
}

func GetServerImage(url string, progressCallbackFn func(string, int)) error {
	err := downloadServerImage(url, progressCallbackFn)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to download server image from '%s': %v", url, err))
		return err
	}

	err = unZipServerImage(progressCallbackFn)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to unZipServerImage: %v", err))
		return err
	}

	return nil
}

func GetAvailableVersion(versionURL string) (string, error) {
	ensureCloudStatusCacheInitialised()

	var version string

	cachedVersion, errCacheGet := cloudStatusCache.Get(versionURL)

	if errCacheGet == nil {
		version = fmt.Sprintf("%v", cachedVersion)
	} else {
		response, err := http.Get(versionURL) // #nosec
		if err != nil {
			log.Debug(fmt.Sprintf("Getting available version from '%s' resulted an error: %v", versionURL, err))
			return "", err
		}

		defer response.Body.Close()

		if response.StatusCode != 200 {
			log.Debug(fmt.Sprintf("Getting available version from '%s' gives a status code %d", versionURL, response.StatusCode))
			return "", fmt.Errorf("%d", response.StatusCode)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Debug(fmt.Sprintf("Reading available version from '%s' resulted and error: %v", versionURL, err))
			return "", err
		}

		version = sanitizeBoxVersionString(string(body))

		errCacheSet := cloudStatusCache.Set(versionURL, version, constants.CloudStatusTimeout)
		if errCacheSet != nil {
			log.Debug(fmt.Sprintf("Could not set cloud status cache: %v", errCacheSet))
		}

		log.Debug(fmt.Sprintf("Box version from '%s' is '%s'", versionURL, version))
	}

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
			log.Debug(fmt.Sprintf("Fatal error: Failed to initialise memory cache: %v", err))
			panic(err)
		}
	}
}
