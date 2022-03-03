//go:build !solution
// +build !solution

package otp

import (
	"io"
)

type reader struct {
	r    io.Reader
	prng io.Reader
}

func (r reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	rng := make([]byte, n)
	r.prng.Read(rng)
	for i := 0; i < n; i++ {
		p[i] = p[i] ^ rng[i]
	}
	return
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	res := &reader{r: r, prng: prng}
	return res
}

type writer struct {
	w    io.Writer
	prng io.Reader
}

func (w writer) Write(p []byte) (n int, err error) {
	rng := make([]byte, len(p))
	w.prng.Read(rng)

	for i := 0; i < len(p); i++ {
		rng[i] = rng[i] ^ p[i]
	}
	return w.w.Write(rng)
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	res := &writer{w: w, prng: prng}
	return res
}
