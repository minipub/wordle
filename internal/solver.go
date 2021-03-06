package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func IsTheStart(b []byte) bool {
	return bytes.HasPrefix(b, []byte(PreText))
}

func IsRetError(b []byte) bool {
	return bytes.HasPrefix(b, []byte(RetERR))
}

func IsTheEnd(b []byte) bool {
	return bytes.HasSuffix(b, []byte(ByeText))
}

// read next input or the end
func ReadLoop(r io.Reader, p Writer) (rs []byte) {
	var n int
	var err error
	var b [512]byte

	var rbuf []byte
	bbuf := bytes.NewBuffer(rbuf)

	for {
		n, err = r.Read(b[:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Err:", err)
			os.Exit(2)
		}

		seg := b[0:n]
		bbuf.Write(seg)
		rs = bbuf.Bytes()

		p.Write(fmt.Sprintf(`resp: {{ %+v }}, {{ %s }}

`, rs, rs))

		if IsTheEnd(rs) {
			break
		} else if !bytes.HasSuffix(rs, []byte(Prompt)) {
			// continue to read if Prompt not direct after PreText or Colored Response
			p.Write(fmt.Sprintln("Prompt not afterwards!"))
			continue
		} else {
			break
		}
	}

	return
}

// calculate the position according to the last response
func CalcPosition(b []byte) (pos [5]int) {
	if IsTheStart(b) {
		return
	}

	var vs []byte
	for i, j, k := 0, 0, 0; ; {
		v := b[i]

		if j > 0 && j%ColoredByteNum == 0 {
			// fmt.Printf("vs: %+v\n", vs)

			if bytes.HasPrefix(vs, []byte(ColorRed)) {
				pos[k] = Miss
			} else if bytes.HasPrefix(vs, []byte(ColorYellow)) {
				pos[k] = Appear
			} else if bytes.HasPrefix(vs, []byte(ColorGreen)) {
				pos[k] = Hit
			}
			k++

			if k == 5 {
				break
			} else {
				vs = make([]byte, 0)
				i += (ColorResetByteNum + 1)
				j = 0
			}
		} else {
			vs = append(vs, v)
			i++
			j++
		}
	}

	return
}
