# Design system

`design.yaml` is both documentation for the agent and the enforcement table for
validation. The agent may only use style letters and characters defined here.

## Schema

```yaml
screen: list-view
sizes:
  - {cols: 100, rows: 30}
  - {cols: 80, rows: 24}
palette:
  bg: "#1a1b26"
  fg: "#c0caf5"
  muted: "#565f89"
  accent: "#7aa2f7"
  warn: "#e0af68"
  sel_bg: "#283457"
styles:
  n: {fg: fg}
  m: {fg: muted}
  t: {fg: accent, bold: true}
  a: {fg: accent}
  w: {fg: warn}
  b: {fg: muted}
  s: {fg: fg, bg: sel_bg, bold: true}
  k: {fg: bg, bg: accent}
borders: [rounded]
symbols: ["▸", "●", "○", "↑", "↓", "…"]
rules:
  - One accent colour; warn only for things needing attention
```

| Field | Type | Rules |
| --- | --- | --- |
| `screen` | string | Required. Default screen name for submissions. |
| `sizes` | list of `{cols, rows}` | At least one. Positive, unique. Submissions must match one entry. |
| `palette` | map name to CSS colour | At least one entry. Style colours must name an entry. |
| `styles` | map letter to style | Letter is one printable, non-space rune of width 1. At least one entry. |
| `borders` | list of set names | `rounded`, `sharp`, `heavy`, `double`, `ascii`. |
| `symbols` | list of strings | Each is one rune of width 1. |
| `rules` | list of strings | Passed to the agent as design guidance. |

Style fields:

| Field | Meaning |
| --- | --- |
| `fg` | Palette name for the foreground. |
| `bg` | Palette name for the background. |
| `bold` | Boolean. |
| `italic` | Boolean. |
| `underline` | Boolean. |

Colours are palette names, not raw values. This keeps the palette as the single
place to change a colour.

## Border sets

| Name | Characters |
| --- | --- |
| `rounded` | `╭ ╮ ╰ ╯ ─ │ ├ ┤ ┬ ┴ ┼` |
| `sharp` | `┌ ┐ └ ┘ ─ │ ├ ┤ ┬ ┴ ┼` |
| `heavy` | `┏ ┓ ┗ ┛ ━ ┃ ┣ ┫ ┳ ┻ ╋` |
| `double` | `╔ ╗ ╚ ╝ ═ ║ ╠ ╣ ╦ ╩ ╬` |
| `ascii` | `+ - |` |

Every rune from every configured set is accepted in text. Configuring more
than one set allows a mix, so use one set per design system.

## Allowed characters

A text cell may contain:

- printable ASCII from U+0020 (space) to U+007E;
- any rune from a configured border set;
- any configured symbol.

Everything else is rejected. Wide characters (CJK, emoji) and zero-width
characters (combining marks) are rejected so one cell always equals one
character.

## Validation messages

Validation collects all independent failures in a fixed order. Examples:

```
row 7: text is 101 cells wide, expected 100
row 7: style has 101 letters, expected 100
row 1: text is 5 cells wide, expected 4
row 1, column 4: style letter "x" is not defined in the design system
row 2: character "🙂" (U+1F642) is 2 cells wide; wide characters are not allowed
size 120x40 is not allowed; allowed sizes: 100x30, 80x24
```

The width error is suppressed when a row contains an invalid character,
because the display width is then undefined.

## Loading behaviour

- Unknown YAML fields are errors. This catches typos such as `colour:` or
  `border:`.
- A missing palette reference in a style is an error.
- A style letter must be one rune. Multi-character keys are rejected.
- Symbols must be one rune of width 1. Emoji symbols are rejected.

`tuimock validate --dir DIR` loads the same file the MCP server and gallery
use. `examples/beandex/design.yaml` is the canonical example.
