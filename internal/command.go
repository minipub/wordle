package internal

import (
	"fmt"
	"regexp"
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

		for _, v := range m.hitLetters {
			isExist = isExist && IsIn([]byte(v), []byte(v)[0])
		}

		for _, v := range m.appearLetters {
			isExist = isExist && IsIn([]byte(v), []byte(v)[0])
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
func (m *message) ChooseWord() (w [5]byte) {
	// layer step down internal mutually exclusive
	layerDownMutexWords := make(map[int][]string)
	for i := 1; i <= 5; i++ {
		layerDownMutexWords[i] = []string{}
	}

	for _, v := range m.nowWords {
		vv := make(map[byte][]int)
		for i, b := range []byte(v) {
			vv[b] = append(vv[b], i)
		}

		mcnt := len(vv)
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

type pipeCmd interface {
	execute(*message) bool
	next() pipeCmd
}

// DoPipeCmds Implements a set of command using pipe to filter
// cat /tmp/words.txt | grep -v "[aplehi]" | grep t | grep n | grep s | grep "^[^t]\w\w[^n][^s]$"
func DoPipeCmds(m *message, pos [5]int, iWord [5]byte, verbose bool) [5]byte {
	m.Input(pos, iWord)

	cc := &cmdConf{
		verbose,
		NewSolverPrinter(verbose),
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
	*SolverPrinter
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
