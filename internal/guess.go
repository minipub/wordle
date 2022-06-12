package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const (
	// status positon of inputted char at in-plan word
	// hit: green, appear: yellow, miss: gray
	hit = iota
	appear
	miss

	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"
)

var Colors map[int]string

func init() {
	Colors = map[int]string{
		hit:    ColorGreen,
		appear: ColorYellow,
		miss:   ColorRed,
	}
}

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

// Equal Implement check inputted & in-plan word is equal
// x: in-plan word
// y: inputted word
func Equal(x, y [5]byte) bool {
	return x == y
}

// FindPos Implement finding inputted word position one-by-one char through in-plan word
// x: in-plan word
// y: inputted word
func FindPos(x, y [5]byte) (pos [5]int) {
	// map store x's char & positons
	xps := make(map[byte][]int)
	for k, v := range x {
		xps[v] = append(xps[v], k)
	}

	for k, v := range y {
		if _, ok := xps[v]; !ok || len(xps[v]) == 0 {
			pos[k] = miss
		} else {
			if IsIn(xps[v], k) {
				pos[k] = hit
			} else {
				pos[k] = appear
			}
		}
	}

	return
}
