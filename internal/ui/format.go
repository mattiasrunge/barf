package ui

import (
	"strings"

	"github.com/c2h5oh/datasize"
)

func byteCountSI(b int64) string {
	bs := datasize.ByteSize(b)

	return strings.ReplaceAll(bs.HumanReadable(), " ", "")
}
