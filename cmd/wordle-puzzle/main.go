// Wordle! Puzzle
package puzzle

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/minipub/wordle/internal"
	"github.com/spf13/cobra"
)

const (
	ModeShell = 1 + iota
	ModeCS
)

var (
	mode int
	port int

	Cmd = &cobra.Command{
		Use:   "puzzle",
		Short: "Wordle Puzzle",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
)

func init() {
	Cmd.Flags().IntVarP(&mode, "mode", "m", 1, `puzzle run mode: 
1. Interactive shell
2. C-S
`)
	Cmd.Flags().IntVarP(&port, "port", "p", 8080, "listen port")
}

func main() {
	var rw internal.ReadWriter

	if mode == ModeCS {
		addr := fmt.Sprintf(":%d", port)

		ln, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Err:", err)
			os.Exit(1)
		}
		fmt.Println("Listening on", addr)

		for {
			conn, err := ln.Accept()
			if err != nil {
				os.Exit(2)
			}

			go func() {
				defer conn.Close()

				r := bufio.NewReader(conn)
				w := bufio.NewWriter(conn)
				rw = internal.NewBotReadWriter(r, w,
					func() {
						r.Reset(conn)
					},
					func() {
						w.Flush()
					})
				internal.DoPuzzle(rw)
			}()
		}
	} else {
		r := bufio.NewReaderSize(os.Stdin, 5)
		w := bufio.NewWriter(os.Stdout)
		rw = internal.NewBotReadWriter(r, w,
			func() {
				r.Reset(os.Stdin)
			},
			func() {
				w.Flush()
			})
		internal.DoPuzzle(rw)
	}

}
