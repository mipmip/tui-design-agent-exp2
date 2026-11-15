# References

Each reference is a pair:

```
refs/lazygit.txt   plain-text screen capture
refs/lazygit.md    one short note: what to learn from it
```

The MCP server pairs files by base name, sorts them by name and returns every
capture with its note in `get_brief`.

## Capture

Capture a real screen with tmux:

```sh
tmux capture-pane -p > refs/lazygit.txt
```

A capture wider than the largest allowed screen is fine for a note, but it is
not a valid mockup. Keep reference captures at or below the widest allowed
`cols` so the example project can be checked as a whole.

## Note

One or two lines. Say what the screen does well, not what it contains:

```
Dense but calm: muted borders, one accent colour, status glyphs in the left
gutter, and every useful key in the footer.
```

## Constraints

- Only printable ASCII, characters from the configured border set, and the
  configured symbols may appear.
- No tabs. Use spaces.
- No wide characters (CJK, emoji) and no zero-width characters.
- The note is plain Markdown text.

The e2e suite loads every pair, checks that the note is non-empty, checks the
width against the widest allowed screen and checks every character against the
design system.

## Current set

| Name | Lesson |
| --- | --- |
| `lazygit` | Muted borders, one accent, status glyphs in the left gutter, keys in the footer. |
| `k9s` | Dense columns, active sort marker, strong selected row. |
| `btop` | Boxed panels, compact numbers, colour used for meaning. |
| `beans` | Ticket counts and status language that the example brief reuses. |

These captures are reconstructed for the example project rather than copied
from a live session. Replace them with real captures when tuning the design
system.
