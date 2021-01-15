package journal

import (
	"path"

	"barf/internal/op"
	"barf/internal/utils/fs"
)

func writeStatus(dir string, operationID op.OperationID, status *op.OperationStatus) error {
	filename := path.Join(dir, string(operationID), "status.json")

	return fs.WriteJSONFile(filename, status)
}

func readStatus(dir string, operationID op.OperationID) *op.OperationStatus {
	filename := path.Join(dir, string(operationID), "status.json")

	status := op.NewStatus()

	_ = fs.ReadJSONFile(filename, status)

	return status
}
