package multisearch

import (
	"testing"
	"unicode"
)

type FakeStemmer struct{}

func (f *FakeStemmer) StemString(input string) string {
	return input
}

func TestTokenize(t *testing.T) {
	input := " pełnoziarnista mąka ryżowa "
	tokenizer := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsNumber(r)
	}
	type capture struct {
		content  string
		captured bool
	}
	captures := make([]capture, 0)
	onBoundary := func(start, end int, captured bool) {
		captures = append(captures, capture{input[start:end], captured})
	}
	tokenize(input, tokenizer, onBoundary)
	for _, i := range []int{0, 2, 4, 6} {
		expected, actual := " ", captures[i].content
		if expected != actual {
			t.Fatalf("captures[%d].content = %q, expected %q", i, expected, actual)
		}
		if captures[i].captured {
			t.Fatalf("captures[%d].captured = true, expected false", i)
		}
	}
	expected, actual := "pełnoziarnista", captures[1].content
	if expected != actual {
		t.Fatalf("captures[1].content = %q, expected %q", expected, actual)
	}
}

func TestIgnore(t *testing.T) {
	eng := NewEngine(&FakeStemmer{})
	ignore := "mąka"
	if err := eng.Ignore(ignore); err != nil {
		t.Fatalf("eng.Ignore(%q) err = %v, expected nil", ignore, err)
	}
}

func TestEngine(t *testing.T) {
	eng := NewEngine(&FakeStemmer{})
	needle := "pełnoziarnista ryżowa"
	if _, err := eng.Match(needle); err != nil {
		t.Fatalf("eng.Add(%q) err = %v, expected nil", needle, err)
	}
	needle = "ryżowa jak"
	if _, err := eng.Match(needle); err != nil {
		t.Fatalf("eng.Add(%q) err = %v, expected nil", needle, err)
	}
	ignore := "mąka"
	if err := eng.Ignore(ignore); err != nil {
		t.Fatalf("eng.Ignore(%q) err = %v, expected nil", ignore, err)
	}
	haystack := "Pełnoziarnista mąka ryżowa – jak zrobić?"
	tokens := eng.Process(haystack)
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
		}
	}
}
