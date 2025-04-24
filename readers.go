package main

import (
	"io"

	idanmadmonReader "github.com/idanmadmon/rate-limited-reader"
)

type ReaderFactory func(io.Reader, int64) io.Reader

func NoLimitReaderFactory(reader io.Reader, limit int64) io.Reader {
	return reader
}

func IdanMadmonRateLimitReaderFactory(reader io.Reader, limit int64) io.Reader {
	return idanmadmonReader.NewRateLimitedReader(reader, limit)
}
