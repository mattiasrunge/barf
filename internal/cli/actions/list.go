package actions

import (
	"fmt"

	"barf/internal/cli"
	"barf/internal/com/client"
)

// List will print a list of active operations
func List(_ map[string]interface{}) error {
	cli.Start()

	operations, err := client.ListOperations()

	if err != nil {
		return err
	}

	for _, o := range operations {
		fmt.Printf(" %d  %s  %s  ", o.Index, o.ID, o.Title)

		for key, value := range o.Args {
			fmt.Printf("%s=%v ", key, value)
		}

		fmt.Printf("\n")
	}

	cli.Finish()

	return nil
}
