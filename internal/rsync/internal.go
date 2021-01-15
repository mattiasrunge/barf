package rsync

import "barf/internal/utils"

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

func (r *Rsync) handleStderrLine(line string) {
	if r.stderrHandler != nil {
		r.stderrHandler(line)
	}
}

func (r *Rsync) handleStdoutLine(line string) {
	if r.stdoutHandler != nil {
		r.stdoutHandler(line)
	}
}

func (r *Rsync) emitStatus() {
	if r.statusHandler != nil {
		r.statusHandler(&r.status)
	}
}
