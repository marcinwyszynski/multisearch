package search_tree

import (
	"testing"
)

func TestNodeAdd(t *testing.T) {
	var err error
	head := newNode("", 0)
	err = head.add([]string{"one", "two"})
	if err != nil {
		t.Fatalf("node add err %v, expected nil", err)
	}
	err = head.add([]string{"one", "two"})
	expectedErr := "duplicate entry: \"one two\""
	if err == nil || err.Error() != expectedErr {
		t.Fatalf("node add err %v, expected %q", err, expectedErr)
	}
}
