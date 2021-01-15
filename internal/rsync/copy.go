package rsync

import "barf/internal/cmd"

func (r *Rsync) parseProgressLine(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fileName, isDir, ok := parseFileName(line)

	if ok {
		r.BytesDoneWholeFiles += r.BytesDoneCurrentFile
		r.BytesDoneCurrentFile = 0

		r.status.BytesDone = r.BytesDoneWholeFiles

		if !isDir {
			r.status.FileName = fileName
			r.status.FilesDone++
			r.status.Message = "Processing " + fileName + "..."
		}

		r.emitStatus()
	} else {
		fileBytesDone, speed, ok := parseProgress(line)

		if ok {
			r.BytesDoneCurrentFile = fileBytesDone
			r.status.Speed = r.getMedianSpeed(speed)
			r.status.BytesDone = r.BytesDoneWholeFiles + r.BytesDoneCurrentFile
			r.status.Progress = (float64(r.status.BytesDone) / float64(r.status.BytesDiffTotal)) * 100
			r.status.SecondsLeft = int64(float64(r.status.BytesDiffTotal-r.status.BytesDone) / r.status.Speed)

			r.emitStatus()
		}
	}

	r.handleStdoutLine(line)
}

func (r *Rsync) doCopy(operationArgs []string) (int, error) {
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

	return exitCode, err
}
