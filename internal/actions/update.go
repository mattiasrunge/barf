package actions

import "barf/internal/update"

// Update checks for an update
func Update(args map[string]interface{}) error {
	err := update.Update()

	if err != nil {
		return err
	}

	// TODO: Restart daemon

	return nil
}
