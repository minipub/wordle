// Wordle! Solver
package solver

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/minipub/wordle/internal"
	"github.com/spf13/cobra"
)

var (
	host    string
	port    int
	verbose bool

	Cmd = &cobra.Command{
		Use:   "solver",
		Short: "Wordle Solver(Cheater)",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
)

func init() {
	Cmd.Flags().StringVar(&host, "host", "127.0.0.1", "dial host")
	Cmd.Flags().IntVar(&port, "port", 8080, "dial port")
	Cmd.Flags().BoolVar(&verbose, "v", false, "enable verbose or debug log")
}

func main() {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Err:", err)
		os.Exit(1)
	}

	r := bufio.NewReaderSize(conn, 1)
	w := bufio.NewWriterSize(conn, 1)

	m := internal.NewMessage()
	p := internal.NewSolverPrinter(verbose)

	var iWord [5]byte // input word

	for {
		rs := internal.ReadLoop(r, func() {
			r.Reset(conn)
		}, p)

		// print server response
		fmt.Printf("%s", rs)

		if internal.IsTheEnd(rs) {
			break
		}

		pos := internal.CalcPosition(rs)

		if bytes.HasSuffix(rs, []byte(internal.Prompt)) {
			time.Sleep(time.Second)
			iWord = internal.DoPipeCmds(m, pos, iWord, verbose, p)
			// print client request
			fmt.Printf("%s\n", iWord)
			w.Write(iWord[:])
		}

	}

}
