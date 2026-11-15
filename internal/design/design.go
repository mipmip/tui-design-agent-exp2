package design

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
)

var borderSets = map[string]string{
	"rounded": "╭╮╰╯─│├┤┬┴┼",
	"sharp":   "┌┐└┘─│├┤┬┴┼",
	"heavy":   "┏┓┗┛━┃┣┫┳┻╋",
	"double":  "╔╗╚╝═║╠╣╦╩╬",
	"ascii":   "+-|",
}

type Size struct {
	Cols int `yaml:"cols" json:"cols"`
	Rows int `yaml:"rows" json:"rows"`
}

type Style struct {
	FG        string `yaml:"fg" json:"fg"`
	BG        string `yaml:"bg" json:"bg"`
	Bold      bool   `yaml:"bold" json:"bold"`
	Italic    bool   `yaml:"italic" json:"italic"`
	Underline bool   `yaml:"underline" json:"underline"`
}

type Design struct {
	Screen  string            `yaml:"screen" json:"screen"`
	Sizes   []Size            `yaml:"sizes" json:"sizes"`
	Palette map[string]string `yaml:"palette" json:"palette"`
	Styles  map[string]Style  `yaml:"styles" json:"styles"`
	Borders []string          `yaml:"borders" json:"borders"`
	Symbols []string          `yaml:"symbols" json:"symbols"`
	Rules   []string          `yaml:"rules" json:"rules"`

	borderRunes map[rune]bool
	symbolRunes map[rune]bool
}

type Resolved struct {
	Color      string `json:"color"`
	Background string `json:"background"`
	Bold       bool   `json:"bold"`
	Italic     bool   `json:"italic"`
	Underline  bool   `json:"underline"`
}

func Load(path string) (*Design, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read design: %w", err)
	}
	return Parse(data)
}

func Parse(data []byte) (*Design, error) {
	var d Design
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	dec.KnownFields(true)
	if err := dec.Decode(&d); err != nil {
		return nil, fmt.Errorf("design.yaml: %w", err)
	}
	if err := d.prepare(); err != nil {
		return nil, fmt.Errorf("design.yaml: %w", err)
	}
	return &d, nil
}

func (d *Design) prepare() error {
	if strings.TrimSpace(d.Screen) == "" {
		return fmt.Errorf("screen must not be empty")
	}
	if len(d.Sizes) == 0 {
		return fmt.Errorf("sizes must list at least one size")
	}
	seen := map[string]bool{}
	for i, s := range d.Sizes {
		if s.Cols <= 0 || s.Rows <= 0 {
			return fmt.Errorf("sizes[%d]: cols and rows must be positive", i)
		}
		key := fmt.Sprintf("%dx%d", s.Cols, s.Rows)
		if seen[key] {
			return fmt.Errorf("sizes[%d]: duplicate size %s", i, key)
		}
		seen[key] = true
	}
	if len(d.Palette) == 0 {
		return fmt.Errorf("palette must not be empty")
	}
	for name, value := range d.Palette {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("palette.%s must not be empty", name)
		}
	}
	if len(d.Styles) == 0 {
		return fmt.Errorf("styles must not be empty")
	}
	for letter, style := range d.Styles {
		runes := []rune(letter)
		if len(runes) != 1 {
			return fmt.Errorf("styles.%q: letter must be exactly one character", letter)
		}
		if runes[0] == ' ' || !unicode.IsPrint(runes[0]) || runewidth.RuneWidth(runes[0]) != 1 {
			return fmt.Errorf("styles.%q: letter must be one printable width-1 character", letter)
		}
		if style.FG != "" {
			if _, ok := d.Palette[style.FG]; !ok {
				return fmt.Errorf("styles.%q: fg refers to unknown palette entry %q", letter, style.FG)
			}
		}
		if style.BG != "" {
			if _, ok := d.Palette[style.BG]; !ok {
				return fmt.Errorf("styles.%q: bg refers to unknown palette entry %q", letter, style.BG)
			}
		}
	}
	if len(d.Borders) == 0 {
		return fmt.Errorf("borders must list at least one border set")
	}
	d.borderRunes = map[rune]bool{}
	for _, name := range d.Borders {
		set, ok := borderSets[name]
		if !ok {
			return fmt.Errorf("borders: unknown set %q, known sets: %s", name, strings.Join(BorderSetNames(), ", "))
		}
		for _, r := range set {
			d.borderRunes[r] = true
		}
	}
	d.symbolRunes = map[rune]bool{}
	for i, sym := range d.Symbols {
		runes := []rune(sym)
		if len(runes) != 1 {
			return fmt.Errorf("symbols[%d]: must be exactly one character", i)
		}
		if runewidth.RuneWidth(runes[0]) != 1 {
			return fmt.Errorf("symbols[%d]: %q is not width 1", i, sym)
		}
		d.symbolRunes[runes[0]] = true
	}
	return nil
}

func BorderSetNames() []string {
	names := make([]string, 0, len(borderSets))
	for name := range borderSets {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (d *Design) HasStyle(letter string) bool {
	_, ok := d.Styles[letter]
	return ok
}

func (d *Design) Resolve(letter string) (Resolved, bool) {
	style, ok := d.Styles[letter]
	if !ok {
		return Resolved{}, false
	}
	return Resolved{
		Color:      d.Palette[style.FG],
		Background: d.Palette[style.BG],
		Bold:       style.Bold,
		Italic:     style.Italic,
		Underline:  style.Underline,
	}, true
}

func (d *Design) AllowedRune(r rune) bool {
	if r >= 0x20 && r <= 0x7e {
		return true
	}
	return d.borderRunes[r] || d.symbolRunes[r]
}

func (d *Design) AllowedSize(cols, rows int) bool {
	for _, s := range d.Sizes {
		if s.Cols == cols && s.Rows == rows {
			return true
		}
	}
	return false
}

func (d *Design) SizesText() string {
	parts := make([]string, 0, len(d.Sizes))
	for _, s := range d.Sizes {
		parts = append(parts, fmt.Sprintf("%dx%d", s.Cols, s.Rows))
	}
	return strings.Join(parts, ", ")
}

func (d *Design) SmallestSize() Size {
	best := d.Sizes[0]
	for _, s := range d.Sizes[1:] {
		if s.Cols*s.Rows < best.Cols*best.Rows {
			best = s
		}
	}
	return best
}

func (d *Design) BorderRunes() []rune {
	runes := make([]rune, 0, len(d.borderRunes))
	for r := range d.borderRunes {
		runes = append(runes, r)
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
	return runes
}

func (d *Design) FrameRunes() (tl, tr, bl, br, h, v rune, ok bool) {
	if len(d.Borders) == 0 {
		return 0, 0, 0, 0, 0, 0, false
	}
	name := d.Borders[0]
	if name == "ascii" {
		return '+', '+', '+', '+', '-', '|', true
	}
	set, found := borderSets[name]
	if !found {
		return 0, 0, 0, 0, 0, 0, false
	}
	runes := []rune(set)
	if len(runes) < 6 {
		return 0, 0, 0, 0, 0, 0, false
	}
	return runes[0], runes[1], runes[2], runes[3], runes[4], runes[5], true
}

func (d *Design) StyleLetters() []string {
	letters := make([]string, 0, len(d.Styles))
	for letter := range d.Styles {
		letters = append(letters, letter)
	}
	sort.Strings(letters)
	return letters
}
