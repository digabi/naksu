package log

import (
	"fmt"
	"os"
  "io"
  "time"
  "log"
	"path/filepath"

  "github.com/mitchellh/go-homedir"
)

var isDebug bool
var debugFilename string

// Close gracefully handles closing of closable item. defer Close(item)
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func appendLogFile(message string) {
	if debugFilename != "" {
		// Append only if the logfile has been set

		// Current timestamp
		t := time.Now()

		/* #nosec */
		f, err := os.OpenFile(debugFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
		if err != nil {
			panic(fmt.Sprintf("Could not append to log file %s: %s", debugFilename, err))
		}
		defer Close(f)

		_, err = f.WriteString(fmt.Sprintf("[%s] %s\n", t.Format("2006-01-02 15:04:05"), message))
		if err != nil {
			if f.Sync() != nil {
				defer Close(f)
			}
		}
	}
}

func existsDir(path string) bool {
	fi, err := os.Lstat(path)

	if err == nil && fi.Mode().IsDir() {
		return true
	}

	return false
}

func existsFile(path string) bool {
	fi, err := os.Lstat(path)

	if err == nil && fi.Mode().IsRegular() {
		return true
	}

	return false
}

// SetDebug enables debug printing if set to true
func SetDebug(newValue bool) {
	isDebug = newValue
}

// SetDebugFilename sets debug log path
func SetDebugFilename(newFilename string) {
	debugFilename = newFilename

	if debugFilename != "" && existsFile(debugFilename) {
		// Re-create the log file
		err := os.Remove(debugFilename)
		if err != nil {
			panic(fmt.Sprintf("Could not open log file %s: %s", debugFilename, err))
		}
	}
}

// GetNewDebugFilename suggests a new debug log filename
func GetNewDebugFilename() string {
	newDebugFilename := ""

  homeDir, err := homedir.Dir()

  if err != nil {
    panic("Could not get home directory")
  }

	if existsDir(filepath.Join(homeDir, "ktp")) {
		newDebugFilename = filepath.Join(homeDir, "ktp", "naksu_lastlog.txt")
	} else {
		newDebugFilename = filepath.Join(os.TempDir(), "naksu_lastlog.txt")
	}

	return newDebugFilename
}

// IsDebug returns true if we need to log debug information
func IsDebug() bool {
	return isDebug
}

// LogDebug logs debug information to log file
func LogDebug(message string) {
	if IsDebug() {
		fmt.Printf("DEBUG: %s\n", message)
	}

	appendLogFile(fmt.Sprintf("DEBUG: %s", message))
}
