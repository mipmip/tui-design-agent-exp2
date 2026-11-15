package brief

import (
	"fmt"
	"strings"

	"github.com/mipmip/tuimock/internal/design"
	"github.com/mipmip/tuimock/internal/project"
)

const Instruction = "Design the screen described in brief.md. Submit one mockup per requested variant, using only the style letters, border characters and symbols of this design system. Choose each size from the allowed sizes. Cite the reference each variant borrows from in its rationale. When a submission is rejected, fix every listed error and submit again."

type Example struct {
	Cols      int      `json:"cols"`
	Rows      int      `json:"rows"`
	Text      []string `json:"text"`
	Style     []string `json:"style"`
	Rationale string   `json:"rationale"`
}

type Format struct {
	Shape   string   `json:"shape"`
	Rules   []string `json:"rules"`
	Example Example  `json:"example"`
}

type Payload struct {
	Screen      string                     `json:"screen"`
	Sizes       []string                   `json:"sizes"`
	Palette     map[string]string          `json:"palette"`
	Styles      map[string]design.Resolved `json:"styles"`
	Borders     []string                   `json:"borders"`
	Symbols     []string                   `json:"symbols"`
	Rules       []string                   `json:"rules"`
	Brief       string                     `json:"brief"`
	References  []project.Ref              `json:"references"`
	Format      Format                     `json:"format"`
	Instruction string                     `json:"instruction"`
}

func Build(p *project.Project) Payload {
	d := p.Design
	sizes := make([]string, 0, len(d.Sizes))
	for _, s := range d.Sizes {
		sizes = append(sizes, fmt.Sprintf("%dx%d", s.Cols, s.Rows))
	}
	styles := make(map[string]design.Resolved, len(d.Styles))
	for _, letter := range d.StyleLetters() {
		resolved, _ := d.Resolve(letter)
		styles[letter] = resolved
	}
	return Payload{
		Screen:     d.Screen,
		Sizes:      sizes,
		Palette:    d.Palette,
		Styles:     styles,
		Borders:    d.Borders,
		Symbols:    d.Symbols,
		Rules:      d.Rules,
		Brief:      p.Brief,
		References: p.Refs,
		Format: Format{
			Shape:   "text and style are parallel arrays with exactly rows entries. Each text row is exactly cols display cells wide. Each style row has exactly one style letter per cell. Only printable ASCII, the configured border characters and the configured symbols are allowed; wide and zero-width characters are rejected.",
			Rules:   d.Rules,
			Example: example(d),
		},
		Instruction: Instruction,
	}
}

func example(d *design.Design) Example {
	size := d.SmallestSize()
	letters := d.StyleLetters()
	letter := "n"
	if len(letters) > 0 {
		letter = letters[0]
	}
	tl, tr, bl, br, h, v, ok := d.FrameRunes()
	if !ok {
		tl, tr, bl, br, h, v = '+', '+', '+', '+', '-', '|'
	}
	text := make([]string, size.Rows)
	style := make([]string, size.Rows)
	for row := 0; row < size.Rows; row++ {
		var line []rune
		switch {
		case row == 0:
			line = append(line, tl)
			line = append(line, repeat(h, size.Cols-2)...)
			line = append(line, tr)
		case row == size.Rows-1:
			line = append(line, bl)
			line = append(line, repeat(h, size.Cols-2)...)
			line = append(line, br)
		default:
			line = append(line, v)
			for i := 0; i < size.Cols-2; i++ {
				line = append(line, ' ')
			}
			line = append(line, v)
		}
		text[row] = string(line)
		style[row] = strings.Repeat(letter, size.Cols)
	}
	label := " tuimock example "
	if size.Cols-6 >= len(label) {
		row := []rune(text[1])
		for i, r := range label {
			row[2+i] = r
		}
		text[1] = string(row)
	}
	return Example{
		Cols:      size.Cols,
		Rows:      size.Rows,
		Text:      text,
		Style:     style,
		Rationale: "Shape example only. Submit variants at the sizes listed for this screen.",
	}
}

func repeat(r rune, n int) []rune {
	if n <= 0 {
		return nil
	}
	out := make([]rune, n)
	for i := range out {
		out[i] = r
	}
	return out
}
