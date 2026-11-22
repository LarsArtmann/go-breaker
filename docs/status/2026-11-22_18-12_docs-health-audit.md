# Status Report — docs-health AUDIT execution (all seven fleet docs)

**2026-11-22 18:12 CET** · session scope: execute the docs-health skill
(AUDIT: BUILD + HARVEST + VERIFY) over TODO_LIST, CHANGELOG, AGENTS, README,
ROADMAP, FEATURES, docs/DOMAIN_LANGUAGE — plus SECURITY.md, pulled in by
fix-on-sight after discovery of critical drift.

**Scope line — what this report is backed by.** VERIFIED with live runs:
`go test ./... -race -count=1` (ok), `go vet ./...` (ok), `go test -cover`
and a fresh coverprofile (98.4%, single uncovered block `breaker.go:265`),
`buildflow -s markdown-lint` (not applicable: no markdownlint config),
`buildflow -s codespell` (0 findings), `docs-gates purge .` dry-run then
`--apply` (8 rows, exit 0 measured unpiped), `docs-gates run .` battery
(5/5 PASS, exit 0, unpiped), `dprint check` over 7 md files (clean), buildflow
true exit code (69, unpiped — 2 known error findings), `git status` (4
modified + 5 untracked). NOT verified, and why: whether sibling repos carry
markdownlint configs (fleet-normalcy of the markdown-lint no-op —
unresearched); whether securitymd regeneration would preserve the manual
SECURITY.md corrections (tool usage not checked); pkg.go.dev rendering and
consumer-repo internals (out of scope); the pma daemon's actual pickup
behavior (not observed).

## a) FULLY DONE

1. **Skill executed properly**: full docs-health SKILL.md + build-guide,
   verify-checklist, and health-report-format references loaded and followed;
   the 2026-11-22_16-24 report (the only `2026-0*` file) viewed in full
   first, per instruction.
