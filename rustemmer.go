// Package rustemmer implements Porter stemmer for Russian language.
package rustemmer

import (
	"regexp"
	"strings"
	"sync"
)

var (
	st = New()
)

const (
	VOWEL = "аеёиоуыэюя"
	REGEX_ADJECTIVE = "(ее|ие|ые|ое|ими|ыми|ей|ий|ый|ой|ем|им|ым|ом|его|ого|ему|ому|их|ых|ую|юю|ая|яя|ою|ею)$"
	REGEX_SOFT_SIGN = "ь$"
	REGEX_NN = "нн$"
	REGEX_I = "и$"
	REGEX_REFLEXIVES = "(ся|сь)$"
	REGEX_NOUN = "(а|ев|ов|ие|ье|е|иями|ями|ами|еи|ии|и|ией|ей|ой|ий|й|иям|ям|ием|ем|ам|ом|о|у|ах|иях|ях|ы|ь|ию|ью|ю|ия|ья|я)$"
	REGEX_SUPERLATIVE = "(ейш|ейше)$"
	REGEX_DERIVATIONAL = "(ост|ость)$"
	REGEX_PERFECTIVE_GERUNDS_1 = "(в|вши|вшись)$"
	REGEX_PERFECTIVE_GERUNDS_2 = "(ив|ивши|ившись|ыв|ывши|ывшись)$"
	REGEX_PARTICIPLE_1 = "(ем|нн|вш|ющ|щ)"
	REGEX_PARTICIPLE_2 = "(ивш|ывш|ующ)"
	REGEX_VERB_1 = "(ла|на|ете|йте|ли|й|л|ем|н|ло|но|ет|ют|ны|ть|ешь|нно)$"
	REGEX_VERB_2 = "(ила|ыла|ена|ейте|уйте|ите|или|ыли|ей|уй|ил|ыл|им|ым|ен|ило|ыло|ено|ят|ует|уют|ит|ыт|ены|ить|ыть|ишь|ую|ю)$"
)


type RuStemmer struct {
	sync.Mutex
	word []rune
	RV int
	R2 int

}

func New() *RuStemmer {
	return &RuStemmer{
		word: []rune(""),
		RV: 0,
		R2: 0,
	}
}

// GetWordBase returns the base word.
func (r *RuStemmer) GetWordBase(word string) string {
	r.Lock()
	defer r.Unlock()
	r.word = []rune(word)
	r.RV = 0
	r.R2 = 0
	r.findRegions()

	// Step 1
	// Find ending PERFECTIVE GERUND. If it exists - delete it and complete this step
	if !r.removeEndings([]string{REGEX_PERFECTIVE_GERUNDS_1, REGEX_PERFECTIVE_GERUNDS_2}, r.RV) {
		// Otherwise, remove ending REFLEXIVE (if it exists)
		r.removeEndings([]string{REGEX_REFLEXIVES}, r.RV)
		// Then try the following procedure to remove ending: ADJECTIVE, VERB, NOUN.
		// As soon as one of them is found - a step ends
		ife := r.removeEndings(
			[]string{
				REGEX_PARTICIPLE_1 + REGEX_ADJECTIVE,
				REGEX_PARTICIPLE_2 + REGEX_ADJECTIVE,
			},
			r.RV,
		) || r.removeEndings([]string{REGEX_ADJECTIVE}, r.RV)

		if !ife && !r.removeEndings([]string{REGEX_VERB_1, REGEX_VERB_2}, r.RV) {
			r.removeEndings([]string{REGEX_NOUN}, r.RV)
		}
	}

	// Step 2
	// If a word ends with "и" - remove the "и"
	r.removeEndings([]string{REGEX_I}, r.RV)
	// Step 3
	// If in "R2" there DERIVATIONAL ending - delete it
	r.removeEndings([]string{REGEX_DERIVATIONAL}, r.R2)
	// Step 4
	// Possible is one of the three variants:
	// If a word ending in "нн" - delete the last letter
	if r.removeEndings([]string{REGEX_NN}, r.RV) {
		r.word = []rune(string(r.word) + "н")
	}
	// If a word ending in SUPERLATIVE - remove it and remove the last letter again if the word ending in "нн"
	r.removeEndings([]string{REGEX_SUPERLATIVE}, r.RV)
	// If a word ending in "ь" - delete it
	r.removeEndings([]string{REGEX_SOFT_SIGN}, r.RV)

	return string(r.word)
}

func (r *RuStemmer) removeEndings(regex []string, region int) bool {
	if region > len(r.word) {
		region = len(r.word)
	}

	prefix := r.word[:region]
	word_ := r.word[len(prefix):]
	if len(regex) > 1 {
		if regexp.MustCompile(".+[а|я]" + regex[0]).MatchString(string(word_)) {
			r.word = []rune(string(prefix) + regexp.MustCompile(regex[0]).ReplaceAllString(string(word_), ""))
			return true
		}
		regex = []string{regex[1]}
	}

	if regexp.MustCompile(".+" + regex[0]).MatchString(string(word_)) {
		r.word = []rune(string(prefix) + regexp.MustCompile(regex[0]).ReplaceAllString(string(word_), ""))
		return true
	}
	return false
}

func (r *RuStemmer) findRegions() {
	state := 0
	wordLength := len(r.word)
	for i := 1; i < wordLength; i++ {
		prevChar := r.word[i - 1]
		char     := r.word[i]
		switch state {
			case 0:
				if r.isVowel(char) {
					r.RV = i + 1
					state = 1
				}
				break
			case 1:
				if r.isVowel(prevChar) && !r.isVowel(char) {
					state = 2
				}
				break
			case 2:
				if r.isVowel(prevChar) && !r.isVowel(char) {
					r.R2 = i + 1
					return
				}
				break
		}
	}
}

func (r *RuStemmer) isVowel(char rune) bool {
	return strings.Contains(VOWEL, string(char))
}