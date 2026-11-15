# Tasks

## 1. Server and brief tool

- [x] 1.1 Implement the stdio server with server name and version and verify a client can connect
- [x] 1.2 Build the `get_brief` payload from design, refs and brief and verify it includes every reference name and note
- [x] 1.3 Add the format specification and a valid example to the brief payload and verify the example validates
- [x] 1.4 Verify `get_brief` creates no state file

## 2. Submission tool

- [x] 2.1 Implement `submit_mockup` input and output types with schema tags and verify tool discovery shows the required fields
- [x] 2.2 Wire validation and state versioning and verify v1 and v2 are returned
- [x] 2.3 Return ordered validation lines as a tool error and verify row, measured and expected widths appear
- [x] 2.4 Verify a rejected submission leaves the inventory unchanged

## 3. Feedback and inventory tools

- [x] 3.1 Implement `get_feedback` with optional screen and verify picks and comments are grouped by variant
- [x] 3.2 Implement `list_mockups` with optional screen and verify sizes and rationales are present
- [x] 3.3 Verify an empty project returns valid empty documents

## 4. Protocol hygiene

- [x] 4.1 Route server diagnostics to stderr and verify a stdio client completes a full session with stdout reserved for JSON-RPC
- [x] 4.2 Verify an unknown tool returns a protocol error and the server stays usable
- [x] 4.3 Document the tool contract in `docs/mcp.md` and verify the documented request fields match the schema

## Workflow follow-up

- Archive the change after all tasks and quality gates pass.
