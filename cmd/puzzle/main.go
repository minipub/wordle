package main

import (
	"fmt"
	"os"
	"wordle/internal"
)

func main() {
	pWord := internal.RandOneWord() // in-plan word
	// fmt.Printf("word: %s\n", pWord)

	for i := 0; i < len(internal.CheerWords); i++ {
		iWord, err := internal.Guess() // inputted word
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
			os.Exit(1)
		}
		// fmt.Printf("iWord: %s\n", iWord)

		ok := internal.Equal(pWord, iWord)
		if ok {
			fmt.Println(internal.CheerWords[i])
			os.Exit(0)
		}
		// fmt.Printf("iWord == pWord: %t\n", ok)

		pos := internal.FindPos(pWord, iWord)
		for _, j := range pos {
			// fmt.Printf("j: %d\n", j)
			fmt.Print(string(internal.Colors[j]), "◼")
		}
		fmt.Println()
	}

	fmt.Println("Come on! Take a break & get another round.")
}
