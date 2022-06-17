package internal

import (
	"bufio"
	"io"
)

type Reader interface {
	Read([]byte) (int, error)
}

func read(conn io.Reader, b []byte) (n int, err error) {
	r := bufio.NewReader(conn)
	n, err = r.Read(b[:])
	return
}

type StdReader struct {
	io.Reader
}

func (s *StdReader) Read(b []byte) (n int, err error) {
	return read(s.Reader, b)
}

type BotReader struct {
	io.Reader
}

func (s *BotReader) Read(b []byte) (n int, err error) {
	return read(s.Reader, b)
}
