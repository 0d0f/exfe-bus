package model

import (
	"database/sql"
	"io"
)

type ErrWriter struct {
	err error
	w   io.Writer
}

func NewErrWriter(w io.Writer) *ErrWriter {
	return &ErrWriter{
		w: w,
	}
}

func (w *ErrWriter) Write(p []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write(p)
}

func (w *ErrWriter) WriteString(s string) {
	w.Write([]byte(s))
}

func (w ErrWriter) Err() error {
	return w.err
}

type ErrPrepare struct {
	err error
	db  *sql.DB
}

func NewErrPrepare(db *sql.DB) *ErrPrepare {
	return &ErrPrepare{
		db: db,
	}
}

func (p *ErrPrepare) Prepare(query string) *sql.Stmt {
	if p.err != nil {
		return nil
	}
	ret, err := p.db.Prepare(query)
	if err != nil {
		p.err = err
		return nil
	}
	return ret
}

func (p *ErrPrepare) Err() error {
	return p.err
}
