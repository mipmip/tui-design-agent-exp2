# Design

## Context

Agents produce text. A grid is only trustworthy if width is measured the way a
terminal measures it. See `proposal.md` for motivation.

## Goals / Non-Goals

**Goals:**

- One validator used by the CLI, the MCP server and the tests.
- Error messages an agent can act on without a human.
- The design system is data, not code.

**Non-Goals:**

- Rendering. The browser resolves styles; the validator never draws.
- Grapheme clustering beyond single-width acceptance.
- Per-project style inheritance or themes.

## Decisions

- **Style letters are single ASCII runes.** A letter is the only thing the agent
  writes, so one byte per cell keeps the parallel `text`/`style` arrays trivial
  to reason about. Multi-byte letters would make column counting ambiguous.
- **`go-runewidth` over a hand-rolled table.** The briefing requires a
  display-width library. `runewidth.RuneWidth` returns 0, 1 or 2 and matches
  terminal behaviour closely enough for the PoC.
- **Reject width != 1 instead of padding.** Silently dropping emoji would make
  the agent believe a cell exists that does not. Rejection is the experiment's
  point: the agent must learn the medium.
- **Collect all errors, not fail-fast.** The agent fixes more per round trip and
  the feedback loop converges faster.
- **Border sets are named groups.** `rounded`, `sharp`, `heavy`, `double` and
  `ascii` each map to a character set. A project opts in by name; the allow-list
  is the union of the chosen sets plus symbols plus printable ASCII.
- **Sizes live in `design.yaml`.** `brief.md` is prose for a human and an agent;
  putting dimensions in YAML keeps them machine-checkable and lets the brief
  narrate them without becoming a parser target.

## Risks / Trade-offs

- [Terminal width differs from browser width for some symbols] -> the PoC only
  allows the curated symbol list, all width 1 in common terminals.
- [ASCII-only style letters limit future themes] -> a theme can reuse letters;
  the palette is the variable layer.
- [Collecting errors costs a second pass] -> irrelevant at 100x30.

## Migration Plan

Greenfield. Archive creates the `mockup-validation` main spec.
