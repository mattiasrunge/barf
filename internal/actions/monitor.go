package actions

import (
	"barf/internal/com/client"
	"barf/internal/typeconv"
	"barf/internal/ui"
	"barf/internal/utils"
)

// Monitor starts monitoring all or specified operations
func Monitor(args map[string]interface{}) error {
	idxArray, _ := typeconv.ToArray(args["ids"])
	strIdx := typeconv.ToStringArray(idxArray)
	idx := typeconv.StringArray2IntArray(strIdx)

	operations, err := client.ListOperations()

	if err != nil {
		return err
	}

	for _, operation := range operations {
		if len(idx) == 0 || utils.IntArrayContains(idx, int(operation.Index)) {
			err = ui.AddOperation(operation)

			if err != nil {
				return err
			}
		}
	}

	// TODO: if len(idx) == 0 listen for new operations and don't exit

	return nil
}
