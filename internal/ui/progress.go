package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/logrusorgru/aurora"
	"github.com/mattiasrunge/goterminal"
)

var writer *goterminal.Writer
var width = 80
var mu sync.Mutex

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
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

func getStep(o *operationWithStatus) string {
	return o.status.Step
}

func getFileName(o *operationWithStatus) string {
	return o.status.FileName
}

func getFinished(o *operationWithStatus) string {
	if o.status.Finished {
		if o.status.ExitCode > 0 {
			return fmt.Sprintf("Failed [code=%d]: %s", o.status.ExitCode, o.status.Error)
		}

		return "Completed successfully!"
	}

	return ""
}

func getByteProgress(o *operationWithStatus) string {
	totalDiffStr := byteCountSI(o.status.BytesDiffTotal)
	totalStr := byteCountSI(o.status.BytesTotal)
	doneStr := aurora.BrightGreen(byteCountSI(o.status.BytesDone)).String()

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
	return fmt.Sprintf("%s/s", byteCountSI(int64(o.status.Speed)))
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
		bars = aurora.BgBrightRed(barsDone).Black().String() + aurora.BgGray(4, barsLeft).String()
	}

	return fmt.Sprintf("%s%s%s", barPrefix, bars, barSuffix)
}

func update() {
	mu.Lock()
	writer.Clear()
	for n, o := range operations {
		index := getIndex(o)
		title := getTitle(o)
		step := getStep(o)
		fileName := getFileName(o)
		finished := getFinished(o)
		sizeProgress := getByteProgress(o)
		fileProgress := getFileProgress(o)
		progress := getProgress(o)
		speed := getSpeed(o)
		timeInfo := getTimeInfo(o)

		barText := ""

		if len(step) > 0 {
			barText = fmt.Sprintf(" %s...", step)
		}

		if len(fileName) > 0 {
			barText = fmt.Sprintf(" %s %s...", step, fileName)
		}

		if len(finished) > 0 {
			barText = fmt.Sprintf(" %s", finished)
		}

		progressPrefix := fmt.Sprintf(" %s. %s ", index, title)
		progressSuffix := fmt.Sprintf(" %s | %s | %s | %s | %s ", progress, speed, sizeProgress, fileProgress, timeInfo)
		progressWidth := width - sLen(progressPrefix) - sLen(progressSuffix)
		progressBar := getProgressBar(o, progressWidth, barText)

		fmt.Fprintf(writer, "%s%s%s\n", progressPrefix, progressBar, progressSuffix)

		if n < len(operations)-1 {
			fmt.Fprintf(writer, "\n")
		}
	}

	writer.Print()
	mu.Unlock()
}
