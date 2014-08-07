package multisearch

import (
	"fmt"
	"strings"
)

type engineImpl struct {
	// Root of the search tree.
	root *matchImpl

	// Collection of ignored words. Stored in a map for greater retrieval
	// efficiency.
	ignores map[string]struct{}

	// Mapping of original terms to matched matchImpls.
	originals map[*matchImpl]string

	// Stemmer to be used for sanitization purposes.
	stemmer Stemmer

	// Tokenizer to be used for tokenization purposes.
	tokenizer Tokenizer
}

func NewEngine(stemmer Stemmer, tokenizer Tokenizer) Engine {
	return &engineImpl{
		root:      newMatchImpl("", 0),
		ignores:   make(map[string]struct{}),
		originals: make(map[*matchImpl]string),
		stemmer:   stemmer,
		tokenizer: tokenizer,
	}
}

func (e *engineImpl) Ignore(input string) error {
	chunks := e.sanitize(input)
	if len(chunks) == 0 {
		return fmt.Errorf("duplicate or empty ignore: %q", input)
	}
	if len(chunks) != 1 {
		return fmt.Errorf("ignore not a single word: %q", input)
	}
	e.ignores[chunks[0]] = struct{}{}
	return nil
}

func (e *engineImpl) Match(needle string) (Match, error) {
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

func (e *engineImpl) Process(input string) []Token {
	cursors, tokens := make(map[*matchImpl]struct{}), make([]Token, 0)
	cursors[e.root] = struct{}{}
	var lastToken *tokenImpl = nil
	tokenize(input, e.tokenizer, func(start, end int, captured bool) {
		t := newTokenImpl()
		tokens = append(tokens, t)
		t.content, t.ignored = input[start:end], !captured
		if t.ignored {
			return
		}
		stem := strings.ToLower(e.stemmer.StemString(t.content))
		if _, isIgnored := e.ignores[stem]; isIgnored {
			t.ignored = true
			return
		}
		t.previous, lastToken = lastToken, t
		cDel, cAdd := make([]*matchImpl, 0), make([]*matchImpl, 0)
		for cursor, _ := range cursors {
			nextCursor, exists := cursor.children[stem]
			if exists && nextCursor.terminal {
				t.recordMatch(nextCursor)
			}
			if exists && !nextCursor.terminal {
				cAdd = append(cAdd, nextCursor)
			}
			if cursor != e.root {
				cDel = append(cDel, cursor)
			}
		}
		for _, cursor := range cDel {
			delete(cursors, cursor)
		}
		for _, cursor := range cAdd {
			cursors[cursor] = struct{}{}
		}
	})
	return tokens
}

func (e *engineImpl) sanitize(input string) []string {
	retVal := make([]string, 0)
	tokenize(input, e.tokenizer, func(start, end int, captured bool) {
		if !captured {
			return
		}
		word := e.stemmer.StemString(strings.ToLower(input[start:end]))
		if _, isStopword := e.ignores[word]; !isStopword {
			retVal = append(retVal, word)
		}
	})
	return retVal
}
