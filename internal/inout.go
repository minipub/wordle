package internal

import (
	"io"
	"os"
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

type SolverPrinter struct {
	Writer
}

func NewSolverPrinter(verbose bool) *SolverPrinter {
	var w Writer
	if verbose {
		// redirect to Stderr (it looks strange to print to Stdout)
		w = &BotWriter{os.Stderr}
	} else {
		w = &NoopWriter{}
	}
	return &SolverPrinter{
		w,
	}
}
