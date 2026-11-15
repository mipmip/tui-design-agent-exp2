# Tasks

## 1. Harness

- [x] 1.1 Add `e2e/TestMain` that builds the binary into a temp dir and verify the path exists
- [x] 1.2 Add a project-copy helper and verify two tests get separate `.tuimock` directories
- [x] 1.3 Include stderr in connection failures and response bodies in HTTP failures

## 2. MCP and gallery coverage

- [x] 2.1 Test `get_brief` over stdio and verify design, refs and brief are present
- [x] 2.2 Test four valid submissions and verify four versions in the gallery payload
- [x] 2.3 Test a one-cell-too-wide submission, fix it and verify the retry succeeds
- [x] 2.4 Test pick and comment over HTTP and verify `get_feedback` returns both
- [x] 2.5 Test a v2 submission and verify v1 and v2 are both listed
- [x] 2.6 Test a second MCP client and verify the same inventory
- [x] 2.7 Test that the repository root gains no `.tuimock` directory

## 3. Quality gate and handoff

- [x] 3.1 Add `doCheck = true` to the Nix package and verify `nix flake check` runs the tests
- [x] 3.2 Add `docs/testing.md` with the gate commands and a map from test to alpha criterion
- [x] 3.3 Fill the evidence column in `docs/alpha-checklist.md` and verify every item names a test
- [x] 3.4 Run the full gate on a clean tree and record the result in the README

## Workflow follow-up

- Archive the change after all tasks and quality gates pass.
