package multisearch

import (
	"fmt"
	"strings"
)

// matchImpl represents a matchImpl in the text tree.
type matchImpl struct {
	content  string
	depth    int
	children map[string]*matchImpl
	parent   *matchImpl
	terminal bool
}

// add places the text in the tree, returning the number of operations
// performed.
func (m *matchImpl) add(text []string) *matchImpl {
	if len(text) == 0 {
		m.terminal = true
		return m
	}
	if potentialChild, exists := m.children[text[0]]; exists {
		return potentialChild.add(text[1:])
	}
	newChild := newMatchImpl(text[0], m.depth+1)
	m.children[text[0]] = newChild
	newChild.parent = m
	return newChild.add(text[1:])
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

// newMatchImpl creates a new matchImpl.
func newMatchImpl(content string, depth int) *matchImpl {
	return &matchImpl{
		content:  content,
		depth:    depth,
		children: make(map[string]*matchImpl),
		terminal: false,
	}
}
