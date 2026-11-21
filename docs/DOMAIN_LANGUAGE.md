# Domain Language — go-breaker

A vocabulary for the failure side of resilience. Grounded in this repo's code.

## Glossary

| Term      | Definition                                                                                                                              | Where it lives        |
| --------- | --------------------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| Circuit   | The per-downstream gate deciding whether calls are attempted at all                                                                     | `breaker.Breaker`     |
| State     | The circuit's position: Closed (calls pass, failures count), Open (calls refused cheaply), HalfOpen (trial calls probe recovery)        | `breaker.State`       |
| Trip      | A failure that counts toward opening the circuit — only Transient/Infrastructure families by default; deterministic failures never trip | `tripsByFamily`       |
| Threshold | Consecutive trip-worthy failures that open the circuit (default 5); half-open successes that close it (default 2)                       | `breaker.Config`      |
| Cooldown  | How long the circuit stays open before the next Allow probes half-open (default 30s)                                                    | `Config.ResetTimeout` |
| ErrOpen   | The cheap refusal while open; classified Transient so outer retry loops back off instead of aborting                                    | `breaker.ErrOpen`     |
