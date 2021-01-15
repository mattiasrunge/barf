package runner

import (
	"barf/internal/coordinator"
)

// Start reads the journal and starts operations
func Start() error {
	coordinator.OnOperationStart(start)
	coordinator.OnOperationAbort(abort)

	return coordinator.Initialize()
}
