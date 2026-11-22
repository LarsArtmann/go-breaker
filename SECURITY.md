# Security Policy

## Supported Versions

| Version | Supported  |
| ------- | ---------- |
| v0.2.0  | ✅ Current |

Only the latest release receives security fixes. Please update to the latest
version before reporting issues found in older releases.

## Reporting a Vulnerability

Please do **not** open public issues for security problems. Report privately instead:

- **GitHub Security Advisory**: https://github.com/LarsArtmann/go-breaker/security/advisories/new

### What to include

- Vulnerability type and severity
- Steps to reproduce
- Potential impact
- Proof of concept (if available)

### Response commitments

| Phase            | Timeline                    |
| ---------------- | --------------------------- |
| Initial response | Within 48 hours             |
| Triage           | Within 72 hours             |
| Disclosure       | Coordinated, within 90 days |

## Security Practices

- All changes are code-reviewed before merge
- Quality gate (lint, format, tests with `-race`) runs via BuildFlow before
  commits land; repository has no CI pipeline yet (tracked in
  [ROADMAP.md](ROADMAP.md))
- Dependencies are kept up to date (weekly Dependabot gomod updates)

---

_Maintained with [securitymd](https://github.com/LarsArtmann/securitymd). Last
updated: 2026-11-22 — corrected for go-breaker (advisory URL, supported
version, CI claim)._
