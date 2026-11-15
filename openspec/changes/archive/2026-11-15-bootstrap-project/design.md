# Design

## Context

The repository is empty apart from licence and ticket tooling. The alpha must be
buildable with one toolchain on Linux and macOS, and the PoC has a hard
one-evening size limit. See `proposal.md` for motivation.

## Goals / Non-Goals

**Goals:**

- One static Go binary with no service dependencies.
- Reproducible builds and a dev shell on four systems.
- Failure messages that tell an agent exactly which file is wrong.

**Non-Goals:**

- Packaging for Homebrew, Docker or nixpkgs upstream.
- Windows support.
- A plugin or configuration system.

## Decisions

- **Go with the standard `flag` package over a CLI framework.** Four subcommands
  do not justify a dependency, and the binary stays small. A hand-rolled
  `switch os.Args[1]` keeps exit codes explicit.
- **Plain `lib.genAttrs` over `flake-utils`.** The briefing forbids flake-utils.
  `eachDefaultSystem` saves seven lines and costs a dependency that hides the
  per-system structure.
- **`buildGoModule` with a pinned `vendorHash`.** This is the standard Nix path
  for Go and works without `go mod vendor` in the tree.
- **Project folder as a value object.** `internal/project.Project` carries the
  parsed design, brief and refs. Every command receives it, so no command
  re-implements path checks.
- **`--addr 127.0.0.1:0` for tests.** Binding port zero and reading the printed
  address avoids port collisions in parallel e2e runs.

## Risks / Trade-offs

- [Nix `vendorHash` drifts when dependencies change] -> the `nix flake check`
  gate fails loudly; refresh with the documented command in the README.
- [Darwin builds are unverified locally] -> `nix flake check` covers evaluation
  on all systems; the CI-free alpha documents that runtime proof is Linux only.
- [A global `flag` set is easy to get wrong] -> each subcommand owns a fresh
  `flag.FlagSet`, tested through the binary rather than through package internals.

## Migration Plan

Greenfield. No migration. The first archive creates the `cli` main spec.
