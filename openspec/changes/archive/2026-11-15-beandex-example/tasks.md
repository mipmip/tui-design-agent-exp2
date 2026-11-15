# Tasks

## 1. Design system

- [x] 1.1 Write `examples/beandex/design.yaml` with palette, eight styles, rounded borders, symbols, rules and both sizes
- [x] 1.2 Verify `tuimock validate --dir examples/beandex` accepts the design

## 2. Brief

- [x] 2.1 Write the list-view brief with grouping, columns, sorting, overdue indicator, filter, key hints and selected row
- [x] 2.2 Add the four layout directions and both sizes and verify an agent reading `get_brief` sees all of them
- [x] 2.3 Add realistic sample data guidance with three groups and twelve repos

## 3. References

- [x] 3.1 Add a lazygit-style capture and note and verify every glyph is allowed
- [x] 3.2 Add a k9s-style capture and note and verify every glyph is allowed
- [x] 3.3 Add a btop-style density capture and note and verify every glyph is allowed
- [x] 3.4 Add a beans TUI capture and note and verify every glyph is allowed
- [x] 3.5 Document the reference format in `examples/beandex/refs/README.md`

## 4. Fixtures

- [x] 4.1 Generate four valid beandex variants in the e2e suite and verify each passes validation
- [x] 4.2 Verify `get_brief` returns brief, design and every reference note in one payload

## Workflow follow-up

- Archive the change after all tasks and quality gates pass.
