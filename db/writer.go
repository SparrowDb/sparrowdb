package db

import (
	"io"

	"github.com/sparrowdb/util"
)

type dbWriter struct {
	writer io.Writer
}

func (w *dbWriter) Append(key string, value []byte) error {
	bout := util.NewByteStream()
	bout.PutUInt32(uint32(len(value)))
	b := bout.Bytes()

	if _, err := w.writer.Write(b); err != nil {
		return err
	}

	if _, err := w.writer.Write(value); err != nil {
		return err
	}

	return nil
}

func (w *dbWriter) Close() error {
	return w.writer.(io.WriteCloser).Close()
}

func newWriter(f io.Writer) *dbWriter {
	return &dbWriter{f}
}

type indexWriter struct {
	writer io.Writer
}

func (iw *indexWriter) Append(value []byte) error {
	bout := util.NewByteStream()
	bout.PutUInt32(uint32(len(value)))
	b := bout.Bytes()

	if _, err := iw.writer.Write(b); err != nil {
		return err
	}

	if _, err := iw.writer.Write(value); err != nil {
		return err
	}

	return nil
}

func (iw *indexWriter) Close() error {
	return iw.writer.(io.WriteCloser).Close()
}

func newIndexWriter(f io.Writer) *indexWriter {
	return &indexWriter{f}
}
