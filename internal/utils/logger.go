package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logFile   *os.File
	logger    *log.Logger
	logMutex  sync.Mutex
	isInitialized bool
)

// InitLogger initializes the logger with a log file
func InitLogger(logPath string) error {
    logMutex.Lock()

    if isInitialized {
        logMutex.Unlock()
        return nil
    }

    // Open log file in append mode, create if doesn't exist
    file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        logMutex.Unlock()
        return fmt.Errorf("failed to open log file: %w", err)
    }

    logFile = file
    logger = log.New(logFile, "", 0) // No prefix, we'll add timestamp manually
    isInitialized = true

    // Release the mutex before logging to avoid re-entrancy deadlock
    logMutex.Unlock()

    LogInfo("Logger initialized successfully")
    return nil
}

// CloseLogger closes the log file
func CloseLogger() error {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// formatTimestamp returns current timestamp in RFC3339 format
func formatTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// LogInfo logs an informational message
func LogInfo(format string, v ...interface{}) {
	logMessage("INFO", format, v...)
}

// LogError logs an error message
func LogError(format string, v ...interface{}) {
	logMessage("ERROR", format, v...)
}

// LogDebug logs a debug message
func LogDebug(format string, v ...interface{}) {
	logMessage("DEBUG", format, v...)
}

// LogMethodInit logs method initialization
func LogMethodInit(methodName string) {
	LogInfo("[METHOD_INIT] %s", methodName)
}

// LogMethodSuccess logs successful method completion
func LogMethodSuccess(methodName string) {
	LogInfo("[METHOD_SUCCESS] %s", methodName)
}

// LogMethodError logs unsuccessful method completion
func LogMethodError(methodName string, err error) {
	LogError("[METHOD_ERROR] %s: %v", methodName, err)
}

// LogMongoTransaction logs MongoDB transaction start
func LogMongoTransaction(operation string, details string) {
	LogInfo("[MONGO_TRANSACTION] %s: %s", operation, details)
}

// LogExcelInit logs Excel document creation initialization
func LogExcelInit(filename string) {
	LogInfo("[EXCEL_INIT] Starting Excel document creation: %s", filename)
}

// LogExcelComplete logs Excel document creation completion
func LogExcelComplete(filename string) {
	LogInfo("[EXCEL_COMPLETE] Excel document created successfully: %s", filename)
}

// logMessage is the internal function that handles all logging
func logMessage(level string, format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	timestamp := formatTimestamp()
	message := fmt.Sprintf(format, v...)
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)

	// Write to log file if initialized
	if logger != nil {
		logger.Println(logLine)
	}

	// Also write to stdout for console visibility
	fmt.Println(logLine)
}

