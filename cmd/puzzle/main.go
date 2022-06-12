package main

import (
	"fmt"
	"os"
	"wordle/internal"
)

func main() {
	pWord := internal.RandOneWord() // in-plan word

	fmt.Println("A Wordle Game!\nPlease input a five-letter word and Press <Enter> to confirm.\n")

	for i := 0; i < len(internal.CheerWords); {
		fmt.Print("input: ")
		iWord, err := internal.InputWord() // inputted word
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
			continue
		}

		ok := internal.Equal(pWord, iWord)
		if ok {
			fmt.Println()
			fmt.Println(internal.CheerWords[i])
			os.Exit(0)
		}

		pos := internal.FindPos(pWord, iWord)
		for m, n := range pos {
			fmt.Print(internal.Colors[n], string(iWord[m]), internal.ColorReset)
		}
		fmt.Println()

		i++
	}

	fmt.Println("\nOut of Chance!")
	fmt.Printf("\nThe Word is <%s>\n", pWord)
	fmt.Println("\nTake a break or get another round.")
}
