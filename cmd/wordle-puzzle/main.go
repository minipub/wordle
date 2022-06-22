// Wordle! Puzzle
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/minipub/wordle/internal"
)

func main() {
	// TODO
	var isBot bool
	var port = ":8080"

	var rw internal.ReadWriter

	if !isBot {

		ln, err := net.Listen("tcp", port)
		if err != nil {
			os.Exit(1)
		}
		fmt.Println("Listening on :8080")

		for {
			conn, err := ln.Accept()
			if err != nil {
				os.Exit(2)
			}

			go func() {
				defer conn.Close()

				r := bufio.NewReader(conn)
				w := bufio.NewWriter(conn)
				rw = internal.NewBotReadWriter(r, w, func() { r.Reset(conn) })
				internal.DoPuzzle(rw, func() { w.Flush() })
			}()
		}
	} else {
		r := bufio.NewReaderSize(os.Stdin, 5)
		rw = internal.NewStdReadWriter(r, func() { r.Reset(os.Stdin) })
		internal.DoPuzzle(rw, func() {})
	}

}
