// Wordle! Solver
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"wordle/internal"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		os.Exit(1)
	}

	r := bufio.NewReader(conn)
	w := bufio.NewWriterSize(conn, 1)

	var iWord [5]byte // only words

	for {
		var b [512]byte // colored text words
		n, err := r.Read(b[:])
		if err != nil {
			os.Exit(2)
		}

		// fmt.Printf("bbbb: %+v, %s\n", b[0:n], b[0:n])

		// s := string(b[0:n])
		// print server response
		fmt.Printf("%s", b[0:n])
		// fmt.Printf("resp: {{{ %s }}}", b[0:n])

		// if !bytes.HasPrefix(b[0:n], []byte(internal.PreText)) {

		if bytes.HasPrefix(b[0:n], []byte(internal.ColorRed)) ||
			bytes.HasPrefix(b[0:n], []byte(internal.ColorYellow)) ||
			bytes.HasPrefix(b[0:n], []byte(internal.ColorGreen)) {

			var vs []byte
			var pos [5]int
			for i, j, k := 0, 0, 0; ; {
				v := b[0:n][i]

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

					// if len(pos) == 5 {
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

			// fmt.Printf("pos: %+v\n", pos)

			if bytes.HasSuffix(b[0:n], []byte(internal.Prompt)) {
				// solve word
				iWord = internal.SolveWord(pos, iWord)
				// print client request
				fmt.Printf("%s\n", iWord)
				w.Write(iWord[:])
			}
		} else if bytes.HasSuffix(b[0:n], []byte(internal.Prompt)) {
			time.Sleep(time.Second)
			copy(iWord[:], "great")
			// print client request
			fmt.Printf("%s\n", iWord)
			w.Write(iWord[:])
		}

	}

}
