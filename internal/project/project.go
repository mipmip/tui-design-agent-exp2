package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mipmip/tuimock/internal/design"
)

type Ref struct {
	Name string `json:"name"`
	Text string `json:"text"`
	Note string `json:"note"`
}

type Project struct {
	Dir    string
	Design *design.Design
	Brief  string
	Refs   []Ref
}

func Load(dir string) (*Project, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("project %s: %w", dir, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("project %s: folder does not exist", abs)
		}
		return nil, fmt.Errorf("project %s: %w", abs, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("project %s: not a directory", abs)
	}

	designPath := filepath.Join(abs, "design.yaml")
	d, err := design.Load(designPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("project %s: design.yaml is missing", abs)
		}
		return nil, fmt.Errorf("project %s: %w", abs, err)
	}

	briefPath := filepath.Join(abs, "brief.md")
	brief, err := os.ReadFile(briefPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("project %s: brief.md is missing", abs)
		}
		return nil, fmt.Errorf("project %s: read brief: %w", abs, err)
	}

	refs, err := loadRefs(filepath.Join(abs, "refs"))
	if err != nil {
		return nil, fmt.Errorf("project %s: %w", abs, err)
	}

	return &Project{Dir: abs, Design: d, Brief: string(brief), Refs: refs}, nil
}

func (p *Project) RefNames() []string {
	names := make([]string, 0, len(p.Refs))
	for _, ref := range p.Refs {
		names = append(names, ref.Name)
	}
	return names
}

func loadRefs(dir string) ([]Ref, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read refs: %w", err)
	}
	var refs []Ref
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}
		text, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read refs/%s: %w", entry.Name(), err)
		}
		name := strings.TrimSuffix(entry.Name(), ".txt")
		note, err := os.ReadFile(filepath.Join(dir, name+".md"))
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("read refs/%s.md: %w", name, err)
		}
		refs = append(refs, Ref{Name: name, Text: string(text), Note: string(note)})
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Name < refs[j].Name })
	return refs, nil
}
