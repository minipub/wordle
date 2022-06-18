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

// func read(r io.Reader, b []byte) (n int, err error) {
// 	n, err = r.Read(b[:])
// 	return
// }

// type StdReader struct {
// 	io.Reader
// }

// func (s *StdReader) Read(b []byte) (n int, err error) {
// 	return read(s.Reader, b)
// }

// type BotReader struct {
// 	io.Reader
// }

// func (s *BotReader) Read(b []byte) (n int, err error) {
// 	return read(s.Reader, b)
// }
