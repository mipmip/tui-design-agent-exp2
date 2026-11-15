# Proposal

## Why

The product has no user interface for the agent. The only coupling is MCP:
whatever client the reviewer uses must be able to fetch the design constraints,
submit a grid, read the verdict and enumerate versions without the tool knowing
which model is on the other end.

## What Changes

- Add an MCP server over stdio using the official Go SDK.
- Expose exactly four tools: `get_brief`, `submit_mockup`, `get_feedback` and
  `list_mockups`.
- Return validation failures as readable text so the agent can correct itself.
- Keep protocol output clean: all diagnostics go to standard error.
- Make the server work from any MCP client with no client-specific code.

## Capabilities

### New Capabilities

- `agent-protocol`: the MCP tool contract, payloads and error semantics.

### Modified Capabilities

None.

## Impact

- New: `internal/mcpserver`, `internal/brief`.
- Dependency: `github.com/modelcontextprotocol/go-sdk`.
- Consumes `design`, `mockup` and `state`.
