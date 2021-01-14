package utils

import (
	"strings"

	"github.com/c2h5oh/datasize"
)

func ByteCountSI(b int64) string {
	bs := datasize.ByteSize(b)

	return strings.ReplaceAll(bs.HumanReadable(), " ", "")
}
