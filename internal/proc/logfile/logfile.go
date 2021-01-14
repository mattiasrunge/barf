package logfile

import (
	"log"
	"os"

	"barf/internal/config"
)

// StartLogging redirectes stdout and stderr to a log file
func StartLogging() {
	f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	os.Stdout = f
	os.Stderr = f
}
