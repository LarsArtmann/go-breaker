# Todo List

Short-term, **actionable** open work for `go-breaker`. Each item is bounded and
cites its evidence. This file lists open work only — completed items move to
`CHANGELOG.md`; long-term/unbounded ideas live in `ROADMAP.md`; questions that
need a human decision live in `ROADMAP.md` → Open questions.

Priority: **P1** = high impact, do first; **P2** = valuable, not blocking;
**P3** = polish. The **Source** column is the trail back to where a task was
recorded. Seeded 2026-11-22 from the buildflow findings-gate baseline;
extended same day by harvest from the 2026-11-22 status reports (T4-T12,
then T13-T18 from the docs-health audit report).

## Open work

| #   | Task                                                                                                                                               | P  | Evidence / why                                                                                                                                   | Source                        |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------- | -- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------- |
| T1  | Split `breaker_test.go` (354 lines, file-size cap 350): move a coherent test group (e.g. the OnTransition tests) into a second file                | P2 | file-size-check error finding: `breaker_test.go: 354 lines (max 350)`; growing it further stacks another finding                                 | buildflow gate 2026-11-22     |
| T2  | Adopt a `.golangci.yml` via `buildflow -s golangci-lint-auto-configure --fix`, review the generated diff before accepting it                       | P2 | go-structure-linter error finding: "Missing golangci-lint config"; repo currently lints with defaults only                                       | buildflow gate 2026-11-22     |
| T8  | Regenerate `reports/coverage.out` from the current tree or trash it                                                                                | P3 | Stale 2026-11-21 artifact, gitignored, wired to nothing; nobody regenerates it                                                                   | 2026-11-22 report §c.7        |
| T9  | Guard test pinning the go.mod `go 1.27` directive (pattern: go-retry's `TestModuleGoDirectiveStaysPinned`)                                         | P2 | AGENTS.md claims "older toolchains will not compile" with no mechanical gate; the sibling repo already solved this                               | 2026-11-22 report §c.4        |
| T13 | markdownlint config decision: check sibling repos for `.markdownlint*` configs; add one or record the absence as a deliberate non-fix in AGENTS.md | P3 | `buildflow -s markdown-lint` is not applicable without a config — the docs pass is checked by no markdown linter; fleet-normalcy unresearched    | 2026-11-22 audit report §f.5  |
| T14 | Standing guard (or an accepted hand-maintenance verdict recorded in FEATURES.md) for the coverage claim                                            | P3 | FEATURES.md hardcodes 98.4%; verify-checklist prefers a gate that re-derives doc numbers over hand-maintained counts                             | 2026-11-22 audit report §f.6  |
| T15 | Annotate the 16-24 report's eight completed §f rows with inline done-markers (lifecycle pass)                                                      | P3 | The report's (f) table still reads as 22 open items though 8 are done and routed; annotate-canon says ask first — owner question 3 pending       | 2026-11-22 audit report §f.9  |
| T16 | Trim CHANGELOG `[Unreleased]` prose bullets that re-tell the eight mechanical purge lines (optional polish)                                        | P3 | Same work narrated twice (prose + purge trail) in one section; append-only instinct kept both, trim is safe pre-release                          | 2026-11-22 audit report §f.10 |
| T17 | Check securitymd usage: does regeneration accept repo args that would preserve the manual SECURITY.md corrections?                                 | P3 | SECURITY.md drifted once via shared hygiene generation (wrong repo URL, nonexistent version, false CI claim); a blind regen would reintroduce it | 2026-11-22 audit report §f.11 |
