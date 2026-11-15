---
# tui-design-agent-exp2-vzy0
title: 03.2 E2E and alpha handoff
status: todo
type: epic
created_at: 2026-11-15T19:27:17Z
updated_at: 2026-11-15T19:27:17Z
parent: tui-design-agent-exp2-a0hk
---

OpenSpec change: e2e-verification.

Black-box e2e suite driving the real binary over stdio MCP and HTTP, covering four variants, validation recovery, pick, comment, v2 and a second client. Documentation and alpha checklist.
Exit: go test ./..., nix flake check and the e2e loop all pass.
