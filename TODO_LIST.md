# Todo List

Short-term, **actionable** open work for `go-breaker`. Each item is bounded and
cites its evidence. This file lists open work only — completed items move to
`CHANGELOG.md`; long-term/unbounded ideas live in `ROADMAP.md`; questions that
need a human decision live in `ROADMAP.md` → Open questions.

Priority: **P1** = high impact, do first; **P2** = valuable, not blocking;
**P3** = polish. The **Source** column is the trail back to where a task was
recorded. Seeded 2026-11-22 from the buildflow findings-gate baseline;
extended same day by harvest from the 2026-11-22 status reports (T4-T12,
then T13-T18 from the docs-health audit report). T8 struck 18:46 — errata:
premise dead (buildflow auto-regenerates reports/coverage.out, observed).

## Open work

| #   | Task                                                                                                                                | P  | Evidence / why                                                                                                                                | Source                       |
| --- | ----------------------------------------------------------------------------------------------------------------------------------- | -- | --------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------- |
| T19 | Archive both fully-annotated 2026-11-22 status reports to `docs/status/archived/` and add the README manifest (go-retry convention) | P3 | Every §f row in both reports is struck with hashes; they are settled history living in the active directory                                   | 2026-11-22 23-43 report §f.1 |
| T20 | Add plain `dprint check` to the standing session battery, next to `buildflow`                                                       | P3 | BuildFlow's dprint step skips dot-directories; CI checks them — a tree can pass locally and fail CI (hit 2026-11-22, `.config/metadata.yaml`) | 2026-11-22 23-43 report §f.2 |
| T21 | actionlint-in-CI decision: tools module (go-retry pattern), pinned `go install`, or record a deliberate non-fix                     | P3 | go-retry lints its workflows in CI; this repo checks them locally only (`actionlint` on PATH) — a small coverage gap vs the sibling           | 2026-11-22 23-43 report §f.3 |
| T22 | First PR-lane CI validation: confirm the `pull_request` trigger runs green (scratch PR or the first dependabot bump PR)             | P3 | Only `push` events have ever exercised the workflow; the PR lane, dispatch, and dependabot-PR paths are untested                              | 2026-11-22 23-43 report §f.4 |

Seeded 2026-11-22 23:43 by harvest from the gate-to-green execution report
(the previous eight rows T1-T17 closed the same day — verdict trails in
`CHANGELOG.md`).
