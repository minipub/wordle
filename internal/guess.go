package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func Guess() (rs [5]byte, err error) {
	r := bufio.NewReader(os.Stdin)

	var b [5]byte
	n, err := r.Read(b[:])
	if err != nil {
		return
	}

	if len(b[0:n]) < 5 {
		err = errors.New("no enough letters")
		return
	}

	// fmt.Printf("s: %s", b[0:n])

	for k, v := range b[0:n] {
		if v >= 'A' && v <= 'Z' {
			rs[k] = v + 32
		} else if v >= 'a' && v <= 'z' {
			rs[k] = v
		} else {
			err = errors.New("no enough letters")
			return
		}
	}

	s := string(rs[:])
	if !isWord(s) {
		err = fmt.Errorf("<%s> not a word", s)
		return
	}

	return
}

// Equal Implement finding inputted word one-by-one char through in-plan word
// x: in-plan word
// y: inputted word
func Equal(x, y [5]byte) bool {
	return x == y
}
