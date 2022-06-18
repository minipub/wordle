package internal

import (
	"fmt"
	"io"
)

type Writer interface {
	Write(string)
}

type StdWriter struct {
}

func (*StdWriter) Write(s string) {
	fmt.Print(s)
}

type BotWriter struct {
	w io.Writer
}

func NewBotWriter(w io.Writer) BotWriter {
	return BotWriter{
		w,
	}
}

func (b *BotWriter) Write(s string) {
	b.w.Write([]byte(s))
}
