# Proposal

## Why

The MCP server and the browser gallery are separate processes. Without one
locked, versioned state file they overwrite each other's picks and comments, and
the agent cannot tell which version a reviewer approved.

## What Changes

- Add `.tuimock/state.json` with a schema version, mockups, picks and comments.
- Add exclusive file locking around every read-modify-write cycle.
- Version a mockup per screen and variant: a resubmission becomes v2, v3, and so
  on, while older versions remain readable.
- Store one pick per screen and comments per mockup version.
- Add atomic writes through a temporary file and rename.

## Capabilities

### New Capabilities

- `mockup-state`: persistence, concurrency, versioning, picks and comments.

### Modified Capabilities

None.

## Impact

- New: `internal/state`.
- Dependency: `github.com/gofrs/flock`.
- The state file is ignored by version control.
