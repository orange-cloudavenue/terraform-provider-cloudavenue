# AGENTS.md

## Project overview

- Terraform provider for CloudAvenue, written in Go on Terraform Plugin Framework.
- Ground truth is `Taskfile.yml`, `.taskfiles/`, and `.github/workflows/`; `README.md` and `CONTRIBUTING.md` contain some stale commands/version notes.

## Commands

- Use these exact commands:

```sh
go test ./...
```

```sh
go test ./internal/testsacc/ -list 'TestAcc'
```

```sh
task generate
```

```sh
task lint
```

```sh
task provider:install
```

```sh
task provider:test TF_TEST_NAME="TestAccIAMUserResource"
```

- `task provider:test` requires `.env` or `~/.cloud_avenue.source` with at least `CLOUDAVENUE_ORG`, `CLOUDAVENUE_USERNAME`, and `CLOUDAVENUE_PASSWORD`.
- `task provider:test` sets `TF_ACC=1` and runs `go test -v ./internal/testsacc/ -count=1 -timeout 30m -run <name>`.
- Useful acceptance-test toggles already wired into the task: `TF_ACC_RUN_TEST`, `TF_ACC_ONLY_PRINT=true`, `TF_LOG=trace|debug|info|warn|error`, `CLOUDAVENUE_DEBUG`, `NETBACKUP_DEBUG`, `S3_DEBUG`.
- `task generate` mutates the repo: `go mod tidy`, `go generate ./...`, `terraform fmt` in `examples/`, and license header fixes.
- `task lint` also mutates the repo: `golangci-lint run --fix ./...` plus license, vuln, module, workflow, and Terraform formatting checks.

## Architecture map

- `main.go` is both the provider entrypoint and the `go generate` hook for docs generation.
- `internal/provider/` is the provider surface. Start with:
  - `provider.go` for provider configuration/client wiring
  - `provider_schema.go` for provider config schema
  - `provider_resources.go` and `provider_datasources.go` for authoritative registration
- Feature code is split by domain under `internal/provider/` (`edgegw`, `elb`, `iam`, `s3`, `vapp`, `vdc`, `vdcg`, `vrf`, etc.). Keep changes local to one domain when possible.
- Typical feature file set is spread across multiple files, not one:
  - `*_schema.go`
  - `*_types.go`
  - `*_resource.go` and/or `*_datasource.go`
  - nearby `*_schema_test.go`
  - `internal/testsacc/*_<feature>_*_test.go`
  - `templates/resources/*.md.tmpl` or `templates/data-sources/*.md.tmpl`
  - `examples/resources/...` or `examples/data-sources/...`
- `docs/` is generated output. Edit schema metadata, `templates/`, or `examples/`, then regenerate.
- `.agents/skills/` contains repo-specific workflows for provider resources, docs, config, migration, and acceptance tests. Check there before inventing a new pattern.

## Conventions that deviate from defaults

- Write all code comments, commit messages, changelog entries, docs, examples, diagnostics, and user-facing strings in English.
- Keep repo text clear, concise, technical and functional. Do not add explanatory filler for obvious behavior.
- Follow existing Terraform Plugin Framework patterns in this repo. Do not introduce foreign patterns in new provider code.
- Give every user-facing schema attribute a `MarkdownDescription`; generated registry docs depend on it.
- Prefer schema validators and plan modifiers over ad hoc runtime validation when the rule is static.
- Preserve compatibility markers already present in schema, especially `DeprecationMessage` and env-var fallback behavior.
- Build resource/data source type names from `req.ProviderTypeName` plus the local `categoryName`; do not hardcode full provider type names.
- Use `tflog` with subsystem context for provider/resource logging. Match existing diagnostic wording like `Error retrieving ...` or `Failed to create ...`.
- Before adding helpers, check whether the domain package already has shared constants/converters; avoid creating a second local pattern.

## Testing guidance

- Keep unit tests next to the code they cover; many domains already use `*_schema_test.go` to validate both schema diagnostics and implementation diagnostics.
- Put acceptance coverage in `internal/testsacc/`, with names aligned to the resource or data source under test.
- Use `t.Setenv(...)` in tests instead of mutating process env manually.
- For an operator-visible behavior change, done usually means: unit tests updated, relevant acceptance coverage updated or justified, `go test ./...` clean, and `task generate` run if schema/templates/examples changed.

## Gotchas and landmines

- Trust `Taskfile.yml` and CI over prose docs. `CONTRIBUTING.md` still mentions `make generate`; the real command is `task generate`.
- Do not hand-edit `docs/` for normal work. CI runs generation and fails on drift.
- `examples/` is excluded from Go linting but still validated through `terraform fmt`; broken examples fail checks.
- If you add or touch Go files, expect SPDX/license header rewrites during `task generate` or `task lint`.
- If you touch schema, resource/data source metadata, `templates/`, or `examples/`, run `task generate` before finishing.
- Provider/test code expects `CLOUDAVENUE_USERNAME`. Some older CI/docs still say `CLOUDAVENUE_USER`.
- Provider configuration still exposes deprecated `vdc`; do not remove or repurpose it without explicit compatibility review.
- `.github/workflows/*` and `.pre-commit-config.yaml` are centrally synced from `orange-cloudavenue/workflows`; local edits are likely to be overwritten.

## Boundaries

- Always:
  - Keep changes small and local to the affected domain package.
  - Check `.agents/skills/` and use the closest repo skill before inventing a new workflow for provider work.
  - Update registration files, tests, templates, and examples when a new resource/data source needs them.
  - Add `.changelog/<PR-number>.txt` for operator-visible provider/resource/data source changes.
  - Run `go test ./...` for code changes.
  - Run `task generate` after schema/template/example changes.
  - Run `task lint` last when you need full repo verification, because it mutates files.
- Ask first:
  - Dependency changes in `go.mod`.
  - Changes to `.github/workflows/*`, `.pre-commit-config.yaml`, or other centrally synced files.
  - Public schema changes, removals, or behavior changes that may need changelog and acceptance-test updates.
  - Changes touching deprecated provider config like `vdc`.
- Never:
  - Hand-edit generated docs in `docs/`.
  - Invent commands from `README.md` or `CONTRIBUTING.md` without checking `Taskfile.yml`/CI.
  - Commit secrets from `.env`, test logs, or local binaries.

## Git/PR workflow

- Use conventional commits. Commit-msg hooks enforce them strictly.
- Add `.changelog/<PR-number>.txt` for every operator-visible provider/resource/data source change before opening or updating a PR. Skip changelog entries for doc-only, test-only, or refactor-only changes.
- History shows common automated commit subjects; do not mimic them unless that is the actual work:
  - `chore: Update CHANGELOG.md`
  - `chore: synced file(s) with orange-cloudavenue/workflows`
