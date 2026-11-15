# example-project Specification

## Purpose
Defines the bundled beandex reference project that gives agents a realistic
brief, a tight design system and reference screens to imitate.

## Requirements

### Requirement: Example project is a valid project folder

`examples/beandex` SHALL contain `design.yaml`, `brief.md` and a `refs/`
directory, and `tuimock validate --dir examples/beandex` SHALL exit 0.

#### Scenario: Validate example

- **WHEN** the user runs `tuimock validate --dir examples/beandex`
- **THEN** the command exits 0 and reports the screen and both sizes

### Requirement: Design system covers the starter palette

The example design SHALL declare the palette from the briefing, style letters
`n`, `m`, `t`, `a`, `w`, `b`, `s` and `k`, the rounded border set, the symbol
list and the five layout rules. It SHALL allow sizes 100x30 and 80x24.

#### Scenario: Style lookup

- **WHEN** a submission uses letter `s`
- **THEN** the validator resolves it to the selected-row style with a background

#### Scenario: Border lookup

- **WHEN** a submission draws a rounded corner
- **THEN** the validator accepts it

### Requirement: Brief describes the beandex list view

`brief.md` SHALL describe repos grouped by parent directory, per-repo ticket
counts and last update time, sortable columns with the active sort visible, an
overdue indicator, a `/` filter, key hints for enter, slash, sort, refresh and
quit, and a selected row. It SHALL ask for four variants that differ in layout
and SHALL specify both allowed sizes and realistic sample volume.

#### Scenario: Brief completeness

- **WHEN** an agent reads the brief through `get_brief`
- **THEN** it contains the grouping, sorting, filter, key-hint and size requirements

### Requirement: References carry notes

The `refs/` directory SHALL contain at least three plain-text captures and, for
each, a markdown note naming the lesson the agent should take.

#### Scenario: Reference list

- **WHEN** `get_brief` is called
- **THEN** every capture is returned together with its note

### Requirement: Reference captures stay within the design system

Every reference capture SHALL use only characters and colours that the example
design system allows, so copying from a reference cannot produce an invalid
mockup by accident.

#### Scenario: No exotic glyphs

- **WHEN** a capture contains a symbol
- **THEN** that symbol is in the example symbol list
