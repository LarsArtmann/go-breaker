# Roadmap — go-breaker

Long-term direction and **raw ideas not yet refined into actionable
[TODO_LIST.md](TODO_LIST.md) rows**. Graduation rule: an idea becomes a
TODO_LIST row when it is bounded and short-term. Seeded 2026-11-22 from the
status-report harvest; pruned the same day by the docs-health audit (two
questions answered, one raw idea shipped into FEATURES.md).

## Raw ideas

- Sliding-window failure counting — consecutive-only today; a window would
  smooth bursty partial-failure patterns (behavior change, needs its own
  design pass)
- Per-breaker name/label — identify breakers in `OnTransition` metrics
  without consumers wrapping them
- Context-aware variant (`ExecuteCtx`) if consumers need cancellation
  propagated into the recorded call
- Metrics adapter (e.g. Prometheus) built on `OnTransition` — core stays
  dependency-free, adapters ship as separate modules
- CI: GitHub Actions running buildflow — today the quality gate is
  local-only (go-retry runs CI; this repo has none)
- ~~Verify pkg.go.dev rendering for v0.2.0 (go-release skill owns the
  ritual)~~ done 2026-11-22: proxy lists both tags; pkg.go.dev indexes v0.2.0
  (published Nov 21); License UNKNOWN is accurate for the Proprietary
  LICENSE, which also means pkg.go.dev shows no rendered docs by policy —
  `go doc` verifies locally instead (package doc, state machine, and all
  Config fields render clean); "Imported by: 0" reflects private consumers
- CONTRIBUTING.md: pointer to AGENTS.md for AI-assisted contributions
- Lychee/link check of README consumer links
- Dependabot: decide a major-version grouping policy

## Non-goals

Deliberately excluded (decisions recorded in [STACK.md](STACK.md)):

- **Third-party breaker libraries** (sony/gobreaker, afex/hystrix-go) — no
  error-family awareness; fleet policy prefers first-party when the semantics
  are fleet-specific; hystrix-go is abandoned upstream.
- **Hosting inside go-retry** — retry and breaker compose as siblings, not
  layers; nesting forecloses independent adoption.
- **Timers/goroutines/background state** — the lazy, call-evaluated design is
  the point; lifecycle cost buys zero semantics.

## Open questions

_None pending._ Both were resolved on 2026-11-22: the license is
**Proprietary** (owner commit `dad3511` corrected the README badge to match
`LICENSE`; pkg.go.dev's UNKNOWN rendering is accurate), and commit
ownership is **manual on request** — the owner commits directly (all
commits in this repo's history are owner-authored; no pma daemon pickup
observed), and agents commit only on an explicit user word.
