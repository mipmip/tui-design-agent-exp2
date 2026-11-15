# Spec Delta

## Purpose

Defines the single shared state file that the MCP server and the gallery read and
write, including versioning, picks, comments and safe concurrent access.

## ADDED Requirements

### Requirement: State file location and schema version

The system SHALL store state in `<project>/.tuimock/state.json` and SHALL include
an integer `version` field. A state file with an unknown major version SHALL be
rejected with an actionable error rather than silently rewritten.

#### Scenario: Fresh state

- **WHEN** no state file exists
- **THEN** reading returns empty mockups, empty picks and empty comments without error

#### Scenario: Unknown schema

- **WHEN** the file declares `"version": 99`
- **THEN** reading fails and the error names the supported version

### Requirement: Writes are atomic and exclusive

Every mutation SHALL hold an exclusive lock and SHALL write through a temporary
file followed by a rename, so a concurrent reader never observes a partial file.

#### Scenario: Concurrent submissions

- **WHEN** two processes submit mockups at the same time
- **THEN** both mockups are present afterwards and the file remains valid JSON

#### Scenario: Failed write

- **WHEN** serialisation fails before rename
- **THEN** the previous state file is unchanged

### Requirement: Mockups are versioned per screen and variant

A submission SHALL store `screen`, `variant`, `cols`, `rows`, `text`, `style`,
`rationale` and a creation timestamp. The first submission for a screen and
variant SHALL receive version 1; each later submission SHALL receive the next
integer version. Earlier versions SHALL remain readable.

#### Scenario: First submission

- **WHEN** variant `a` of `list-view` is submitted
- **THEN** it is stored as version 1

#### Scenario: Revision

- **WHEN** variant `a` of `list-view` is submitted again
- **THEN** a version 2 record is added and version 1 is still present

### Requirement: Each mockup has a stable id

Every stored version SHALL have a unique id that is stable for that version and
SHALL be returned to the caller.

#### Scenario: Id uniqueness

- **WHEN** four variants are submitted
- **THEN** all four ids differ

### Requirement: Picks are per screen

The system SHALL store at most one picked mockup id per screen. Picking a second
mockup SHALL replace the previous pick. A pick SHALL refer to an existing
version.

#### Scenario: Replace pick

- **WHEN** variant `a` is picked and then variant `c` is picked
- **THEN** feedback reports `c` as picked and `a` as not picked

#### Scenario: Unknown pick

- **WHEN** a pick names an id that does not exist
- **THEN** the write is rejected and the previous pick is unchanged

### Requirement: Comments attach to a version

A comment SHALL record the target mockup id, its variant, its version and the
comment text. Multiple comments per version SHALL be allowed. Empty comments
SHALL be rejected.

#### Scenario: Comment on v1

- **WHEN** a reviewer comments on v1 of variant `a`
- **THEN** `get_feedback` returns that comment under variant `a` with version 1

#### Scenario: Empty comment

- **WHEN** the comment text is empty or whitespace
- **THEN** the system rejects it

### Requirement: Reads never block writers unnecessarily

A read SHALL use a shared lock and SHALL return a consistent snapshot. A read
SHALL NOT create the state file when none exists.

#### Scenario: Read on empty project

- **WHEN** the gallery polls before any submission
- **THEN** it receives an empty snapshot and no file is created
