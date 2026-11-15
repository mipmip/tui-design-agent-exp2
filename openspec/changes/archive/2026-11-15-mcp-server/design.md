# Design

## Context

The reviewer picks the MCP client, not the other way round. The server is a thin
adapter over `brief`, `mockup` and `state`. See `proposal.md` for motivation.

## Goals / Non-Goals

**Goals:**

- Four tools that map one-to-one onto the briefing's table.
- Typed inputs and outputs so the published schema matches reality.
- Validation errors that read like compiler errors.

**Non-Goals:**

- MCP resources, prompts, sampling or notifications.
- HTTP or SSE transport. Stdio only for the alpha.
- Authentication.

## Decisions

- **Official Go SDK with `mcp.AddTool` and typed structs.** The SDK derives the
  JSON schema from the struct and validates input before the handler runs. This
  keeps the tool surface honest without hand-written schemas.
- **Structured output plus text.** `get_brief` returns a structured document and
  the SDK also serialises it as text, so clients that only read content still
  work.
- **Validation errors as tool errors, not protocol errors.** `AddTool` packs a
  returned Go error into `CallToolResult` with `IsError`, which is exactly the
  behaviour the agent needs: it sees a failed tool call with readable text and
  can retry.
- **State writes only in `submit_mockup`.** `get_brief` and `list_mockups` are
  read-only; `get_feedback` is read-only. This keeps the experiment reproducible.
- **Screen defaults to the design's screen.** A tool call without `screen` uses
  `design.yaml`'s `screen`, so a one-screen project needs no ceremony.
- **Logging to stderr only.** The stdio transport owns stdout.

## Risks / Trade-offs

- [Structured output duplicates the text payload] -> acceptable size at 100x30
  and four variants; clients need both forms.
- [SDK schema inference cannot express cross-field rules] -> the validator is the
  authority; schema checks are only a first line.
- [A client that ignores tool errors sees only text] -> text contains the same
  information.

## Migration Plan

Greenfield. Archive creates the `agent-protocol` main spec.
