package coordinator

import (
	"barf/internal/op"
)

// WriteOperationStdout writes to the operations stdout log
func WriteOperationStdout(operationID op.OperationID, line string) {
	e := getEntry(operationID)

	if e == nil {
		return // Just drop
	}

	e.WriteStdout(line)
}

// WriteOperationStderr writes to the operations stderr log
func WriteOperationStderr(operationID op.OperationID, line string) {
	e := getEntry(operationID)

	if e == nil {
		return // Just drop
	}

	e.WriteStderr(line)
}
