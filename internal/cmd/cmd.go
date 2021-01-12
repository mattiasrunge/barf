package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"rft/internal/utils"
)

// LogHandler is the definition for a log line handler func
type LogHandler func(string)
type finishhandler func(int, string)

// Cmd represents a command
type Cmd struct {
	wg            sync.WaitGroup
	cmd           *exec.Cmd
	stdoutHandler LogHandler
	stderrHandler LogHandler
	finishHandler finishhandler
}

// NewCmd creates a new cmd
func NewCmd() *Cmd {
	return &Cmd{}
}

// Start starts the execution of the process returns error and exitCode
func (cmd *Cmd) Start(args []string) (int, error) {
	if len(args) < 1 {
		err := errors.New("args must at least include the executable name")
		cmd.handleFinish(255, err.Error())
		return 255, err
	}

	cmd.cmd = exec.Command(args[0], args[1:]...)
	var wg sync.WaitGroup

	stdout, _ := cmd.cmd.StdoutPipe()
	stderr, _ := cmd.cmd.StderrPipe()

	go cmd.handleLog(stdout, cmd.stdoutHandler)
	go cmd.handleLog(stderr, cmd.stderrHandler)

	err := cmd.cmd.Start()

	if err != nil {
		cmd.cmd = nil
		cmd.handleFinish(255, err.Error())
		return 255, err
	}

	wg.Wait()

	err = cmd.cmd.Wait()

	cmd.cmd = nil

	if err != nil {
		// https://stackoverflow.com/questions/10385551/get-exit-code-go
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode := status.ExitStatus()
				cmd.handleFinish(exitCode, exiterr.Error())
				return exitCode, err
			}
		}
	}

	cmd.handleFinish(0, "")
	return 0, nil
}

// Abort tries to abort the process
func (cmd *Cmd) Abort() error {
	if cmd.cmd != nil {
		return cmd.cmd.Process.Signal(syscall.SIGTERM)
	}

	return nil
}

// OnStdout register handler for lines on stdout
func (cmd *Cmd) OnStdout(handler LogHandler) {
	cmd.stdoutHandler = handler
}

// OnStderr registers handler for lines on stderr
func (cmd *Cmd) OnStderr(handler LogHandler) {
	cmd.stderrHandler = handler
}

// OnFinish registers handler for when cmd has finished
func (cmd *Cmd) OnFinish(handler finishhandler) {
	cmd.finishHandler = handler
}

func (cmd *Cmd) handleFinish(exitCode int, errorMsg string) {
	if cmd.finishHandler != nil {
		cmd.finishHandler(exitCode, errorMsg)
	}
}

func (cmd *Cmd) handleLog(r io.Reader, handler LogHandler) {
	if handler == nil {
		return
	}

	cmd.wg.Add(1)
	scanner := bufio.NewScanner(r)

	scanner.Split(utils.ScanLines)

	for scanner.Scan() {
		handler(scanner.Text())
	}

	err := scanner.Err()

	if err != nil && !strings.Contains(err.Error(), "file already closed") {
		fmt.Println("Cmd:", err)
	}

	cmd.wg.Done()
}
