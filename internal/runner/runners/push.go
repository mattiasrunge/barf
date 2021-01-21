package runners

import (
	"barf/internal/cmd"
	"barf/internal/op"
	"barf/internal/rsync"
)

type pushRunner struct {
	operation     *op.Operation
	rsync         *rsync.Rsync
	stdoutHandler cmd.LogHandler
	stderrHandler cmd.LogHandler
	statusHandler statushandler
}

func (r *pushRunner) init(operation *op.Operation) {
	r.operation = operation
	r.rsync = rsync.NewRsync()
	r.rsync.OnStdout(r.handleStdout)
	r.rsync.OnStderr(r.handleStderr)
	r.rsync.OnStatus(r.handleStatus)
}

func (r *pushRunner) Start() {
	args := []string{
		"--delete-after",
	}
	src := []string{
		r.operation.Args["src"].(string),
	}
	dst := r.operation.Args["dst"].(string)

	r.rsync.Copy(args, src, dst)
}

func (r *pushRunner) Abort() error {
	return r.rsync.Abort()
}

func (r *pushRunner) OperationID() op.OperationID {
	return r.operation.ID
}

func (r *pushRunner) OnStdout(handler cmd.LogHandler) {
	r.stdoutHandler = handler
}

func (r *pushRunner) OnStderr(handler cmd.LogHandler) {
	r.stderrHandler = handler
}

func (r *pushRunner) OnStatus(handler statushandler) {
	r.statusHandler = handler
}

func (r *pushRunner) handleStdout(line string) {
	r.stdoutHandler(line)
}

func (r *pushRunner) handleStderr(line string) {
	r.stderrHandler(line)
}

func (r *pushRunner) handleStatus(status *op.OperationStatus) {
	r.statusHandler(status)
}
