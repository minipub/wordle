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
	for {
		var b [512]byte
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
			var pos []int
			for i, j := 0, 0; ; {
				v := b[0:n][i]

				if j > 0 && j%internal.ColoredByteNum == 0 {
					// fmt.Printf("vs: %+v\n", vs)

					if bytes.HasPrefix(vs, []byte(internal.ColorRed)) {
						pos = append(pos, internal.Miss)
					} else if bytes.HasPrefix(vs, []byte(internal.ColorYellow)) {
						pos = append(pos, internal.Appear)
					} else if bytes.HasPrefix(vs, []byte(internal.ColorGreen)) {
						pos = append(pos, internal.Hit)
					}

					if len(pos) == 5 {
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
		}

		time.Sleep(time.Second)

		if bytes.HasSuffix(b[0:n], []byte(internal.Prompt)) {
			// TODO chosen word
			word := "great"
			// print client request
			fmt.Println(word)
			w.Write([]byte(word))
		}
	}

}
