package multisearch

// Tokenizer is a function capable of taking a rune as an input and deciding
// whether the rune should be captured (true) or ignored (false).
type Tokenizer func(rune) bool

// bounaryCallack is a function called by the tokenize method on each token it
// identifies.
type bounaryCallack func(start, end int, matched bool)

// Engine represents a multisearch engine capable of registering new matches,
// ignores (stopwords) and - most importantly - processing new inputs.
type Engine interface {
	// Ignore registers a new ignored word.
	Ignore(input string) error

	// Match registers and returns a new Match.
	Match(needle string) (Match, error)

	// Process breaks up the input into a slice of Tokens based on Matches
	// and ignores.
	Process(input string) []Token
}

// Match represents a matched term.
type Match interface {
	// Size provides the total length of the Match.
	Size() int

	// String returns a full string representation of the Match.
	String() string
}

// Token represents a piece of processed text. From user's perspective this is
// behaves like a single linked list.
type Token interface {
	// Ignored tells whether this given piece of text has been ignored.
	Ignored() bool

	// Matches returns an array of Matches that were matched against
	// this particular piece of text.
	Matches() []Match

	// String returns the content of the token.
	String() string
}

// Stemmer does the job of stemming words.
type Stemmer interface {
	// StemString stems a string to produce another string. The original
	// is not changed in the process.
	StemString(word string) string
}
