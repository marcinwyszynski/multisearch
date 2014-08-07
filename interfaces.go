package multisearch

type Match interface {
	Size()   int
	String() string
	Weight() int
}

type Token interface {
	Ignored()  bool
	Next()     Token
	Matches()  []Match
	TopMatch() Match
	String()   string
}

// Stemmer does the job of stemming words.
type Stemmer interface {
	// StemString stems a string to produce another string. The original
	// is not changed in the process.
	StemString(word string) string
}