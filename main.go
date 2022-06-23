package main

import (
	puzzle "github.com/minipub/wordle/cmd/wordle-puzzle"
	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{Use: "wordle"}

	root.AddCommand(
		puzzle.Cmd,
	)

	root.Execute()
}
