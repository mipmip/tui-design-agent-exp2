# Spec Delta

## Purpose

Defines how the design system is declared and the exact character-grid contract a
submitted mockup must satisfy before it may enter state or reach the gallery.

## ADDED Requirements

### Requirement: Design system declares styles and allowed characters

The system SHALL read `design.yaml` and expose a named screen, one or more
allowed sizes, a palette, a map of single-character style letters to text
attributes, a list of border sets and a list of allowed symbols. It SHALL reject
a design where a style letter is not exactly one printable character or refers
to a palette entry that does not exist.

#### Scenario: Valid design

- **WHEN** `design.yaml` declares styles `n` and `b` and both resolve
- **THEN** the design loads and exposes both letters

#### Scenario: Unknown palette reference

- **WHEN** a style refers to palette entry `accent` that is not declared
- **THEN** loading fails and the error names the style letter and `accent`

### Requirement: Grid dimensions are exact

A mockup SHALL be rejected unless `text` and `style` each contain exactly
`rows` entries.

#### Scenario: Short style array

- **WHEN** `rows` is 24 and `style` has 23 entries
- **THEN** validation reports that row 24 of `style` is missing

### Requirement: Text rows use display width

Every `text` row SHALL have a display width of exactly `cols`, measured with a
display-width algorithm. The error SHALL state the measured width and the
expected width.

#### Scenario: Row too wide

- **WHEN** row 7 measures 82 cells and `cols` is 80
- **THEN** validation reports row 7 as 82 wide when 80 was expected

#### Scenario: Combining character

- **WHEN** a row contains a combining mark that adds zero display width
- **THEN** validation rejects the row and names the offending character

### Requirement: Style rows map one letter per cell

Every `style` row SHALL contain exactly one style letter per cell, and every
letter SHALL exist in the design system.

#### Scenario: Unknown style letter

- **WHEN** a style row uses letter `z` and the design has no `z`
- **THEN** validation reports the row, the column and the unknown letter

#### Scenario: Style row too short

- **WHEN** a style row has fewer letters than `cols`
- **THEN** validation reports the measured letter count

### Requirement: Character allow-list

A mockup SHALL only contain printable ASCII, characters from a configured border
set, and declared symbols. Any other character SHALL be rejected with its row,
column and code point.

#### Scenario: Disallowed glyph

- **WHEN** a row contains a snowman outside the symbol list
- **THEN** validation reports the row, column and code point

#### Scenario: Border glyph allowed

- **WHEN** a row contains `╭` and the `rounded` border set is configured
- **THEN** validation accepts it

### Requirement: Wide characters are rejected

The system SHALL reject any character with display width other than one,
including emoji, CJK characters, zero-width joiners and combining marks.

#### Scenario: Emoji

- **WHEN** a row contains an emoji whose display width is two
- **THEN** validation rejects it and states that wide characters are not allowed

### Requirement: Size must be allowed

The `cols` and `rows` pair SHALL match one of the sizes declared in
`design.yaml`. The error SHALL list the allowed sizes.

#### Scenario: Unsupported size

- **WHEN** a mockup declares 120x40 and the design allows 100x30 and 80x24
- **THEN** validation rejects it and lists `100x30` and `80x24`

### Requirement: Validation reports every failure

Validation SHALL collect every independent failure and return them as separate
lines, ordered by row. A mockup with three bad rows SHALL produce three lines,
not one.

#### Scenario: Multiple failures

- **WHEN** two rows are too wide and one style row has an unknown letter
- **THEN** the result contains three error lines
