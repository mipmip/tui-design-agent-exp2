package mockup

import (
	"strings"
	"testing"

	"github.com/mipmip/tuimock/internal/design"
)

const designYAML = `
screen: test
sizes:
  - {cols: 4, rows: 2}
palette:
  bg: "#000000"
  fg: "#ffffff"
  accent: "#7aa2f7"
styles:
  n: {fg: fg}
  a: {fg: accent}
borders: [rounded]
symbols: ["▸"]
`

func testDesign(t *testing.T) *design.Design {
	t.Helper()
	d, err := design.Parse([]byte(designYAML))
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func valid() Submission {
	return Submission{
		Screen:    "test",
		Variant:   "a",
		Rationale: "borrows from lazygit",
		Grid: Grid{
			Cols:  4,
			Rows:  2,
			Text:  []string{"╭──╮", "╰──╯"},
			Style: []string{"nnnn", "aann"},
		},
	}
}

func TestValidateOK(t *testing.T) {
	if errs := Validate(testDesign(t), valid()); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestValidateWidth(t *testing.T) {
	s := valid()
	s.Text[0] = "╭───╮"
	s.Style[0] = "nnnnn"
	errs := Validate(testDesign(t), s)
	joined := strings.Join(errs, "\n")
	for _, want := range []string{
		"row 1: text is 5 cells wide, expected 4",
		"row 1: style has 5 letters, expected 4",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q in:\n%s", want, joined)
		}
	}
}

func TestValidateWideCharacter(t *testing.T) {
	s := valid()
	s.Text[0] = "╭🙂─╮"
	errs := Validate(testDesign(t), s)
	joined := strings.Join(errs, "\n")
	if !strings.Contains(joined, "is 2 cells wide; wide characters are not allowed") {
		t.Fatalf("got:\n%s", joined)
	}
	if strings.Contains(joined, "cells wide, expected") {
		t.Fatalf("width error should be suppressed when a character is invalid:\n%s", joined)
	}
}

func TestValidateDisallowedCharacter(t *testing.T) {
	s := valid()
	s.Text[0] = "╭☃─╮"
	errs := Validate(testDesign(t), s)
	joined := strings.Join(errs, "\n")
	if !strings.Contains(joined, "is not allowed") {
		t.Fatalf("got:\n%s", joined)
	}
}

func TestValidateCJKCharacter(t *testing.T) {
	s := valid()
	s.Text[0] = "╭漢─╮"
	errs := Validate(testDesign(t), s)
	joined := strings.Join(errs, "\n")
	if !strings.Contains(joined, "is 2 cells wide; wide characters are not allowed") {
		t.Fatalf("got:\n%s", joined)
	}
}

func TestValidateBeandexSizes(t *testing.T) {
	d, err := design.Parse([]byte(`
screen: list-view
sizes:
  - {cols: 100, rows: 30}
  - {cols: 80, rows: 24}
palette: {fg: "#c0caf5"}
styles: {n: {fg: fg}}
borders: [rounded]
`))
	if err != nil {
		t.Fatal(err)
	}
	for _, size := range d.Sizes {
		s := Submission{
			Screen:    "list-view",
			Variant:   "flat",
			Rationale: "grouped table",
			Grid:      Grid{Cols: size.Cols, Rows: size.Rows},
		}
		for row := 0; row < size.Rows; row++ {
			text := strings.Repeat(" ", size.Cols)
			switch row {
			case 0:
				text = "╭" + strings.Repeat("─", size.Cols-2) + "╮"
			case size.Rows - 1:
				text = "╰" + strings.Repeat("─", size.Cols-2) + "╯"
			default:
				text = "│" + text[1:size.Cols-1] + "│"
			}
			s.Text = append(s.Text, text)
			s.Style = append(s.Style, strings.Repeat("n", size.Cols))
		}
		if errs := Validate(d, s); len(errs) != 0 {
			t.Fatalf("%dx%d: %v", size.Cols, size.Rows, errs)
		}
	}
}

func TestValidateZeroWidthCharacter(t *testing.T) {
	s := valid()
	s.Text[0] = "╭─\u0301─╮"
	errs := Validate(testDesign(t), s)
	if !strings.Contains(strings.Join(errs, "\n"), "zero width") {
		t.Fatalf("got:\n%v", errs)
	}
}

func TestValidateUnknownStyle(t *testing.T) {
	s := valid()
	s.Style[0] = "nnnx"
	errs := Validate(testDesign(t), s)
	if !strings.Contains(strings.Join(errs, "\n"), `row 1, column 4: style letter "x" is not defined in the design system`) {
		t.Fatalf("got:\n%v", errs)
	}
}

func TestValidateSizeAndRows(t *testing.T) {
	s := valid()
	s.Cols = 5
	s.Rows = 3
	s.Text = []string{"╭───╮", "╰───╯"}
	s.Style = []string{"nnnnn", "aaaaa"}
	errs := Validate(testDesign(t), s)
	joined := strings.Join(errs, "\n")
	for _, want := range []string{
		"size 5x3 is not allowed; allowed sizes: 4x2",
		"text: got 2 rows, expected 3",
		"style: got 2 rows, expected 3",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q in:\n%s", want, joined)
		}
	}
}

func TestValidateIdentity(t *testing.T) {
	s := valid()
	s.Screen = ""
	s.Variant = ""
	errs := Validate(testDesign(t), s)
	joined := strings.Join(errs, "\n")
	if !strings.Contains(joined, "screen must not be empty") || !strings.Contains(joined, "variant must not be empty") {
		t.Fatalf("got:\n%s", joined)
	}
}

func TestValidateInvalidSize(t *testing.T) {
	s := valid()
	s.Cols = 0
	errs := Validate(testDesign(t), s)
	if len(errs) != 1 || !strings.Contains(errs[0], "cols and rows must be positive") {
		t.Fatalf("got:\n%v", errs)
	}
}
