# Current Mission
update go.mod to use the latest commit on `main` for `github.com/orange-cloudavenue/cloudavenue-sdk-go`

## Plan
1. [in_progress] Inspect the current dependency state and determine whether a local replace is still active.
2. [pending] Update `go.mod`/`go.sum` to the latest SDK commit on `main`.
3. [pending] Review and summarize the result.

## Active Task
inspect current module dependency state before applying the update

### Sub-tasks
- [x] previous changelog mission completed and closed
- [ ] determine current `require` version for `cloudavenue-sdk-go`
- [ ] determine whether a local `replace` directive is still present
- [ ] determine the latest commit on SDK `main` and the resulting Go pseudo-version

### Files Being Modified
- `.opencode/scratchpad.md` only for now

### Context for Resume
- user explicitly asked: the `cloudavenue-sdk-go` is updated, use the latest commit on `main` and update `go.mod`
- likely work includes updating `go.mod` and `go.sum`, and possibly removing a local `replace` if it would prevent using the remote main commit
- next step is to inspect the current dependency declaration and the SDK upstream state

## Agent Results
- prior mission completed: added changelog fragment for issue `#1220`

## Decisions
- none yet beyond following the user's request to track the latest SDK `main` commit in `go.mod`

## Open Questions
- is a local `replace` still present and should it be removed as part of the update
- what exact pseudo-version corresponds to the latest SDK `main` commit

## Parked Scopes
- verify official upstream state of `core_api` versus local modifications
