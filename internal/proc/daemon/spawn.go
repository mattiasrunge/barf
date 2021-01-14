package daemon

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"barf/internal/config"
	"barf/internal/proc/pidfile"
)

func isAlive(pid int) bool {
	process, err := os.FindProcess(pid)

	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))

	return err == nil
}

func waitForSocket() error {
	for tries := 0; tries < 10; tries++ {
		if _, err := os.Stat(config.SocketFile); err == nil {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return errors.New("Could not find socket")
}

// IsRunning checks if daemon process is running
func IsRunning() bool {
	pid, err := pidfile.Read()

	if err != nil {
		return false
	}

	return isAlive(pid)
}

// Spawn spawns a new daemon process if none is running
func Spawn() error {
	if IsRunning() {
		return nil
	}

	err := os.RemoveAll(config.SocketFile)

	if err != nil {
		return err
	}

	fmt.Println("No daemon process found, will try to start...")

	newEnv := append(os.Environ(), config.DaemonVariable+"="+config.DaemonVariable)
	sysproc := &syscall.SysProcAttr{Setsid: true}
	attr := os.ProcAttr{
		Dir: ".",
		Env: newEnv,
		Files: []*os.File{
			nil,
			nil,
			nil,
		},
		Sys: sysproc,
	}
	dirname, _ := os.Getwd()
	binary := os.Args[0]

	fmt.Println("Binary:", binary)
	if strings.Contains(binary, "/go-build") {
		binary = "barf.sh" // For development
	}

	var executable = path.Join(dirname, binary)

	process, err := os.StartProcess(executable, []string{executable, "background"}, &attr)

	if err != nil {
		return err
	}

	err = process.Release()

	if err != nil {
		return err
	}

	return waitForSocket()
}
