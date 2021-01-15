package op

import (
	"strings"

	"github.com/rs/xid"
)

// Operation types
const (
	OpCopy    OperationType = "copy"
	OpMove    OperationType = "move"
	OpPull    OperationType = "pull"
	OpPush    OperationType = "push"
	OpBackup  OperationType = "backup"
	OpRestore OperationType = "restore"
	OpDummy   OperationType = "dummy"
)

// OperationIndex represents unique index for the operation among the active operations
type OperationIndex int

// OperationID represents unique id for the operation
type OperationID string

// OperationTitle should be a human readable name for the operation
type OperationTitle string

// OperationType should be set to one of the Operation types: OpCopy, OpMove etc.
type OperationType string

// OperationArgs holds the arguments for the operation
type OperationArgs map[string]interface{}

// Operation defines the parameters for an operation
type Operation struct {
	ID    OperationID    `json:"id"`
	Index OperationIndex `json:"index"`
	Title OperationTitle `json:"title"`
	Type  OperationType  `json:"type"`
	Args  OperationArgs  `json:"args"`
}

// NewOperation creates an operation object
func NewOperation(opType OperationType, args OperationArgs, index OperationIndex) *Operation {
	id := xid.New()

	return &Operation{
		ID:    OperationID(id.String()),
		Index: index,
		Title: OperationTitle(strings.Title(string(opType))),
		Type:  opType,
		Args:  args,
	}
}
