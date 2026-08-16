package write

import (
	"bufio"
	"fmt"
	"io"

	"ini-merge/internal/ini"
)

// File writes the INI document. The writer is flushed before return.
func File(w io.Writer, f *ini.File) (err error) {
	bw := bufio.NewWriter(w)
	defer func() {
		if ferr := bw.Flush(); ferr != nil && err == nil {
			err = ferr
		}
	}()
	for i, s := range f.Sections {
		if i > 0 {
			if _, err := fmt.Fprintln(bw); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(bw, "[%s]\n", s.Name); err != nil {
			return err
		}
		for _, p := range s.Keys {
			if _, err := fmt.Fprintf(bw, "%s=%s\n", p.Key, p.Value); err != nil {
				return err
			}
		}
	}
	return nil
}
