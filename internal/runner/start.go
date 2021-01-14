package runner

import (
	"barf/internal/op"
)

func start(operation *op.Operation) error {
	return createRunner(operation)
}
