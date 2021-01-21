package actions

import "barf/internal/proc/daemon"

// Stop stops the background process if it is running
func Stop(args map[string]interface{}) error {
	daemon.Stop()

	return nil
}
