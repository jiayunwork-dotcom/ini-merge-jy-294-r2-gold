package write

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"ini-merge/internal/ini"
)

type failWriter struct{}

func (failWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write fail")
}

func TestFileOK(t *testing.T) {
	f := &ini.File{Sections: []ini.Section{
		{Name: "app", Keys: []ini.Pair{{Key: "name", Value: "demo"}}},
	}}
	var buf bytes.Buffer
	if err := File(&buf, f); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "[app]") || !strings.Contains(got, "name=demo") {
		t.Fatalf("missing content: %q", got)
	}
}

func TestFileFlushError(t *testing.T) {
	f := &ini.File{Sections: []ini.Section{
		{Name: "app", Keys: []ini.Pair{{Key: "n", Value: "1"}}},
	}}
	if err := File(failWriter{}, f); err == nil {
		t.Fatal("expected flush/write error to propagate")
	}
}
