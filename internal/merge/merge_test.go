package merge

import (
	"strings"
	"testing"

	"ini-merge/internal/ini"
)

func mustParse(t *testing.T, s string) *ini.File {
	t.Helper()
	f, err := ini.Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestOverlayReplacesKey(t *testing.T) {
	base := mustParse(t, "[app]\nname=old\nport=80\n")
	over := mustParse(t, "[app]\nname=new\n")
	out := Overlay(base, over)
	v, ok := out.Lookup("app", "name")
	if !ok || v != "new" {
		t.Fatalf("name=%q ok=%v", v, ok)
	}
	v, ok = out.Lookup("app", "port")
	if !ok || v != "80" {
		t.Fatalf("port=%q ok=%v", v, ok)
	}
}

func TestOverlayAddsSection(t *testing.T) {
	base := mustParse(t, "[app]\nname=demo\n")
	over := mustParse(t, "[db]\nhost=localhost\n")
	out := Overlay(base, over)
	if len(out.Sections) != 2 {
		t.Fatalf("sections=%d", len(out.Sections))
	}
	v, ok := out.Lookup("db", "host")
	if !ok || v != "localhost" {
		t.Fatalf("host=%q ok=%v", v, ok)
	}
}

func TestOverlayDoesNotMutateBase(t *testing.T) {
	base := mustParse(t, "[app]\nname=old\n")
	over := mustParse(t, "[app]\nname=new\n")
	out := Overlay(base, over)
	out.Sections[0].Keys[0].Value = "CHANGED"
	v, _ := base.Lookup("app", "name")
	if v != "old" {
		t.Fatalf("overlay mutated base: %q", v)
	}
}

func TestOverlayEmptyOver(t *testing.T) {
	base := mustParse(t, "[app]\nname=demo\n")
	over := mustParse(t, "")
	out := Overlay(base, over)
	if len(out.Sections) != 1 {
		t.Fatalf("sections=%d", len(out.Sections))
	}
}
