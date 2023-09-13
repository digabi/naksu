package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/natefinch/lumberjack.v2"

	"naksu/constants"
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
		const maxBackups = 5
		const maxLogSizeInMegabytes = 10

		lumberLog := lumberjack.Logger{
			Compress:   false,
			Filename:   debugFilename,
			LocalTime:  false,
			MaxAge:     0,
			MaxBackups: maxBackups,
			MaxSize:    maxLogSizeInMegabytes,
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

	ktpDir := filepath.Join(homeDir, "ktp")

	if !existsDir(ktpDir) {
		err = os.Mkdir(ktpDir, constants.FilePermissionsOwnerRWX)
		if err != nil {
			fmt.Printf("Warning: log.GetNewDebugFilename() could not create directory '%s': %v\n", ktpDir, err)
		}
	}

	if existsDir(ktpDir) {
		newDebugFilename = filepath.Join(ktpDir, "naksu_lastlog.txt")
	} else {
		newDebugFilename = filepath.Join(os.TempDir(), "naksu_lastlog.txt")
	}

	return newDebugFilename
}

// IsDebug returns true if we need to log debug information
func IsDebug() bool {
	return isDebug
}

// writeLogMessage writes log entries with the specified prefix
func writeLogMessage(prefix string, message string, vars ...interface{}) {
	formattedMessage := fmt.Sprintf(message, vars...)
	if (prefix == "DEBUG" && IsDebug()) || prefix != "DEBUG" {
		fmt.Printf("%s: %s\n", prefix, formattedMessage)
	}

	appendLogFile(fmt.Sprintf("%s: %s", prefix, formattedMessage))
}

// Debug logs debug information to log file
func Debug(message string, vars ...interface{}) {
	writeLogMessage("DEBUG", message, vars...)
}

// Error logs error message to log file
func Error(message string, vars ...interface{}) {
	writeLogMessage("ERROR", message, vars...)
}

// Warning logs warning information to log file
func Warning(message string, vars ...interface{}) {
	writeLogMessage("WARNING", message, vars...)
}

// Warning logs info message to log file
func Info(message string, vars ...interface{}) {
	writeLogMessage("INFO", message, vars...)
}

// Action logs action information (i.e. user action) to log file
func Action(message string, vars ...interface{}) {
	writeLogMessage("ACTION", message, vars...)
}
