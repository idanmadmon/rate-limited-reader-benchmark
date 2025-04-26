package main

import (
	"io"

	idanmadmonReader "github.com/idanmadmon/rate-limited-reader"
)

type ReaderFactory func(io.ReadCloser, int64) io.ReadCloser

func NoLimitReaderFactory(reader io.ReadCloser, limit int64) io.ReadCloser {
	return reader
}

func IdanMadmonRateLimitReaderFactory(reader io.ReadCloser, limit int64) io.ReadCloser {
	return idanmadmonReader.NewRateLimitedReadCloser(reader, limit)
}
