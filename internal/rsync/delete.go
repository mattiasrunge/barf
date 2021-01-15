package rsync

import (
	"os"
	"strings"

	"barf/internal/cmd"
)

func (r *Rsync) doDelete(entities []string) (int, error) {
	r.status.Message = "Cleanup..."
	r.emitStatus()

	// TODO: Sort and look for top directories, remove those and skip stuff inside them

	for _, entity := range entities {
		if isRemote(entity) {
			parts := strings.Split(entity, ":")

			args := []string{
				"ssh",
				parts[0],
				"rm -rf " + parts[1],
			}

			r.cmd = cmd.NewCmd()
			r.cmd.OnStdout(r.handleStdoutLine)
			r.cmd.OnStderr(r.handleStderrLine)

			exitCode, err := r.cmd.Start(args)
			r.cmd = nil

			if err != nil {
				return exitCode, err
			}
		} else {
			r.handleStdoutLine("Removing " + entity)
			err := os.RemoveAll(entity)

			if err != nil {
				return 255, err
			}
		}
	}

	return 0, nil
}
