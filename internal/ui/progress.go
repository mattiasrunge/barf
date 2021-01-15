package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"barf/internal/utils"

	"github.com/acarl005/stripansi"
	"github.com/logrusorgru/aurora"
	"github.com/mattiasrunge/goterminal"
)

var writer *goterminal.Writer
var width = 132
var mu sync.Mutex

func replaceAtIndex(in string, r rune, i int) string {
	if i < len(in) {
		out := []rune(in)
		out[i] = r
		return string(out)
	}

	return in
}

func sLen(s string) int {
	return len(stripansi.Strip(s))
}

func getTitle(o *operationWithStatus) string {
	return aurora.BrightMagenta(string(o.operation.Title)).Bold().String()
}

func getIndex(o *operationWithStatus) string {
	return fmt.Sprintf("%d", o.operation.Index)
}

func getMessage(o *operationWithStatus) string {
	return o.status.Message
}

func getFileName(o *operationWithStatus) string {
	return o.status.FileName
}

func getByteProgress(o *operationWithStatus) string {
	totalDiffStr := utils.ByteCountSI(o.status.BytesDiffTotal)
	totalStr := utils.ByteCountSI(o.status.BytesTotal)
	doneStr := aurora.BrightGreen(utils.ByteCountSI(o.status.BytesDone)).String()

	if o.status.Finished && o.status.ExitCode == 0 {
		totalDiffStr = aurora.BrightGreen(totalDiffStr).String()
		totalStr = aurora.BrightGreen(totalStr).String()
	}

	result := fmt.Sprintf("%s / %s", doneStr, totalDiffStr)

	if o.status.BytesDiffTotal != o.status.BytesTotal {
		return fmt.Sprintf("%s(%s)", result, totalStr)
	}

	return result
}

func getFileProgress(o *operationWithStatus) string {
	totalDiffStr := strconv.Itoa(int(o.status.FilesDiffTotal))
	totalStr := strconv.Itoa(int(o.status.FilesTotal))
	doneStr := aurora.BrightGreen(strconv.Itoa(int(o.status.FilesDone))).String()

	if o.status.Finished && o.status.ExitCode == 0 {
		totalDiffStr = aurora.BrightGreen(totalDiffStr).String()
		totalStr = aurora.BrightGreen(totalStr).String()
	}

	result := fmt.Sprintf("%s / %s", doneStr, totalDiffStr)

	if o.status.FilesDiffTotal != o.status.FilesTotal {
		return fmt.Sprintf("%s(%s)", result, totalStr)
	}

	return result
}

func getProgress(o *operationWithStatus) string {
	return fmt.Sprintf("%d%%", int(math.Ceil(o.status.Progress)))
}

func getSpeed(o *operationWithStatus) string {
	return fmt.Sprintf("%s/s", utils.ByteCountSI(int64(o.status.Speed)))
}

func getTimeInfo(o *operationWithStatus) string {
	if o.status.Finished {
		duration := time.Duration((o.status.FinishedTimestamp - o.status.StartedTimestamp) * int64(time.Second))

		return fmt.Sprintf("%s", aurora.Bold(duration.String()).String())
	}

	leftDuration := time.Duration(o.status.SecondsLeft * int64(time.Second))

	return fmt.Sprintf("ETA: %s", aurora.Bold(leftDuration.String()).String())
}

func getProgressBar(o *operationWithStatus, width int, text string) string {
	barDoneChar := " "
	barLeftChar := " "
	barPrefix := ""
	barSuffix := ""
	barMax := width - sLen(barPrefix) - sLen(barSuffix)
	barCurrent := int(math.Ceil(float64(o.status.Progress) / 100 * float64(barMax)))
	barsDone := strings.Repeat(barDoneChar, barCurrent)
	barsLeft := strings.Repeat(barLeftChar, barMax-barCurrent)

	// TODO: Don't do this one rune at a time
	for n := 0; n < len(text); n++ {
		if n < len(barsDone) {
			barsDone = replaceAtIndex(barsDone, rune(text[n]), n)
		} else {
			barsLeft = replaceAtIndex(barsLeft, rune(text[n]), n-len(barsDone))
		}
	}

	bars := aurora.BgBrightGreen(barsDone).Black().String() + aurora.BgGray(4, barsLeft).String()

	if o.status.Finished && o.status.ExitCode > 0 {
		bars = aurora.BgBrightRed(barsDone).Black().String() + aurora.BgGray(4, barsLeft).BrightRed().String()
	}

	return fmt.Sprintf("%s%s%s", barPrefix, bars, barSuffix)
}

func update() {
	mu.Lock()
	defer mu.Unlock()

	writer.Clear()
	for n, o := range operations {
		index := getIndex(o)
		title := getTitle(o)
		message := getMessage(o)
		sizeProgress := getByteProgress(o)
		fileProgress := getFileProgress(o)
		progress := getProgress(o)
		speed := getSpeed(o)
		timeInfo := getTimeInfo(o)

		progressPrefix := fmt.Sprintf(" %s. %s ", index, title)
		progressSuffix := fmt.Sprintf(" %s | %s | %s | %s | %s ", progress, speed, sizeProgress, fileProgress, timeInfo)
		progressWidth := width - sLen(progressPrefix) - sLen(progressSuffix)
		progressBar := getProgressBar(o, progressWidth, message)

		fmt.Fprintf(writer, "%s%s%s\n", progressPrefix, progressBar, progressSuffix)

		if n < len(operations)-1 {
			fmt.Fprintf(writer, "\n")
		}
	}
	writer.Print()
}
