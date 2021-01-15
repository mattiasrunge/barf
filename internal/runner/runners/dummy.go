package runners

import (
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"barf/internal/cmd"
	"barf/internal/op"
)

type dummyRunner struct {
	operation     *op.Operation
	cmd           *cmd.Cmd
	stdoutHandler cmd.LogHandler
	stderrHandler cmd.LogHandler
	statusHandler statushandler
}

func (r *dummyRunner) init(operation *op.Operation) {
	r.operation = operation
	r.cmd = cmd.NewCmd()
	r.cmd.OnStdout(r.handleStdout)
	r.cmd.OnStderr(r.handleStderr)
	r.cmd.OnFinish(r.handleFinish)
}

func (r *dummyRunner) Start() {
	dirname, _ := os.Getwd()
	args := []string{
		path.Join(dirname, "scripts", "dummy.sh"),
		r.operation.Args["iterations"].(string),
	}

	r.cmd.Start(args)
}

func (r *dummyRunner) Abort() error {
	return r.cmd.Abort()
}

func (r *dummyRunner) OperationID() op.OperationID {
	return r.operation.ID
}

func (r *dummyRunner) OnStdout(handler cmd.LogHandler) {
	r.stdoutHandler = handler
}

func (r *dummyRunner) OnStderr(handler cmd.LogHandler) {
	r.stderrHandler = handler
}

func (r *dummyRunner) OnStatus(handler statushandler) {
	r.statusHandler = handler
}

func (r *dummyRunner) handleStdout(line string) {
	parts := strings.Split(line, "/")
	current, _ := strconv.Atoi(parts[0])
	total, _ := strconv.Atoi(parts[1])

	bytesDone := int64(current) * 1024
	speed := float64(1024)
	progress := (float64(current) / float64(total)) * 100
	bytesTotal := int64(float64(bytesDone) / (float64(progress) / 100))
	bytesLeft := bytesTotal - bytesDone
	left := int64(float64(bytesLeft) / speed)

	r.statusHandler(&op.OperationStatus{
		BytesTotal:     bytesTotal,
		BytesDiffTotal: bytesTotal,
		BytesDone:      bytesDone,
		Progress:       math.Round(progress*100) / 100,
		Speed:          math.Round(speed*100) / 100,
		FilesTotal:     int64(total),
		FilesDiffTotal: int64(total),
		FilesDone:      int64(current),
		SecondsLeft:    left,
	})

	r.stdoutHandler(line)
}

func (r *dummyRunner) handleStderr(line string) {
	r.stderrHandler(line)
}

func (r *dummyRunner) handleFinish(exitCode int, errorMsg string) {
	r.statusHandler(&op.OperationStatus{
		Finished: true,
		ExitCode: exitCode,
		Message:  errorMsg,
	})
}
