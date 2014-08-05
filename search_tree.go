package search_tree

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	splitter = regexp.MustCompile("[^\\w]+")
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

func (n *node) walk(callback func(*node) bool) {
	if !callback(n) {
		return
	}
	for _, child := range n.children {
		child.walk(callback)
	}
}

func (n *node) debugString() string {
	retVal := ""
	for i := 0; i < n.depth; i++ {
		retVal += "  "
	}
	retVal += n.String()
	if retVal == "" {
		retVal += "[root]"
	}
	if n.terminal {
		retVal += " [t]"
	}
	retVal += "\n"
	for _, child := range n.children {
		retVal += child.debugString()
	}
	return retVal
}

func (n *node) String() string {
	if n.parent == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", n.parent.String(), n.content))
}

// SearchTree is a data structure capable of finding the ratio of needles
// (added through the Add) method to the size of the haystack (passed via the
// Compute method.
type SearchTree interface {
	// Add adds a search term to the search tree.
	Add(term string) error

	// Compute computes the ratio of matched expressions to the total length
	// of the text.
	Compute(text string) float64
}

// searchTreeImpl is an implementation of the SearchTree interface.
type searchTreeImpl struct {
	head      *node
	stopwords map[string]struct{}
}

func NewSearchTree() *searchTreeImpl {
	return &searchTreeImpl{
		head:      newNode("", 0),
		stopwords: make(map[string]struct{}),
	}
}

func (s *searchTreeImpl) Add(term string) error {
	sanitized := s.sanitize(term)
	if len(sanitized) == 0 {
		return fmt.Errorf("term only consists of stopwords: %q", term)
	}
	s.head.add(sanitized)
	return nil
}

// compute computes the ratio of matched expressions to the total length of the
// text. It is assumed that both expressions in the tree and the text being
// analyzed are sanitized (downcase, no stopwords).
func (s *searchTreeImpl) Compute(text string) float64 {
	sanitizedText := s.sanitize(text)
	if len(sanitizedText) == 0 {
		return 0.0 // As to avoid division by 0.
	}
	matchedWords := 0
	cursors := make(map[*node]struct{})
	cursors[s.head] = struct{}{}
	for _, word := range sanitizedText {
		for cursor, _ := range cursors {
			nextCursor, exists := cursor.children[word]
			if exists && nextCursor.terminal {
				matchedWords += nextCursor.depth
			}
			if exists && !nextCursor.terminal {
				cursors[nextCursor] = struct{}{}
			}
			if cursor != s.head {
				delete(cursors, cursor)
			}
		}
	}
	return float64(matchedWords) / float64(len(sanitizedText))
}

func (s *searchTreeImpl) addStopword(stopword string) error {
	if len(splitter.Split(stopword, -1)) > 1 {
		return fmt.Errorf("stopword not a single word: %q", stopword)
	}
	stopword = strings.ToLower(stopword)
	if _, exists := s.stopwords[stopword]; exists {
		return fmt.Errorf("duplicate stopword: %q", stopword)
	}
	s.stopwords[stopword] = struct{}{}
	return nil
}

func (s *searchTreeImpl) sanitize(text string) []string {
	retVal := make([]string, 0)
	for _, token := range splitter.Split(text, -1) {
		word := strings.ToLower(token)
		if _, isStopword := s.stopwords[word]; !isStopword {
			retVal = append(retVal, word)
		}
	}
	return retVal
}
