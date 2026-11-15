# Tasks

## 1. HTTP server and API

- [x] 1.1 Implement the loopback listener and static handler and verify `/` returns the page
- [x] 1.2 Implement `GET /api/state` with design resolution and verify the payload shape
- [x] 1.3 Implement `POST /api/pick` and verify a valid pick and an unknown id with 400
- [x] 1.4 Implement `POST /api/comment` and verify empty text is rejected
- [x] 1.5 Set `no-store` and verify the response headers
- [x] 1.6 Verify the handler serves no project files such as `design.yaml` or `refs/`

## 2. Gallery page

- [x] 2.1 Render mockups side by side with variant, version, size and rationale and verify four cards in rendered HTML
- [x] 2.2 Render every cell with resolved colour and background and verify a space in a selected row keeps its background
- [x] 2.3 Implement the version switcher and verify v1 and v2 show different grids
- [x] 2.4 Implement Pick and comment controls and verify the API calls and the picked marker
- [x] 2.5 Implement one-second polling with draft preservation and verify the page contains the interval and draft map

## 3. Visual checks and docs

- [x] 3.1 Add a golden HTML fixture for one 4x2 mockup and verify it matches the renderer output
- [x] 3.2 Document the gallery and its API in `docs/gallery.md` and verify the endpoints through the e2e tests

## Workflow follow-up

- Archive the change after all tasks and quality gates pass.
