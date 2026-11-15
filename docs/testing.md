# Testing

## Gates

```sh
go test ./...
go vet ./...
nix flake check
```

`nix flake check` builds `packages.default` with `doCheck = true`, so it runs
the same tests inside the Nix build. A plain `go test ./...` is the fast local
loop. `nix develop -c go test ./...` proves the dev shell is complete.

## Unit and package tests

| Package | Covers |
| --- | --- |
| `internal/design` | YAML parsing, five-style palettes, all five border sets, allowed characters, sizes, frame runes, unknown fields. |
| `internal/mockup` | Row counts, display width, wide and zero-width characters, unknown style letters, allowed sizes, valid 100x30 and 80x24 grids. |
| `internal/project` | Design, brief and reference loading, missing files, paired notes, invalid design. |
| `internal/state` | Empty and missing files, schema version, versioning, unique IDs, picks, comments, failed writes, parallel submissions, read during write, unsupported schema. |
| `internal/brief` | Payload contents, example validity, example without a frame. |
| `internal/web` | Page contract, stylesheet, escaped cells, selected-row backgrounds, version switcher, fragment, API, status codes, golden grid. |
| `cmd/tuimock` | Version, help, dispatch, validate, missing project, empty brief, broken mockup JSON. |

## End-to-end tests

`e2e` builds the real binary once in `TestMain`, copies `examples/beandex` to a
temporary folder, and drives it over stdio and HTTP.

| Test | Proves |
| --- | --- |
| `TestBriefAndValidation` | `get_brief` returns both sizes, four references with notes, the four layout directions, a valid example and creates no state. A 101-cell row is rejected with row, measured and expected widths; nothing is stored; empty list and feedback documents are valid. |
| `TestToolDiscovery` | Exactly four tools exist, `submit_mockup` declares its required fields, an unknown tool returns a protocol error, and the session stays usable. |
| `TestAlphaLoop` | Four variants submit, the gallery serves the same state and colours, pick and comment work over HTTP, a second client sees the same inventory and feedback, a v2 is created and both versions appear, pick and comment survive a gallery restart and a fresh client. |
| `TestExampleProject` | Every reference has a note, fits the screen width and uses only allowed characters; the brief names both sizes. |
| `TestRepositoryStaysClean` | The test run leaves no `.tuimock` in the repository root. |
| `TestNoCodeComments` | No Go source file contains comments other than the `go:embed` directive. |

## Alpha criterion map

| Criterion | Evidence |
| --- | --- |
| Four MCP tools, no more | `TestToolDiscovery` |
| Gallery and MCP share state | `TestAlphaLoop` |
| Four valid variants with self-correction | `TestBriefAndValidation`, `TestAlphaLoop` |
| Side by side, resolved colours, true size | `TestIndex`, `TestFragmentAndComments`, `TestSelectedRowSpaceKeepsBackground`, `TestAlphaLoop` |
| Pick and comment in the browser | `TestAPIPickAndComment`, `TestAlphaLoop` |
| Feedback and v2 comparison | `TestAlphaLoop`, `TestVersionSwitcher` |
| Second MCP client, no code change | `TestAlphaLoop` |
| Nix build and flake check | `nix build .#default`, `nix flake check` |
| Full test suite | `go test ./...`, `nix develop -c go test ./...` |
| Concurrent state safety | `TestConcurrentSubmit`, `TestConcurrentReadDuringWrite` |
| No code comments | `TestNoCodeComments` |
