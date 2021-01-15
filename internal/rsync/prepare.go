package rsync

import (
	"barf/internal/cmd"
)

func (r *Rsync) parsePreparationLine(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.status.FilesTotal = parseNumberOfFiles(line, r.status.FilesTotal)
	r.status.FilesDiffTotal = parseNumberOfCreatedFiles(line, r.status.FilesDiffTotal)
	r.status.BytesTotal = parseTotalFileSize(line, r.status.BytesTotal)
	r.status.BytesDiffTotal = parseTotalTransferredFileSize(line, r.status.BytesDiffTotal)

	r.handleStdoutLine(line)
}

func (r *Rsync) doPreparation(operationArgs []string) (int, error) {
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

	return exitCode, err
}