2. **FEATURES.md BUILT** (was missing — the doc set's one hole): sibling
   go-retry house style, provenance line, every row evidence-cited to
   code + passing tests, coverage claim verified against a fresh profile
   (the one uncovered block is the no-op return in `transition`).
3. **README repaired**: Installation section added (Go 1.27+ floor, `go
   get`); "Zero dependencies" lie replaced by "One dependency, by design".
4. **docs/DOMAIN_LANGUAGE.md extended**: `TripOn`, `Probe`,
   `OnTransition`, `Clock` terms + three language rules (trip ≠ any error;
   probe reserved for half-open; lazy = observed, not timer-fired).
5. **SECURITY.md corrected** (critical drift found by reading it — T11):
   advisory URL pointed at **go-must** (wrong repository), supported
   version cited **v0.1.2** (never released here), practices claimed **CI
   security scanning** (no CI exists). All fixed; regeneration gotcha
   recorded in AGENTS.md.
6. **AGENTS.md endurance-rewritten**: dated baseline section replaced by
   enduring quality-gate gotchas (T10 executed early — the three-way split
   brain is dead); FEATURES/ROADMAP/SECURITY wired into files-that-move-
   together; internal/ non-fix verdict recorded (T12); 6.8 KB, zero
   dates/hashes.
7. **ROADMAP.md restructured**: non-goals added from STACK.md's deliberate
   exclusions; two owner-answered questions removed; **license
   contradiction added as Open question 2** (new critical finding: README
   MIT badge vs PROPRIETARY LICENSE, present since the initial commit;
   sibling go-retry is MIT).
8. **CHANGELOG `[Unreleased]`** written (Added/Changed/Fixed) **and** the
   mechanical purge ritual executed per the active-table canon: 8 completed
   rows struck whole-line → `docs-gates purge --apply` → verdict lines in
   CHANGELOG. TODO_LIST pruned 12 → 4 open rows (T1, T2, T8, T9 — all
   code/tooling tasks, correctly still open).
9. **On-demand gates run** (T5): codespell 0 findings; markdown-lint not
   applicable (no config — a no-op by prerequisite, reported as such, not
   claimed as a pass).
10. **Inline health report delivered** per the skill format: Accuracy
    5.0/10, Fitness 8.25/10, AS-FOUND scoring, visible math, residuals
    named. First audit — no baseline claimed.

## b) PARTIALLY DONE

1. **T5's markdown-lint half** — ran, but as a no-op: this repo has no
   markdownlint config. Whether that is fleet-normal (sibling repos
   unchecked) or a gap to close is unresolved (new TODO T13).
2. **CHANGELOG `[Unreleased]`** now carries both prose entries and 8
   verbose purge lines covering the same work — two lenses, one file;
   deliberate keep (append-only instinct), but trimmable.
3. **TODO_LIST** holds 4 genuinely open rows — the doc work is done; the
   code/tooling work (T1/T2/T8/T9) is out of docs-health's lane by design.
4. **Quality gate** still exits 69 on the same 2 tracked findings —
   unchanged, correctly tracked, not doc work.

## c) NOT STARTED

1. T1 (split `breaker_test.go`), T2 (`.golangci.yml`), T8
   (regenerate/trash `reports/coverage.out`), T9 (go.mod directive guard).
2. License resolution (ROADMAP Open question 2 — owner).
3. CI decision, pkg.go.dev verification, CONTRIBUTING pointer, link check,
   dependabot grouping (ROADMAP raw ideas).
4. securitymd regeneration path check (would it preserve the manual
   SECURITY.md fixes?).
5. Annotating the 16-24 report's §f table for its 8 now-completed rows —
   deliberately skipped: same-day report, harvest trail (TODO → CHANGELOG
   purge lines) intact, and the annotate canon says ask before touching
   historical files nobody named.

## d) TOTALLY FUCKED UP!

Nothing destructive — battery green, tree exactly as intended, no data
lost. The sloppiness ledger (all caught and corrected in-flight):

1. **Exit-code-through-pipe, AGAIN.** First purge attempt printed
   "PURGE-EXIT=0" beside a FATAL — the exact mistake documented as a
   process fix in the 16-24 report, repeated one hour later. Corrected on
   the retry (unpiped redirect). The lesson needed a second application to
   stick.
2. **First purge apply FATALed** — I hand-wrote the CHANGELOG structure
   without the `### Changed` section the purge tool requires, discovering
   the requirement by failure instead of by checking first.
3. **Two stale-read edit rejections** — dprint auto-fix ran between my
   read and my edit (TODO_LIST, 16-24 report). Process fix now explicit:
   re-view immediately before editing anything a formatter touched.
4. **Freebie missed**: the audit generated a fresh coverage profile at
   `/tmp/cov.out` — T8 only needed it written to `reports/coverage.out`.
   Left on the table; kept as a one-command TODO.

## e) WHAT WE SHOULD IMPROVE!

- **Repeat-offender discipline**: the pipe-exit-code lesson recidivated
  within the hour. Standing rule from now on: any command whose exit code
  matters gets `> file; echo $?` — no `| tail` in the same breath, ever.
- **Small self-created duplication**: FEATURES' `WORTH_CONSIDERING`
  section and ROADMAP's raw ideas list overlap (sliding window, labels,
  ExecuteCtx, metrics adapter). Mitigated by an explicit pointer ("refined
  in ROADMAP.md") and sibling precedent (go-retry does the same), but it is
  the drift shape the skill warns about — monitor, or collapse FEATURES'
  section to a one-line pointer.
- **Hardcoded coverage number** in FEATURES.md (98.4%) — command-cited but
  not gate-protected; verify-checklist prefers a standing guard over
  hand-maintained counts (new TODO T14).
- **markdownlint fleet-normalcy unresearched** — the "no config" verdict
  would be stronger with one `ls` across sibling repos (new TODO T13).
- **Annotation restraint** — skipping the 16-24 annotate was the skill's
  default (ask before touching), but the next session that touches docs
  should do it; leaving stale-open §f rows to age is how reports rot.

## f) Up to 50 things to get done next

17 real items (routes: TODO = bounded; ROADMAP = raw/owner; the ceiling is
not a target — padding would be dishonest filler):

| #  | Task                                                                                                   | P  | Route   |
| -- | ------------------------------------------------------------------------------------------------------ | -- | ------- |
| 1  | T1: split `breaker_test.go` (354 > 350 cap) — OnTransition tests to a second file                      | P2 | TODO    |
| 2  | T2: adopt `.golangci.yml` via `buildflow -s golangci-lint-auto-configure --fix`, review diff           | P2 | TODO    |
| 3  | T8: regenerate `reports/coverage.out` (the audit's profile was one `cp` away) or trash it              | P3 | TODO    |
| 4  | T9: guard test pinning go.mod `go 1.27` (go-retry's `TestModuleGoDirectiveStaysPinned` pattern)        | P2 | TODO    |
| 5  | T13 (new): markdownlint config decision — check sibling repos, add config or record deliberate non-fix | P3 | TODO    |
| 6  | T14 (new): standing guard (or accepted hand-maintenance verdict) for FEATURES' coverage claim          | P3 | TODO    |
| 7  | License contradiction resolution — MIT badge vs PROPRIETARY LICENSE (owner, Q1 below)                  | P1 | OWNER   |
| 8  | Commit ownership — 4 modified + 5 untracked files waiting (owner, Q2 below)                            | P1 | OWNER   |
| 9  | Annotate the 16-24 report's 8 completed §f rows (next docs session)                                    | P3 | TODO    |
| 10 | Trim CHANGELOG `[Unreleased]` prose that re-tells the 8 purge lines (polish, optional)                 | P3 | TODO    |
| 11 | securitymd regeneration path: verify correct args exist so future regens keep the fixes                | P3 | TODO    |
| 12 | Verify pma daemon actually picks up these files (Q2-adjacent observation task)                         | P3 | TODO    |
| 13 | CI decision: GitHub Actions running buildflow vs local-gate-only                                       | P2 | ROADMAP |
| 14 | Verify pkg.go.dev rendering for v0.2.0 (go-release skill)                                              | P3 | ROADMAP |
| 15 | CONTRIBUTING.md pointer to AGENTS.md for AI-assisted contributions                                     | P3 | ROADMAP |
| 16 | Lychee/link check of README consumer links                                                             | P3 | ROADMAP |
| 17 | Dependabot major-version grouping policy                                                               | P3 | ROADMAP |

Items 1-4 are existing TODO rows; 5-6 are new rows harvested now; 9-12 are
new bounded rows harvested now; 7-8 live in ROADMAP Open questions; 13-17
are ROADMAP raw ideas.

## g) Questions I cannot figure out myself

1. **License.** The README badge says MIT; the LICENSE file is PROPRIETARY
   ("All rights reserved", present since the initial commit); sibling
   go-retry is MIT and seven fleet consumers already import this library.
   Which is intended? I will not pick a side on a legal artifact.
2. **Commit ownership.** Ten files (4 modified, 5 untracked, this report's
   directory) are uncommitted. The harness forbids committing without your
   word; the pma daemon reportedly has a blind spot. Say "commit" and I
   will — properly messaged — or confirm the daemon owns it.
3. **Annotate the 16-24 report now?** Its (f) table still reads as 22 open
   items though 8 are done and routed. I left it untouched (same-day,
   annotate-canon says ask). Want the inline done-markers applied, or leave
   it for the next docs session's lifecycle pass?

---

_Point-in-time snapshot. When stale, ANNOTATE — do not rewrite (docs-health
policy)._
