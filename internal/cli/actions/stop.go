package actions

import "barf/internal/proc/daemonctrl"

// Stop stops the background process if it is running
func Stop(args map[string]interface{}) error {
	daemonctrl.Stop()

	return nil
}
