# Tasks

## 1. State model

- [x] 1.1 Define the state, mockup, pick and comment structs with JSON tags and verify a round trip test
- [x] 1.2 Implement schema-version validation and verify an unknown version is rejected
- [x] 1.3 Implement empty-state reads and verify no state file is created

## 2. Locking and persistence

- [x] 2.1 Implement the shared read lock and verify a concurrent read during a write returns valid JSON
- [x] 2.2 Implement the exclusive write lock with temp-file rename and verify the previous file survives a failed serialisation
- [x] 2.3 Implement twenty parallel submissions and verify all twenty are present with unique ids

## 3. Versioning, picks and comments

- [x] 3.1 Implement screen-and-variant version numbering and verify v1 and v2 coexist
- [x] 3.2 Implement pick replacement and verify the previous pick is replaced
- [x] 3.3 Reject a pick for an unknown id and verify the stored pick is unchanged
- [x] 3.4 Implement comment append with version metadata and verify empty text is rejected
- [x] 3.5 Implement a per-screen feedback projection and verify it groups comments by variant and version
- [x] 3.6 Document the state schema in `docs/state.md` and verify it matches a state file produced by the tests

## Workflow follow-up

- Archive the change after all tasks and quality gates pass.
