# Design

## Context

The gallery is a reviewer tool, not a product for end users. It runs on the same
machine as the MCP server and reads the same project folder. See `proposal.md`
for motivation.

## Goals / Non-Goals

**Goals:**

- One embedded page with no build step and no external assets.
- Correct cell colours and predictable character metrics.
- Writes go through the same locked state path as the MCP server.

**Non-Goals:**

- Editing mockup text in the browser.
- Multiple screens in one view.
- Server-sent events. Polling is explicitly allowed for the PoC.

## Decisions

- **`embed.FS` for HTML, CSS and JS.** One binary remains the deliverable, and
  `nix build` needs no asset pipeline.
- **Server-rendered cell spans.** Go renders one escaped `<span>` per cell into
  the fragment. The agent's text is never interpreted as HTML, which avoids an
  injection class and keeps the browser script small. The script toggles
  versions and swaps the fragment after a poll.
- **Styles resolved server-side into CSS declarations.** The client receives
  `letter -> {color, background, bold, italic, underline}`. Palette resolution
  stays in Go and the browser never re-implements YAML semantics.
- **Polling every second with `fetch` and a whole-fragment swap.** The page keeps
  the previous fragment text and replaces `#gallery` only when the response
  changed. Comment drafts live in a map keyed by variant so refreshes preserve
  them, and focus and selection are restored after the swap.
- **Pick and comment over POST JSON.** Form encoding would need CSRF thoughts;
  a same-origin JSON endpoint is smaller and testable with `curl`.
- **`Cache-Control: no-store` on the API.** The state is live; caching would
  hide submissions.

## Risks / Trade-offs

- [Polling costs CPU on an idle page] -> one request per second on loopback is
  negligible; the page pauses polling when the tab is hidden.
- [Rebuilding a card loses scroll position] -> the horizontal scroller is on the
  container, not the card; cards keep a fixed width per version.
- [Font differences change the advance width] -> the page uses a bundled
  monospace stack and sets explicit `font-size`, `line-height` and letter
  spacing; the e2e test asserts cell count rather than pixel width.

## Migration Plan

Greenfield. Archive creates the `gallery` main spec.
