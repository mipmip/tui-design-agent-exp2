# Proposal

## Why

A PoC is only an alpha base if the next agent can trust it. Claimed behaviour
that no test exercises decays immediately, and the briefing's done-criteria are
exactly a sequence of observable steps that can be automated.

## What Changes

- Add unit tests for design parsing, validation, state and brief composition.
- Add integration tests that start the real binary over stdio and call the MCP
  tools through the official client.
- Add black-box e2e tests for the complete loop: brief, four variants,
  validation recovery, gallery data, pick, comment, feedback and v2.
- Add a second-client test that uses the same server with different client
  metadata, proving there is no client-specific coupling.
- Add `nix flake check` coverage for tests and a documented single quality gate.

## Capabilities

### New Capabilities

- `verification`: the automated evidence that the alpha checklist is true.

### Modified Capabilities

None.

## Impact

- New: tests beside each package, `e2e/**`, `docs/testing.md`.
- Consumes the built binary; no new runtime dependency.
