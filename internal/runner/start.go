package runner

import (
	"rft/internal/op"
)

func start(operation *op.Operation) error {
	return createRunner(operation)
}
