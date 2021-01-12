package runner

import (
	"rft/internal/journal"
)

func init() {
	journal.OnOperationStart(start)
	journal.OnOperationAbort(abort)
}

// StartRunner reads the journal and starts operations
func StartRunner() error {
	return journal.ReadFromDisk()
}
