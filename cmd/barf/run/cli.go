package run

import (
	"fmt"
	"os"

	"barf/internal/com/socket"
	"barf/internal/proc/daemon"
	"barf/internal/ui"
)

// StartCLI starts the CLI process
func StartCLI(width int, action func() error) {
	normalClose := false
	err := daemon.Spawn()

	if err != nil {
		fmt.Println(err)
		os.Exit(255)
		return
	}

	err = socket.Connect()

	if err != nil {
		fmt.Println(err)
		os.Exit(255)
		return
	}

	socket.OnClose(func() {
		if !normalClose {
			fmt.Println("Lost connection to backend")
			os.Exit(1)
		}
	})

	ui.Start(width)

	err = action()

	if err != nil {
		fmt.Println(err)
		os.Exit(255)
		return
	}

	exitCode := ui.Wait()

	normalClose = true
	socket.Close()
	socket.WaitOnClose()

	os.Exit(exitCode)
}
