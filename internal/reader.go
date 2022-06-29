package internal

import (
	"io"
)

type Reader interface {
	Read([]byte) (int, error)
}

type StdReader struct {
	io.Reader
	f func()
}

func NewStdReader(r io.Reader, f func()) *StdReader {
	return &StdReader{r, f}
}

func (s *StdReader) Read(b []byte) (n int, err error) {
	n, err = s.Reader.Read(b[:])
	defer s.f()
	return
}
