package main

import (
	puzzle "github.com/minipub/wordle/cmd/wordle-puzzle"
	solver "github.com/minipub/wordle/cmd/wordle-solver"
	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{Use: "wordle"}

	root.AddCommand(
		puzzle.Cmd,
		solver.Cmd,
	)

	root.Execute()
}
