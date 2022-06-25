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

func (s *StdReader) Read(b []byte) (n int, err error) {
	n, err = s.Reader.Read(b[:])
	defer s.f()
	return
}
