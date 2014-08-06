package multisearch

import (
	"testing"
)

type FakeStemmer struct{}

func (f *FakeStemmer) StemString(input string) string {
	return input
}

func TestEngine(t *testing.T) {
	eng := NewEngine(&FakeStemmer{})
	needle := "pełnoziarnista ryżowa"
	if err := eng.Add(needle, 1); err != nil {
		t.Fatalf("eng.Add(%q) err = %v, expected nil", needle, err)
	}
	needle = "ryżowa jak"
	if err := eng.Add(needle, 2); err != nil {
		t.Fatalf("eng.Add(%q) err = %v, expected nil", needle, err)
	}
	ignore := "mąka"
	if err := eng.Ignore(ignore); err != nil {
		t.Fatalf("eng.Ignore(%q) err = %v, expected nil", ignore, err)
	}
	haystack := "Pełnoziarnista mąka ryżowa – jak zrobić?"
	tokens := make([]*token, 0)
	for this := eng.Process(haystack); this != nil; this = this.next {
		tokens = append(tokens, this)
	}
	if len(tokens) != 10 {
		t.Errorf("len(tokens) = %d, expected 10", len(tokens))
	}
	for _, k := range tokens {
		if k.content == "Pełnoziarnista" || k.content == "jak" {
			if len(k.matchedBy) != 1 {
				t.Errorf("len(k.matchedBy) = %d for %q, expected 1", len(k.matchedBy), k.content)
			}
		}
		if k.content == "ryżowa" {
			if len(k.matchedBy) != 2 {
				t.Errorf("len(k.matchedBy) = %d for %q, expected 2", len(k.matchedBy), k.content)
			}
		}
	}
}
