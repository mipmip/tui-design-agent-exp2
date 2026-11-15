# TUI mockup PoC — one-evening briefing

Nov 15, 2026 · @Pim Snel (prive)

## Goal

Build in one evening a tiny tool that lets any MCP-capable agent produce TUI mockups as strict character grids, shown side by side in a browser, so I can pick and comment before any code is written. Working name: `tuimock`.

The question this PoC answers: **do a tight design system plus reference screens make agent-made TUI designs good enough to build from?** Everything else is secondary.

What changed since the first experiment:

- **Much smaller.** No terminal engine, no PTY, no Bubble Tea adapter, no golden files. Mockups are data, not rendered programs.
- **Quality first.** The agent can only use the colours, styles and border characters the design system defines, and it gets real reference screens as examples.
- **Same provider-agnostic coupling.** The tool never calls an LLM; any agent drives it through MCP.

## How the loop works

One Go binary runs an MCP server (stdio) and a small web page on localhost; both share one JSON file as state.

1. I run `tuimock serve` in a project folder and open the browser page. The folder holds `design.yaml`, `refs/` and `brief.md`.
2. I tell my agent: "design the screen in the brief, 4 variants".
3. The agent calls `get_brief` and receives the design system, the reference screens and the screen brief in one go.
4. The agent calls `submit_mockup` per variant. Invalid mockups are rejected with exact errors (row 7 is 82 wide, expected 80) and the agent fixes them.
5. Valid mockups appear live in the browser, side by side at true character size.
6. I mark one as picked and write comments per mockup.
7. The agent calls `get_feedback`, revises the picked one and submits it as a new version.

State lives in `.tuimock/state.json`: mockups with versions, picks and comments. No database.

## Mockup format

A mockup is a strict grid: `text` holds each row as a string of exactly `cols` cells, and `style` holds a parallel string with one letter per cell. Each letter is a style defined in the design system. The agent never writes colours, only letters, so it cannot leave the design system.

```json
{
  "screen": "list-view",
  "variant": "a",
  "cols": 24,
  "rows": 4,
  "text": [
    "╭─ beandex ────────────╮",
    "│ ~/code         12  3 │",
    "│ ~/work/nivis   4   1 │",
    "╰──────────────────────╯"
  ],
  "style": [
    "bbtttttttttbbbbbbbbbbbbb",
    "bnnnnnnnnnnnnnnnnnnnnanb",
    "bssssssssssssssssssssssb",
    "bbbbbbbbbbbbbbbbbbbbbbbb"
  ],
  "rationale": "One line on the main design choice"
}
```

Validation in `submit_mockup`, each failure returned as one precise line:

- `text` and `style` both have exactly `rows` entries
- every `text` row is exactly `cols` cells wide (measure with a display-width library, not byte or rune length)
- every `style` row has exactly one letter per cell, and every letter exists in the design system
- only allowed characters: printable ASCII, the border sets listed in the design system, and its symbol list
- wide characters (emoji, CJK) are rejected in this PoC, which keeps one cell equal to one character
- `cols` × `rows` matches one of the sizes the brief asks for

## Design system and references

These two inputs are the experiment: the tool's job is to hand them to the agent in full and enforce them.

**`design.yaml`** defines the style letters, the allowed characters and the layout rules. Starter version (values are placeholders to tune):

```yaml
palette:
  bg: "#1a1b26"
  fg: "#c0caf5"
  muted: "#565f89"
  accent: "#7aa2f7"
  warn: "#e0af68"
  sel_bg: "#283457"
styles:            # letter -> style
  n: { fg: fg }
  m: { fg: muted }
  t: { fg: accent, bold: true }
  a: { fg: accent }
  w: { fg: warn }
  b: { fg: muted }
  s: { fg: fg, bg: sel_bg, bold: true }
  k: { fg: bg, bg: accent }       # key hints in the footer
borders: [rounded]                 # ╭╮╰╯─│ plus ├┤┬┴┼
symbols: ["▸", "●", "○", "↑", "↓", "…"]
rules:
  - One accent colour; warn only for things needing attention
  - Selection is a full-width bar, never just a cursor glyph
  - Footer row shows the 4-6 most useful keys as k-styled chips
  - Numbers right-aligned; truncate long paths with … in the middle
  - Never fill more than about 70% of cells with text
```

**`refs/`** holds 3 to 5 real TUI screens I admire, each as a plain-text capture plus a note on what to learn from it. Capture with `tmux capture-pane -p > refs/lazygit.txt` and write `refs/lazygit.md` ("dense but calm: muted borders, one accent, keys in footer"). Good candidates: lazygit, k9s, btop, and the beans TUI itself for consistency.

`get_brief` returns the design system, all refs with their notes, and `brief.md` as one block of text, with an instruction to cite which reference each variant borrows from in its `rationale`.

## MCP tools and browser gallery

Four MCP tools, nothing more:

| Tool | Input | Returns |
| --- | --- | --- |
| `get_brief` | — | design system, refs with notes, screen brief, required sizes, format spec with the example |
| `submit_mockup` | the mockup JSON | `ok` with id and version, or the list of validation errors |
| `get_feedback` | screen? | per mockup: picked yes/no and my comments |
| `list_mockups` | screen? | ids, variants, versions, sizes, rationale |

Use the official MCP Go SDK, stdio transport. `tuimock serve` starts the web page and `tuimock mcp` the MCP server; both read and write `.tuimock/state.json` with a file lock, so they can run as two processes.

**Browser gallery**, one static page embedded in the binary:

- all mockups of the current screen side by side, each drawn as a grid of monospace cells with colours resolved from `design.yaml`, at true size with horizontal scroll
- under each: variant name, version, rationale, a **Pick** button and a comment box
- a version switcher per mockup to compare v1 against v2
- live refresh by polling the state file every second (simpler than SSE for tonight)
- no build step: plain HTML, CSS and a small inline script

## Test case: beandex list view

The first `brief.md` is the main list view of beandex, the planned Go/Bubble Tea navigator over git repos that use beans. Sizes: 100×30 and 80×24.

What the screen must show:

- repos found under the configured base paths, grouped by parent directory
- per repo: name, ticket count, open ticket count, time of the last ticket update; columns sortable, with the active sort visible
- an indicator for repos with overdue or stale tickets
- a quick filter opened with `/` (substring by default, fuzzy optional), visible when active
- key hints for at least: enter (open in beans), `/` (filter), sort, `r` (refresh), `q` (quit)
- a selected row, and enough sample data to fill the screen: about 3 groups and 12 repos, with realistic paths and counts

Ask the agent for 4 variants that differ in layout, not just colour: for example a flat table with group header rows, a two-pane list plus detail, a compact one-line-per-repo view, and one free choice.

## Done tonight when

- [ ] an agent calls `get_brief` and submits 4 valid beandex variants, fixing its own validation errors without my help
- [ ] the variants show up side by side in the browser in the right colours and at true size
- [ ] I pick one and comment; the agent reads that via `get_feedback` and submits a v2 I can compare with v1
- [ ] the loop also works with a second MCP agent (for example Codex), with no code change
- [ ] I write down a verdict: are these designs good enough to build from, and did the design system or the refs make the difference?

**Out of scope:** rendering real TUIs, PTY or terminal engines, Bubble Tea code generation, golden files or tests, multiple screens or projects, auth, packaging beyond `go build`. If the verdict is positive, the next step is turning a picked mockup into a Bubble Tea `View()`; that is a separate briefing.
