package mockup

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/mipmip/tuimock/internal/design"
)

type Grid struct {
	Cols  int      `json:"cols" yaml:"cols"`
	Rows  int      `json:"rows" yaml:"rows"`
	Text  []string `json:"text" yaml:"text"`
	Style []string `json:"style" yaml:"style"`
}

type Submission struct {
	Screen    string `json:"screen"`
	Variant   string `json:"variant"`
	Rationale string `json:"rationale"`
	Grid
}

func Validate(d *design.Design, s Submission) []string {
	var errs []string
	if s.Screen == "" {
		errs = append(errs, "screen must not be empty")
	}
	if s.Variant == "" {
		errs = append(errs, "variant must not be empty")
	}
	if s.Cols <= 0 || s.Rows <= 0 {
		errs = append(errs, fmt.Sprintf("size %dx%d is invalid; cols and rows must be positive", s.Cols, s.Rows))
		return errs
	}
	if !d.AllowedSize(s.Cols, s.Rows) {
		errs = append(errs, fmt.Sprintf("size %dx%d is not allowed; allowed sizes: %s", s.Cols, s.Rows, d.SizesText()))
	}
	if len(s.Text) != s.Rows {
		errs = append(errs, fmt.Sprintf("text: got %d rows, expected %d", len(s.Text), s.Rows))
	}
	if len(s.Style) != s.Rows {
		errs = append(errs, fmt.Sprintf("style: got %d rows, expected %d", len(s.Style), s.Rows))
	}
	for i := 0; i < s.Rows; i++ {
		if i < len(s.Text) {
			errs = append(errs, validateTextRow(d, i+1, s.Cols, s.Text[i])...)
		}
		if i < len(s.Style) {
			errs = append(errs, validateStyleRow(d, i+1, s.Cols, s.Style[i])...)
		}
	}
	return errs
}

func validateTextRow(d *design.Design, row, cols int, text string) []string {
	var errs []string
	invalid := false
	col := 0
	for _, r := range text {
		col++
		width := runewidth.RuneWidth(r)
		if width == 0 {
			invalid = true
			errs = append(errs, fmt.Sprintf("row %d, column %d: character %q (U+%04X) is zero width; zero-width characters are not allowed", row, col, string(r), r))
			continue
		}
		if width != 1 {
			invalid = true
			errs = append(errs, fmt.Sprintf("row %d, column %d: character %q (U+%04X) is %d cells wide; wide characters are not allowed", row, col, string(r), r, width))
			continue
		}
		if !d.AllowedRune(r) {
			invalid = true
			errs = append(errs, fmt.Sprintf("row %d, column %d: character %q (U+%04X) is not allowed", row, col, string(r), r))
		}
	}
	if !invalid {
		if got := runewidth.StringWidth(text); got != cols {
			errs = append(errs, fmt.Sprintf("row %d: text is %d cells wide, expected %d", row, got, cols))
		}
	}
	return errs
}

func validateStyleRow(d *design.Design, row, cols int, style string) []string {
	var errs []string
	letters := []rune(style)
	if len(letters) != cols {
		errs = append(errs, fmt.Sprintf("row %d: style has %d letters, expected %d", row, len(letters), cols))
	}
	for i, letter := range letters {
		if !d.HasStyle(string(letter)) {
			errs = append(errs, fmt.Sprintf("row %d, column %d: style letter %q is not defined in the design system", row, i+1, string(letter)))
		}
	}
	return errs
}
