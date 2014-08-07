package multisearch

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Word splitter is essentially a reverse of character matcer.
	wordSplitter = regexp.MustCompile("[^\\w\\pL]+")
)

type Engine struct {
	// Root of the search tree.
	root *matchImpl

	// Collection of ignored words. Stored in a map for greater retrieval
	// efficiency.
	ignores map[string]struct{}

	// Mapping of original terms to matched matchImpls.
	originals map[*matchImpl]string

	// Stemmer to be used for sanitization purposes.
	stemmer Stemmer
}

func NewEngine(stemmer Stemmer) *Engine {
	return &Engine{
		root:      newMatchImpl("", 0),
		ignores:   make(map[string]struct{}),
		originals: make(map[*matchImpl]string),
		stemmer:   stemmer,
	}
}

func (e *Engine) Ignore(stopword string) error {
	if len(wordSplitter.Split(stopword, -1)) > 1 {
		return fmt.Errorf("ignore not a single word: %q", stopword)
	}
	stopword = e.stemmer.StemString(strings.ToLower(stopword))
	if _, exists := e.ignores[stopword]; exists {
		return fmt.Errorf("duplicate ignore: %q", stopword)
	}
	e.ignores[stopword] = struct{}{}
	return nil
}

func (e *Engine) Match(needle string) (Match, error) {
	sanitized := e.sanitize(needle)
	if len(sanitized) == 0 {
		return nil, fmt.Errorf("only consists of ignores: %q", needle)
	}
	newMatch := e.root.add(sanitized)
	if original, existed := e.originals[newMatch]; existed {
		return nil, fmt.Errorf("duplicate of %q: %q", original, needle)
	}
	e.originals[newMatch] = needle
	return newMatch, nil
}

func (e *Engine) Process(input string) Token {
	cursors := make(map[*matchImpl]struct{})
	cursors[e.root] = struct{}{}
	return tokenize(input, func(t *tokenImpl) {
		if t.ignored {
			return
		}
		stem := strings.ToLower(e.stemmer.StemString(t.content))
		if _, isIgnored := e.ignores[stem]; isIgnored {
			t.ignored = true
			return
		}
		cursorsToDelete, cursorsToAdd := make([]*matchImpl, 0), make([]*matchImpl, 0)
		for cursor, _ := range cursors {
			nextCursor, exists := cursor.children[stem]
			if exists && nextCursor.terminal {
				t.recordMatch(nextCursor)
			}
			if exists && !nextCursor.terminal {
				cursorsToAdd = append(cursorsToAdd, nextCursor)
			}
			if cursor != e.root {
				cursorsToDelete = append(cursorsToDelete, cursor)
			}
		}
		for _, cursor := range cursorsToDelete {
			delete(cursors, cursor)
		}
		for _, cursor := range cursorsToAdd {
			cursors[cursor] = struct{}{}
		}
	})
}

func (e *Engine) sanitize(text string) []string {
	retVal := make([]string, 0)
	for _, token := range wordSplitter.Split(text, -1) {
		word := e.stemmer.StemString(strings.ToLower(token))
		if _, isStopword := e.ignores[word]; !isStopword {
			retVal = append(retVal, word)
		}
	}
	return retVal
}
