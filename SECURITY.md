# Security Policy

## Supported versions

templ-charts follows semantic versioning. Security fixes are applied to the latest
released minor version. Older versions are not maintained.

| Version | Supported |
| ------- | --------- |
| 1.x     | ✅        |
| < 1.0   | ❌        |

## Reporting a vulnerability

Please **do not** open a public issue for security vulnerabilities.

Instead, report privately using GitHub's
[private vulnerability reporting](https://github.com/geoffjay/templ-charts/security/advisories/new)
("Report a vulnerability" under the repository's **Security** tab). If that is
unavailable, email the maintainer at geoff.jay@gmail.com with details.

Please include:

- A description of the vulnerability and its impact.
- Steps to reproduce, or a proof-of-concept.
- Affected version(s).

You can expect an acknowledgement within a few days. Once a fix is available, a patched
release will be published and the advisory disclosed.

## Scope

This library renders SVG server-side from data supplied by the embedding application.
When charting **untrusted input**, note:

- All text is escaped before rendering, but you remain responsible for validating and
  sanitizing data appropriate to your context.
- Report any input that causes a panic, an unbounded allocation, or produces markup that
  escapes the SVG context — these are treated as security-relevant.
