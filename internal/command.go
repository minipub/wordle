package internal

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type chars []string

// message be sent between pipe commmands
type message struct {
	// Command Input
	//
	// position of last round word state
	pos [5]int
	// last round word
	iWord [5]byte

	// Intermediate State
	//
	// store iWord chars
	hitLetters    chars
	appearLetters chars
	missLetters   chars
	// store iWord map: pos as key, letter as val
	hitIpl    map[int]chars
	appearIpl map[int]chars
	missIpl   map[int]chars
	// last round words
	lastWords []string
	// current round words
	nowWords []string

	// Command Output
	oWord [5]byte
}

func NewMessage() *message {
	return &message{
		hitLetters:    chars{},
		appearLetters: chars{},
		missLetters:   chars{},
		hitIpl:        make(map[int]chars),
		appearIpl:     make(map[int]chars),
		missIpl:       make(map[int]chars),
		lastWords:     words,
		nowWords:      []string{},
	}
}

func (m *message) Input(pos [5]int, iWord [5]byte) {
	m.pos = pos
	m.iWord = iWord
}

func (m *message) Output() [5]byte {
	return m.oWord
}

func (m *message) filter(notPattern, posPattern string) {
	for _, v := range m.lastWords {
		// discard those have missed letters
		match, _ := regexp.MatchString(notPattern, v)
		if match {
			continue
		}

		// discard those not in hited & appeared letters
		isExist := true

		for _, l := range m.hitLetters {
			isExist = isExist && IsIn([]byte(v), []byte(l)[0])
		}

		for _, l := range m.appearLetters {
			isExist = isExist && IsIn([]byte(v), []byte(l)[0])
		}

		if !isExist {
			continue
		}

		// save words which matched char position
		match, _ = regexp.MatchString(posPattern, v)
		if match {
			m.nowWords = append(m.nowWords, v)
		}
	}
}

// choose word
// Horizontal order: make chars inside word more mutually exclusive
// Vertical order: choose word which those chars occur more possibly
func (m *message) ChooseWord() (w [5]byte) {
	horders := make(map[int][]string)
	for i := 1; i <= 5; i++ {
		horders[i] = []string{}
	}

	// alphabet occur count
	var showCnts [26]int64

	for _, word := range m.nowWords {
		v := make(map[byte][]int)

		for i, b := range []byte(word) {
			v[b] = append(v[b], i)
			showCnts[b-'a'] += 1
		}

		mcnt := len(v)
		horders[mcnt] = append(horders[mcnt], word)
	}

	// descend
	sort.SliceStable(showCnts[:], func(i, j int) bool { return showCnts[i] > showCnts[j] })

	// number of digits
	digitNum := len(strconv.Itoa(len(words) * 5))
	// power base
	powBase := NewXBigInt(int64(math.Pow(10, float64(digitNum-1))))

	vorders := make(map[int][]string)
	for i := 1; i <= 5; i++ {
		vorders[i] = []string{}
	}

	rankWord := make(map[string][]string)

	for i := 5; i > 0; i-- {
		if len(horders[i]) > 0 {
			var topV []*Int

			for _, w := range horders[i] {
				var t [5]int64

				for k, v := range w {
					t[k] = showCnts[byte(v)-'a']
				}

				sort.SliceStable(t[:], func(i, j int) bool { return t[i] < t[j] })

				total := NewXBigInt(0)

				var j int64
				for j = 0; j < int64(len(t)); j++ {
					y := Mul(NewXBigInt(t[j]), Pow(powBase, NewXBigInt(j)))
					total.Add(y)
				}

				topV = append(topV, total)
				rankWord[fmt.Sprint(total)] = append(rankWord[fmt.Sprint(total)], w)
			}

			sort.SliceStable(topV, func(i, j int) bool { return topV[i].Gt(topV[j]) })

			for _, v := range topV {
				vorders[i] = append(vorders[i], rankWord[fmt.Sprint(v)]...)
			}
		}
	}

	if true {
		for i := 5; i > 0; i-- {
			if len(vorders[i]) > 0 {
				for k, v := range vorders[i][0] {
					w[k] = byte(v)
				}
				return
			}
		}
	} else {
		for i := 5; i > 0; i-- {
			if len(horders[i]) > 0 {
				w = RandOneWord(horders[i])
				return
			}
		}
	}

	return
}

type pipeCmd interface {
	execute(*message) bool
	next() pipeCmd
}

