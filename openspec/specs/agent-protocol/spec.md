# agent-protocol Specification

## Purpose
Defines the MCP tool surface through which any agent reads the design brief,
submits mockups and reads reviewer feedback, independent of the client used.

## Requirements

### Requirement: Four tools only

The MCP server SHALL expose exactly the tools `get_brief`, `submit_mockup`,
`get_feedback` and `list_mockups`. Tool names and input schemas SHALL be stable
within the alpha.

#### Scenario: Tool discovery

- **WHEN** a client lists tools
- **THEN** it receives those four names and no others

### Requirement: get_brief returns the whole design context

`get_brief` SHALL return, in one response, the design system including palette,
style letters, border sets, symbols, layout rules and allowed sizes; every
reference capture with its note; the full `brief.md`; the mockup format with a
valid example; and an instruction to cite which reference each variant borrows
from in its rationale.

#### Scenario: Complete payload

- **WHEN** an agent calls `get_brief`
- **THEN** the response contains the brief text, every reference name and the allowed sizes

#### Scenario: Read-only

- **WHEN** `get_brief` is called before any submission
- **THEN** the state file is not created

### Requirement: submit_mockup accepts a grid

`submit_mockup` SHALL accept `screen`, `variant`, `cols`, `rows`, `text`,
`style` and `rationale`, validate the grid, and on success return `ok` with the
assigned id and version.

#### Scenario: Successful submission

- **WHEN** a valid 100x30 variant is submitted
- **THEN** the response contains `ok`, an id and version 1

#### Scenario: Revision

- **WHEN** the same screen and variant is submitted again
- **THEN** the response contains version 2

### Requirement: submit_mockup explains every rejection

On validation failure `submit_mockup` SHALL return a failure result whose text
contains one line per validation error and no stored state change. The response
SHALL NOT be a protocol error or a panic.

#### Scenario: Bad width

- **WHEN** row 7 is 82 cells wide in an 80-column mockup
- **THEN** the response is a tool failure naming row 7, 82 and 80

#### Scenario: No partial write

- **WHEN** validation fails
- **THEN** `list_mockups` shows the same versions as before the call

### Requirement: get_feedback reports picks and comments

`get_feedback` SHALL return, for the current screen by default or the screen
passed as input, every variant with its latest version, whether it is picked,
and all comments with their version and text.

#### Scenario: Picked variant with comment

- **WHEN** variant `b` is picked and has a comment on v1
- **THEN** the response marks `b` picked and includes the comment text

#### Scenario: No feedback yet

- **WHEN** nothing is picked and there are no comments
- **THEN** the response is a valid empty feedback document

### Requirement: list_mockups is a compact inventory

`list_mockups` SHALL return, per stored version, the id, screen, variant,
version, size and rationale. It SHALL accept an optional screen filter.

#### Scenario: Four variants

- **WHEN** four variants are stored
- **THEN** the response contains four entries with their sizes and rationales

#### Scenario: Screen filter

- **WHEN** the input names `list-view`
- **THEN** only mockups of that screen are returned

### Requirement: Standard streams stay clean

The server SHALL write only MCP protocol messages to standard output. Logs and
diagnostics SHALL go to standard error. A malformed unknown tool call SHALL
return a protocol error without terminating the server.

#### Scenario: Unknown tool

- **WHEN** a client calls `delete_everything`
- **THEN** the server returns an error and remains usable

#### Scenario: Clean stdout

- **WHEN** the server starts and answers one call
- **THEN** standard output parses as newline-delimited JSON-RPC only
