# Screen: beandex list view

Build the main list view of beandex, a Go/Bubble Tea navigator over git
repositories that use beans.

Sizes: 100x30 and 80x24.

The screen must show:

- repositories found under the configured base paths, grouped by parent
  directory
- per repo: name, ticket count, open ticket count, time of the last ticket
  update; columns are sortable and the active sort is visible
- an indicator for repos with overdue or stale tickets
- a quick filter opened with `/` (substring by default, fuzzy optional),
  visible when active
- key hints for at least: enter (open in beans), `/` (filter), sort, `r`
  (refresh), `q` (quit)
- a selected row, and enough sample data to fill the screen: about 3 groups
  and 12 repos, with realistic paths and counts

Deliver four variants that differ in layout, not just colour. Use these
directions:

1. a flat table with group header rows
2. a two-pane list plus detail
3. a compact one-line-per-repo view
4. one free choice

Cite the reference each variant borrows from in its rationale.
