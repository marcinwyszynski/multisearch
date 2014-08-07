package multisearch

import (
	"unicode"
)

type tokenImpl struct {
	content        string
	ignored        bool
	previous, next *tokenImpl
	matchedBy      []*matchImpl
}

// newTokenImpl is a constructor for the tokenImpl object.
func newTokenImpl() *tokenImpl {
	return &tokenImpl{
		matchedBy: make([]*matchImpl, 0),
	}
}

func (t *tokenImpl) Ignored() bool {
	return t.ignored
}

func (t *tokenImpl) Next() Token {
	if t.next == nil {
		return nil
	}
	return t.next
}

func (t *tokenImpl) Matches() []Match {
	retVal := make([]Match, len(t.matchedBy), len(t.matchedBy))
	for i, node := range t.matchedBy {
		retVal[i] = node
	}
	return retVal
}

func (t *tokenImpl) String() string {
	return t.content
}

// recordMatch back-propagates a successful match. Please note that ignored
// tokens are also marked with the match.
func (t *tokenImpl) recordMatch(match *matchImpl) {
	this, recorded := t, 0
	for {
		this.matchedBy = append(this.matchedBy, match)
		if !this.ignored {
			recorded++
		}
		if recorded == match.depth {
			return
		}
		this = this.previous
	}
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func tokenize(input string, callback func(*tokenImpl)) *tokenImpl {
	currentToken := newTokenImpl()
	firstToken := currentToken
	var previousToken *tokenImpl
	for _, r := range input {
		isWord := unicode.IsLetter(r) || unicode.IsNumber(r)
		if currentToken.content == "" || isWord != currentToken.ignored {
			currentToken.content += string(r)
			currentToken.ignored = !isWord
			continue
		}
		if isWord == currentToken.ignored {
			previousToken = currentToken
			currentToken = newTokenImpl()
			previousToken.next = currentToken
			currentToken.content = string(r)
			currentToken.previous = previousToken
			currentToken.ignored = !isWord
			callback(previousToken)
		}
	}
	if currentToken.content != "" {
		callback(currentToken)
	}
	return firstToken
}
