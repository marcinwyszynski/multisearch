package multisearch

type tokenImpl struct {
	content   string
	ignored   bool
	previous  *tokenImpl
	matchedBy []*matchImpl
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

func tokenize(input string, tok Tokenizer, cb bounaryCallack) {
	var captured bool
	start, inLength := 0, len(input)
	for pos, r := range input {
		thisCaptured := tok(r)
		if pos == 0 {
			captured = thisCaptured
			continue
		}
		if captured != thisCaptured {
			cb(start, pos, captured)
			captured, start = thisCaptured, pos
		}
		if pos == inLength-1 {
			cb(start, inLength, captured)
		}
	}
}
