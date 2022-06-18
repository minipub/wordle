package internal

import (
	"fmt"
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
	hitIpc    = map[int]chars{}
	appearIpc = map[int]chars{}
	missIpc   = map[int]chars{}

	lastWords = words
	nowWords  = []string{}
)

// cat /tmp/words.txt | grep -v "[aplehi]" | grep t | grep n | grep s | grep "^[^t]\w\w[^n][^s]$"
func ChooseWord(pos [5]int, iWord [5]byte) {
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

}

func sillyFilter(notPattern, posPattern string) {

	for _, v := range lastWords {
		// discard those have missed letters
		match, _ := regexp.MatchString(notPattern, v)
		if match {
			continue
		}

		// discard those not in hited & appeared letters
		for _, m := range hitLetters {
			if !IsIn([]byte(v), []byte(m)[0]) {
				continue
			}
		}

		for _, m := range appearLetters {
			if !IsIn([]byte(v), []byte(m)[0]) {
				continue
			}
		}

		// save words which matched char position
		match, _ = regexp.MatchString(posPattern, v)
		if match {
			nowWords = append(nowWords, v)
		}
	}

}
