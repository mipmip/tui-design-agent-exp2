---
# tui-design-agent-exp2-9o6k
title: 02.1 MCP tools
status: todo
type: epic
created_at: 2026-11-15T19:27:17Z
updated_at: 2026-11-15T19:27:17Z
parent: tui-design-agent-exp2-88d4
---

OpenSpec change: mcp-server.

Four tools over stdio with the official Go SDK: get_brief, submit_mockup, get_feedback, list_mockups. Errors are returned as readable text, never panics.
Exit: an MCP client can complete the submit and feedback loop.
