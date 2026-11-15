# State

All state lives in one JSON file inside the project folder:

```
.tuimock/state.json
.tuimock/state.json.lock
```

The lock file exists so the MCP server and the gallery can run as separate
processes and still share state safely.

## Schema

```json
{
  "version": 1,
  "mockups": [
    {
      "id": "list-view-flat-4f2a9c1e",
      "screen": "list-view",
      "variant": "flat",
      "version": 1,
      "cols": 100,
      "rows": 30,
      "text": ["╭─ beandex ──╮", "│ ~/code     │", "╰────────────╯"],
      "style": ["tttttttttttttt", "bnnnnnnnnnnnnn", "bbbbbbbbbbbbbb"],
      "rationale": "flat table with group headers, borrowed from lazygit",
      "created_at": "2026-11-15T20:12:41Z"
    }
  ],
  "picks": {
    "list-view": "list-view-flat-4f2a9c1e"
  },
  "comments": [
    {
      "id": "cmt-3b71d0a4",
      "mockup_id": "list-view-flat-4f2a9c1e",
      "screen": "list-view",
      "variant": "flat",
      "version": 1,
      "text": "keep the group headers",
      "created_at": "2026-11-15T20:14:02Z"
    }
  ]
}
```

| Field | Meaning |
| --- | --- |
| `version` | Schema version. Only `1` is accepted. |
| `mockups` | Append-only list of submitted grids. |
| `picks` | One mockup ID per screen. |
| `comments` | Append-only comments with mockup, variant and version metadata. |

`text` and `style` are parallel arrays: one string per row, one style letter per
cell. The grid fields are flat on the mockup object so the file stays readable
and easy to inspect with `jq`.

## Ordering

`list_mockups` and the gallery sort by variant, then version. The order makes
v1 appear before v2 for the same variant.

`get_feedback` groups comments by variant in append order. A comment written
against v1 stays attached to v1 after v2 arrives.

## Version numbers

A submission gets version 1 plus the highest existing version for the same
screen and variant. Variants are independent: `flat` v2 and `compact` v1 can
exist side by side.

IDs use the form `<screen>-<variant>-<8 hex>`. The suffix is random and unique
per submission, so two submissions in the same second do not collide.

## Picks

A pick replaces the previous pick for that screen. Picking a mockup ID that does
not exist returns an error and leaves the stored pick unchanged. Feedback marks
the variant as picked and reports the latest mockup ID and version for that
variant.

## Comments

Comment text is trimmed. Empty or whitespace-only text is rejected. A comment
can target any stored mockup, including an older version.

## Concurrency and durability

- Reads take a shared lock; writes take an exclusive lock.
- Writes serialise to a temporary file in `.tuimock/` and rename it over the
  state file. A crash or a failed update leaves the previous file intact.
- A missing state file reads as an empty state and is not created until the
  first write.
- An empty file reads as an empty state.
- An unknown schema version is an error rather than a silent reset.

Tests cover parallel submissions, a read that arrives during a write, a failed
update that keeps the previous bytes, pick replacement and unknown picks.

## Reset

Delete `.tuimock/` to start over. The gallery and MCP server recreate it on the
next submission.
