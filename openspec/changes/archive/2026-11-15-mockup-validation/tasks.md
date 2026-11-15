# Tasks

## 1. Design system parsing

- [x] 1.1 Implement palette and style parsing and verify a fixture with five styles resolves
- [x] 1.2 Implement border-set lookup for rounded, sharp, heavy, double and ascii and verify every set's characters are accepted
- [x] 1.3 Implement size parsing and verify duplicate and empty sizes are rejected
- [x] 1.4 Implement an `AllowedRune` predicate and verify printable ASCII, borders, symbols and a disallowed glyph

## 2. Grid validation

- [x] 2.1 Implement row-count checks for `text` and `style` and verify short and long arrays
- [x] 2.2 Implement display-width checks with `go-runewidth` and verify a row one cell too wide
- [x] 2.3 Implement wide and zero-width rejection and verify emoji, CJK and combining marks
- [x] 2.4 Implement one-letter-per-cell checks and verify unknown letters report row and column
- [x] 2.5 Implement the size allow-list and verify an unsupported size
- [x] 2.6 Implement ordered error collection and verify three independent failures produce three lines
- [x] 2.7 Add a table-driven test with a valid 100x30 and a valid 80x24 fixture

## 3. Wiring

- [x] 3.1 Use the validator in `tuimock validate` and verify a broken mockup file is reported
- [x] 3.2 Document the design-system schema in `docs/design-system.md` and verify the example matches `examples/beandex/design.yaml`

## Workflow follow-up

- Archive the change after all tasks and quality gates pass.
