# Contributing to templ-charts

Thanks for your interest in contributing! This document explains how to set up your
environment, the conventions the project follows, and what a good pull request looks
like.

## Prerequisites

- **Go 1.25 or newer** (the declared minimum; see `go.mod`).
- The [`templ`](https://templ.guide) CLI, pinned to the version in the `Makefile`.
  Running `make templ` installs the pinned version into `$HOME/go/bin` if it is missing.

## Getting started

```sh
git clone https://github.com/geoffjay/templ-charts
cd templ-charts
make templ   # generate Go from .templ sources
make ci      # lint + test — should pass on a clean checkout
```

To see charts rendered live, run the demo app:

```sh
make run-demo   # serves the demo on http://localhost:8000
```

## Development workflow

Common `make` targets:

| Target        | What it does                                              |
| ------------- | --------------------------------------------------------- |
| `make templ`  | Regenerate `*_templ.go` from `.templ` sources             |
| `make build`  | Build all packages                                        |
| `make test`   | Run the test suite                                        |
| `make lint`   | `go vet` + `gofmt` check                                  |
| `make cover`  | Coverage profile + HTML report (`coverage.html`)          |
| `make bench`  | Run benchmarks (not part of CI; timing-sensitive)         |
| `make golden` | Regenerate golden SVG/path snapshots (see below)          |
| `make ci`     | Lint + test — the same gate CI enforces                   |

### Generated code

`*_templ.go` files are generated from `.templ` sources. If you change a `.templ` file,
run `make templ` and commit the regenerated Go alongside your change.

### Golden snapshots

Many chart packages assert their rendered output against golden snapshots in
`testdata/golden/`. If your change **intentionally** alters rendered output, regenerate
the snapshots with `make golden` and review the diff carefully before committing — an
unexpected diff usually means an unintended rendering change.

## Coding conventions

- Code must pass `make lint` (`go vet` clean, `gofmt` formatted).
- Public API changes should be considered carefully: this project follows semantic
  versioning, so breaking changes to anything under `charts/...` require a major version
  bump. Prefer additive changes.
- Exported symbols should have godoc comments.
- Keep the dependency footprint minimal — the library depends only on `templ` at runtime.
- Add or update tests for behavior changes; new chart features should include golden
  snapshots and, where useful, an `ExampleXxx` test.

## Commit messages

The history follows [Conventional Commits](https://www.conventionalcommits.org/)
(`feat:`, `fix:`, `chore:`, `docs:`, …). Please match that style.

## Pull requests

1. Fork and create a topic branch.
2. Make your change with tests; ensure `make ci` passes locally.
3. Regenerate templ/golden output if applicable and commit it.
4. Open a PR describing the change and its motivation. Link any related issue.

CI runs lint, tests (across supported Go versions), a vulnerability scan, and reports
coverage on the PR.

## License

By contributing, you agree that your contributions are licensed under the project's
[MIT License](LICENSE).
