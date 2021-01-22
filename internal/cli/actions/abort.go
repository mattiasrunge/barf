package actions

import (
	"errors"
	"fmt"
	"strconv"

	"barf/internal/cli"
	"barf/internal/com/client"
)

// Abort aborts an active operation
func Abort(args map[string]interface{}) error {
	cli.Start()

	id := args["id"].(*string)
	idInt, _ := strconv.Atoi(*id)

	operations, err := client.ListOperations()

	if err != nil {
		return err
	}

	for _, o := range operations {
		if int(o.Index) == idInt || *id == string(o.ID) {
			err = client.AbortOperation(o.ID)

			if err != nil {
				return err
			}

			fmt.Println("Operation aborted!")

			cli.Finish()

			return nil
		}
	}

	return errors.New("no operation matching that ID found")
}
