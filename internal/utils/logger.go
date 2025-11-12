package utils

import (
	"io"
	"log"
	"os"
)

var (
	// InfoLogger logs informational messages
	InfoLogger *log.Logger

	// ErrorLogger logs error messages
	ErrorLogger *log.Logger
)

// InitLogger initializes the application loggers
// It creates log files in the logs/ directory and sets up loggers with appropriate prefixes
func InitLogger() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// Open or create the log file
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Create multi-writers that write to both file and stdout
	infoMulti := io.MultiWriter(logFile, os.Stdout)
	errorMulti := io.MultiWriter(logFile, os.Stderr)

	// Initialize loggers with prefixes and flags
	InfoLogger = log.New(infoMulti, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorMulti, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger.Println("Logger initialized successfully")

	return nil
}
