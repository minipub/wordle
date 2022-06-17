package internal

type ReadWriter interface {
	Reader
	Writer
}

type StdReadWriter struct {
	StdReader
	StdWriter
}

type BotReadWriter struct {
	BotReader
	BotWriter
}
