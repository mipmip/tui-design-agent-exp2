---
# tui-design-agent-exp2-aw3w
title: 01.1 Bootstrap and Nix
status: completed
type: epic
priority: normal
created_at: 2026-11-15T19:27:17Z
updated_at: 2026-11-15T19:58:05Z
parent: tui-design-agent-exp2-dj8m
---

OpenSpec change: bootstrap-project.

Go module, CLI skeleton with serve/mcp/validate/version, plain Nix flake for four systems, dev shell, build and check outputs. No flake-utils.
Exit: nix build . produces a working binary and nix flake check passes.
