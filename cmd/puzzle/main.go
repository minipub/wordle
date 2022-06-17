// Wordle! Puzzle
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"wordle/internal"
)

func main() {
	var isBot bool

	var rw internal.ReadWriter
	if isBot {
		rw = &internal.BotReadWriter{}
	} else {
		rw = &internal.StdReadWriter{}
	}

	// handle different reader
	ln, err := net.Listen("tcp", ":8080")
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
			// defer conn.Close()

			for {
				fmt.Println("heheh")
				iWord, err := internal.HandleInput(rw)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
					continue
				}

				fmt.Println("eeeey")
				// conn.Write(iWord[:])
				// conn.Write([]byte("gogo"))
				w := bufio.NewWriterSize(conn, 5)
				w.Write([]byte(iWord[:]))
				w.Write([]byte("gogo"))
			}
		}()
	}

	return

	pWord := internal.RandOneWord() // in-plan word

	fmt.Print(`A Wordle Game!

Please input a five-letter word and Press <Enter> to confirm.

`)

	for i := 0; i < len(internal.CheerWords); {
		// handle different writer
		fmt.Print("input: ")
		iWord, err := internal.InputWord() // inputted word
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
			continue
		}

		ok := internal.Equal(pWord, iWord)
		if ok {
			fmt.Printf("\n%s\n", internal.CheerWords[i])
			os.Exit(0)
		}

		pos := internal.FindPos(pWord, iWord)
		for m, n := range pos {
			fmt.Print(internal.Colors[n], string(iWord[m]), internal.ColorReset)
		}
		fmt.Println()

		i++
	}

	fmt.Printf(`
Out of Chance!

The Word is <%s>

Take a break or get another round.
`, pWord)
}
