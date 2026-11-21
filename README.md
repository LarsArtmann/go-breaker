# go-breaker

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A dependency-light, error-family-aware circuit breaker for Go — the sibling of [go-retry](https://github.com/Larsartmann/go-retry) for the failure side of resilience. Retry answers "should I try again?"; the breaker answers "should I even try?". Extracted 2026-11-21 from the domain-language census, which counted 12+ bespoke circuit-breaker implementations across the fleet.

```go
import "github.com/larsartmann/go-breaker"
```

## Why error-family aware

Only **Transient** and **Infrastructure** failures count toward opening the circuit — a Rejection (bad input), Conflict (state mismatch), or Corruption (integrity) failure is deterministic, and breaking the circuit would not make the next identical call succeed. The default predicate is one line of [go-error-family](https://github.com/larsartmann/go-error-family) classification, overridable per breaker.

## Usage

```go
b := breaker.New(breaker.Config{
    FailureThreshold: 5,             // consecutive trip-worthy failures before opening
    ResetTimeout:     30 * time.Second, // open -> half-open cooldown
    SuccessThreshold: 2,             // half-open successes to close
})

err := b.Execute(func() error {
    return callDownstream(ctx)
})

if errors.Is(err, breaker.ErrOpen) {
    // refused cheaply, downstream never touched
}
```

State machine: `Closed --threshold--> Open --cooldown--> HalfOpen --successes--> Closed` (a trip-worthy failure in half-open re-opens immediately). `ErrOpen` is classified Transient on purpose: outer retry loops back off rather than treat the refusal as terminal.

Lower-level API when the call and the verdict are separated: `Allow() error` + `Record(err)`. Concurrency-safe; clock injectable via `Config.Now` for deterministic tests. `Config.OnTransition func(from, to State)` fires on every state change, serialized in transition order (under the lock) — for metrics and logging; observers must not call back into the same Breaker.

## Zero dependencies

One dependency by design: go-error-family (the classification contract the trip semantics ride on). No timers, no goroutines, no background state — the circuit evaluates lazily on each call.

## Consumers

Adopted: **file-and-image-renamer** (v0.2.0, 2026-11-21 — provider-call gate; `OnTransition` feeds the resilience metrics; grew from this consumer's need).

Planned migration rows (filed 2026-11-21 per ADR-0001, wire-or-delete for each bespoke breaker): Cyberdom, Kernovia, RedditParse, Rolls-Royce-mtuGoHelpCenter-golang, StopTube, dnsblockd, superb-gh-milestone-extention (and the config-declared breakers in ActaFlow, ai-router-implementation, ai-task-prioritizer, BuildFlow, ast-state-analyzer where the verdict is a row, not a migration).
