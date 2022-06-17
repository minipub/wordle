package internal

import (
	"fmt"
	"net"
)

type Writer interface {
	Write(string)
}

type StdWriter struct {
}

func (StdWriter) Write(s string) {
	fmt.Print(s)
}

type BotWriter struct {
	conn net.Conn
}

func NewBotWriter(conn net.Conn) BotWriter {
	return BotWriter{
		conn,
	}
}

func (b *BotWriter) Write(s string) {
	b.conn.Write([]byte(s))
}
