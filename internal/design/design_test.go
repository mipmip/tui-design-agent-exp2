package design

import "testing"

const validYAML = `
screen: test
sizes:
  - {cols: 10, rows: 3}
  - {cols: 4, rows: 2}
palette:
  bg: "#000000"
  fg: "#ffffff"
  accent: "#7aa2f7"
  sel: "#283457"
styles:
  n: {fg: fg}
  a: {fg: accent, bold: true}
  s: {fg: fg, bg: sel}
borders: [rounded]
symbols: ["▸"]
rules: ["one accent"]
`

func parse(t *testing.T, yaml string) *Design {
	t.Helper()
	d, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return d
}

func TestParseValid(t *testing.T) {
	d := parse(t, validYAML)
	if d.Screen != "test" {
		t.Fatalf("screen = %q", d.Screen)
	}
	if got := d.SizesText(); got != "10x3, 4x2" {
		t.Fatalf("sizes = %q", got)
	}
	if got := d.SmallestSize(); got.Cols != 4 || got.Rows != 2 {
		t.Fatalf("smallest = %+v", got)
	}
	if got := d.StyleLetters(); len(got) != 3 || got[0] != "a" || got[1] != "n" || got[2] != "s" {
		t.Fatalf("letters = %v", got)
	}
	resolved, ok := d.Resolve("a")
	if !ok || resolved.Color != "#7aa2f7" || !resolved.Bold {
		t.Fatalf("resolve a = %+v %v", resolved, ok)
	}
	if _, ok := d.Resolve("z"); ok {
		t.Fatal("z should not resolve")
	}
}

func TestParseErrors(t *testing.T) {
	cases := map[string]string{
		"missing screen": `
sizes: [{cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [rounded]
`,
		"no sizes": `
screen: test
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [rounded]
`,
		"duplicate size": `
screen: test
sizes: [{cols: 4, rows: 2}, {cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [rounded]
`,
		"bad size": `
screen: test
sizes: [{cols: 0, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [rounded]
`,
		"unknown palette reference": `
screen: test
sizes: [{cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: missing}}
borders: [rounded]
`,
		"multi character letter": `
screen: test
sizes: [{cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {nn: {fg: fg}}
borders: [rounded]
`,
		"unknown border set": `
screen: test
sizes: [{cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [wavy]
`,
		"wide symbol": `
screen: test
sizes: [{cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [rounded]
symbols: ["🙂"]
`,
		"unknown field": `
screen: test
sizes: [{cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [rounded]
colours: {}
`,
	}
	for name, yaml := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Parse([]byte(yaml)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestAllowedRune(t *testing.T) {
	d := parse(t, validYAML)
	for _, r := range []rune{'A', ' ', '~', '╭', '─', '▸'} {
		if !d.AllowedRune(r) {
			t.Fatalf("%q should be allowed", r)
		}
	}
	for _, r := range []rune{'☃', '🙂', '\t'} {
		if d.AllowedRune(r) {
			t.Fatalf("%q should be rejected", r)
		}
	}
}

func TestFrameRunes(t *testing.T) {
	d := parse(t, validYAML)
	tl, tr, bl, br, h, v, ok := d.FrameRunes()
	if !ok || tl != '╭' || tr != '╮' || bl != '╰' || br != '╯' || h != '─' || v != '│' {
		t.Fatalf("frame = %c %c %c %c %c %c %v", tl, tr, bl, br, h, v, ok)
	}
	ascii, err := Parse([]byte(`
screen: test
sizes: [{cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [ascii]
`))
	if err != nil {
		t.Fatal(err)
	}
	tl, tr, bl, br, h, v, ok = ascii.FrameRunes()
	if !ok || tl != '+' || tr != '+' || bl != '+' || br != '+' || h != '-' || v != '|' {
		t.Fatalf("ascii frame = %c %c %c %c %c %c %v", tl, tr, bl, br, h, v, ok)
	}
}

func TestAllowedSize(t *testing.T) {
	d := parse(t, validYAML)
	if !d.AllowedSize(10, 3) || d.AllowedSize(11, 3) {
		t.Fatal("AllowedSize mismatch")
	}
}

func TestParseFiveStyles(t *testing.T) {
	d := parse(t, `
screen: test
sizes: [{cols: 4, rows: 2}]
palette:
  bg: "#000000"
  fg: "#ffffff"
  accent: "#7aa2f7"
  warn: "#e0af68"
  sel: "#283457"
styles:
  n: {fg: fg}
  m: {fg: bg}
  a: {fg: accent, bold: true}
  w: {fg: warn, italic: true}
  s: {fg: fg, bg: sel, underline: true}
borders: [rounded]
`)
	if got := len(d.Styles); got != 5 {
		t.Fatalf("styles = %d", got)
	}
	if got := d.StyleLetters(); len(got) != 5 || got[0] != "a" || got[4] != "w" {
		t.Fatalf("letters = %v", got)
	}
	resolved, ok := d.Resolve("s")
	if !ok || resolved.Background != "#283457" || !resolved.Underline {
		t.Fatalf("resolve s = %+v %v", resolved, ok)
	}
}

func TestAllBorderSetsAccepted(t *testing.T) {
	names := BorderSetNames()
	if len(names) != 5 {
		t.Fatalf("border sets = %v", names)
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			d, err := Parse([]byte(`
screen: test
sizes: [{cols: 4, rows: 2}]
palette: {fg: "#fff"}
styles: {n: {fg: fg}}
borders: [` + name + `]
`))
			if err != nil {
				t.Fatal(err)
			}
			runes := d.BorderRunes()
			if len(runes) == 0 {
				t.Fatalf("no border runes for %s", name)
			}
			for _, r := range runes {
				if !d.AllowedRune(r) {
					t.Fatalf("%s border rune %q is not allowed", name, r)
				}
			}
		})
	}
}
