package multisearch

// Match represents a matched term.
type Match interface {
	// Size provides the total length of the Match.
	Size() int

	// String returns a full string representation of the Match.
	String() string

	// Weight represents a user-defined weight of the Match.
	Weight() int
}

// Token represents a piece of processed text. From user's perspective this is
// behaves like a single linked list.
type Token interface {
	// Ignored tells whether this given piece of text has been ignored.
	Ignored() bool

	// Next provides a link to the subsequent Token.
	Next() Token

	// Matches returns an array of Matches that were matched against
	// this particular piece of text.
	Matches() []Match

	// TopMatch represents the highest scoring Match with weight considered
	// first, and size considered second.
	TopMatch() Match

	// String returns the content of the token.
	String() string
}

// Stemmer does the job of stemming words.
type Stemmer interface {
	// StemString stems a string to produce another string. The original
	// is not changed in the process.
	StemString(word string) string
}
