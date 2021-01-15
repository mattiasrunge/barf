package journal

import (
	"os"
	"path"

	"barf/internal/config"
)

var activeDir = ""
var historyDir = ""

// Initialize initializes the journal, ensures directories and loads active entries
func Initialize() ([]*JournalEntry, error) {
	activeDir = path.Join(config.JournalDir, "active")
	historyDir = path.Join(config.JournalDir, "history")

	err := os.MkdirAll(activeDir, 0700)

	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(historyDir, 0700)

	if err != nil {
		return nil, err
	}

	return loadActiveJournalEntries()
}
