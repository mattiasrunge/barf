package actions

import (
	"fmt"

	"barf/internal/proc/daemonctrl"
	"barf/internal/update"
)

// Update checks for an update
func Update(args map[string]interface{}) error {
	err := update.Update()

	if err != nil {
		return err
	}

	fmt.Println("Stopping daemon process...")
	daemonctrl.Stop()

	fmt.Println("Starting daemon process...")
	return daemonctrl.Spawn()
}
