package search_tree

import (
	"testing"
)

// Check if adding a node works.
func TestNodeAdd(t *testing.T) {
	head := newNode("", 0)
	node := head.add([]string{"one", "two"})
	if node == nil {
		t.Fatal("node add = nil")
	}
}

// Check if walking the tree bredth-first works. While at it, check string
// serialization just as well.
func TestNodeWalkAndString(t *testing.T) {
	head := newNode("", 0)
	head.add([]string{"one", "two", "three"})
	nodes := make([]*node, 0)
	head.walk(func(n *node) bool {
		if n.String() == "one two three" {
			return false
		}
		nodes = append(nodes, n)
		return true
	})
	if len(nodes) != 3 {
		t.Fatalf("len(nodes) = %d, expected 3", len(nodes))
	}
	for i, expString := range []string{"", "one", "one two"} {
		if nodes[i].depth != i {
			t.Errorf("nodes[%d].depth = %d, expected %d", nodes[i].depth, i)
		}
		if nodes[i].String() != expString {
			t.Errorf("nodes[%d].String() = %q, expected %q", i, nodes[i].String(), expString)
		}
	}
}
