# Gallery

`tuimock serve --dir DIR` starts a local gallery and prints the bound address:

```
tuimock gallery: http://127.0.0.1:7373
```

The default address is `127.0.0.1:7373`. Use `--addr 127.0.0.1:0` to let the
kernel pick a port, which is what the e2e tests do. The server shuts down on
SIGINT or SIGTERM.

## What the page shows

- All variants of the current screen side by side in one horizontal scroller.
- Each cell rendered as a span with the resolved design colour, background and
  font attributes, at true character size in a monospace font.
- A version switcher per variant. All versions are present in the page; the
  switcher hides the others.
- Variant name, version, size and rationale under each grid.
- Comments per version.
- A **Pick** button and a comment box per variant, with the picked card
  highlighted.

The page is embedded in the binary and needs no build step. It polls the state
file every second and replaces the gallery fragment only when the HTML changed.
A hidden tab skips polling. The version selection and unfinished comment text
survive a refresh.

`--dir` must point at a folder with `design.yaml` and `brief.md`. The gallery
does not need the MCP server, but both can run at once.

## HTTP API

All responses use `Cache-Control: no-store`.

| Method and path | Purpose |
| --- | --- |
| `GET /` | Gallery page. |
| `GET /fragment` | Gallery cards only, used by polling. |
| `GET /api/state` | State plus resolved design. |
| `POST /api/pick` | Set the pick for a screen. |
| `POST /api/comment` | Append a comment to one mockup. |

### State

```sh
curl -s http://127.0.0.1:7373/api/state | jq
```

```json
{
  "screen": "list-view",
  "design": {
    "sizes": ["100x30", "80x24"],
    "styles": {"n": {"color": "#c0caf5"}},
    "palette": {"fg": "#c0caf5"}
  },
  "mockups": [
    {
      "id": "list-view-flat-4f2a9c1e",
      "screen": "list-view",
      "variant": "flat",
      "version": 1,
      "cols": 100,
      "rows": 30,
      "text": ["..."],
      "style": ["..."],
      "rationale": "flat table",
      "created_at": "2026-11-15T20:12:41Z"
    }
  ],
  "picks": {"list-view": "list-view-flat-4f2a9c1e"},
  "comments": []
}
```

### Fragment

```sh
curl -s http://127.0.0.1:7373/fragment
```

Returns the cards without the page shell. An empty state returns a short
`No mockups yet` block.

### Pick

```sh
curl -s -X POST http://127.0.0.1:7373/api/pick \
  -H 'content-type: application/json' \
  -d '{"mockup_id":"list-view-flat-4f2a9c1e"}'
```

A pick replaces the previous pick for that screen.

### Comment

```sh
curl -s -X POST http://127.0.0.1:7373/api/comment \
  -H 'content-type: application/json' \
  -d '{"mockup_id":"list-view-flat-4f2a9c1e","text":"keep the group headers"}'
```

The comment is trimmed and stored against the mockup's variant and version.

## Status codes

| Code | Meaning |
| --- | --- |
| 200 | Success. |
| 400 | Invalid JSON, unknown mockup ID, empty comment, or body over 1 MB. |
| 404 | Unknown path. |
| 405 | Wrong method for a known path. |
