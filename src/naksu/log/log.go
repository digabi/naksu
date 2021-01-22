package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/natefinch/lumberjack.v2"
)

var isDebug bool
var debugFilename string
var logger *log.Logger
var loggerWriter io.WriteCloser

func appendLogFile(message string) {
	if debugFilename != "" {
		// Append only if the logfile has been set
		logger.Print(message)
	}
}

func existsDir(path string) bool {
	fi, err := os.Lstat(path)

	if err == nil && fi.Mode().IsDir() {
		return true
	}

	return false
}

// SetDebug enables debug printing if set to true
func SetDebug(newValue bool) {
	isDebug = newValue
}

// SetDebugFilename sets debug log path
// Setting filename of "-" prints errors to standard error
func SetDebugFilename(newFilename string) {
	debugFilename = newFilename

	if loggerWriter != nil {
		err := loggerWriter.Close()
		if err != nil {
			// nolint: gosec, errcheck
			os.Stderr.WriteString(fmt.Sprintf("Could not close log file: %v", err))
		}
	}

	if newFilename == "-" {
		loggerWriter = os.Stderr
	} else {
		lumberLog := lumberjack.Logger{
			Filename:   debugFilename,
			MaxSize:    3, // megabytes
			MaxBackups: 3,
		}

		loggerWriter = &lumberLog
	}

	logger = log.New(loggerWriter, "", log.Ldate|log.Ltime)
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

// Debug logs debug information to log file
func Debug(message string, vars ...interface{}) {
	formattedMessage := fmt.Sprintf(message, vars...)
	if IsDebug() {
		fmt.Printf("DEBUG: %s\n", formattedMessage)
	}

	appendLogFile(fmt.Sprintf("DEBUG: %s", formattedMessage))
}
