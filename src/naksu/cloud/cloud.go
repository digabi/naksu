package cloud

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

  humanize "github.com/dustin/go-humanize"

	"naksu/log"
	"naksu/mebroutines"
)

// Package "cloud" can be used to query info about box versions available
// at the cloud (Abitti, yo?) and retrieve these box images. The images
// can be installed using package "box".

func unZipAbittiImage(zipPath string, destinationImagePath string, progressCallbackFn func(string)) error {
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
					progressCallbackFn(fmt.Sprintf("Unzipping image: %d %%", uint64(100*imageFileStat.Size()) / file.UncompressedSize64))
				}

        select {
        case unzipDoneValue := <-unzipDone:
          waitForZip = ! unzipDoneValue
        default:
        }
			}

			progressCallbackFn("unzip finished")
		}
	}

	return nil
}

func GetAbittiImage(destinationImagePath string, progressCallbackFn func(string)) error {
	progressCallbackFn("started")

	err := unZipAbittiImage(filepath.Join(mebroutines.GetHomeDirectory(), "ktp", "ktp-etcher.zip"), destinationImagePath, progressCallbackFn)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to unZipAbittiImage: %v", err))
		progressCallbackFn("failed")
		return err
	}

	progressCallbackFn("ended")

	return nil
}
