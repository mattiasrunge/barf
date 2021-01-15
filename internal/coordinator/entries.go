package coordinator

import (
	"errors"
	"sync"

	"barf/internal/journal"
	"barf/internal/op"
)

var entries []*journal.JournalEntry
var entriesMu sync.Mutex

func addEntry(e *journal.JournalEntry) {
	entriesMu.Lock()
	defer entriesMu.Unlock()

	entries = append(entries, e)
}

func removeEntry(e *journal.JournalEntry) {
	entriesMu.Lock()
	defer entriesMu.Unlock()

	for i, s := range entries {
		if s == e {
			copy(entries[i:], entries[i+1:])
			entries = entries[:len(entries)-1]

			return
		}
	}
}

func getEntry(operationID op.OperationID) *journal.JournalEntry {
	entriesMu.Lock()
	defer entriesMu.Unlock()

	for _, e := range entries {
		if e.Operation.ID == operationID {
			return e
		}
	}

	return nil
}

func getEntries() []*journal.JournalEntry {
	entriesMu.Lock()
	defer entriesMu.Unlock()

	result := make([]*journal.JournalEntry, len(entries))
	copy(result, entries)

	return result
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
