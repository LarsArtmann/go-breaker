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
- ~~CI: GitHub Actions running buildflow — today the quality gate is
  local-only (go-retry runs CI; this repo has none)~~ done 2026-11-22 at
  `4786db3`/`efa78f9` (owner-approved; go-retry's shape with SHA-pinned
  actions replicating the buildflow verdict — race+govulncheck, 95%
  coverage floor, dprint+golangci-lint; buildflow stays the local
  orchestrator since it is nix-hosted, not Actions-installable)
- ~~Verify pkg.go.dev rendering for v0.2.0 (go-release skill owns the
  ritual)~~ done 2026-11-22: proxy lists both tags; pkg.go.dev indexes v0.2.0
  (published Nov 21); License UNKNOWN is accurate for the Proprietary
  LICENSE, which also means pkg.go.dev shows no rendered docs by policy —
  `go doc` verifies locally instead (package doc, state machine, and all
  Config fields render clean); "Imported by: 0" reflects private consumers
- ~~CONTRIBUTING.md: pointer to AGENTS.md for AI-assisted contributions~~
  done 2026-11-22 at `efa78f9` (pointer plus the fleet BuildFlow-vs-raw-tools
  workflow note; the lint-command divergence closed — `.golangci.yml` is
  in-repo so external `golangci-lint run` matches CI)
- ~~Lychee/link check of README consumer links~~ done 2026-11-22 (all 7
  URLs across the living docs answered 200 — go-retry, go-error-family,
  shields badge, advisories, securitymd, keepachangelog, semver; all 4
  relative doc links resolve; zero rot, no tool adoption warranted)
- ~~Dependabot: decide a major-version grouping policy~~ done 2026-11-22
  at `efa78f9` (majors stay ungrouped — breaking changes deserve individual
  PRs, matching go-retry; a `github-actions` ecosystem entry was added for
  the new CI's pinned actions)

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

Three pending since 2026-11-22 23:43 (from the gate-to-green execution
report):

1. **Dprint pin ownership** — CI here checks with 0.64.3 (matches the local
   toolchain that formats the tree); go-retry pins 0.57.4. Keep the split,
   align the fleet on 0.64.3, or revert here and reformat to the older
   opinion? This also decides whether commit `b0037d5` stands as-is.
2. **Fleet fixer-gotcha ledger** — the dupword-auto-fix-deletes-diagram-words
   lesson is pinned in cqrs-htmx's agents-notes; it was hit again here
   regardless. Commit it to crush-config `references/lessons.md` (loaded by
   every session; needs a commit to that repo), or is per-repo AGENTS.md
   placement enough?

Both earlier questions were resolved on 2026-11-22: the license is
**Proprietary** (owner commit `dad3511` corrected the README badge to match
`LICENSE`; pkg.go.dev's UNKNOWN rendering is accurate), and commit
ownership is **manual on request** — the owner commits directly (all
commits in this repo's history are owner-authored; no pma daemon pickup
observed), and agents commit only on an explicit user word.
