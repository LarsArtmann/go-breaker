# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versioning follows [Semantic Versioning](https://semver.org/).

## [Unreleased]

Nothing yet.

## [0.2.0] - 2026-11-21

### Added

- `Config.OnTransition func(from, to State)`: observer fired on every state change (closed→open, open→half_open, half_open→closed, half_open→open), serialized in transition order under the breaker's lock. Observers must be cheap and must never call back into the same Breaker. Grew from the first consumer migration (file-and-image-renamer's resilience metrics feed). Tests: lifecycle order goldens, half-open reopen sequence.

## [0.1.0] - 2026-11-21

### Added

- Initial extraction (domain-language census plan M18): error-family-aware circuit breaker — Closed/Open/HalfOpen machine, `Execute`/`Allow`+`Record`/`State` API, `ErrOpen` (Transient, so outer retry loops back off), family-aware default trip predicate (Transient + Infrastructure trip; Rejection/Conflict/Corruption never; unclassified errors inherit go-error-family's fail-open Transient), concurrency-safe, injectable clock. Tests: state-transition goldens, family table, `-race` concurrency.
