package journal

import (
	"log"
	"os"
	"path"

	"barf/internal/op"
)

type operationLog struct {
	close  func() error
	stdout func(...interface{})
	stderr func(...interface{})
}

func openLog(dir string, operationID op.OperationID) (*operationLog, error) {
	filename := path.Join(dir, string(operationID), "output.log")

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	if err != nil {
		return nil, err
	}

	stdout := log.New(f, " OUT  ", log.LstdFlags|log.LUTC|log.Lmsgprefix)
	stderr := log.New(f, " ERR  ", log.LstdFlags|log.LUTC|log.Lmsgprefix)

	return &operationLog{
		close:  f.Close,
		stdout: stdout.Println,
		stderr: stderr.Println,
	}, nil
}
