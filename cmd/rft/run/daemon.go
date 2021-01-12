package run

import (
	"fmt"
	"os"

	"rft/internal/com/socket"
	"rft/internal/config"
	"rft/internal/proc/daemon"
	"rft/internal/proc/life"
	"rft/internal/proc/pidfile"
	"rft/internal/runner"
)

// CheckDaemon starts the background daemon process if that is set
func CheckDaemon() {
	if os.Getenv(config.DaemonVariable) == config.DaemonVariable {
		err := startDaemon()

		if err != nil {
			fmt.Println(err)
			life.RunShutdownHooks()
			os.Exit(1)
		}

		os.Exit(0)
	}
}

func startDaemon() error {
	if daemon.IsRunning() {
		fmt.Println("Daemon process already running...")
		return nil
	}

	life.Start()

	// logfile.StartLogging()
	err := pidfile.Write(os.Getpid())

	if err != nil {
		return nil
	}

	life.AddShutdownHook(func() {
		pidfile.Delete()
	})

	fmt.Println("Daemon started")

	err = socket.Listen()

	if err != nil {
		return err
	}

	fmt.Println("Listening for connections")

	err = runner.StartRunner()

	if err != nil {
		return err
	}

	fmt.Println("Runner started")

	socket.WaitOnClose()

	fmt.Println("Daemon finished")

	return nil
}
