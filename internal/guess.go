package internal

import (
	"errors"
	"fmt"
)

const (
	// status positon of inputted char at in-plan word
	// Hit: green, Appear: yellow, Miss: gray
	Hit = iota
	Appear
	Miss

	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorRed    = "\033[31m"

	ColoredByteNum    = len(ColorRed)
	ColorResetByteNum = len(ColorReset)

	Prompt = "input: "

	PreText = `A Wordle Game!

Please input a five-letter word and Press <Enter> to confirm.

`
)

var (
	Colors     map[int]string
	CheerWords map[int]string
)

func init() {
	Colors = map[int]string{
		Hit:    ColorGreen,
		Appear: ColorYellow,
		Miss:   ColorRed,
	}

	CheerWords = map[int]string{
		0: "God Like!!!!",
		1: "Holy Shit!!!",
		2: "Splendid!!",
		3: "Great Job!",
		4: "Well Done",
		5: "Not Bad",
	}
}

func DoPuzzle(rw ReadWriter) {
	pWord := RandOneWord(words) // in-plan word
	PreWord(rw)
	i := GuessWord(rw, pWord)
	PostWord(rw, i, pWord)
}

func GuessWord(rw ReadWriter, pWord [5]byte) (cnt int) {
	cnt = -1
	for i := 0; i < len(CheerWords); {
		// handle different writer
		rw.Write(Prompt)
		iWord, err := HandleInput(rw) // inputted word
		if err != nil {
			rw.Write(fmt.Sprintf("Error: %+v\n", err))
			continue
		}

		ok := Equal(pWord, iWord)
		if ok {
			cnt = i
			return
		}

		pos := FindPos(pWord, iWord)
		var s string
		for m, n := range pos {
			s += fmt.Sprint(Colors[n], string(iWord[m]), ColorReset)
		}
		rw.Write(fmt.Sprintln(s))

		i++
	}
	return
}

func PreWord(rw ReadWriter) {
	rw.Write(PreText)
}

func PostWord(rw ReadWriter, i int, pWord [5]byte) {
	if i > -1 {
		rw.Write(fmt.Sprintf("\n%s\n", CheerWords[i]))
	} else {
		rw.Write(fmt.Sprintf(`
Out of Chance!

The Word is <%s>

Take a break or get another round.
`, pWord))
	}
}

func HandleInput(rw ReadWriter) (rs [5]byte, err error) {
	var b [5]byte
	n, err := rw.Read(b[:])
	if err != nil {
		return
	}

	// fmt.Printf("b read: %+v\n", b[:])

	if isCRLF(b[n-1]) {
		err = errors.New("no enough letters")
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
			err = errors.New("letters must be between a-zA-z")
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
// pos: status position at inputted word
func FindPos(x, y [5]byte) (pos [5]int) {
	// map store x's char & positons
	xps := make(map[byte][]int)
	for k, v := range x {
		xps[v] = append(xps[v], k)
	}

	for k, v := range y {
		if _, ok := xps[v]; !ok || len(xps[v]) == 0 {
			pos[k] = Miss
		} else {
			if IsIn(xps[v], k) {
				pos[k] = Hit
			} else {
				pos[k] = Appear
			}
		}
	}

	return
}
