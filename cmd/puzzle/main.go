package main

import (
	"fmt"
	"os"
	"wordle/internal"
)

func main() {
	pWord := internal.RandOneWord() // in-plan word

	for i := 0; i < len(internal.CheerWords); {
		iWord, err := internal.Guess() // inputted word
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
			continue
		}
		// fmt.Printf("iWord: %s\n", iWord)

		i++

		ok := internal.Equal(pWord, iWord)
		if ok {
			fmt.Println(internal.CheerWords[i])
			os.Exit(0)
		}
		// fmt.Printf("iWord == pWord: %t\n", ok)

		pos := internal.FindPos(pWord, iWord)
		for _, j := range pos {
			fmt.Print(string(internal.Colors[j]), "◼", internal.ColorReset)
		}
		fmt.Println()
	}

	fmt.Printf("\nThe Plan Word is <%s>\n", pWord)
	fmt.Println("Come on! Take a break & get another round.")
}
