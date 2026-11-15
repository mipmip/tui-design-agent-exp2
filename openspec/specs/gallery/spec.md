# gallery Specification

## Purpose
Defines the local browser gallery in which mockups are reviewed side by side,
versions are compared and a pick and comments are recorded.

## Requirements

### Requirement: Gallery shows mockups side by side

The gallery SHALL render every stored variant of the selected screen in one
horizontally scrollable row, each mockup drawn at true character size with a
monospace font. Each mockup SHALL show its variant name, current version,
rationale and size.

#### Scenario: Four variants

- **WHEN** four variants are stored for the screen
- **THEN** four mockups appear in the row with their rationales

#### Scenario: True size

- **WHEN** a mockup is 80 columns wide
- **THEN** its rendered width equals 80 times the character advance of the font

### Requirement: Styles resolve from the design system

Each cell SHALL use the foreground, background, bold and other attributes of its
style letter as declared in `design.yaml`. Spaces SHALL still receive the
background colour of their style.

#### Scenario: Selected row

- **WHEN** a cell uses style `s` with a background
- **THEN** the cell renders with that background including its spaces

#### Scenario: Unknown palette entry

- **WHEN** a design cannot resolve a colour
- **THEN** the gallery shows a readable error instead of a broken page

### Requirement: Version switcher

Each variant SHALL offer every stored version and SHALL default to the latest.
Switching versions SHALL update the rendered grid, the rationale and the
comments for that version without a page load.

#### Scenario: Compare v1 and v2

- **WHEN** variant `a` has v1 and v2
- **THEN** the switcher lists both and selecting v1 shows the v1 grid and rationale

### Requirement: Pick one mockup per screen

The gallery SHALL render a Pick button per variant. Picking SHALL replace the
previous pick for that screen, and the picked variant SHALL be visibly marked.

#### Scenario: Pick replaces

- **WHEN** the reviewer picks `b` after `a`
- **THEN** `b` is marked picked and `a` is not

### Requirement: Comments are per version

The gallery SHALL provide a comment box per variant. Submitting a non-empty
comment SHALL append it to the displayed version and clear the box. An empty
comment SHALL be ignored.

#### Scenario: Comment persists

- **WHEN** the reviewer comments on v2 of `a`
- **THEN** the comment shows under v2 after the next poll

### Requirement: Live refresh by polling

The page SHALL fetch the current state at most once per second and update the
gallery in place. It SHALL preserve the selected version and any text being
typed in a comment box across refreshes.

#### Scenario: Agent submits while reviewing

- **WHEN** the MCP server stores a new variant while the page is open
- **THEN** the new variant appears within two seconds without a reload

#### Scenario: Typing is preserved

- **WHEN** a poll completes while the reviewer is typing
- **THEN** the text in the comment box is unchanged

### Requirement: JSON API

The server SHALL expose `GET /api/state`, `POST /api/pick` and
`POST /api/comment`. `GET /api/state` SHALL return the design, screen, mockups,
picks and comments as JSON. Invalid JSON or an unknown mockup id SHALL return a
4xx status with a readable message.

#### Scenario: State payload

- **WHEN** a client requests `/api/state`
- **THEN** the response contains `design`, `screen` and `mockups`

#### Scenario: Bad pick

- **WHEN** a pick names a missing id
- **THEN** the response status is 400 and the message names the id

### Requirement: Local only

The server SHALL bind to `127.0.0.1` by default and SHALL serve no file from
outside its embedded assets. It SHALL NOT require network access.

#### Scenario: Default bind

- **WHEN** the server starts without an address
- **THEN** it listens on a loopback address only

#### Scenario: Path traversal

- **WHEN** a request asks for `../go.mod`
- **THEN** the server does not return that file
