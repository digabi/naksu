package cloud

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	humanize "github.com/dustin/go-humanize"

	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
)

// Suppress progress messages if there has been less than 2 seconds from a message
const progressLastMessageTimeout = 2 * time.Second

// writeCounter defines  io.Writer interface (see downloadServerImage, unZipServerImage)
type writeCounter struct {
	Total              uint64
	FileSize           uint64
	ProgressCallbackFn func(string)
	ProgressString     string
}

var progressLastMessageTime = time.Now()

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)

	if time.Now().After(progressLastMessageTime.Add(progressLastMessageTimeout)) {
		wc.ProgressCallbackFn(fmt.Sprintf(wc.ProgressString, (100*wc.Total)/wc.FileSize))
		progressLastMessageTime = time.Now()
	}

	return n, nil
}

// GetServerImagePath returns path to cached server image (~/ktp/naksu_last_image.zip)
func GetServerImagePath() string {
	return filepath.Join(mebroutines.GetVagrantDirectory(), "naksu_last_image.zip")
}

func downloadServerImage(url string, destinationZipPath string, progressCallbackFn func(string)) error {
	progressCallbackFn("Contacting server")
	response, errHTTPGet := http.Get(url) // #nosec
	if errHTTPGet != nil {
		log.Debug(fmt.Sprintf("HTTP GET from url '%s' gives an error: %v", url, errHTTPGet))
		return errHTTPGet
	}
	defer response.Body.Close()

	fileSize := uint64(response.ContentLength)

	progressCallbackFn("Opening file")
	zipFile, errFile := os.OpenFile(destinationZipPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if errFile != nil {
		log.Debug(fmt.Sprintf("Could not open file '%s' for server image zip: %v", destinationZipPath, errFile))
		return errFile
	}
	defer zipFile.Close()

	progressCallbackFn("Downloading server image")

	counter := &writeCounter{}
	counter.ProgressCallbackFn = progressCallbackFn
	counter.FileSize = fileSize
	counter.ProgressString = "Downloading image: %d %%"

	var errCopy error
	if _, errCopy = io.Copy(zipFile, io.TeeReader(response.Body, counter)); errCopy != nil {
		return errCopy
	}

	progressCallbackFn("Server image downloaded")

	return nil
}

// Package "cloud" can be used to query info about box versions available
// at the cloud (Abitti, yo?) and retrieve these box images. The images
// can be installed using package "box".

func unZipServerImage(zipPath string, destinationImagePath string, progressCallbackFn func(string)) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("could not open zip %s: %v", zipPath, err)
	}
	defer r.Close()

	for _, file := range r.File {
		log.Debug(fmt.Sprintf("Etcher zip contains file %s, size %s", file.Name, humanize.Bytes(file.UncompressedSize64)))

		if file.Name == "ytl/ktp.img" {

			fImage, errImage := os.OpenFile(destinationImagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if errImage != nil {
				return fmt.Errorf("could not create image file %s: %v", destinationImagePath, errImage)
			}

			fZipped, errZipped := file.Open()
			if errZipped != nil {
				return fmt.Errorf("could not open file inside the zip: %v", errZipped)
			}

			progressCallbackFn("Starting to unzip raw image")

			counter := &writeCounter{}
			counter.ProgressCallbackFn = progressCallbackFn
			counter.FileSize = file.UncompressedSize64
			counter.ProgressString = "Uncompressing image: %d %%"

			var errCopy error
			if _, errCopy = io.Copy(fImage, io.TeeReader(fZipped, counter)); errCopy != nil {
				return errCopy
			}

			progressCallbackFn("unzip finished")
		}
	}

	return nil
}

func getAndUnzipCloudImage(url string, destinationImagePath string, progressCallbackFn func(string)) error {
	zipPath := GetServerImagePath()

	errDownload := downloadServerImage(url, zipPath, progressCallbackFn)
	if errDownload != nil {
		log.Debug(fmt.Sprintf("Failed to download server image from '%s' to '%s': %v", url, destinationImagePath, errDownload))
		return errDownload
	}

	errUnzip := unZipServerImage(zipPath, destinationImagePath, progressCallbackFn)
	if errUnzip != nil {
		log.Debug(fmt.Sprintf("Failed to unZipServerImage %s: %v", zipPath, errUnzip))
		return errUnzip
	}

	return nil
}

func GetAbittiImage(destinationImagePath string, progressCallbackFn func(string)) error {
	return getAndUnzipCloudImage(constants.AbittiEtcherURL, destinationImagePath, progressCallbackFn)
}
