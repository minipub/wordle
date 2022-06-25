package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	m := message{
		hitLetters:    chars{"y", "i"},
		appearLetters: chars{"d", "m"},
		missLetters:   chars{"c", "a", "v", "e", "r", "s", "o", "u", "p", "b", "n", "g", "w", "f", "t", "l"},

		lastWords: []string{
			"diddy",
			"dilly",
			"dimly",
			"dizzy",
			"hilly",
			"jimmy",
			"kiddy",
			"middy",
			"milky",
			"mizzy",
		},
	}

	notPattern := "[caversoupbngwftl]"
	posPattern := `^[^d]i[^m]\wy$`

	m.filter(notPattern, posPattern)

	assert.Equal(t, 1, len(m.nowWords))
	assert.Equal(t, "middy", m.nowWords[0])
}
