package main

import "io"

type ReaderFactory func(io.Reader, int64) io.Reader

func NoLimitReaderFactory(reader io.Reader, limit int64) io.Reader {
	return reader
}
