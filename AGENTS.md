# AGENTS.md — go-breaker

Concise, enduring context for AI sessions working in this repo. Change history
lives in `CHANGELOG.md`; open work in `TODO_LIST.md`; the feature inventory
in `FEATURES.md`; long-term ideas and owner questions in `ROADMAP.md`. Domain
vocabulary lives in `docs/DOMAIN_LANGUAGE.md`.

## What this is

A single-file, zero-goroutine, error-family-aware circuit breaker library —
the sibling of go-retry ("should I even try?" vs "should I try again?").
Entire implementation is `breaker.go`; entire test suite is `breaker_test.go`.
No CLI, no service, no packages, no flake.nix (nix-hermetic by triviality —
STACK.md). One dependency by design: `go-error-family` (the trip semantics
ride on its classification contract; do not add more).

## Commands

| Task             | Command                                | Notes                                                   |
| ---------------- | -------------------------------------- | ------------------------------------------------------- |
| Test (with race) | `go test ./... -race`                  | CONTRIBUTING.md's canonical command; passes clean       |
| Quality gate     | `buildflow`                            | Auto-detects this repo (no `.buildflow.yml` needed)     |
| One tool         | `buildflow -s <tool> --format finding` | e.g. `-s go-structure-linter`; `-s` does not accumulate |
| Vet              | `go vet ./...`                         | Passes                                                  |

- BuildFlow owns formatting/lint/tidy/dependabot-sync — do not run
  `golangci-lint --fix`, `go mod tidy`, or formatter wrappers by hand, and do
  not hand-edit the `buildflow-managed` block in `.gitignore`.
- Formatting split: Go via gofmt (tabs); everything else (md/json/yaml/toml)
  via `dprint.json`. `CHANGELOG.md` is dprint-excluded on purpose.
- A green-looking `buildflow` run can still exit non-zero: the findings gate
  fails when error-severity findings remain even if all steps succeeded.
  Read the summary's findings section, not just the exit code.

### Quality-gate gotchas

- The findings gate can be red on pre-existing findings before you touch
  anything — check `TODO_LIST.md` for tracked gate findings before assuming
  you broke something.
- `file-size-check` caps files at 350 lines — put new tests in a second
  `_test.go` file rather than growing `breaker_test.go`.
- The `internal/`-directory suggestion from go-structure-linter is a
  deliberate non-fix: this is a public single-package library by design
  (STACK.md).
- `docs-gates` hard-requires `TODO_LIST.md` to exist — delete or rename it
  and the docs-gates step FATALs ("reading TODO_LIST.md: no such file").
  Its freshness gate also requires every open item in a `docs/status/`
  report to be represented in a living doc (TODO_LIST/ROADMAP/CHANGELOG).
- `SECURITY.md` is securitymd-generated and drifted once (wrong repo URL,
  nonexistent version, false CI claim) — re-verify its repo-specific fields
  after any shared-hygiene regeneration.

## Design contracts (breaking any of these is a breaking change)

- **Lazy state machine, no background anything.** Open→HalfOpen happens only
  when a call observes the elapsed cooldown. No timers, no goroutines —
  `Config.Now` is the only clock.
- **`State()` is a view, `Allow()` is the commit.** After the cooldown
  elapses, `State()` already reports `half_open`, but the internal state
  stays `open` (and `OnTransition` has NOT fired) until the next `Allow()`
  performs the transition. Tests rely on this.
- **`ErrOpen` is Transient on purpose.** Outer retry loops must back off, not
  abort. Changing its family breaks the fleet contract (see
  `TestErrOpenIsTransientFamily`).
- **Default trip predicate = error families.** Transient + Infrastructure
  trip; Rejection/Conflict/Corruption never do; unclassified errors inherit
  go-error-family's fail-open Transient and therefore trip.
- **`Execute` returns fn's error unmodified** — recorded, never wrapped.
- **HalfOpen does not limit concurrent probes.** Every `Allow()` passes in
  half-open (unlike sony/gobreaker's single-probe semantics); recovery is
  proven by `SuccessThreshold` consecutive successes. Consumers migrating
  from gobreaker may expect single-probe — that is not this library.
- **`OnTransition` runs under the breaker's lock**, serialized in transition
  order. Observers must be cheap, non-blocking, and must never call back
  into the same Breaker (deadlock).

## Conventions

- Go 1.27 (`go.mod`); the code uses `for range N`, `wg.Go`, table tests with
  `t.Run`. Older toolchains will not compile.
- Stdlib `testing` only — no testify. Error fixtures are tiny constructors
  (`transientErr()`, `rejectionErr()`, …) with namespaced error-family codes
  (`test.transient`, `breaker.open`).
- Deterministic time in tests: capture `now` in a variable, inject
  `Config.Now: func() time.Time { return now }`, then advance `now` — no
  sleeps anywhere in the suite.
- Doc comments are the public API docs (pkg.go.dev) — package comment in
  `breaker.go` carries the state-machine diagram; keep it in sync.
- Error codes are dotted lowercase namespaced (`breaker.open`).
- Tabs in Go, 2-space in yml/json/toml/nix (`.editorconfig`).

## Files that must move together

- `README.md` — consumer-facing; its **Consumers table** lists fleet adopters
  (wire-or-delete migrations). Update it when another repo adopts/drops this
  library; git history shows dedicated "consumer table" commits for this.
- `CHANGELOG.md` — Keep-a-Changelog + SemVer; `## [Unreleased]` section on
  top. Releases are tagged (`v0.1.0`, `v0.2.0` exist) — follow the go-release
  skill for cutting versions.
- `FEATURES.md` — honest feature inventory by status; never round up (a
  status without evidence is a guess); when a feature ships, update the row
  and delete the matching TODO row (no split brains).
- `ROADMAP.md` — raw ideas, non-goals, and the owner Open questions; an idea
  graduates to `TODO_LIST.md` when it becomes bounded.
- `STACK.md` — stack decision record (why stdlib state machine, why not
  sony/gobreaker); re-derive before extending, don't append blindly.
- `docs/DOMAIN_LANGUAGE.md` — glossary (Circuit, Trip, Threshold, Cooldown,
  ErrOpen); update when introducing new domain terms.
- `.github/dependabot.yml` — weekly gomod updates; BuildFlow's
  `dependabot-auto-configure` syncs it with real modules.

## Repo facts

- Default branch is `master` (git-town config; `main` is remapped).
- No CI workflows — `.github/` contains only `dependabot.yml`.
- `reports/` is gitignored (coverage artifacts; `reports/coverage.out` may
  exist locally, untracked).
- `dprint.json` fetches its plugins over the network on first use — a cold
  run can hang; cached under `~/.cache/dprint` afterwards.
