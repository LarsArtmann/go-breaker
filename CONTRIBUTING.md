# Contributing

Thanks for your interest in contributing!

## How to Contribute

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request (CI runs race tests, a 95% coverage floor,
   govulncheck, and lint on every PR)

## Development Setup

Run the following commands to set up your development environment:

    go test ./... -race
    golangci-lint run ./...

The linter configuration lives in-repo (`.golangci.yml`), so a plain
`golangci-lint run` matches what CI enforces.

## Working with AI assistants

This repo keeps [AGENTS.md](AGENTS.md) as the session-context file for AI
contributors — it records the quality-gate commands, deliberate non-fixes,
and the design contracts that breaking counts as a breaking change. Read it
before changing `breaker.go`. Fleet-internal sessions use
[BuildFlow](https://github.com/LarsArtmann/BuildFlow) as the local quality
gate instead of invoking the underlying tools directly.

## Reporting Issues

Please use GitHub Issues to report bugs or request features.
