package runners

import (
	"barf/internal/cmd"
	"barf/internal/op"
	"barf/internal/rsync"
)

type pullRunner struct {
	operation     *op.Operation
	rsync         *rsync.Rsync
	stdoutHandler cmd.LogHandler
	stderrHandler cmd.LogHandler
	statusHandler statushandler
}

func (r *pullRunner) init(operation *op.Operation) {
	r.operation = operation
	r.rsync = rsync.NewRsync()
	r.rsync.OnStdout(r.handleStdout)
	r.rsync.OnStderr(r.handleStderr)
	r.rsync.OnStatus(r.handleStatus)
}

func (r *pullRunner) Start() {
	args := []string{
		"--delete-after",
	}
	src := []string{
		r.operation.Args["dst"].(string),
	}
	dst := r.operation.Args["src"].(string)

	r.rsync.Copy(args, src, dst)
}

func (r *pullRunner) Abort() error {
	return r.rsync.Abort()
}

func (r *pullRunner) OperationID() op.OperationID {
	return r.operation.ID
}

func (r *pullRunner) OnStdout(handler cmd.LogHandler) {
	r.stdoutHandler = handler
}

func (r *pullRunner) OnStderr(handler cmd.LogHandler) {
	r.stderrHandler = handler
}

func (r *pullRunner) OnStatus(handler statushandler) {
	r.statusHandler = handler
}

func (r *pullRunner) handleStdout(line string) {
	r.stdoutHandler(line)
}

func (r *pullRunner) handleStderr(line string) {
	r.stderrHandler(line)
}

func (r *pullRunner) handleStatus(status *op.OperationStatus) {
	r.statusHandler(status)
}
