package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const designYAML = `
screen: test
sizes:
  - {cols: 4, rows: 2}
palette:
  fg: "#ffffff"
styles:
  n: {fg: fg}
borders: [rounded]
`

func projectDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "design.yaml"), []byte(designYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "brief.md"), []byte("brief"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func invoke(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestVersion(t *testing.T) {
	code, stdout, _ := invoke("version")
	if code != 0 || stdout == "" {
		t.Fatalf("code=%d stdout=%q", code, stdout)
	}
}

func TestUsage(t *testing.T) {
	code, _, stderr := invoke()
	if code != 2 || !strings.Contains(stderr, "tuimock <command>") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
	code, stdout, _ := invoke("help")
	if code != 0 || !strings.Contains(stdout, "serve") {
		t.Fatalf("code=%d stdout=%q", code, stdout)
	}
	code, _, stderr = invoke("nope")
	if code != 2 || !strings.Contains(stderr, "unknown subcommand") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestValidateProject(t *testing.T) {
	dir := projectDir(t)
	code, stdout, stderr := invoke("validate", "--dir", dir)
	if code != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
	for _, want := range []string{"ok: screen test", "sizes 4x2", "0 references", "brief 5 bytes"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("missing %q in %q", want, stdout)
		}
	}
}

func TestValidateMissingProject(t *testing.T) {
	code, _, stderr := invoke("validate", "--dir", filepath.Join(t.TempDir(), "missing"))
	if code != 1 || !strings.Contains(stderr, "does not exist") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestValidateEmptyBrief(t *testing.T) {
	dir := projectDir(t)
	if err := os.WriteFile(filepath.Join(dir, "brief.md"), []byte("  \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderr := invoke("validate", "--dir", dir)
	if code != 1 || !strings.Contains(stderr, "brief.md is empty") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestValidateMockup(t *testing.T) {
	dir := projectDir(t)
	good := filepath.Join(dir, "good.json")
	if err := os.WriteFile(good, []byte(`{"variant":"a","cols":4,"rows":2,"text":["╭──╮","╰──╯"],"style":["nnnn","nnnn"],"rationale":"x"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	code, stdout, stderr := invoke("validate", "--dir", dir, "--mockup", good)
	if code != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
	if !strings.Contains(stdout, "ok: mockup test variant a 4x2") {
		t.Fatalf("stdout=%q", stdout)
	}

	bad := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(bad, []byte(`{"variant":"a","cols":4,"rows":2,"text":["╭───╮","╰──╯"],"style":["nnnnn","nnnn"],"rationale":"x"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderr = invoke("validate", "--dir", dir, "--mockup", bad)
	if code != 1 || !strings.Contains(stderr, "row 1: text is 5 cells wide, expected 4") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}

func TestValidateMockupUnreadable(t *testing.T) {
	dir := projectDir(t)
	code, _, stderr := invoke("validate", "--dir", dir, "--mockup", filepath.Join(dir, "nope.json"))
	if code != 1 || !strings.Contains(stderr, "read mockup") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
	if err := os.WriteFile(filepath.Join(dir, "broken.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, _, stderr = invoke("validate", "--dir", dir, "--mockup", filepath.Join(dir, "broken.json"))
	if code != 1 || !strings.Contains(stderr, "parse mockup") {
		t.Fatalf("code=%d stderr=%q", code, stderr)
	}
}
