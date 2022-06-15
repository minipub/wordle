// Wordle! Puzzle
package main

import (
	"fmt"
	"os"
	"wordle/internal"
)

func main() {
	pWord := internal.RandOneWord() // in-plan word

	fmt.Print(`A Wordle Game!

Please input a five-letter word and Press <Enter> to confirm.

`)

	// handle different reader
	// ln, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	// handle error
	// }
	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		// handle error
	// 	}
	// 	go handleConnection(conn)
	// }

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
