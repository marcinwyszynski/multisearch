package multisearch

import "unicode"

// WordFinder is an implementation of a Tokenizer designed to break up text on
// word boundaries.
func WordFinder(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}
