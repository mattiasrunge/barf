package utils

import "bytes"

// https://golang.org/src/bufio/scan.go

// ScanLines tweaks the default ScanLines to split on \r also
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	for {
		i := bytes.IndexByte(data, '\r')

		if i == -1 {
			break
		}

		data[i] = '\n'
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
