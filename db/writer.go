package db

import (
	"io"

	"github.com/SparrowDb/sparrowdb/util"
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

type bufWriter struct {
	writer io.Writer
}

func (bw *bufWriter) Append(value []byte) error {
	bout := util.NewByteStream()
	bout.PutUInt32(uint32(len(value)))
	b := bout.Bytes()

	if _, err := bw.writer.Write(b); err != nil {
		return err
	}

	if _, err := bw.writer.Write(value); err != nil {
		return err
	}

	return nil
}

func (bw *bufWriter) Close() error {
	return bw.writer.(io.WriteCloser).Close()
}

func newBufWriter(f io.Writer) *bufWriter {
	return &bufWriter{f}
}
