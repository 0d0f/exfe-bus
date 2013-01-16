package streaming

import (
	"bytes"
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"testing"
	"time"
)

type FakeWriter struct {
	isClosed bool
	buf      *bytes.Buffer
}

func NewFakeWriter() *FakeWriter {
	return &FakeWriter{
		isClosed: false,
		buf:      bytes.NewBuffer(nil),
	}
}

func (w *FakeWriter) Close() error {
	w.isClosed = true
	return nil
}

func (w *FakeWriter) Write(p []byte) (int, error) {
	if w.isClosed {
		return -1, fmt.Errorf("closed")
	}
	return w.buf.Write(p)
}

func TestStreaming(t *testing.T) {
	streaming := New()
	id := "user123"

	buf := NewFakeWriter()

	go func() {
		streaming.Connect(id, buf)
	}()

	time.Sleep(time.Second)

	err := streaming.Feed(id, "abcde")
	assert.Equal(t, err, nil)
	err = streaming.Feed(id, "123")
	assert.Equal(t, err, nil)
	err = streaming.Feed("user789", "123")
	assert.NotEqual(t, err, nil)
	err = streaming.Feed(id, "xyz")
	assert.Equal(t, err, nil)

	time.Sleep(time.Second)

	buf.Close()
	streaming.Feed(id, "")
	err = streaming.Feed(id, "")
	assert.NotEqual(t, err, nil)

	assert.Equal(t, buf.buf.String(), "abcde\n123\nxyz\n")
}
