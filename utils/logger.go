package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// Set up loggers with custom prefixes and flags
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogInfo logs an info message with optional fields
func LogInfo(message string, fields ...interface{}) {
	logMessage := formatLogMessage("INFO", message, fields...)
	InfoLogger.Output(2, logMessage)
}

// LogError logs an error message with optional fields
func LogError(message string, fields ...interface{}) {
	logMessage := formatLogMessage("ERROR", message, fields...)
	ErrorLogger.Output(2, logMessage)
}

// formatLogMessage formats the log message with timestamp and fields
func formatLogMessage(level, message string, fields ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	baseMsg := fmt.Sprintf("[%s] %s", timestamp, message)

	if len(fields) > 0 {
		// Format additional fields
		fieldStr := "{"
		for i := 0; i < len(fields); i += 2 {
			if i > 0 {
				fieldStr += ", "
			}
			if i+1 < len(fields) {
				fieldStr += fmt.Sprintf("%v: %v", fields[i], fields[i+1])
			}
		}
		fieldStr += "}"
		baseMsg += " " + fieldStr
	}

	return baseMsg
}
