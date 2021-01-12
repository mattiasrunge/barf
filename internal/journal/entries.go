package journal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"rft/internal/op"
)

type entry struct {
	Operation *op.Operation       `json:"operation"`
	Status    *op.OperationStatus `json:"status"`
}

var entries []*entry

func addEntry(e *entry) error {
	index, err := getNextIndex()

	if err != nil {
		return err
	}

	e.Operation.Index = index
	entries = append(entries, e)

	return writeEntry(e)
}

func removeEntry(e *entry) error {
	for i, s := range entries {
		if s == e {
			copy(entries[i:], entries[i+1:])
			entries = entries[:len(entries)-1]

			return finishEntry(e)
		}
	}

	return nil
}

func getEntry(operationID op.OperationID) (*entry, error) {
	for _, e := range entries {
		if e.Operation.ID == operationID {
			return e, nil
		}
	}

	return nil, errors.New("No such operation found")
}

func writeEntry(e *entry) error {
	bytes, err := json.Marshal(e)

	if err != nil {
		return err
	}

	file := path.Join(activeDir, string(e.Operation.ID)+".json")

	return ioutil.WriteFile(file, bytes, 0600)
}

func readEntry(operationID op.OperationID) (*entry, error) {
	var e entry
	file := path.Join(activeDir, string(operationID)+".json")
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return &e, err
	}

	err = json.Unmarshal(data, &e)

	return &e, err
}

func finishEntry(e *entry) error {
	from := path.Join(activeDir, string(e.Operation.ID)+".json")
	to := path.Join(historyDir, string(e.Operation.ID)+".json")

	return os.Rename(from, to)
}

func getNextIndex() (op.OperationIndex, error) {
	var index op.OperationIndex = 1

	indexFree := func(index op.OperationIndex) bool {
		for _, e := range entries {
			if e.Operation.Index == index {
				return false
			}
		}

		return true
	}

	for ; !indexFree(index); index++ {
		if index == 0 {
			return 0, errors.New("Could not get new index, overflow")
		}
	}

	return index, nil
}

// ReadFromDisk reads active entries from disk
func ReadFromDisk() error {
	if registerdStartHandler == nil {
		return errors.New("No start handler registered")
	}

	entries = nil
	fileInfo, err := ioutil.ReadDir(activeDir)

	if err != nil {
		return err
	}

	for _, file := range fileInfo {
		operationID := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		e, err := readEntry(op.OperationID(operationID))

		if err != nil {
			return err
		}

		entries = append(entries, e)
	}

	for _, e := range entries {
		err = registerdStartHandler(e.Operation)

		if err != nil {
			UpdateOperationStatus(e.Operation.ID, &op.OperationStatus{
				Finished: true,
				ExitCode: 255,
				Error:    err.Error(),
			})
			return err
		}
	}

	return nil
}
