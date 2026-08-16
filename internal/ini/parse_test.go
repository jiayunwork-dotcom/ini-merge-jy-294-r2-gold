package ini

import (
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	in := "[app]\nname = demo\n"
	f, err := Parse(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Sections) != 1 || f.Sections[0].Name != "app" {
		t.Fatalf("sections=%v", f.Sections)
	}
	v, ok := f.Lookup("app", "name")
	if !ok || v != "demo" {
		t.Fatalf("lookup=%q ok=%v", v, ok)
	}
}

func TestParseStripsBOM(t *testing.T) {
	in := "\xef\xbb\xbf[app]\nname=demo\n"
	f, err := Parse(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Sections) == 0 || f.Sections[0].Name != "app" {
		t.Fatalf("BOM not stripped, sections=%v", f.Sections)
	}
}

func TestParseRejectsBareLine(t *testing.T) {
	if _, err := Parse(strings.NewReader("[app]\nnot-a-pair\n")); err == nil {
		t.Fatal("expected error for line without =")
	}
}

func TestParseEmptyNonNil(t *testing.T) {
	f, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if f == nil || f.Sections == nil {
		t.Fatal("empty input must return non-nil File with empty Sections")
	}
}

func TestLookupLastWins(t *testing.T) {
	in := "[app]\nname=first\nname=second\n"
	f, err := Parse(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	v, ok := f.Lookup("app", "name")
	if !ok || v != "second" {
		t.Fatalf("got %q ok=%v, want second", v, ok)
	}
}
