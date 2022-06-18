package internal

import (
	"io"
)

type ReadWriter interface {
	Reader
	Writer
}

type StdReadWriter struct {
	*StdReader
	*StdWriter
}

func NewStdReadWriter(r io.Reader, f func()) StdReadWriter {
	return StdReadWriter{
		&StdReader{r, f},
		&StdWriter{},
	}
}

type BotReadWriter struct {
	*StdReader
	*BotWriter
}

func NewBotReadWriter(r io.Reader, w io.Writer, f func()) BotReadWriter {
	return BotReadWriter{
		&StdReader{r, f},
		&BotWriter{w},
	}
}
