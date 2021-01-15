package op

import "time"

// OperationStatus represents the currents status for an operation
type OperationStatus struct {
	BytesDiffTotal int64 `json:"bytesDiffTotal"`
	BytesTotal     int64 `json:"bytesTotal"`
	BytesDone      int64 `json:"bytesDone"`

	FilesDiffTotal int64 `json:"filesDiffTotal"`
	FilesTotal     int64 `json:"filesTotal"`
	FilesDone      int64 `json:"filesDone"`

	Progress float64 `json:"progress"`
	Speed    float64 `json:"speed"`

	StartedTimestamp  int64 `json:"startedTimestamp"`
	FinishedTimestamp int64 `json:"finishedTimestamp"`
	SecondsLeft       int64 `json:"secondsLeft"`

	FileName string `json:"fileName"`
	Finished bool   `json:"finished"`
	Message  string `json:"message"`
	ExitCode int    `json:"exitCode"`
}

// NewStatus creates a operation status object with default values
func NewStatus() *OperationStatus {
	now := time.Now()

	return &OperationStatus{
		StartedTimestamp: now.Unix(),
		Finished:         false,
		ExitCode:         -1,
	}
}

// UpdateStatus updates the first supplied status with values from the second
func UpdateStatus(a *OperationStatus, b *OperationStatus) {
	if b.BytesTotal > 0 {
		a.BytesTotal = b.BytesTotal
	}

	if b.BytesDiffTotal > 0 {
		a.BytesDiffTotal = b.BytesDiffTotal
	}

	if b.BytesDone > 0 {
		a.BytesDone = b.BytesDone
	}

	if b.Progress > 0 {
		a.Progress = b.Progress
	}

	if b.Speed > 0 {
		a.Speed = b.Speed
	}

	if b.FilesTotal > 0 {
		a.FilesTotal = b.FilesTotal
	}

	if b.FilesDiffTotal > 0 {
		a.FilesDiffTotal = b.FilesDiffTotal
	}

	if b.FilesDone > 0 {
		a.FilesDone = b.FilesDone
	}

	if b.SecondsLeft > 0 {
		a.SecondsLeft = b.SecondsLeft
	}

	if len(b.FileName) > 0 {
		a.FileName = b.FileName
	}

	if len(b.Message) > 0 {
		a.Message = b.Message
	}

	if b.Finished {
		a.Finished = b.Finished
	}

	if a.Finished && a.FinishedTimestamp == 0 {
		now := time.Now()
		a.SecondsLeft = 0
		a.FinishedTimestamp = now.Unix()
	}

	if b.ExitCode > a.ExitCode {
		a.ExitCode = b.ExitCode
	}
}
