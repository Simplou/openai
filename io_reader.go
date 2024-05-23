package openai

import (
	"bytes"
	"io"
)

var ioReader = func(b []byte) io.Reader {
	return bytes.NewBuffer(b)
}
