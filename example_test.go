package rustemmer_test

import (
	"fmt"
	"github.com/liderman/rustemmer"
)

func ExampleGetWordBase() {
	word := "вазы"
	wordBase := rustemmer.GetWordBase(word)
	fmt.Printf("%s => %s!\n", word, wordBase)
}

func ExampleNormalizeText() {
	text := "г. Москва, ул. Полярная, д. 31А, стр. 1"
	fmt.Printf(
		"%s => %s!\n",
		text,
		rustemmer.NormalizeText(text),
	)
}
