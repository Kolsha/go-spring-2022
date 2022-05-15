//go:build !solution

package gzep

import (
	"compress/gzip"
	"io"
	"sync"
)

var pool sync.Pool

func Encode(data []byte, dst io.Writer) error {
	var ww *gzip.Writer
	if w := pool.Get(); w != nil {
		ww = w.(*gzip.Writer)
		ww.Reset(dst)
	} else {
		ww, _ = gzip.NewWriterLevel(dst, gzip.DefaultCompression)
	}
	defer func() {
		_ = ww.Close()
		pool.Put(ww)
	}()
	if _, err := ww.Write(data); err != nil {
		return err
	}
	return ww.Flush()
}
