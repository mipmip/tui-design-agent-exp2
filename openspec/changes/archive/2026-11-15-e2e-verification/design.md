# Design

## Context

Every previous change ships its own tests. This change adds the black-box layer
and the gate that makes the alpha claim checkable. See `proposal.md` for
motivation.

## Goals / Non-Goals

**Goals:**

- Prove the briefing's done-criteria without a browser in the loop.
- Keep the suite hermetic and parallel-safe.
- Fail with the information needed to fix the failure.

**Non-Goals:**

- Browser automation. The gallery's data and rendered HTML are tested separately.
- Mutation testing or coverage thresholds.
- CI configuration. The alpha is local-first.

## Decisions

- **`TestMain` builds the binary once.** `go build -o <tmp>/tuimock ./cmd/tuimock`
  runs before the e2e tests and every test uses that path. This proves the same
  artefact users install.
- **`t.TempDir()` project copies.** Each test copies `examples/beandex` into a
  temp folder, so tests never share `.tuimock` and can run in parallel.
- **MCP over `CommandTransport`.** The official client launches the binary
  itself, which exercises the stdio transport and process lifecycle.
- **HTTP over `net/http` against the printed address.** `tuimock serve --addr
  127.0.0.1:0` prints the bound URL; the test parses it instead of guessing a
  port.
- **A second client is a new `mcp.Client` with different implementation
  metadata.** It is the cheapest faithful proof that nothing depends on the
  first client's identity.
- **Nix check runs `go test ./...` in the sandbox.** Tests bind loopback and use
  no external network, so the sandbox is sufficient.
- **Failures include stderr and response text.** Connection failures print the
  server's stderr and HTTP failures print the status and body, so a black-box
  failure is debuggable without a rerun.

## Risks / Trade-offs

- [Subprocess tests are slower] -> two processes and a handful of calls; the
  whole suite stays under a few seconds.
- [Port parsing is brittle if `serve` changes its banner] -> the banner is part
  of the `cli` spec and covered by its own test.
- [Nix sandbox may restrict loopback] -> verified during this change; if it does,
  the check runs the unit and integration layers and e2e stays in `go test`.

## Migration Plan

Greenfield. Archive creates the `verification` main spec.
