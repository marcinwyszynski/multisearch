package multisearch

import (
	"fmt"
	"strings"
)

// node represents a node in the text tree.
type node struct {
	content       string
	depth, weight int
	children      map[string]*node
	parent        *node
	terminal      bool
}

// add places the text in the tree, returning the number of operations
// performed.
func (n *node) add(text []string, weight int) *node {
	if len(text) == 0 {
		n.terminal = true
		return n
	}
	if potentialChild, exists := n.children[text[0]]; exists {
		return potentialChild.add(text[1:], weight)
	}
	newChild := newNode(text[0], n.depth+1, weight)
	n.children[text[0]] = newChild
	newChild.parent = n
	return newChild.add(text[1:], weight)
}

func (n *node) Size() int {
	return n.depth
}

func (n *node) String() string {
	if n.parent == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", n.parent.String(), n.content))
}

func (n *node) Weight() int {
	return n.weight
}

// newNode creates a new node.
func newNode(content string, depth, weight int) *node {
	return &node{
		content:  content,
		depth:    depth,
		weight:   weight,
		children: make(map[string]*node),
		terminal: false,
	}
}

// byWeight provides an auxilary data structure for sorting nodes in
// *decreasing* order by their weight first and (then) their depth.
type byWeight []*node

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
