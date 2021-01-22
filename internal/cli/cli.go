package cli

import (
	"fmt"
	"os"

	"barf/internal/cli/ui"
	"barf/internal/com/socket"
	"barf/internal/proc/daemonctrl"
)

// Action is the action function type
type Action func(args map[string]interface{}) error

// Start initializes the CLI and connects to the daemon
func Start() {
	err := daemonctrl.Spawn()

	if err != nil {
		fmt.Println(err)
		os.Exit(255)
	}

	err = socket.Connect()

	if err != nil {
		fmt.Println(err)
		os.Exit(255)
	}

	socket.OnClose(func(normalClose bool) {
		if !normalClose {
			fmt.Println("Lost connection to backend")
			os.Exit(1)
		}
	})

	ui.Start()
}

// Finish waits for the CLI to finish and exits
func Finish() {
	exitCode := ui.Wait()

	socket.Close()
	socket.WaitOnClose()

	os.Exit(exitCode)
}
