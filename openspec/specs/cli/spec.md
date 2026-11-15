# cli Specification

## Purpose
Defines the command surface of the single `tuimock` binary and the project folder
it operates on, so every later capability has one predictable entry point.

## Requirements

### Requirement: Binary exposes four subcommands

The `tuimock` binary SHALL expose exactly `serve`, `mcp`, `validate` and
`version` as top-level subcommands. An unknown subcommand SHALL print usage to
standard error and exit non-zero.

#### Scenario: Known subcommand

- **WHEN** the user runs `tuimock version`
- **THEN** the process prints the version string and exits with status 0

#### Scenario: Unknown subcommand

- **WHEN** the user runs `tuimock frobnicate`
- **THEN** the process prints usage to standard error and exits with status 2

### Requirement: Project folder is explicit

The `serve`, `mcp` and `validate` subcommands SHALL resolve the project folder
from `--dir`, defaulting to the current working directory. They SHALL reject a
folder that lacks `design.yaml` or `brief.md`, naming the missing file and the
resolved folder in the error.

#### Scenario: Missing brief

- **WHEN** the project folder contains `design.yaml` but no `brief.md`
- **THEN** the command exits non-zero and the error names `brief.md` and the folder

#### Scenario: Explicit directory

- **WHEN** the user passes `--dir /tmp/example` from another directory
- **THEN** the command operates on `/tmp/example`

### Requirement: Serve reports its address

`tuimock serve` SHALL print the gallery URL to standard output before it starts
accepting requests, and SHALL keep running until interrupted. It SHALL bind to
`127.0.0.1` by default.

#### Scenario: Default bind

- **WHEN** the user runs `tuimock serve` without `--addr`
- **THEN** the printed URL starts with `http://127.0.0.1:`

#### Scenario: Custom port

- **WHEN** the user runs `tuimock serve --addr 127.0.0.1:0`
- **THEN** the printed URL contains the port the kernel assigned

### Requirement: Validate is a one-shot check

`tuimock validate` SHALL parse `design.yaml`, count reference files, confirm that
`brief.md` is non-empty, and exit 0 with a one-line summary. It SHALL exit
non-zero with one precise error otherwise.

#### Scenario: Valid project

- **WHEN** every project input parses
- **THEN** the command prints a summary and exits 0

#### Scenario: Invalid design

- **WHEN** `design.yaml` is not valid YAML
- **THEN** the command exits non-zero and the error names the parse failure
