package journal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"barf/internal/op"
)

// JournalEntry represents an operation and its status in the journal
type JournalEntry struct {
	Operation *op.Operation
	Status    *op.OperationStatus
	log       *operationLog
	mutex     sync.Mutex
}

// NewJournalEntry creates a new journal entry based on the supplied operation
func NewJournalEntry(operation *op.Operation) (*JournalEntry, error) {
	err := writeOperation(activeDir, operation)

	if err != nil {
		return nil, err
	}

	log, err := openLog(activeDir, operation.ID)

	if err != nil {
		return nil, err
	}

	e := JournalEntry{
		Operation: operation,
		Status:    op.NewStatus(),
		log:       log,
	}

	return &e, nil
}

func readJournalEntry(operationID op.OperationID) (*JournalEntry, error) {
	operation, err := readOperation(activeDir, operationID)

	if err != nil {
		return nil, err
	}

	status := readStatus(activeDir, operationID)

	log, err := openLog(activeDir, operation.ID)

	if err != nil {
		return nil, err
	}

	e := JournalEntry{
		Operation: operation,
		Status:    status,
		log:       log,
	}

	return &e, nil
}

func loadActiveJournalEntries() ([]*JournalEntry, error) {
	var entries []*JournalEntry
	fileInfo, err := ioutil.ReadDir(activeDir)

	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		e, err := readJournalEntry(op.OperationID(file.Name()))

		if err != nil {
			fmt.Printf("Could not read active journal entry for %s, will skip it..\n", file.Name())
			fmt.Println(err)
			continue
		}

		entries = append(entries, e)
	}

	return entries, nil
}

// UpdateStatus updates the status of the operation
func (e *JournalEntry) UpdateStatus(status *op.OperationStatus) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	op.UpdateStatus(e.Status, status)

	err := writeStatus(activeDir, e.Operation.ID, e.Status)

	if err != nil {
		return err
	}

	if e.Status.Finished {
		return e.archive()
	}

	return nil
}

// WriteStdout writes to the log of the operation
func (e *JournalEntry) WriteStdout(line string) {
	e.log.stdout(line)
}

// WriteStderr writes to the log of the operation
func (e *JournalEntry) WriteStderr(line string) {
	e.log.stderr(line)
}

func (e *JournalEntry) archive() error {
	err := e.log.close()

	if err != nil {
		return err
	}

	from := path.Join(activeDir, string(e.Operation.ID))
	to := path.Join(historyDir, string(e.Operation.ID))

	return os.Rename(from, to)
}
