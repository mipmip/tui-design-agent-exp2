# Tasks

## 1. Module and Nix foundation

- [x] 1.1 Create `go.mod` with module `github.com/mipmip/tuimock` and verify `go build ./...` succeeds
- [x] 1.2 Add `flake.nix` with `lib.genAttrs` over the four supported systems and verify the flake evaluates
- [x] 1.3 Add `packages.default` with `buildGoModule` and verify `nix build .` produces `result/bin/tuimock`
- [x] 1.4 Generate `flake.lock` and verify `nix flake check` passes

## 2. CLI surface

- [x] 2.1 Implement `version` and verify `tuimock version` exits 0 and prints a semver
- [x] 2.2 Implement subcommand dispatch with exit status 2 for unknown commands and verify it in a unit test
- [x] 2.3 Implement `--dir` resolution and missing-file errors and verify both cases in a table test
- [x] 2.4 Implement `validate` as a one-shot project check and verify valid and broken fixtures
- [x] 2.5 Implement `serve` startup that prints the bound URL and verify the e2e client reaches the listener
- [x] 2.6 Verify `nix develop -c go test ./...` passes

## 3. Documentation

- [x] 3.1 Write README build and run instructions and verify every command runs as written
- [x] 3.2 Document the `vendorHash` refresh command and verify it appears in the README

## Workflow follow-up

- Archive the change after all tasks and quality gates pass.
