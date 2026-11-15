package brief

import (
	"strings"
	"testing"

	"github.com/mipmip/tuimock/internal/design"
	"github.com/mipmip/tuimock/internal/mockup"
	"github.com/mipmip/tuimock/internal/project"
)

const designYAML = `
screen: beandex
sizes:
  - {cols: 30, rows: 5}
palette:
  bg: "#1a1b26"
  fg: "#c0caf5"
  accent: "#7aa2f7"
styles:
  n: {fg: fg}
  a: {fg: accent, bold: true}
borders: [rounded]
symbols: ["▸"]
rules: ["one accent"]
`

func testProject(t *testing.T) *project.Project {
	t.Helper()
	d, err := design.Parse([]byte(designYAML))
	if err != nil {
		t.Fatal(err)
	}
	return &project.Project{
		Dir:    t.TempDir(),
		Design: d,
		Brief:  "Design the beandex list view.",
		Refs: []project.Ref{
			{Name: "lazygit", Text: "capture", Note: "dense but calm"},
		},
	}
}

func TestBuild(t *testing.T) {
	p := testProject(t)
	payload := Build(p)
	if payload.Screen != "beandex" || len(payload.Sizes) != 1 || payload.Sizes[0] != "30x5" {
		t.Fatalf("payload = %+v", payload)
	}
	if payload.Brief != p.Brief || len(payload.References) != 1 || payload.References[0].Note != "dense but calm" {
		t.Fatalf("brief/refs = %+v", payload)
	}
	if payload.Styles["a"].Color != "#7aa2f7" || !payload.Styles["a"].Bold {
		t.Fatalf("styles = %+v", payload.Styles)
	}
	if payload.Instruction == "" || payload.Format.Shape == "" || len(payload.Format.Rules) != 1 {
		t.Fatalf("format = %+v", payload.Format)
	}
}

func TestExampleIsValid(t *testing.T) {
	p := testProject(t)
	payload := Build(p)
	ex := payload.Format.Example
	if ex.Cols != 30 || ex.Rows != 5 {
		t.Fatalf("example size = %dx%d", ex.Cols, ex.Rows)
	}
	if len(ex.Text) != ex.Rows || len(ex.Style) != ex.Rows {
		t.Fatalf("example rows = %d %d", len(ex.Text), len(ex.Style))
	}
	if !strings.Contains(strings.Join(ex.Text, "\n"), " tuimock example ") {
		t.Fatal("example label missing")
	}
	errs := mockup.Validate(p.Design, mockup.Submission{
		Screen:    p.Design.Screen,
		Variant:   "example",
		Rationale: ex.Rationale,
		Grid:      mockup.Grid{Cols: ex.Cols, Rows: ex.Rows, Text: ex.Text, Style: ex.Style},
	})
	if len(errs) != 0 {
		t.Fatalf("example invalid: %v", errs)
	}
}

func TestExampleWithoutFrame(t *testing.T) {
	d, err := design.Parse([]byte(`
screen: test
sizes: [{cols: 10, rows: 3}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [rounded]
`))
	if err != nil {
		t.Fatal(err)
	}
	ex := Build(&project.Project{Design: d}).Format.Example
	if len(ex.Text) != 3 || len([]rune(ex.Text[0])) != 10 {
		t.Fatalf("example = %+v", ex)
	}
}
