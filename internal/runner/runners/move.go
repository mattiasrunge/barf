package runners

import (
	"barf/internal/cmd"
	"barf/internal/op"
	"barf/internal/rsync"
	"barf/internal/typeconv"
)

type moveRunner struct {
	operation     *op.Operation
	rsync         *rsync.Rsync
	stdoutHandler cmd.LogHandler
	stderrHandler cmd.LogHandler
	statusHandler statushandler
}

func (r *moveRunner) init(operation *op.Operation) {
	r.operation = operation
	r.rsync = rsync.NewRsync()
	r.rsync.OnStdout(r.handleStdout)
	r.rsync.OnStderr(r.handleStderr)
	r.rsync.OnStatus(r.handleStatus)
}

func (r *moveRunner) Start() {
	args := []string{}
	srcArray, _ := typeconv.ToArray(r.operation.Args["src"])
	src := typeconv.ToStringArray(srcArray)
	dst := r.operation.Args["dst"].(string)

	r.rsync.Move(args, src, dst)
}

func (r *moveRunner) Abort() error {
	return r.rsync.Abort()
}

func (r *moveRunner) OperationID() op.OperationID {
	return r.operation.ID
}

func (r *moveRunner) OnStdout(handler cmd.LogHandler) {
	r.stdoutHandler = handler
}

func (r *moveRunner) OnStderr(handler cmd.LogHandler) {
	r.stderrHandler = handler
}

func (r *moveRunner) OnStatus(handler statushandler) {
	r.statusHandler = handler
}

func (r *moveRunner) handleStdout(line string) {
	r.stdoutHandler(line)
}

func (r *moveRunner) handleStderr(line string) {
	r.stderrHandler(line)
}

func (r *moveRunner) handleStatus(status *op.OperationStatus) {
	r.statusHandler(status)
}
