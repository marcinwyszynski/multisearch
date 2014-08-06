package search_tree

import (
	"fmt"
	"strings"
)

// node represents a node in the text tree.
type node struct {
	content  string
	depth    int
	children map[string]*node
	parent   *node
	terminal bool
}

// newNode creates a new node.
func newNode(content string, depth int) *node {
	return &node{
		content:  content,
		depth:    depth,
		children: make(map[string]*node),
		terminal: false,
	}
}

// add places the text in the tree, returning the number of operations
// performed.
func (n *node) add(text []string) *node {
	if len(text) == 0 {
		n.terminal = true
		return n
	}
	if potentialChild, exists := n.children[text[0]]; exists {
		return potentialChild.add(text[1:])
	}
	newChild := newNode(text[0], n.depth+1)
	n.children[text[0]] = newChild
	newChild.parent = n
	return newChild.add(text[1:])
}

// String reconstructs the full content of the node.
func (n *node) String() string {
	if n.parent == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", n.parent.String(), n.content))
}
