package run

import (
	"fmt"

	"rft/internal/com/socket"
	"rft/internal/proc/daemon"
	"rft/internal/ui"

	cli "github.com/jawher/mow.cli"
)

// StartCLI starts the CLI process
func StartCLI(action func() error) {
	normalClose := false
	err := daemon.Spawn()

	if err != nil {
		fmt.Println(err)
		cli.Exit(255)
		return
	}

	err = socket.Connect()

	if err != nil {
		fmt.Println(err)
		cli.Exit(255)
		return
	}

	socket.OnClose(func() {
		if !normalClose {
			fmt.Println("Lost connection to backend")
			cli.Exit(1)
		}
	})

	ui.Start()

	err = action()

	if err != nil {
		fmt.Println(err)
		cli.Exit(255)
		return
	}

	exitCode := ui.Wait()

	normalClose = true
	socket.Close()
	socket.WaitOnClose()

	cli.Exit(exitCode)
}
