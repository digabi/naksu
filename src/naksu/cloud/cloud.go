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

// GetServerImagePath returns path to cached server image (~/ktp/naksu_last_image.zip)
func GetServerImagePath() string {
	return filepath.Join(mebroutines.GetVagrantDirectory(), "naksu_last_image.zip")
}

func downloadServerImage(url string, destinationZipPath string, progressCallbackFn func(string)) error {
	progressCallbackFn("Contacting server")
	response, errHttpGet := http.Get(url)
	if errHttpGet != nil {
		log.Debug(fmt.Sprintf("HTTP GET from url '%s' gives an error: %v", url, errHttpGet))
		return errHttpGet
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

	var bytesReadTotal uint64

	for {
		buffer := make([]byte, 100000)
		bytesRead, errRead := response.Body.Read(buffer)
		if errRead != nil && errRead != io.EOF {
			log.Debug(fmt.Sprintf("Error while reading data from '%s': %v", url, errRead))
			return errRead
		}

		zipFile.Write(buffer[0:bytesRead])

		bytesReadTotal = bytesReadTotal + uint64(bytesRead)
		progressCallbackFn(fmt.Sprintf("Downloading image: %d %%", (100*bytesReadTotal)/fileSize))

		if errRead == io.EOF {
			break
		}
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

			unzipDone := make(chan bool)

			go func() {
				_, err = io.Copy(fImage, fZipped)
				fZipped.Close()
				fImage.Close()
				unzipDone <- true
			}()

			waitForZip := true
			for waitForZip {
				time.Sleep(3 * time.Second)
				imageFileStat, err := fImage.Stat()
				if err != nil {
					log.Debug(fmt.Sprintf("Failed to stat raw image file %s: %v", destinationImagePath, err))
				} else {
					progressCallbackFn(fmt.Sprintf("Unzipping image: %d %%", uint64(100*imageFileStat.Size())/file.UncompressedSize64))
				}

				select {
				case unzipDoneValue := <-unzipDone:
					waitForZip = !unzipDoneValue
				default:
				}
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
