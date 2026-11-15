# Proposal

## Why

The PoC needs a reproducible foundation before any agent-facing behaviour exists.
Without a single Go module, a plain Nix flake and a predictable CLI surface, every
later change would invent its own build and entry points. The alpha must install
and run the same way on Linux and macOS from day one.

## What Changes

- Add a Go module for `tuimock` with a single binary.
- Add the `tuimock` CLI with `serve`, `mcp`, `validate` and `version` subcommands.
- Add a plain Nix flake (no `flake-utils`) for `x86_64-linux`, `aarch64-linux`,
  `x86_64-darwin` and `aarch64-darwin`, with a package, a dev shell and checks.
- Add a project-directory resolver that rejects a folder without `design.yaml`
  or `brief.md` with an actionable error.
- Add `tuimock validate` to check a project folder without starting a server.

## Capabilities

### New Capabilities

- `cli`: command surface, flags, exit codes and project-folder preconditions.

### Modified Capabilities

None.

## Impact

- New: `go.mod`, `go.sum`, `cmd/tuimock`, `internal/project`, `flake.nix`,
  `flake.lock`.
- Toolchain: Go 1.24 or newer, Nix 2.18 or newer.
- No runtime services or external APIs.
