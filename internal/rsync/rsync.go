package rsync

import (
	"strings"
	"sync"

	"barf/internal/cmd"
	"barf/internal/op"
)

type statusHandler func(status *op.OperationStatus)

// Rsync is the rsync execution object
type Rsync struct {
	cmd           *cmd.Cmd
	stdoutHandler cmd.LogHandler
	stderrHandler cmd.LogHandler
	statusHandler statusHandler
	status        op.OperationStatus

	speed                []float64
	BytesDoneWholeFiles  int64
	BytesDoneCurrentFile int64

	mu sync.Mutex
}

// NewRsync creates a new rsync object
func NewRsync() *Rsync {
	r := &Rsync{
		status: op.OperationStatus{
			Finished: false,
			ExitCode: -1,
		},
	}

	return r
}

// Copy will start a copy operation
func (r *Rsync) Copy(extraArgs []string, src []string, dst string) {
	args := []string{}

	for _, value := range extraArgs {
		args = append(args, value)
	}

	for _, value := range src {
		args = append(args, strings.TrimSuffix(value, "/"))
	}

	args = append(args, strings.TrimSuffix(dst, "/"))

	exitCode, err := r.doPreparation(args)

	if err != nil {
		r.status.Finished = true
		r.status.Message = "Error: " + err.Error()
		r.status.ExitCode = exitCode
		r.emitStatus()

		return
	}

	if r.status.BytesDiffTotal == 0 {
		r.status.Progress = 100
		r.status.Finished = true
		r.status.Message = "No work needed!"
		r.status.ExitCode = exitCode
		r.emitStatus()

		return
	}

	exitCode, err = r.doCopy(args)

	r.status.Finished = true
	r.status.ExitCode = exitCode

	if err != nil {
		r.status.Message = "Error: " + err.Error()
	} else {
		r.status.Message = "Finished successfully!"
	}

	r.emitStatus()
}

// Move will start a move operation
func (r *Rsync) Move(extraArgs []string, src []string, dst string) {
	// TODO: Evaluate if there could be a special case for local src and dst on the same disk which does not use rsync but mv. This will not take 2 times the disk and be very much faster... how to make resumable though...

	args := []string{}
	srcTrimmed := []string{}

	for _, value := range extraArgs {
		args = append(args, value)
	}

	for _, value := range src {
		trimmed := strings.TrimSuffix(value, "/")

		srcTrimmed = append(srcTrimmed, trimmed)
		args = append(args, trimmed)
	}

	args = append(args, strings.TrimSuffix(dst, "/"))

	_, err := r.doPreparation(args)

	if err != nil {
		return
	}

	if r.status.BytesDiffTotal > 0 {
		_, err = r.doCopy(args)

		if err != nil {
			return
		}
	}

	r.status.Progress = 100

	exitCode, err := r.doDelete(srcTrimmed)

	r.status.Finished = true
	r.status.ExitCode = exitCode

	if err != nil {
		r.status.Message = "Error: " + err.Error()
	} else {
		r.status.Message = "Finished successfully!"
	}

	r.emitStatus()
}

// Abort will abort the current command
func (r *Rsync) Abort() error {
	if r.cmd != nil {
		return r.cmd.Abort()
	}

	return nil
}

// OnStdout registers a stdout handler for log lines
func (r *Rsync) OnStdout(handler cmd.LogHandler) {
	r.stdoutHandler = handler
}

// OnStderr registers a stderr handler for log lines
func (r *Rsync) OnStderr(handler cmd.LogHandler) {
	r.stderrHandler = handler
}

// OnStatus registers an status handler for the status event
func (r *Rsync) OnStatus(handler statusHandler) {
	r.statusHandler = handler
}
