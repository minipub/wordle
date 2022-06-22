package internal

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type chars []string

var (
	// store iWord chars
	hitLetters    = []string{}
	appearLetters = []string{}
	missLetters   = []string{}

	// store iWord char eg. pos => char
	hitIpc    = make(map[int]chars)
	appearIpc = make(map[int]chars)
	missIpc   = make(map[int]chars)

	lastWords = words
	nowWords  = []string{}
)

// cat /tmp/words.txt | grep -v "[aplehi]" | grep t | grep n | grep s | grep "^[^t]\w\w[^n][^s]$"
func SolveWord(pos [5]int, iWord [5]byte) (rs [5]byte) {
	fmt.Fprintf(os.Stderr, "pos: {{ %+v }}\n", pos)
	fmt.Fprintf(os.Stderr, "iWord: {{ %+v }}\n", iWord)

	for k, v := range iWord {
		w := string(v)

		switch pos[k] {
		case Hit:
			if !IsIn(hitLetters, w) {
				hitLetters = append(hitLetters, w)
				hitIpc[k] = append(hitIpc[k], w)
			}
		case Appear:
			if !IsIn(appearLetters, w) {
				appearLetters = append(appearLetters, w)
				appearIpc[k] = append(appearIpc[k], w)
			}
		case Miss:
			if !IsIn(missLetters, w) {
				missLetters = append(missLetters, w)
				missIpc[k] = append(missIpc[k], w)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "hitLetters: {{ %+v }}\n", hitLetters)
	fmt.Fprintf(os.Stderr, "appearLetters: {{ %+v }}\n", appearLetters)
	fmt.Fprintf(os.Stderr, "missLetters: {{ %+v }}\n", missLetters)

	// not pattern
	notPattern := fmt.Sprintf("[%s]", strings.Join(missLetters, ""))

	// position pattern
	var posPattern string

	for i := 0; i < 5; i++ {
		if v, ok := hitIpc[i]; ok {
			posPattern += v[0]
			continue
		}

		if v, ok := appearIpc[i]; ok {
			posPattern += fmt.Sprintf("[^%s]", strings.Join(v, ""))
			continue
		}

		posPattern += `\w`
	}

	posPattern = fmt.Sprintf("^%s$", posPattern)

	sillyFilter(notPattern, posPattern)

	fmt.Fprintf(os.Stderr, "notPattern: {{ %s }}\n", notPattern)
	fmt.Fprintf(os.Stderr, "posPattern: {{ %s }}\n", posPattern)

	var candiWords string
	for _, v := range nowWords {
		candiWords += fmt.Sprintln(v)
	}
	fmt.Fprintf(os.Stderr, `candiWords: 
%s`, candiWords)

	rs = ChooseWord()
	fmt.Fprintf(os.Stderr, `chosen word: {{ %+v }}

`, rs)

	// at the end
	reset()

	return
}

func reset() {
	lastWords = nowWords
	nowWords = []string{}
}

func sillyFilter(notPattern, posPattern string) {

	for _, v := range lastWords {
		// discard those have missed letters
		match, _ := regexp.MatchString(notPattern, v)
		if match {
			continue
		}

		// discard those not in hited & appeared letters
		isExist := true

		for _, m := range hitLetters {
			isExist = isExist && IsIn([]byte(v), []byte(m)[0])
		}

		for _, m := range appearLetters {
			isExist = isExist && IsIn([]byte(v), []byte(m)[0])
		}

		if !isExist {
			continue
		}

		// save words which matched char position
		match, _ = regexp.MatchString(posPattern, v)
		if match {
			nowWords = append(nowWords, v)
		}
	}

}

// choose word
func ChooseWord() (w [5]byte) {
	// layer step down internal mutually exclusive
	layerDownMutexWords := make(map[int][]string)
	for i := 1; i <= 5; i++ {
		layerDownMutexWords[i] = []string{}
	}

	for _, v := range nowWords {
		m := make(map[byte][]int)
		for i, b := range []byte(v) {
			m[b] = append(m[b], i)
		}

		mcnt := len(m)
		layerDownMutexWords[mcnt] = append(layerDownMutexWords[mcnt], v)
	}

	for i := 5; i > 0; i-- {
		if len(layerDownMutexWords[i]) > 0 {
			w = RandOneWord(layerDownMutexWords[i])
			return
		}
	}

	return
}
