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
	if _, err := eng.Match(needle, 1); err != nil {
		t.Fatalf("eng.Add(%q) err = %v, expected nil", needle, err)
	}
	needle = "ryżowa jak"
	if _, err := eng.Match(needle, 2); err != nil {
		t.Fatalf("eng.Add(%q) err = %v, expected nil", needle, err)
	}
	ignore := "mąka"
	if err := eng.Ignore(ignore); err != nil {
		t.Fatalf("eng.Ignore(%q) err = %v, expected nil", ignore, err)
	}
	haystack := "Pełnoziarnista mąka ryżowa – jak zrobić?"
	tokens := make([]Token, 0)
	for this := eng.Process(haystack); this != nil; this = this.Next() {
		tokens = append(tokens, this)
	}
	if len(tokens) != 10 {
		t.Errorf("len(tokens) = %d, expected 10", len(tokens))
	}
	for _, k := range tokens {
		if k.String() == "Pełnoziarnista" || k.String() == "jak" {
			if len(k.Matches()) != 1 {
				msg := "len(k.Matches()) = %d for %q, expected 1"
				t.Errorf(msg, len(k.Matches()), k.String())
			}
		}
		if k.String() == "ryżowa" {
			if len(k.Matches()) != 2 {
				msg := "len(k.Matches()) = %d for %q, expected 2"
				t.Errorf(msg, len(k.Matches()), k.String())
			}
			expTopMatch := "ryżowa jak"
			actual := k.TopMatch().String()
			if actual != expTopMatch {
				t.Errorf("topMatch weight: %d", k.TopMatch().Weight())
				msg := "k.topMatch for %q = %q, expected %q"
				t.Errorf(msg, k.String(), actual, expTopMatch)
			}
		}
	}
}
