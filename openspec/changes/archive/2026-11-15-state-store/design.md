# Design

## Context

Two processes, one file. The briefing forbids a database and asks for a file lock.
See `proposal.md` for motivation.

## Goals / Non-Goals

**Goals:**

- Never lose a pick or comment under concurrent writes.
- Keep old versions so the reviewer can compare v1 and v2.
- One small, inspectable JSON file.

**Non-Goals:**

- Multi-user identity or auth.
- History beyond what is needed for version comparison.
- Cross-machine sync.

## Decisions

- **Append-only versions, derived latest.** A submission never mutates an
  existing record; it adds a new one. `latest` is computed by screen and variant.
  This makes v1 and v2 comparison free and avoids in-place edits.
- **One pick per screen, stored as a mockup id.** A pick is a decision about a
  variant, and only the latest version of that variant matters for `get_feedback`.
  Storing the id keeps the decision auditable.
- **`gofrs/flock` over hand-rolled `syscall.Flock`.** It handles the platform
  differences and stale-lock semantics. The dependency is tiny and widely used.
- **Temp file plus rename for writes.** `os.CreateTemp` in the state directory
  followed by `os.Rename` gives atomic replacement on the same filesystem.
- **Exclusive lock for writes, shared for reads.** Writes are rare and small;
  readers poll once per second. A shared read lock keeps the gallery responsive.
- **IDs from a random hex suffix with a screen-variant prefix.** Human-readable
  ids help the reviewer and the agent talk about the same mockup.

## Risks / Trade-offs

- [State file grows without bound across sessions] -> v1, acceptable for a PoC;
  the archive change documents pruning as future work.
- [A crash between temp write and rename leaves a temp file] -> temp names are
  ignored by reads and overwritten by the next write.
- [Pick points at a superseded version] -> `get_feedback` maps a pick to its
  variant and reports the latest version as picked.

## Migration Plan

Greenfield. Archive creates the `mockup-state` main spec.
