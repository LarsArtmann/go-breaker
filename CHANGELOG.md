# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versioning follows [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Fleet-standard documentation set: `AGENTS.md` (session context — commands, design contracts, conventions, gotchas), `TODO_LIST.md` (evidence-cited open work), `ROADMAP.md` (raw ideas, non-goals, open questions), `FEATURES.md` (honest status inventory; 98.4% statement coverage cited with its command), and `docs/status/` opened with the 2026-11-22 bootstrap report. `docs/DOMAIN_LANGUAGE.md` extended with `TripOn`, `Probe`, `OnTransition`, and `Clock` terms plus usage rules.
- README: Installation section (Go 1.27+ floor, `go get` command).

### Changed

- Documentation overhauled by the 2026-11-22 docs-health audit: `docs/DOMAIN_LANGUAGE.md` gained the `TripOn`, `Probe`, `OnTransition`, and `Clock` terms plus language rules; `AGENTS.md` rewritten to pass the endurance test (dated baseline section replaced by enduring quality-gate gotchas); `ROADMAP.md` gained non-goals and shed its two answered open questions; completed TODO rows are purged mechanically via `docs-gates purge` (verdict lines below).

- TODO_LIST purge (2026-11-22): | T3  | Add an Installation section to README.md (explicit `go get` step before Usage)                                                                                | P3 | go-structure-linter info finding: "README.md lacks an installation section"; import path is shown but the install command is not | buildflow gate 2026-11-22   | — done (audit 2026-11-22: Installation section added — Go 1.27+ floor + `go get`)
- TODO_LIST purge (2026-11-22): | T4  | Append a `## [Unreleased]` CHANGELOG entry recording the AGENTS.md + TODO_LIST.md doc bootstrap                                                                 | P2 | CHANGELOG [Unreleased] still says "Nothing yet." while two doc files landed this session                                        | 2026-11-22 report §f.4      | — done (audit 2026-11-22: Unreleased Added/Fixed entries written)
- TODO_LIST purge (2026-11-22): | T5  | Run the on-demand gates over the new docs: `buildflow -s markdown-lint`, then `-s codespell`                                                                    | P2 | markdown-lint and codespell never run in pipeline modes; the new docs' prose is checked by no gate                                | 2026-11-22 report §f.5      | — done (audit 2026-11-22: codespell 0 findings; markdown-lint not applicable — repo has no markdownlint config, a no-op by prerequisite)
- TODO_LIST purge (2026-11-22): | T6  | Reconcile the findings-gate count attribution (gate said 4 findings, `--format finding` listed 3 messages); capture buildflow's true exit code without a pipe     | P3 | Discrepancy hand-waved as "attribution"; exit code was measured through `tail` (wrong value shown)                               | 2026-11-22 report §d.1/§d.3 | — done (audit 2026-11-22: true exit code captured unpiped = 69 with 2 error findings; count consistent since run 2 — run-1 4-vs-3 no longer reproducible, evidence discarded)
- TODO_LIST purge (2026-11-22): | T7  | Resolve README "## Zero dependencies" heading vs "One dependency by design" body with the owner, then apply the wording                                         | P3 | Heading contradicts the sentence directly under it on the consumer-facing page; wording is an owner decision (report question 2) | 2026-11-22 report §f.7      | — done (audit 2026-11-22: heading fixed to "One dependency, by design" per owner's truthful-docs instruction)
- TODO_LIST purge (2026-11-22): | T10 | When T1/T2 land: shrink the AGENTS.md "Known gate baseline" section to a TODO_LIST pointer                                                                    | P2 | The two findings would then live in three places (gate, TODO_LIST, AGENTS.md restatement) — a split brain that drifts on landing | 2026-11-22 report §e        | — done (audit 2026-11-22: executed early — section rewritten to enduring gotchas + TODO_LIST pointer, dates/counts removed)
- TODO_LIST purge (2026-11-22): | T11 | Read SECURITY.md / AUTHORS / LICENSE; fold any non-obvious supported-versions fact into AGENTS.md                                                             | P3 | Discovery gap: never opened this session; SECURITY.md may state a support policy worth surfacing                                 | 2026-11-22 report §c.5      | — done (audit 2026-11-22: SECURITY.md corrected — advisory URL named go-must, version v0.1.2 never existed, CI claim false; LICENSE/README MIT contradiction routed to ROADMAP Open questions)
- TODO_LIST purge (2026-11-22): | T12 | Document the go-structure-linter internal/-dir info finding as a deliberate non-fix in AGENTS.md (public single-package library)                              | P3 | Info finding recurs every gate run; AGENTS.md mentions it in passing but not as a verdict                                        | 2026-11-22 report §f.11     | — done (audit 2026-11-22: non-fix verdict recorded in AGENTS.md quality-gate gotchas)

- TODO_LIST purge (2026-11-22): | T18 | Observe whether the pma auto-commit daemon picks up the 2026-11-22 doc set (commit-ownership question 2 adjacent)                                  | P3 | buildflow skill warns of a daemon blind spot; ten files were uncommitted at 18:12 — unverified either way                                        | 2026-11-22 audit report §f.12 | — done (2026-11-22: commit policy resolution by evidence — owner commits directly, dad3511 landed on top of agent work already pushed, all repo history owner-authored, no daemon pickup observed; agents commit on explicit request)

### Fixed

- README license badge corrected from MIT to Proprietary, matching the actual `LICENSE` (owner commit `dad3511`, 2026-11-22 — pkg.go.dev's "UNKNOWN" license rendering is accurate; closes the README-vs-LICENSE contradiction the audit flagged).
- README: the "Zero dependencies" heading contradicted the body's "one dependency by design" — the heading now tells the truth.
- SECURITY.md: drift from the shared hygiene config — the vulnerability advisory URL pointed at `go-must` (wrong repository), the supported-versions table cited `v0.1.2` (never released here), and the practices section claimed CI security scanning (this repository has no CI). All corrected for go-breaker.

## [0.2.0] - 2026-11-21

### Added

- `Config.OnTransition func(from, to State)`: observer fired on every state change (closed→open, open→half_open, half_open→closed, half_open→open), serialized in transition order under the breaker's lock. Observers must be cheap and must never call back into the same Breaker. Grew from the first consumer migration (file-and-image-renamer's resilience metrics feed). Tests: lifecycle order goldens, half-open reopen sequence.

## [0.1.0] - 2026-11-21

### Added

- Initial extraction (domain-language census plan M18): error-family-aware circuit breaker — Closed/Open/HalfOpen machine, `Execute`/`Allow`+`Record`/`State` API, `ErrOpen` (Transient, so outer retry loops back off), family-aware default trip predicate (Transient + Infrastructure trip; Rejection/Conflict/Corruption never; unclassified errors inherit go-error-family's fail-open Transient), concurrency-safe, injectable clock. Tests: state-transition goldens, family table, `-race` concurrency.
