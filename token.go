package multisearch

import (
	"regexp"
	"sort"
	"strings"
)

var (
	// A regular expression matcher accepting all Unicode letters and
	// ASCII word characters (numbers, underscore).
	charRegexp = regexp.MustCompile("[\\pL\\w]")
)

type token struct {
	content        string
	ignored        bool
	previous, next *token
	matchedBy      []*node
}

// newToken is a constructor for the token object.
func newToken() *token {
	return &token{
		matchedBy: make([]*node, 0),
	}
}

func (t *token) Ignored() bool {
	return t.ignored
}

func (t *token) Next() Token {
	if t.next == nil {
		return nil
	}
	return t.next
}

func (t *token) Matches() []Match {
	retVal := make([]Match, len(t.matchedBy), len(t.matchedBy))
	for i, node := range t.matchedBy {
		retVal[i] = node
	}
	return retVal
}

func (t *token) TopMatch() Match {
	if !t.matched() {
		return nil
	}
	sort.Sort(byWeight(t.matchedBy))
	return t.matchedBy[0]
}

func (t *token) String() string {
	return t.content
}

// recordMatch back-propagates a successful match. Please note that ignored
// tokens are also marked with the match.
func (t *token) recordMatch(match *node) {
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

// matched reports whether the token has been matched by one or more terminal
// nodes.
func (t *token) matched() bool {
	return len(t.matchedBy) > 0
}

func tokenize(input string, callback func(*token)) *token {
	currentToken := newToken()
	firstToken := currentToken
	var previousToken *token
	for _, char := range strings.Split(input, "") {
		isWord := charRegexp.MatchString(char)
		if currentToken.content == "" || isWord != currentToken.ignored {
			currentToken.content += char
			currentToken.ignored = !isWord
			continue
		}
		if isWord == currentToken.ignored {
			previousToken = currentToken
			currentToken = newToken()
			previousToken.next = currentToken
			currentToken.content = char
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
