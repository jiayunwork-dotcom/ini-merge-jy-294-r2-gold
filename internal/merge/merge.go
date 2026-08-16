package merge

import "ini-merge/internal/ini"

// Overlay copies base then applies over: matching section keys are replaced,
// new sections and keys are appended. Neither input is mutated.
func Overlay(base, over *ini.File) *ini.File {
	out := &ini.File{Sections: make([]ini.Section, 0, len(base.Sections)+len(over.Sections))}
	idx := map[string]int{}
	for _, s := range base.Sections {
		cp := copySection(s)
		idx[cp.Name] = len(out.Sections)
		out.Sections = append(out.Sections, cp)
	}
	for _, s := range over.Sections {
		i, ok := idx[s.Name]
		if !ok {
			cp := copySection(s)
			idx[cp.Name] = len(out.Sections)
			out.Sections = append(out.Sections, cp)
			continue
		}
		dst := &out.Sections[i]
		kidx := map[string]int{}
		for j, p := range dst.Keys {
			kidx[p.Key] = j
		}
		for _, p := range s.Keys {
			if j, hit := kidx[p.Key]; hit {
				dst.Keys[j].Value = p.Value
				continue
			}
			kidx[p.Key] = len(dst.Keys)
			dst.Keys = append(dst.Keys, p)
		}
	}
	return out
}

func copySection(s ini.Section) ini.Section {
	keys := append([]ini.Pair(nil), s.Keys...)
	return ini.Section{Name: s.Name, Keys: keys}
}