// DoPipeCmds Implements handling a set of pipe commands to solve word
// cat /tmp/words.txt | grep -v "[aplehi]" | grep t | grep n | grep s | grep "^[^t]\w\w[^n][^s]$"
func DoPipeCmds(m *message, pos [5]int, iWord [5]byte, verbose bool, p Writer) [5]byte {
	m.Input(pos, iWord)

	cc := &cmdConf{
		verbose,
		p,
	}

	for pc := NewMacroCmd(cc); pc.execute(m); {
		if pc = pc.next(); pc == nil {
			break
		}
	}

	return m.Output()
}

type cmdConf struct {
	verbose bool
	Writer
}

// the begin of the pipe commands
type macroCmd struct {
	*cmdConf
}

func NewMacroCmd(cc *cmdConf) pipeCmd {
	return &macroCmd{cc}
}

func (c *macroCmd) execute(m *message) bool {
	c.Write(`Thinking...
`)
	c.Write(fmt.Sprintf("pos: {{ %+v }}\n", m.pos))
	c.Write(fmt.Sprintf("iWord: {{ %+v }}\n", m.iWord))

	for k, v := range m.iWord {
		if v == byte(0) {
			continue
		}

		w := string(v)

		switch m.pos[k] {
		case Hit:
			if !IsIn(m.hitLetters, w) {
				m.hitLetters = append(m.hitLetters, w)
				m.hitIpl[k] = append(m.hitIpl[k], w)
			}
		case Appear:
			if !IsIn(m.appearLetters, w) {
				m.appearLetters = append(m.appearLetters, w)
				m.appearIpl[k] = append(m.appearIpl[k], w)
			}
		case Miss:
			if !IsIn(m.missLetters, w) {
				m.missLetters = append(m.missLetters, w)
				m.missIpl[k] = append(m.missIpl[k], w)
			}
		}
	}

	c.Write(fmt.Sprintf("hitLetters: {{ %+v }}\n", m.hitLetters))
	c.Write(fmt.Sprintf("appearLetters: {{ %+v }}\n", m.appearLetters))
	c.Write(fmt.Sprintf("missLetters: {{ %+v }}\n", m.missLetters))

	return true
}

func (c *macroCmd) next() pipeCmd {
	return NewPatternCmd(c.cmdConf)
}

type patternCmd struct {
	*cmdConf
}

func NewPatternCmd(cc *cmdConf) pipeCmd {
	return &patternCmd{cc}
}

func (c *patternCmd) execute(m *message) bool {
	// not pattern
	notPattern := fmt.Sprintf("[%s]", strings.Join(m.missLetters, ""))

	// position pattern
	var posPattern string

	for i := 0; i < 5; i++ {
		if v, ok := m.hitIpl[i]; ok {
			posPattern += v[0]
			continue
		}

		if v, ok := m.appearIpl[i]; ok {
			posPattern += fmt.Sprintf("[^%s]", strings.Join(v, ""))
			continue
		}

		posPattern += `\w`
	}

	posPattern = fmt.Sprintf("^%s$", posPattern)

	m.filter(notPattern, posPattern)

	c.Write(fmt.Sprintf("notPattern: {{ %s }}\n", notPattern))
	c.Write(fmt.Sprintf("posPattern: {{ %s }}\n", posPattern))

	return true
}

func (c *patternCmd) next() pipeCmd {
	return NewCandiCmd(c.cmdConf)
}

// choose the candidate word
type candiCmd struct {
	*cmdConf
}

func NewCandiCmd(cc *cmdConf) pipeCmd {
	return &candiCmd{cc}
}

func (c *candiCmd) execute(m *message) bool {
	var candiWords string
	for _, v := range m.nowWords {
		candiWords += fmt.Sprintln(v)
	}
	c.Write(fmt.Sprintf(`candiWords: 
%s`, candiWords))

	m.oWord = m.ChooseWord()

	c.Write(fmt.Sprintf(`chosen word: {{ %+v }}

`, m.oWord))

	return true
}

func (c *candiCmd) next() pipeCmd {
	return NewResetCmd()
}

// the end of the pipe commands
type resetCmd struct {
}

func NewResetCmd() pipeCmd {
	return &resetCmd{}
}

func (*resetCmd) execute(m *message) bool {
	m.lastWords = m.nowWords
	m.nowWords = []string{}
	return true
}

func (*resetCmd) next() pipeCmd {
	return nil
}
