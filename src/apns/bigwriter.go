package apns

import (
	"encoding/binary"
	"io"
)

type BigWriter struct {
	Error  error
	writer io.Writer
}

func newBigWriter(writer io.Writer) *BigWriter {
	return &BigWriter{
		writer: writer,
	}
}

func (w *BigWriter) Write(v interface{}) {
	if w.Error != nil {
		return
	}
	w.Error = binary.Write(w.writer, binary.BigEndian, v)
}
