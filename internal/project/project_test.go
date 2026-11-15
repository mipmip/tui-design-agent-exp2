package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const designYAML = `
screen: beandex
sizes:
  - {cols: 10, rows: 3}
palette:
  bg: "#1a1b26"
  fg: "#c0caf5"
styles:
  n: {fg: fg}
borders: [rounded]
`

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "design.yaml"), designYAML)
	write(t, filepath.Join(dir, "brief.md"), "# brief\n")
	write(t, filepath.Join(dir, "refs", "lazygit.txt"), "one\ntwo\n")
	write(t, filepath.Join(dir, "refs", "lazygit.md"), "dense but calm")
	write(t, filepath.Join(dir, "refs", "btop.txt"), "three\n")
	write(t, filepath.Join(dir, "refs", "notes.md"), "ignored")

	p, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Design.Screen != "beandex" {
		t.Fatalf("screen = %q", p.Design.Screen)
	}
	if p.Brief != "# brief\n" {
		t.Fatalf("brief = %q", p.Brief)
	}
	if got := p.RefNames(); len(got) != 2 || got[0] != "btop" || got[1] != "lazygit" {
		t.Fatalf("refs = %v", got)
	}
	if p.Refs[1].Note != "dense but calm" {
		t.Fatalf("note = %q", p.Refs[1].Note)
	}
	if p.Refs[0].Note != "" {
		t.Fatalf("btop note = %q", p.Refs[0].Note)
	}
	if !filepath.IsAbs(p.Dir) {
		t.Fatalf("dir not absolute: %s", p.Dir)
	}
}

func TestLoadMissingFiles(t *testing.T) {
	dir := t.TempDir()
	if _, err := Load(dir); err == nil || !strings.Contains(err.Error(), "design.yaml is missing") {
		t.Fatalf("err = %v", err)
	}
	write(t, filepath.Join(dir, "design.yaml"), designYAML)
	if _, err := Load(dir); err == nil || !strings.Contains(err.Error(), "brief.md is missing") {
		t.Fatalf("err = %v", err)
	}
}

func TestLoadNoRefs(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "design.yaml"), designYAML)
	write(t, filepath.Join(dir, "brief.md"), "brief")
	p, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Refs) != 0 {
		t.Fatalf("refs = %v", p.Refs)
	}
}

func TestLoadInvalidDesign(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "design.yaml"), "screen: test\n")
	write(t, filepath.Join(dir, "brief.md"), "brief")
	_, err := Load(dir)
	if err == nil || !strings.Contains(err.Error(), "design.yaml") {
		t.Fatalf("err = %v", err)
	}
}

func TestLoadNotADirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file")
	write(t, path, "x")
	if _, err := Load(path); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("err = %v", err)
	}
}
