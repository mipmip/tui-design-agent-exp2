# Spec Delta

## Purpose

Defines the automated evidence that `tuimock` performs the alpha loop and keeps
performing it as the code changes.

## ADDED Requirements

### Requirement: One gate runs every test

The repository SHALL expose `go test ./...` as the single command that runs unit,
integration and end-to-end tests, and `nix flake check` SHALL pass on every
supported system.

#### Scenario: Clean checkout

- **WHEN** a fresh checkout runs `go test ./...`
- **THEN** every package passes without manual setup

#### Scenario: Nix check

- **WHEN** `nix flake check` runs
- **THEN** it evaluates and runs the checks without network access

### Requirement: E2E drives the real binary

The e2e suite SHALL start the compiled `tuimock` binary as a subprocess, connect
to `tuimock mcp` through stdio and to `tuimock serve` over HTTP. It SHALL NOT
call internal packages to simulate the loop.

#### Scenario: Real MCP transport

- **WHEN** the e2e test calls `get_brief`
- **THEN** the call travels over stdin and stdout to a separate process

#### Scenario: Real HTTP

- **WHEN** the e2e test fetches `/api/state`
- **THEN** the response comes from a listening loopback socket

### Requirement: E2E proves the full loop

The suite SHALL submit four valid variants, receive a validation failure and fix
it, verify all four in the gallery payload, pick one, comment on it, read the
feedback through MCP, and submit a v2 that the gallery reports alongside v1.

#### Scenario: Four variants

- **WHEN** the suite submits four variants
- **THEN** the gallery reports four variants at the requested sizes

#### Scenario: Self-correction

- **WHEN** a row is one cell too wide
- **THEN** the tool error names the row and widths and the corrected submission succeeds

#### Scenario: Revision

- **WHEN** the picked variant is resubmitted
- **THEN** it becomes v2 and both versions are listed

### Requirement: A second client is proven

The suite SHALL connect a second MCP client with different name and version
metadata to the same server and complete a tool call.

#### Scenario: Second client

- **WHEN** the second client calls `list_mockups`
- **THEN** it receives the same inventory as the first client

### Requirement: Failures are diagnosable

A failing e2e test SHALL print the server's stderr and the last tool response.
The suite SHALL use temporary project folders and SHALL NOT touch the repository
state.

#### Scenario: Failure output

- **WHEN** a tool call fails unexpectedly
- **THEN** the test output contains the stderr tail and the response text

#### Scenario: No leaked state

- **WHEN** the suite finishes
- **THEN** no `.tuimock` directory exists in the repository root
