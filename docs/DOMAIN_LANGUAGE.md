# Domain Language — go-breaker

A vocabulary for the failure side of resilience. Grounded in this repo's code.

## Glossary

| Term         | Definition                                                                                                                              | Where it lives        |
| ------------ | --------------------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| Circuit      | The per-downstream gate deciding whether calls are attempted at all                                                                     | `breaker.Breaker`     |
| State        | The circuit's position: Closed (calls pass, failures count), Open (calls refused cheaply), HalfOpen (trial calls probe recovery)        | `breaker.State`       |
| Trip         | A failure that counts toward opening the circuit — only Transient/Infrastructure families by default; deterministic failures never trip | `tripsByFamily`       |
| TripOn       | The configurable predicate deciding whether an error trips; overridable per breaker for domain-specific tripping                        | `Config.TripOn`       |
| Threshold    | Consecutive trip-worthy failures that open the circuit (default 5); half-open successes that close it (default 2)                       | `breaker.Config`      |
| Cooldown     | How long the circuit stays open before the next Allow probes half-open (default 30s)                                                    | `Config.ResetTimeout` |
| Probe        | A trial call admitted in HalfOpen to test recovery; successes accumulate toward closing, one trip-worthy failure re-opens               | `Allow`, `Record`     |
| ErrOpen      | The cheap refusal while open; classified Transient so outer retry loops back off instead of aborting                                    | `breaker.ErrOpen`     |
| OnTransition | Observer fired on every state change with from/to states, serialized in transition order under the breaker's lock                       | `Config.OnTransition` |
| Clock        | The injectable time source; the only time the circuit reads, enabling sleep-free deterministic tests                                    | `Config.Now`          |

## Language rules

- **"The breaker tripped"** means a trip-worthy failure crossed the threshold
  — not merely any error (a Rejection never trips the breaker).
- **"Probe"** is reserved for half-open trial calls; a closed-state call is
  just a call.
- Lazy by design: state moves only when observed — "the cooldown passed"
  means the next `Allow`/`State` observed it, not that a timer fired.
