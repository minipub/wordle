package internal

import (
	"io"
	"os"
)

type ReadWriter interface {
	Reader
	Writer
}

type BotReadWriter struct {
	*StdReader
	*BotWriter
}

func NewBotReadWriter(r io.Reader, w io.Writer, rf, wf func()) BotReadWriter {
	return BotReadWriter{
		&StdReader{r, rf},
		&BotWriter{w, wf},
	}
}

type SolverPrinter struct {
	Writer
}

func NewSolverPrinter(verbose bool) *SolverPrinter {
	var w Writer
	if verbose {
		// redirect to Stderr (it looks strange to print to Stdout)
		w = &BotWriter{os.Stderr, func() {}}
	} else {
		w = &NoopWriter{}
	}
	return &SolverPrinter{
		w,
	}
}
