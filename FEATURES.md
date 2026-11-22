# Features

Honest inventory of what `go-breaker` does, by status. Code is the source of
truth — every entry cites where it lives. Status vocabulary: `FULLY_FUNCTIONAL`
(works, exercised by passing tests), `PARTIALLY_FUNCTIONAL` (ships with known
gaps), `BROKEN` (present but failing), `PLANNED` (no code yet), and
`WORTH_CONSIDERING` (idea, not committed — lives in
[ROADMAP.md](ROADMAP.md)).

_Test status: `go test ./... -race` is green; statement coverage is 98.4%
(`go test -cover ./...` — the `coverage: 98.4% of statements` line; the
uncovered remainder is the no-op branch of `transition` when the target state
equals the current one, `breaker.go:266`). CI enforces a 95% coverage floor
on every push (`.github/workflows/ci.yml`), so this number cannot rot
silently._

**Created:** 2026-11-22 · Crush (glm-5.3), docs-health AUDIT · backlog basis:
full code walk of `breaker.go` + `breaker_test.go`, coverage run above.

---

## FULLY_FUNCTIONAL

### Circuit core

- **Guarded execution (`Execute`)** — refuses cheaply with `ErrOpen` while the
  circuit is open (fn never called), otherwise calls fn exactly once and
  returns its error unmodified — recorded, never wrapped. `breaker.go`
  (`Execute`); `TestExecutePassesThroughWhenClosed`,
  `TestOpenRefusesCheaply`, `TestExecuteReturnsFnErrorUnwrapped`.
- **Separated verdict API (`Allow` + `Record`)** — for callers where the
  attempt and its outcome are separate steps; `Allow` moves an expired open
  circuit to half-open, `Record` classifies the outcome. `breaker.go`
  (`Allow`, `Record`); `TestResetTimeoutMovesToHalfOpen`,
  `TestHalfOpenFailureReopensImmediately`.
- **Closed → Open → HalfOpen → Closed machine** — `FailureThreshold`
  consecutive trip-worthy failures open the circuit; after `ResetTimeout` the
  next `Allow` probes half-open; `SuccessThreshold` consecutive half-open
  successes close it; one trip-worthy failure in half-open re-opens
  immediately. `breaker.go` (`recordSuccess`, `recordFailure`);
  `TestThresholdConsecutiveFailuresOpen`, `TestHalfOpenClosesAfterSuccesses`.
- **Consecutive-counter semantics** — a success in the closed state resets
  the failure count (bursty partial failure does not trip); counters reset on
  every transition. `breaker.go` (`recordSuccess`, `recordFailure`);
  `TestSuccessResetsFailureCount`.
- **Lazy timing, no background anything** — cooldown is evaluated per call
  against the injected clock; no timers, no goroutines, no background state.
  `State()` is a non-mutating view (reports `half_open` once the cooldown
  elapsed, without firing the transition). `breaker.go` (`currentState`,
  `Allow`); `TestResetTimeoutMovesToHalfOpen`.

### Error-family awareness

- **Family-aware default trip predicate** — only Transient and Infrastructure
  failures count toward opening; Rejection/Conflict (deterministic) failures
  never trip. `breaker.go` (`tripsByFamily`);
  `TestDeterministicFamiliesNeverTrip`, `TestInfrastructureAndTransientTrip`.
- **Unclassified errors trip (fail-open)** — plain errors classify Transient
  via go-error-family's fail-open rule and therefore count.
  `TestPlainErrorTripsAsTransient`.
- **`ErrOpen` is Transient by design** — an open circuit is temporary, so
  outer retry loops back off instead of aborting. `breaker.go` (`ErrOpen`);
  `TestErrOpenIsTransientFamily`.
- **Custom trip predicates (`Config.TripOn`)** — override the family default
  for domain-specific tripping (e.g. dnsblockd's always-trip restart gate).
  `breaker.go` (`Config.TripOn`).

### Observability and testability

- **State-change observer (`Config.OnTransition`)** — fired on every
  transition with from/to states, serialized in transition order under the
  breaker's lock; must stay cheap and never call back into the same breaker.
  `breaker.go` (`transition`); `TestOnTransitionFiresInOrderAcrossLifecycle`,
  `TestOnTransitionHalfOpenFailureFiresReopen`.
- **Injectable clock (`Config.Now`)** — deterministic time in tests; the
  entire suite runs without a single sleep. `breaker.go` (`Config.Now`);
  every timing test in `breaker_test.go`.
- **Concurrency safety** — all state behind one mutex; exercised under
  `-race`. `TestConcurrentUseStaysConsistent`.

## WORTH_CONSIDERING

Raw ideas with no code — demand-gated, refined in [ROADMAP.md](ROADMAP.md):

- Sliding-window (rather than consecutive-only) failure counting
- Per-breaker name/label for metrics identification
- Context-aware `ExecuteCtx` variant
- Metrics adapters (e.g. Prometheus) riding `OnTransition`, shipped as
  separate modules to keep the core dependency-light
