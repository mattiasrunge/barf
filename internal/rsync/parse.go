package rsync

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/c2h5oh/datasize"
)

func parseByteStr(str string) int64 {
	cleanStr := strings.ReplaceAll(str, ",", "")
	bytes, err := strconv.ParseInt(cleanStr, 10, 64)

	if err != nil {
		return -1
	}

	return bytes
}

func parseSpeedStr(str string) float64 {
	var v datasize.ByteSize
	err := v.UnmarshalText([]byte("260MB"))

	// b, err := units.ParseBase2Bytes(strings.ToUpper(str))

	if err != nil {
		fmt.Println("str:", str)
		fmt.Println("err:", err)
		return -1
	}

	return float64(v.Bytes())
}

/*
rsync --stats --dry-run from/* to/

Number of files: 4 (reg: 4)
Number of created files: 4 (reg: 4)
Number of deleted files: 0
Number of regular files transferred: 4
Total file size: 4,000,000,000 bytes
Total transferred file size: 4,000,000,000 bytes
Literal data: 0 bytes
Matched data: 0 bytes
File list size: 0
File list generation time: 0.001 seconds
File list transfer time: 0.000 seconds
Total bytes sent: 122
Total bytes received: 28

sent 122 bytes  received 28 bytes  300.00 bytes/sec
total size is 4,000,000,000  speedup is 26,666,666.67 (DRY RUN)
*/

var reNumberOfFiles = regexp.MustCompile(`Number of files: \S \Sreg: ([0-9]+)`)
var reNumberOfCreatedFiles = regexp.MustCompile(`Number of created files: \S \Sreg: ([0-9]+)`)
var reTotalFileSize = regexp.MustCompile(`Total file size: (.*) bytes`)
var reTotalTransferredFileSize = regexp.MustCompile(`Total transferred file size: (.*) bytes`)

func parseNumberOfFiles(line string, def int64) int64 {
	matches := reNumberOfFiles.FindStringSubmatch(line)

	if len(matches) < 2 {
		return def
	}

	number, err := strconv.ParseInt(matches[1], 10, 64)

	if err != nil {
		return def
	}

	return number
}

func parseNumberOfCreatedFiles(line string, def int64) int64 {
	matches := reNumberOfCreatedFiles.FindStringSubmatch(line)

	if len(matches) < 2 {
		return def
	}

	number, err := strconv.ParseInt(matches[1], 10, 64)

	if err != nil {
		return def
	}

	return number
}

func parseTotalFileSize(line string, def int64) int64 {
	matches := reTotalFileSize.FindStringSubmatch(line)

	if len(matches) < 2 {
		return def
	}

	bytes := parseByteStr(matches[1])

	if bytes == -1 {
		return def
	}

	return bytes
}

func parseTotalTransferredFileSize(line string, def int64) int64 {
	matches := reTotalTransferredFileSize.FindStringSubmatch(line)

	if len(matches) < 2 {
		return def
	}

	bytes := parseByteStr(matches[1])

	if bytes == -1 {
		return def
	}

	return bytes
}

/*
rsync -aP --inplace --no-whole-file from/* to/

sending incremental file list
abc1.txt
  1,000,000,000 100%  319.26MB/s    0:00:02 (xfr#1, to-chk=3/4)
abc2.txt
  1,000,000,000 100%  237.00MB/s    0:00:04 (xfr#2, to-chk=2/4)
abc3.txt
  1,000,000,000 100%  325.15MB/s    0:00:02 (xfr#3, to-chk=1/4)
abc4.txt
  1,000,000,000 100%  224.39MB/s    0:00:04 (xfr#4, to-chk=0/4)
*/
var reProgress = regexp.MustCompile(`^\s+(\S+)\s+(\S+)%\s+(\S+)/s`)

func parseProgress(line string) (int64, float64, bool) {
	matches := reProgress.FindStringSubmatch(line)

	if len(matches) < 3 {
		return -1, -1, false
	}

	fileBytesDone := parseByteStr(matches[1])
	speed := parseSpeedStr(matches[3])

	return fileBytesDone, speed, true
}

func parseFileName(line string) (string, bool, bool) {
	if strings.HasPrefix(line, "__file:") {
		fileName := strings.TrimPrefix(strings.Replace(line, "__file:", "", 1), " ")
		isDir := strings.HasPrefix(fileName, "./")

		return fileName, isDir, true
	}

	return "", false, false
}
