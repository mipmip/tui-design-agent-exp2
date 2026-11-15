---
# tui-design-agent-exp2-2biz
title: 01.3 State store
status: completed
type: epic
priority: normal
created_at: 2026-11-15T19:27:17Z
updated_at: 2026-11-15T19:58:35Z
parent: tui-design-agent-exp2-dj8m
---

OpenSpec change: state-store.

Shared .tuimock/state.json with exclusive file locking, schema versioning, per-variant version numbers, picks per screen and comments per mockup version.
Exit: concurrent writers lose no updates.
