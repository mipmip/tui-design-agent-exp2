# Alpha checklist

`tuimock` is an alpha base when every item below is true and proven by the
automated suite. The verdict section is written by Pim after using the tool.

## Product

- [x] `tuimock mcp` speaks MCP over stdio and exposes exactly four tools:
      `get_brief`, `submit_mockup`, `get_feedback`, `list_mockups`.
      Evidence: `TestToolDiscovery`.
- [x] `tuimock serve` serves a gallery on localhost and reads the same state.
      Evidence: `TestAlphaLoop`.
- [x] An agent can call `get_brief`, submit four valid variants, and fix its own
      validation errors from the returned messages.
      Evidence: `TestBriefAndValidation`, `TestAlphaLoop`.
- [x] Variants appear side by side in the browser at true character size with
      colours resolved from `design.yaml`.
      Evidence: `TestIndex`, `TestFragmentAndComments`,
      `TestSelectedRowSpaceKeepsBackground`, `TestAlphaLoop`.
- [x] Picking a variant and writing comments is possible in the browser.
      Evidence: `TestAPIPickAndComment`, `TestAlphaLoop`.
- [x] `get_feedback` returns picks and comments; a revised submission creates a
      v2 that can be compared with v1 in the gallery.
      Evidence: `TestAlphaLoop`, `TestVersionSwitcher`.
- [x] A second MCP client works without code changes.
      Evidence: `TestAlphaLoop` connects a client named `codex` and compares its
      inventory with the first client.

## Engineering

- [x] The flake declares `x86_64-linux`, `aarch64-linux`, `x86_64-darwin` and
      `aarch64-darwin` with plain `lib.genAttrs`, and `nix build .` passes on
      the current system.
      Evidence: `flake.nix`, `nix build .#default` on x86_64-linux.
- [x] `nix flake check` passes.
      Evidence: `nix flake check` on 2026-11-15.
- [x] `go test ./...` passes, including black-box e2e tests.
      Evidence: `go test ./...` and `nix develop -c go test ./...` on
      2026-11-15.
- [x] e2e covers the complete loop in `docs/briefing.md`.
      Evidence: `TestAlphaLoop`, `TestBriefAndValidation`.
- [x] State writes are safe when the MCP server and gallery run concurrently.
      Evidence: `TestConcurrentSubmit`, `TestConcurrentReadDuringWrite`.
- [x] No code comments outside generated or machine-required files.
      Evidence: `TestNoCodeComments`.

## Verdict

Filled in by Pim after the first real session:

1. Are these designs good enough to build from?
2. Did the design system or the reference screens make the difference?

Verdict: _pending_
