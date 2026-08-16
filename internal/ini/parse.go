package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// File is a parsed INI document. Section order is insertion order.
type File struct {
	Sections []Section
}

// Section holds a name and key/value pairs.
type Section struct {
	Name string
	Keys []Pair
}

// Pair is one key=value line.
type Pair struct {
	Key   string
	Value string
}

// Parse reads an INI file. UTF-8 BOM is stripped. Lines starting with
// '#' or ';' are comments. A line that is not [section] or key=value is an error.
func Parse(r io.Reader) (*File, error) {
	br := bufio.NewReader(r)
	head, err := br.Peek(3)
	if err == nil && bytes.Equal(head, []byte{0xEF, 0xBB, 0xBF}) {
		if _, err := br.Discard(3); err != nil {
			return nil, err
		}
	}
	sc := bufio.NewScanner(br)
	out := &File{Sections: []Section{}}
	var cur *Section
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			name := strings.TrimSpace(line[1 : len(line)-1])
			if name == "" {
				return nil, fmt.Errorf("ini: line %d: empty section name", lineNo)
			}
			out.Sections = append(out.Sections, Section{Name: name, Keys: []Pair{}})
			cur = &out.Sections[len(out.Sections)-1]
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			return nil, fmt.Errorf("ini: line %d: expected key=value, got %q", lineNo, line)
		}
		if cur == nil {
			return nil, fmt.Errorf("ini: line %d: key outside section", lineNo)
		}
		k := strings.TrimSpace(line[:eq])
		v := strings.TrimSpace(line[eq+1:])
		if k == "" {
			return nil, fmt.Errorf("ini: line %d: empty key", lineNo)
		}
		cur.Keys = append(cur.Keys, Pair{Key: k, Value: v})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Lookup returns the last value of key in section, or ("", false).
func (f *File) Lookup(section, key string) (string, bool) {
	for i := range f.Sections {
		if f.Sections[i].Name != section {
			continue
		}
		for j := len(f.Sections[i].Keys) - 1; j >= 0; j-- {
			if f.Sections[i].Keys[j].Key == key {
				return f.Sections[i].Keys[j].Value, true
			}
		}
	}
	return "", false
}
