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
	help bool
	host string
	port int
	// step    string
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
	Cmd.Flags().BoolVarP(&help, "help", "", false, "help for wordle solver")
	Cmd.Flags().StringVarP(&host, "host", "h", "127.0.0.1", "dial host")
	Cmd.Flags().IntVarP(&port, "port", "p", 8080, "dial port")
	// Cmd.Flags().StringVarP(&step, "step", "s", "", "enable single step")
	Cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose or debug log")
}

func main() {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Err:", err)
		os.Exit(1)
	}

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	rr := internal.NewStdReader(r, func() { r.Reset(conn) })

	m := internal.NewMessage()
	p := internal.NewSolverPrinter(verbose)

	var iWord [5]byte // input word

	for {
		rs := internal.ReadLoop(rr, p)

		// print server response
		fmt.Printf("%s", rs)

		if internal.IsTheEnd(rs) {
			break
		}

		if internal.IsRetError(rs) {
			fmt.Fprintln(os.Stderr, "\nPlease upgrade server version to the latest.")
			os.Exit(3)
		}

		pos := internal.CalcPosition(rs)

		if bytes.HasSuffix(rs, []byte(internal.Prompt)) {
			time.Sleep(time.Second)
			iWord = internal.DoPipeCmds(m, pos, iWord, verbose, p)
			// print client request
			fmt.Printf("%s\n", iWord)
			w.Write(iWord[:])
			w.Flush()
		}

	}

}
