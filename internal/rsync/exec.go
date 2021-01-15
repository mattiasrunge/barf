package rsync

import (
	"barf/internal/cmd"
	"barf/internal/utils"
	"fmt"
	"sync"
)

type RsyncStatus struct {
	Message string

	BytesDiffTotal int64
	BytesTotal     int64
	FilesDiffTotal int64
	FilesTotal     int64

	BytesDone            int64
	BytesDoneTotal       int64
	CurrentFileIndex     int64
	CurrentFileBytesDone int64
	CurrentFileName      string

	Progress    float64
	Speed       float64
	SecondsLeft int64

	Finished bool
	ExitCode int
}

type StatusHandler func(status *RsyncStatus)

type Rsync struct {
	cmd           *cmd.Cmd
	stdoutHandler cmd.LogHandler
	stderrHandler cmd.LogHandler
	statusHandler StatusHandler
	status        RsyncStatus

	speed []float64
	mu    sync.Mutex
}

// NewRsync creates a new rsync object
func NewRsync() *Rsync {
	r := &Rsync{
		status: RsyncStatus{
			Finished: false,
			ExitCode: -1,
		},
	}

	return r
}

func (r *Rsync) getArgs(stepArgs []string, operationArgs []string) []string {
	args := []string{
		"rsync",
		// "--checksum", // Very slow!
		// "--compress", // TODO: Use for remote
		"--executability",
		"--archive", // -rlptgoD
		"--partial",
		"--inplace",
		"--no-whole-file",
	}

	for _, arg := range stepArgs {
		args = append(args, arg)
	}

	for _, arg := range operationArgs {
		args = append(args, arg)
	}

	return args
}

func (r *Rsync) getMedianSpeed(newSpeed float64) float64 {
	if newSpeed == 0 {
		return 0
	}

	r.speed = append(r.speed, newSpeed)

	if len(r.speed) > 4 {
		r.speed = r.speed[1:]
	}

	return utils.Median(r.speed)
}

func (r *Rsync) parseProgressLine(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fileName, isDir, ok := parseFileName(line)

	if ok {
		r.status.BytesDone += r.status.CurrentFileBytesDone
		r.status.CurrentFileBytesDone = 0

		if !isDir {
			r.status.CurrentFileName = fileName
			r.status.CurrentFileIndex++
			r.status.Message = "Processing " + fileName + "..."
		}

		r.emitStatus()
	} else {
		fileBytesDone, speed, ok := parseProgress(line)

		if ok {
			r.status.CurrentFileBytesDone = fileBytesDone
			r.status.Speed = r.getMedianSpeed(speed)
			r.status.BytesDoneTotal = r.status.BytesDone + r.status.CurrentFileBytesDone
			r.status.Progress = (float64(r.status.BytesDoneTotal) / float64(r.status.BytesDiffTotal)) * 100
			r.status.SecondsLeft = int64(float64(r.status.BytesDiffTotal-r.status.BytesDoneTotal) / r.status.Speed)

			r.emitStatus()
		}
	}

	if r.stdoutHandler != nil {
		r.stdoutHandler(line)
	}
}

func (r *Rsync) parsePreparationLine(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.status.FilesTotal = parseNumberOfFiles(line, r.status.FilesTotal)
	r.status.FilesDiffTotal = parseNumberOfCreatedFiles(line, r.status.FilesDiffTotal)
	r.status.BytesTotal = parseTotalFileSize(line, r.status.BytesTotal)
	r.status.BytesDiffTotal = parseTotalTransferredFileSize(line, r.status.BytesDiffTotal)

	if r.stdoutHandler != nil {
		r.stdoutHandler(line)
	}
}

func (r *Rsync) handleStderrLine(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.stderrHandler != nil {
		r.stderrHandler(line)
	}
}

func (r *Rsync) doPreparation(operationArgs []string) error {
	r.status.Message = "Preparing..."
	r.emitStatus()

	args := r.getArgs([]string{
		"--dry-run",
		"--stats",
	}, operationArgs)

	r.cmd = cmd.NewCmd()
	r.cmd.OnStdout(r.parsePreparationLine)
	r.cmd.OnStderr(r.handleStderrLine)

	exitCode, err := r.cmd.Start(args)
	r.cmd = nil

	if err != nil {
		r.status.Finished = true
		r.status.Message = err.Error()
		r.status.ExitCode = exitCode
		r.emitStatus()
	} else if r.status.BytesDiffTotal == 0 {
		fmt.Println("No bytes found that needs transfer, will do nothing")
		r.status.Progress = 100
		r.status.Finished = true
		r.status.Message = "No work needed!"
		r.status.ExitCode = exitCode
		r.emitStatus()
	}

	return err
}

func (r *Rsync) doSync(operationArgs []string) error {
	r.status.Message = "Processing..."
	r.emitStatus()

	args := r.getArgs([]string{
		"--progress",
		"--out-format=__file:%n",
	}, operationArgs)

	r.cmd = cmd.NewCmd()
	r.cmd.OnStdout(r.parseProgressLine)
	r.cmd.OnStderr(r.handleStderrLine)

	exitCode, err := r.cmd.Start(args)
	r.cmd = nil

	r.status.Finished = true
	r.status.ExitCode = exitCode

	if err != nil {
		r.status.Message = err.Error()
	} else {
		r.status.Message = "Finished successfully!"
	}

	r.emitStatus()

	return err
}

func (r *Rsync) emitStatus() {
	if r.statusHandler != nil {
		r.statusHandler(&r.status)
	}
}

func (r *Rsync) Copy(operationArgs []string) {
	err := r.doPreparation(operationArgs)

	if err != nil || r.status.Finished {
		return
	}

	_ = r.doSync(operationArgs)
}

func (r *Rsync) Abort() error {
	if r.cmd != nil {
		return r.cmd.Abort()
	}

	return nil
}

func (r *Rsync) OnStdout(handler cmd.LogHandler) {
	r.stdoutHandler = handler
}

func (r *Rsync) OnStderr(handler cmd.LogHandler) {
	r.stderrHandler = handler
}

func (r *Rsync) OnStatus(handler StatusHandler) {
	r.statusHandler = handler
}
