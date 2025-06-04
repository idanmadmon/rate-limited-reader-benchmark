package main

import (
	"context"
	"io"

	idanmadmonReader "github.com/idanmadmon/rate-limited-reader"
	"github.com/juju/ratelimit"
	uberratelimit "go.uber.org/ratelimit"
	"golang.org/x/time/rate"
)

type ReaderFactory func(reader io.ReadCloser, bufferSize, limit int) io.ReadCloser

func NoLimitReaderFactory(reader io.ReadCloser, _ int) io.ReadCloser {
	return reader
}

func IdanMadmonDeterministicRateLimitReaderFactory(reader io.ReadCloser, bufferSize, limit int) io.ReadCloser {
	return idanmadmonReader.NewRateLimitedReadCloser(reader, int64(limit))
}

func GolangBurstsRateLimitReaderFactory(reader io.ReadCloser, bufferSize, limit int) io.ReadCloser {
	// limiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(limit/bufferSize)), 1)
	// limiter := rate.NewLimiter(rate.Limit(limit/bufferSize), 1)
	limiter := rate.NewLimiter(rate.Limit(limit/bufferSize), limit/bufferSize)
	return &GolangRateLimitedReader{
		reader:  reader,
		limiter: limiter,
		ctx:     context.Background(),
	}
}

type GolangRateLimitedReader struct {
	reader  io.ReadCloser
	limiter *rate.Limiter
	ctx     context.Context
}

func (r *GolangRateLimitedReader) Read(p []byte) (n int, err error) {
	// fmt.Println("Started", time.Now())
	// err = r.limiter.Wait(r.ctx) // wait until tokens are available
	err = r.limiter.WaitN(r.ctx, 700) // wait until tokens are available
	// fmt.Println("Waited", time.Now())
	if err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}

func (r *GolangRateLimitedReader) Close() error {
	return r.reader.Close()
}

func JujuBurstsRateLimitReaderFactory(reader io.ReadCloser, bufferSize, limit int) io.ReadCloser {
	bucket := ratelimit.NewBucketWithRate(float64(limit), int64(limit))
	return &JujuRateLimitedReader{
		reader: reader,
		bucket: bucket,
	}
}

type JujuRateLimitedReader struct {
	reader io.ReadCloser
	bucket *ratelimit.Bucket
}

func (r *JujuRateLimitedReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	r.bucket.Wait(int64(n))
	return n, err
}

func (r *JujuRateLimitedReader) Close() error {
	return r.reader.Close()
}

func UberDeterministicRateLimitReaderFactory(reader io.ReadCloser, bufferSize, limit int) io.ReadCloser {
	rl := uberratelimit.New(limit / bufferSize) // operations per second
	return &UberRateLimitedReader{
		reader:  reader,
		limiter: rl,
	}
}

type UberRateLimitedReader struct {
	reader  io.ReadCloser
	limiter uberratelimit.Limiter
}

func (r *UberRateLimitedReader) Read(p []byte) (n int, err error) {
	r.limiter.Take()
	return r.reader.Read(p)
}

func (r *UberRateLimitedReader) Close() error {
	return r.reader.Close()
}
