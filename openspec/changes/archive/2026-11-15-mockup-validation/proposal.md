# Proposal

## Why

The experiment only means something if the agent cannot leave the design system
and cannot hand back a grid that renders differently in the browser than in a
terminal. Validation is the product, not an afterthought: it is what turns agent
output into a trustworthy design artefact.

## What Changes

- Parse `design.yaml` into a palette, style letters, border sets, symbols, layout
  rules and the allowed screen sizes.
- Validate submitted mockups against a strict contract and return one precise
  error line per failure.
- Reject double-width and zero-width characters so one cell is always one
  character.
- Expose the design system as structured data for the brief and the gallery.

## Capabilities

### New Capabilities

- `mockup-validation`: the design-system contract and the exact rules a mockup
  must satisfy.

### Modified Capabilities

None.

## Impact

- New: `internal/design`, `internal/mockup`.
- Dependency: `github.com/mattn/go-runewidth`, `gopkg.in/yaml.v3`.
- Later changes consume the parsed design and the validator.
