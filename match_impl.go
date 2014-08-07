package multisearch

import (
	"fmt"
	"strings"
)

// matchImpl represents a matchImpl in the text tree.
type matchImpl struct {
	content       string
	depth, weight int
	children      map[string]*matchImpl
	parent        *matchImpl
	terminal      bool
}

// add places the text in the tree, returning the number of operations
// performed.
func (m *matchImpl) add(text []string, weight int) *matchImpl {
	if len(text) == 0 {
		m.terminal = true
		return m
	}
	if potentialChild, exists := m.children[text[0]]; exists {
		return potentialChild.add(text[1:], weight)
	}
	newChild := newMatchImpl(text[0], m.depth+1, weight)
	m.children[text[0]] = newChild
	newChild.parent = m
	return newChild.add(text[1:], weight)
}

func (m *matchImpl) Size() int {
	return m.depth
}

func (m *matchImpl) String() string {
	if m.parent == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", m.parent.String(), m.content))
}

func (m *matchImpl) Weight() int {
	return m.weight
}

// newMatchImpl creates a new matchImpl.
func newMatchImpl(content string, depth, weight int) *matchImpl {
	return &matchImpl{
		content:  content,
		depth:    depth,
		weight:   weight,
		children: make(map[string]*matchImpl),
		terminal: false,
	}
}

// byWeight provides an auxilary data structure for sorting matchImpls in
// *decreasing* order by their weight first and (then) their depth.
type byWeight []*matchImpl

func (b byWeight) Len() int {
	return len(b)
}

func (b byWeight) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byWeight) Less(i, j int) bool {
	if b[i].weight == b[j].weight {
		return b[i].depth > b[j].depth
	}
	return b[i].weight > b[j].weight
}
