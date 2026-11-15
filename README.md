# tuimock

`tuimock` lets an MCP-capable agent produce TUI mockups as strict character
grids. A local gallery shows the variants side by side so you can pick one and
comment before any application code is written.

One Go binary provides two processes:

- `tuimock mcp` is an MCP server over stdio with four tools.
- `tuimock serve` is a localhost gallery that renders the same state.

Both read and write `.tuimock/state.json` with a file lock, so they can run at
the same time.

## Requirements

- Go 1.25 or newer for a plain build.
- Nix with flakes for the packaged build and the development shell.

## Quick start

```sh
nix build .#default
./result/bin/tuimock validate --dir examples/beandex
./result/bin/tuimock serve --dir examples/beandex
```

Open the printed URL, then connect an agent to the same project folder. The
`examples/beandex` folder contains a design system, four reference screens and
a brief for a repository list view.

Without Nix:

```sh
go build -o tuimock ./cmd/tuimock
./tuimock serve --dir examples/beandex
```

## Agent connection

Claude Code:

```sh
claude mcp add tuimock -- /absolute/path/to/tuimock mcp --dir /absolute/path/to/project
```

Any MCP client that supports stdio can use a JSON entry like this:

```json
{
  "mcpServers": {
    "tuimock": {
      "command": "/absolute/path/to/tuimock",
      "args": ["mcp", "--dir", "/absolute/path/to/project"]
    }
  }
}
```

The project folder must contain `design.yaml` and `brief.md`. `refs/` is
optional but recommended.

## Commands

| Command | Purpose |
| --- | --- |
| `tuimock mcp --dir DIR` | Run the MCP server on stdio. |
| `tuimock serve --dir DIR [--addr HOST:PORT]` | Run the gallery. Default address is `127.0.0.1:7373`. |
| `tuimock validate --dir DIR` | Check the design, brief and references. |
| `tuimock validate --dir DIR --mockup FILE` | Validate one mockup JSON file. |
| `tuimock version` | Print the version. |
| `tuimock help` | Print usage. |

## Project folder

```
project/
  design.yaml        palette, styles, borders, symbols, rules, sizes
  brief.md           the screen to design
  refs/
    lazygit.txt      plain-text capture
    lazygit.md       what to learn from it
```

State is written to `project/.tuimock/state.json`. Add `.tuimock/` to the
project's ignore file.

## Repository layout

```
cmd/tuimock/           CLI entry point
internal/brief/        get_brief payload and format example
internal/design/       design system parser and character rules
internal/mcpserver/    four MCP tools over stdio
internal/mockup/       grid validation
internal/project/      project folder loader
internal/state/        locked JSON state store
internal/web/          gallery server and embedded page
e2e/                   black-box tests against the built binary
examples/beandex/      example project
openspec/              change proposals, specs and archived changes
docs/                  design, state, MCP, gallery and testing notes
```

## Quality gates

Run these before calling a change done:

```sh
go test ./...
go vet ./...
nix flake check
```

`nix flake check` builds the package with `doCheck = true`, so the full test
suite runs inside the Nix build.

Verified on 2026-11-15 on x86_64-linux:

- `go test ./...` passes, including the black-box e2e suite.
- `go vet ./...` passes.
- `nix build .#default` passes and the binary reports `0.1.0`.
- `nix flake check` passes.
- `nix develop -c go test ./...` passes.

The flake declares `x86_64-linux`, `aarch64-linux`, `x86_64-darwin` and
`aarch64-darwin` with `lib.genAttrs`. Only the current system is built locally.

## Updating dependencies

`packages.default` uses `buildGoModule` with a pinned `vendorHash`. After
changing `go.mod` or `go.sum`, set a fake hash and read the real one from the
build error:

```sh
# set vendorHash to lib.fakeHash in flake.nix
nix build .#default
# copy the reported got: sha256-... into flake.nix
nix build .#default
```

## Scope

This is an alpha base. It does not render a real terminal, start a PTY,
generate Bubble Tea code, or support more than one screen or project at a
time. See `docs/alpha-checklist.md` for the acceptance list and
`docs/briefing.md` for the original experiment.

## Documentation

- `docs/design-system.md`: `design.yaml` schema and character rules.
- `docs/state.md`: state file schema and locking.
- `docs/mcp.md`: tool contracts and examples.
- `docs/gallery.md`: gallery and HTTP API.
- `docs/testing.md`: gates and test map.
