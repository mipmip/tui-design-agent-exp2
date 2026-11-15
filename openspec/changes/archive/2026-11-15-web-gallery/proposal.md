# Proposal

## Why

A grid in a terminal transcript is unreadable. The reviewer needs to see variants
side by side at true character size, compare versions and record a decision
without leaving the browser or asking the agent what it thinks.

## What Changes

- Add an HTTP gallery at `127.0.0.1` serving one embedded page.
- Render every mockup of the current screen side by side with colours resolved
  from `design.yaml`.
- Add a per-mockup version switcher, a Pick button and a comment box.
- Poll the state file once per second and refresh the page content in place.
- Keep the page build-free: one HTML file, inline CSS and script.
- Add a small JSON API used by the page and by tests.

## Capabilities

### New Capabilities

- `gallery`: the browser surface, its rendering rules and its write API.

### Modified Capabilities

None.

## Impact

- New: `internal/web`, `internal/web/static`.
- Consumes `design`, `state` and the project loader.
- No frontend toolchain and no external assets.
