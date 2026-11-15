# Proposal

## Why

Validation proves a grid is legal; it does not prove the experiment's question.
That needs a real brief and reference screens that teach the agent restraint. The
beandex list view is the first screen Pim will actually judge.

## What Changes

- Add `examples/beandex/design.yaml` with the starter palette, eight style
  letters, the rounded border set, the allowed symbols, the layout rules and the
  100x30 and 80x24 sizes.
- Add `examples/beandex/brief.md` describing the beandex list view, its required
  content and the four requested layout directions.
- Add three to five reference captures under `refs/`, each with a plain-text
  capture and a note on what to learn from it.
- Add an example expectation file used by the e2e suite: four valid variants that
  the tests submit.

## Capabilities

### New Capabilities

- `example-project`: the bundled reference project and what it must contain.

### Modified Capabilities

None.

## Impact

- New: `examples/beandex/**`, `e2e/fixtures/beandex/**`.
- No Go packages. The binary treats it as an ordinary project folder.
