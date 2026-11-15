# MCP tools

`tuimock mcp --dir DIR` runs an MCP server over stdio. It exposes exactly four
tools and never calls an LLM.

| Tool | Input | Output |
| --- | --- | --- |
| `get_brief` | none | Design system, references, brief, format spec and a valid example. |
| `submit_mockup` | one grid | `ok`, ID and version, or validation errors. |
| `get_feedback` | optional `screen` | Picks and comments grouped by variant. |
| `list_mockups` | optional `screen` | IDs, variants, versions, sizes and rationales. |

`screen` is optional on inputs and defaults to the screen named in
`design.yaml`.

## get_brief

Input: `{}`

The response is one payload:

```json
{
  "screen": "list-view",
  "sizes": ["100x30", "80x24"],
  "palette": {"bg": "#1a1b26", "fg": "#c0caf5"},
  "styles": {
    "n": {"color": "#c0caf5"},
    "s": {"color": "#c0caf5", "background": "#283457", "bold": true}
  },
  "borders": ["rounded"],
  "symbols": ["▸", "●", "○", "↑", "↓", "…"],
  "rules": ["One accent colour; warn only for things needing attention"],
  "brief": "The text of brief.md ...",
  "references": [
    {
      "name": "lazygit",
      "text": "╭─ lazygit ─────╮ ...",
      "note": "Dense but calm: muted borders, one accent colour."
    }
  ],
  "format": {
    "shape": "text and style are parallel arrays ...",
    "rules": ["text has exactly rows entries", "..."],
    "example": {
      "cols": 80,
      "rows": 24,
      "text": ["╭─ tuimock example ──╮", "..."],
      "style": ["nnnnnnnnnnnnnnnnnnnn", "..."],
      "rationale": "reference frame at the smallest allowed size"
    }
  },
  "instruction": "Use only style letters and characters from this design system ..."
}
```

The example is generated from the design system and validates against it. The
instruction tells the agent to cite a reference in each `rationale` and to fix
every validation error.

`get_brief` does not create `.tuimock/`.

## submit_mockup

Input:

```json
{
  "screen": "list-view",
  "variant": "flat",
  "cols": 100,
  "rows": 30,
  "text": ["╭─ beandex ──╮", "│ ~/code     │", "╰────────────╯"],
  "style": ["tttttttttttttt", "bnnnnnnnnnnnnn", "bbbbbbbbbbbbbb"],
  "rationale": "flat table with group headers, borrowed from lazygit"
}
```

The grid in this example is shortened for readability. A real request has
exactly `rows` text entries and exactly `cols` style letters per row.

`screen` may be omitted. Success:

```json
{"ok": true, "id": "list-view-flat-4f2a9c1e", "version": 1}
```

Failure returns a tool error (`isError: true`) with one line per problem:

```
row 7: text is 101 cells wide, expected 100
row 7: style has 101 letters, expected 100
row 1, column 4: style letter "x" is not defined in the design system
```

The same list is available in the structured result as `errors`. The agent
fixes the grid and calls the tool again. Until a submission passes validation,
nothing is stored, so the version counter does not move.

## get_feedback

Input: `{}` or `{"screen": "list-view"}`

```json
{
  "screen": "list-view",
  "variants": [
    {
      "variant": "flat",
      "latest_id": "list-view-flat-9c2b71da",
      "latest_version": 2,
      "picked": true,
      "comments": [
        {
          "id": "cmt-3b71d0a4",
          "mockup_id": "list-view-flat-4f2a9c1e",
          "variant": "flat",
          "version": 1,
          "text": "keep the group headers",
          "created_at": "2026-11-15T20:14:02Z"
        }
      ]
    }
  ]
}
```

An empty project returns `{"screen": "list-view", "variants": []}`.

## list_mockups

Input: `{}` or `{"screen": "list-view"}`

```json
{
  "screen": "list-view",
  "mockups": [
    {
      "id": "list-view-flat-4f2a9c1e",
      "variant": "flat",
      "version": 1,
      "cols": 100,
      "rows": 30,
      "rationale": "flat table with group headers"
    }
  ]
}
```

## Client setup

Claude Code:

```sh
claude mcp add tuimock -- /absolute/path/to/tuimock mcp --dir /absolute/path/to/project
```

Generic stdio client:

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

## Protocol notes

- stdout carries JSON-RPC only. Logs and errors go to stderr.
- An unknown tool returns a protocol error. The session stays usable.
- Tool inputs are validated against their JSON Schema before the handler runs.
- The server has no authentication and is meant for a local project folder.
