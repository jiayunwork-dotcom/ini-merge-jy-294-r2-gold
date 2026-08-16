// Command ini-merge overlays one INI file onto another and writes the result.
package main

import (
	"flag"
	"fmt"
	"os"

	"ini-merge/internal/ini"
	"ini-merge/internal/merge"
	"ini-merge/internal/write"
)

func main() {
	basePath := flag.String("base", "", "base INI path (required)")
	overPath := flag.String("over", "", "overlay INI path (required)")
	out := flag.String("out", "-", "output path, or - for stdout")
	flag.Parse()
	if *basePath == "" || *overPath == "" {
		fatal("missing required -base and -over")
	}
	base, err := load(*basePath)
	if err != nil {
		fatal("base: %v", err)
	}
	over, err := load(*overPath)
	if err != nil {
		fatal("over: %v", err)
	}
	merged := merge.Overlay(base, over)
	w := os.Stdout
	if *out != "-" {
		f, err := os.Create(*out)
		if err != nil {
			fatal("create %s: %v", *out, err)
		}
		defer f.Close()
		w = f
	}
	if err := write.File(w, merged); err != nil {
		fatal("write: %v", err)
	}
}

func load(path string) (*ini.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ini.Parse(f)
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "ini-merge: "+format+"\n", a...)
	os.Exit(1)
}
