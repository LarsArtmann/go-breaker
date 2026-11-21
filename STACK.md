# STACK.md — go-breaker

Decision record per the stack-selection skill. Profile-driven; re-derive before extending.

## Profile

| Axis        | Answer                                                                                        |
| ----------- | --------------------------------------------------------------------------------------------- |
| Shape       | Pure Go library (no CLI, no service, no goroutines)                                           |
| Domain      | Circuit breaking, error-family aware (fleet resilience pair with go-retry)                    |
| Persistence | None (in-memory state, lazily evaluated)                                                      |
| Interfaces  | Go API only (Execute / Allow+Record / State)                                                  |
| Constraints | Minimal dependencies; deterministic testability (clock injection); nix-hermetic by triviality |

## Chosen foundation

| Concern              | Choice                | Why this, not that                                                                                                                                           | Owning skill                        | Revisit trigger                              |
| -------------------- | --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------- | -------------------------------------------- |
| Error classification | go-error-family v0.12 | The trip predicate IS the error-family contract (Transient/Infrastructure trip; Rejection/Conflict/Corruption never); every fleet consumer already speaks it | how-to-golang                       | a classification need Classify cannot answer |
| State machine        | stdlib (mutex + enum) | Three states, lazy evaluation; timer/goroutine designs add lifecycle cost for zero semantics                                                                 | how-to-golang (smallest foundation) | metrics/events hooks needing observers       |

## Deliberately excluded

| Candidate               | Reason                                                                                                                              |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| sony/gobreaker          | Third-party; no error-family awareness (every error counts), fleet policy prefers first-party when the semantics are fleet-specific |
| afex/hystrix-go         | Abandoned upstream; pulls a full concurrency runtime                                                                                |
| hosting inside go-retry | Retry and breaker compose as siblings, not layers; nesting forecloses independent adoption (plan M18 named it a sibling explicitly) |

## Wiring order

1. `meta set --project go-breaker` (registered 2026-11-21)
2. First consumer migration per ADR-0001 rule 5 (rows filed in the docs TODO_LIST, 2026-11-21)
