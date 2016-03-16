package rustemmer_test

import (
	"fmt"
	"github.com/liderman/rustemmer"
)

func Example_GetWordBase() {
	word := "вазы"
	wordBase := rustemmer.GetWordBase(word)
	fmt.Printf("%s => %s!\n", word, wordBase)
}

func Example_NormalizeText() {
	text := "г. Москва, ул. Полярная, д. 31А, стр. 1"
	fmt.Printf(
		"%s => %s!\n",
		text,
		rustemmer.NormalizeText(text),
	)
}
