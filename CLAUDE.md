# tuimock

`tuimock` lets any MCP-capable agent turn a screen brief into strict character-grid
TUI mockups, shown side by side in a local browser gallery for review. One Go
binary provides an MCP server (stdio) and an HTTP gallery; both share
`.tuimock/state.json`. The tool never calls an LLM.

The authoritative product briefing is `docs/briefing.md`. Read it before changing
behaviour. The alpha done-criteria live in `docs/alpha-checklist.md`.

## Non-negotiable rules

- **Beans is the only work tracker.** Never use TodoWrite or ad-hoc todo lists.
  Run `beans prime` at the start of a session. Check admin: milestone and epic
  beans are created and maintained by the agent.
- **Milestone titles start with an incremental two-digit number**: `01 ...`,
  `02 ...`, and so on. Never renumber an existing milestone.
- **OpenSpec owns planning.** Every epic maps to exactly one OpenSpec change.
  Create the change with `/opsx:propose` before writing code. Implement with
  `/opsx:apply`. Archive with `/opsx:archive` when every task is checked.
- **Commit after every archive.** Use `jj` only. Author is Pim Snel. No
  self-promotion, no `Co-authored-by`, no mention of Claude or any AI tool in
  commit messages.
- **No code comments.** Names and structure carry the meaning. This applies to
  Go, shell, Nix, HTML, CSS and JavaScript. Prose belongs in docs.
- **Nix from the start.** `flake.nix` uses plain Nix with
  `lib.genAttrs systems`. Never add `flake-utils`. Supported systems:
  `x86_64-linux`, `aarch64-linux`, `x86_64-darwin`, `aarch64-darwin`.
- **Evidence over claims.** A change is done only when `go test ./...`,
  `nix flake check` and the e2e suite pass. Never mark a task complete on a
  promise.

## Architecture

```
cmd/tuimock/            CLI entry: serve, mcp, validate, version
internal/design/        design.yaml parsing, border sets, allowed characters
internal/mockup/        strict grid validation and display-width rules
internal/state/         state.json schema, file locking, versions, picks, comments
internal/project/       project folder loading (design.yaml, brief.md, refs/)
internal/brief/         the single get_brief payload
internal/mcpserver/     four MCP tools over stdio
internal/web/           embedded gallery, JSON API, polling
internal/width/         centered display-width helpers
examples/beandex/       reference brief, design system and reference captures
e2e/                    black-box tests: real binary, real MCP, real HTTP
```

State lives in `<project>/.tuimock/state.json`. All writes take an exclusive
lock. The gallery only reads state; pick and comment go through the HTTP API
which writes under the same lock.

## Validation contract

`submit_mockup` rejects with one precise line per failure. See
`openspec/specs/mockup-validation/spec.md`. The rules are exact:

- `text` and `style` have exactly `rows` entries.
- Each `text` row has display width exactly `cols`, measured with
  `github.com/mattn/go-runewidth`, never bytes or runes.
- Each `style` row has exactly one style letter per cell.
- Every style letter exists in the design system.
- Only printable ASCII, characters from the configured border sets, and the
  configured symbols are allowed.
- Any zero-width or double-width character is rejected.
- `cols` x `rows` must be one of the sizes in `design.yaml`.

## Quality gates

Run before every commit:

```bash
go test ./...
go vet ./...
nix flake check
```

`e2e` must prove the full loop: MCP `get_brief`, four submitted variants,
gallery rendering data, pick plus comment, `get_feedback`, revised v2.

## Commit style

```
<area>: <imperative summary>

<body explaining why, wrapped at 72 columns>
```

Example: `web: add version switcher to the gallery`.

## Agent loop

1. `beans prime`; read ready beans and existing specs.
2. Pick the next epic bean. Its OpenSpec change must exist and be unarchived.
3. `/opsx:apply <change>`; work tasks one by one, checking them off.
4. Run the quality gates.
5. `/opsx:archive <change>`.
6. `jj commit` as Pim Snel.
7. Update the epic bean, then the next epic.
