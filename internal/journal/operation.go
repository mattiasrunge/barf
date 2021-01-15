package journal

import (
	"path"

	"barf/internal/op"
	"barf/internal/utils/fs"
)

func writeOperation(dir string, operation *op.Operation) error {
	filename := path.Join(dir, string(operation.ID), "operation.json")

	return fs.WriteJSONFile(filename, operation)
}

func readOperation(dir string, operationID op.OperationID) (*op.Operation, error) {
	filename := path.Join(dir, string(operationID), "operation.json")

	var operation op.Operation

	err := fs.ReadJSONFile(filename, &operation)

	return &operation, err
}
