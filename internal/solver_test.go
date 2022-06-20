package internal

import "testing"

func TestSillyFilter(t *testing.T) {
	hitLetters = []string{}
	appearLetters = []string{"g", "e", "a"}
	missLetters = []string{"r", "t"}
	notPattern := "[rt]"
	posPattern := `^[^g]\w[^e][^a]\w$`

	sillyFilter(notPattern, posPattern)

	t.Log("nowWords:")
	t.Log(nowWords)
	t.Log(len(nowWords))
}
