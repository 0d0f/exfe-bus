package streaming

import (
	"bytes"
	"fmt"
	"github.com/stretchrcom/testify/assert"
	"io"
	"net"
	"testing"
	"time"
)

type FakeWriter struct {
	isClosed    bool
	buf         *bytes.Buffer
	readTimeout time.Time
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

func (w *FakeWriter) Flush() error {
	return nil
}

func (w *FakeWriter) Read(b []byte) (n int, err error) {
	if w.isClosed {
		return -1, io.EOF
	}
	time.Sleep(w.readTimeout.Sub(time.Now()))
	return 0, nil
}

func (w *FakeWriter) SetReadDeadline(t time.Time) error {
	w.readTimeout = t
	return nil
}

func (w *FakeWriter) SetWriteDeadline(t time.Time) error {
	return nil
}

func (w *FakeWriter) SetDeadline(t time.Time) error {
	return nil
}

func (w *FakeWriter) LocalAddr() net.Addr {
	return new(net.IPAddr)
}

func (w *FakeWriter) RemoteAddr() net.Addr {
	return new(net.IPAddr)
}

func TestStreaming(t *testing.T) {
	streaming := New(time.Second / 2)
	id := "user123"

	buf1 := NewFakeWriter()
	go func() {
		streaming.Connect(id, buf1, buf1)
	}()

	time.Sleep(time.Second)

	err := streaming.Feed(id, "abcde")
	assert.Equal(t, err, nil)

	buf2 := NewFakeWriter()
	go func() {
		streaming.Connect(id, buf2, buf2)
	}()

	time.Sleep(time.Second)

	err = streaming.Feed(id, "123")
	assert.Equal(t, err, nil)
	err = streaming.Feed("user789", "123")
	assert.NotEqual(t, err, nil)
	err = streaming.Feed(id, "xyz")
	assert.Equal(t, err, nil)

	time.Sleep(time.Second)

	buf1.Close()
	buf2.Close()

	time.Sleep(time.Second)

	err = streaming.Feed(id, "")
	assert.NotEqual(t, err, nil)

	assert.Equal(t, buf1.buf.String(), "abcde\n123\nxyz\n")
	assert.Equal(t, buf2.buf.String(), "123\nxyz\n")
}
