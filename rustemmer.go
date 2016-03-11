// Package rustemmer implements Porter stemmer for Russian language.
package rustemmer

import (
	"regexp"
	"strings"
	"sync"
)

var instance = New()

const VOWEL = "аеёиоуыэюя"

var suffixNN = []string{"нн"}
var suffixPerfectiveGerunds = [][]string{
	{"в", "вши", "вшись", "в", "вши", "вшись"},
	{"ив", "ивши", "ившись", "ыв", "ывши", "ывшись"},
}
var suffixReflexives = []string{"ся", "сь"}
var suffixAdjective = []string{
	"ее", "ие", "ые", "ое", "ими", "ыми", "ей", "ий", "ый", "ой", "ем", "им", "ым", "ом", "его", "ого", "ему",
	"ому", "их", "ых", "ую", "юю", "ая", "яя", "ою", "ею",
}
var suffixVerb = [][]string{
	{"ла", "на", "ете", "йте", "ли", "й", "л", "ем", "н", "ло", "но", "ет", "ют", "ны", "ть", "ешь", "нно"},
	{
		"ила", "ыла", "ена", "ейте", "уйте", "ите", "или", "ыли", "ей", "уй", "ил", "ыл", "им", "ым", "ен",
		"ило", "ыло", "ено", "ят", "ует", "уют", "ит", "ыт", "ены", "ить", "ыть", "ишь", "ую", "ю",
	},
}

var suffixNoun = []string{
	"а", "ев", "ов", "ие", "ье", "е", "иями", "ями", "ами", "еи", "ии", "и", "ией", "ей", "ой", "ий", "й", "иям",
	"ям", "ием", "ем", "ам", "ом", "о", "у", "ах", "иях", "ях", "ы", "ь", "ию", "ью", "ю", "ия", "ья", "я",
}
var suffixSuperlative = []string{"ейш", "ейше"}
var suffixSoftSign = []string{"ь"}
var suffixI = []string{"и"}
var suffixDerivational = []string{"ост", "ость"}
var suffixParticiple = [][]string{
	appendPrefix(suffixAdjective, []string{"ем", "нн", "вш", "ющ", "щ"}),
	appendPrefix(suffixAdjective, []string{"ивш", "ывш", "ующ"}),
}

type RuStemmer struct {
	mu sync.Mutex
	word []rune
	RV int
	R2 int
}

// New creates a new RuStemmer.
func New() *RuStemmer {
	return &RuStemmer{
		word: []rune(""),
		RV: 0,
		R2: 0,
	}
}

// GetWordBase returns the base word.
func GetWordBase(word string) string {
	instance.mu.Lock()
	defer instance.mu.Unlock()
	return instance.GetWordBase(word)
}

// NormalizeText returns normalized text.
// Returns text in which all words will be replaced with the basics of words separated by a space.
// All Special characters except "_" will be removed.
func NormalizeText(text string) string {
	instance.mu.Lock()
	defer instance.mu.Unlock()
	return instance.NormalizeText(text)
}

// GetWordBase returns the base word.
func (r *RuStemmer) GetWordBase(word string) string {
	r.word = []rune(word)
	r.RV = 0
	r.R2 = 0
	r.findRegions()

	// Step 1
	// Find ending PERFECTIVE GERUND. If it exists - delete it and complete this step
	if !r.removeEndings(r.RV, suffixPerfectiveGerunds[0], suffixPerfectiveGerunds[1]) {
		// Otherwise, remove ending REFLEXIVE (if it exists)
		r.removeEndings(r.RV, suffixReflexives)
		// Then try the following procedure to remove ending: ADJECTIVE, VERB, NOUN.
		// As soon as one of them is found - a step ends
		ife := r.removeEndings(
			r.RV,
			suffixParticiple[0],
			suffixParticiple[1],
		) || r.removeEndings(r.RV, suffixAdjective)

		if !ife && !r.removeEndings(r.RV, suffixVerb[0], suffixVerb[1]) {
			r.removeEndings(r.RV, suffixNoun)
		}
	}

	// Step 2
	// If a word ends with "и" - remove the "и"
	r.removeEndings(r.RV, suffixI)

	// Step 3
	// If in "R2" there DERIVATIONAL ending - delete it
	r.removeEndings(r.R2, suffixDerivational)

	// Step 4
	// Possible is one of the three variants:
	// If a word ending in "нн" - delete the last letter
	if r.removeEndings(r.RV, suffixNN) {
		r.word = []rune(string(r.word) + "н")
	}

	// If a word ending in SUPERLATIVE - remove it and remove the last letter again if the word ending in "нн"
	r.removeEndings(r.RV, suffixSuperlative)
	// If a word ending in "ь" - delete it
	r.removeEndings(r.RV, suffixSoftSign)

	return string(r.word)
}

// NormalizeText returns normalized text.
// Returns text in which all words will be replaced with the basics of words separated by a space.
// All Special characters except "_" will be removed.
func (r *RuStemmer) NormalizeText(text string) string {
	regexWords := regexp.MustCompile("[\\p{L}\\d_]+")
	words := regexWords.FindAllString(text, -1)
	for k, word := range words {
		words[k] = r.GetWordBase(word)
	}

	return strings.Join(words, " ")
}

func (r *RuStemmer) removeEndings(region int, suffixesPacks ...[]string) bool {
	if region > len(r.word) {
		region = len(r.word)
	}

	prefix := r.word[:region]
	word_ := string(r.word[len(prefix):])

	suffixes := suffixesPacks[0]
	if len(suffixesPacks) == 2 {
		if result := trimFirstSuffix(word_, suffixes, true); result != word_ {
			r.word = []rune(string(prefix) + result)
			return true
		}
		suffixes = suffixesPacks[1]
	}

	if result := trimFirstSuffix(word_, suffixes, false); result != word_ {
		r.word = []rune(string(prefix) + result)
		return true
	}

	return false
}

func appendPrefix(strs []string, prefixesPacks ...[]string) []string {
	ret := []string{}
	for _, str := range strs {
		for _, prefixes := range prefixesPacks {
			for _, prefix := range prefixes {
				ret = append(ret, prefix + str)
			}
		}
	}

	return ret
}

func trimFirstSuffix(word string, suffixes []string, isAYA bool) string {
	for _, suffix := range suffixes {
		if isAYA && !(strings.HasSuffix(word, "а" + suffix) || strings.HasSuffix(word, "я" + suffix)) {
			continue
		}

		if result := strings.TrimSuffix(word, suffix); result != word {
			return result
		}
	}
	return word
}

func (r *RuStemmer) findRegions() {
	state := 0
	wordLength := len(r.word)
	for i := 1; i < wordLength; i++ {
		prevChar := string(r.word[i - 1])
		char     := string(r.word[i])
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

func (r *RuStemmer) isVowel(char string) bool {
	return strings.Contains(VOWEL, char)
}