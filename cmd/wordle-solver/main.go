// Wordle! Solver
package solver

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
		os.Exit(1)
	}

	r := bufio.NewReaderSize(conn, 1)
	w := bufio.NewWriterSize(conn, 1)

	var iWord [5]byte // only words

	for {
		rs := readLoop(r, func() {
			r.Reset(conn)
		})

		// print server response
		fmt.Printf("%s", rs)

		if internal.IsTheEnd(rs) {
			break
		}

		var pos [5]int
		if !internal.IsTheStart(rs) {
			fmt.Fprintf(os.Stderr, `Thinking...
`)
			var vs []byte
			for i, j, k := 0, 0, 0; ; {
				v := rs[i]

				if j > 0 && j%internal.ColoredByteNum == 0 {
					// fmt.Printf("vs: %+v\n", vs)

					if bytes.HasPrefix(vs, []byte(internal.ColorRed)) {
						pos[k] = internal.Miss
					} else if bytes.HasPrefix(vs, []byte(internal.ColorYellow)) {
						pos[k] = internal.Appear
					} else if bytes.HasPrefix(vs, []byte(internal.ColorGreen)) {
						pos[k] = internal.Hit
					}
					k++

					if k == 5 {
						break
					} else {
						vs = make([]byte, 0)
						i += (internal.ColorResetByteNum + 1)
						j = 0
					}
				} else {
					vs = append(vs, v)
					i++
					j++
				}
			}
		}

		if bytes.HasSuffix(rs, []byte(internal.Prompt)) {
			time.Sleep(time.Second)
			iWord = internal.SolveWord(pos, iWord)
			// print client request
			fmt.Printf("%s\n", iWord)
			w.Write(iWord[:])
		}

	}

}

// read next input or the end
func readLoop(r io.Reader, f func()) (rs []byte) {
	var keepRead bool
	var n int
	var err error
	var b [512]byte

	defer f()

	for {
		if keepRead {
			// fmt.Fprintf(os.Stderr, "keepRead: %t\n", keepRead)
			n, err = r.Read(b[n:])
		} else {
			n, err = r.Read(b[:])
		}
		if err != nil {
			fmt.Printf("readLoop err: %+v\n", err)
			os.Exit(2)
		}

		rs = b[0:n]
		fmt.Fprintf(os.Stderr, `resp: {{ %+v }}, {{ %s }}

`, rs, rs)

		if internal.IsTheEnd(rs) {
			break
		} else if !bytes.HasSuffix(rs, []byte(internal.Prompt)) {
			// continue to read if Prompt not direct after PreText or Colored Response
			fmt.Fprintln(os.Stderr, "Prompt not afterwards!")
			keepRead = true
			continue
		} else {
			break
		}
	}

	return
}
