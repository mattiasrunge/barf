package actions

import (
	"fmt"

	"barf/internal/proc/daemon"
	"barf/internal/update"
)

// Update checks for an update
func Update(args map[string]interface{}) error {
	err := update.Update()

	if err != nil {
		return err
	}

	fmt.Println("Stopping daemon process...")
	daemon.Stop()

	fmt.Println("Starting daemon process...")
	return daemon.Spawn()
}
