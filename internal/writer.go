package internal

import (
	"fmt"
	"io"
)

type Writer interface {
	Write(string)
}

type NoopWriter struct {
}

func (*NoopWriter) Write(s string) {
}

type StdWriter struct {
}

func (*StdWriter) Write(s string) {
	fmt.Print(s)
}

type BotWriter struct {
	w io.Writer
	f func()
}

// func NewBotWriter(w io.Writer) BotWriter {
// 	return BotWriter{
// 		w,
// 	}
// }

func (b *BotWriter) Write(s string) {
	b.w.Write([]byte(s))
	b.f()
}
