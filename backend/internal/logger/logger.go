package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	debugMode   bool
	currentLevel LogLevel
)

func init() {
	// Check if debug mode is enabled
	debugMode = os.Getenv("DEBUG") == "true"

	if debugMode {
		currentLevel = DEBUG
	} else {
		currentLevel = INFO
	}

	// Configure log output format
	log.SetFlags(0) // We'll handle our own formatting
}

func formatMessage(level, color, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// ANSI color codes
	colorReset := "\033[0m"

	return fmt.Sprintf("%s[%s] %s%-5s%s %s",
		colorReset, timestamp, color, level, colorReset, message)
}

func Debug(format string, v ...interface{}) {
	if currentLevel <= DEBUG {
		colorCyan := "\033[36m"
		message := fmt.Sprintf(format, v...)
		log.Println(formatMessage("DEBUG", colorCyan, message))
	}
}

func Info(format string, v ...interface{}) {
	if currentLevel <= INFO {
		colorGreen := "\033[32m"
		message := fmt.Sprintf(format, v...)
		log.Println(formatMessage("INFO", colorGreen, message))
	}
}

func Warn(format string, v ...interface{}) {
	if currentLevel <= WARN {
		colorYellow := "\033[33m"
		message := fmt.Sprintf(format, v...)
		log.Println(formatMessage("WARN", colorYellow, message))
	}
}

func Error(format string, v ...interface{}) {
	if currentLevel <= ERROR {
		colorRed := "\033[31m"
		message := fmt.Sprintf(format, v...)
		log.Println(formatMessage("ERROR", colorRed, message))
	}
}

func Fatal(format string, v ...interface{}) {
	colorRed := "\033[31m"
	message := fmt.Sprintf(format, v...)
	log.Println(formatMessage("FATAL", colorRed, message))
	os.Exit(1)
}

func IsDebugMode() bool {
	return debugMode
}
