# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0]

First stable release. From this release on, the public API under `github.com/geoffjay/templ-charts/charts/...`
follows semantic versioning: breaking changes will require a new major version.

### Added

- **28 chart families** rendered as server-side SVG via [templ](https://templ.guide):
  bar, line, pie, scatterplot, heatmap, waffle, calendar, radar, radial-bar, stream,
  bullet, funnel, boxplot, bump, marimekko, parallel-coordinates, polar-bar, treemap,
  sunburst, icicle, circle-packing, tree, voronoi, network, swarmplot, sankey, chord,
  and geo.
- **Foundation packages**: core, axes, arcs, legends, theming, colors, scales, render,
  samples, static, text, tooltip, annotations, grid, rects, interact, htmx, and canvas.
- **Theming and color palettes** with a documented palette catalog (see `docs/PALETTES.md`).
- **Interactivity** via HTMX endpoints and a lightweight client hover script
  (`charts/interact`), including series toggling and hierarchy zoom.
- **Accessibility and responsiveness** support documented in `docs/USAGE.md`.
- **Pure-Go ports of D3 primitives** under `internal/d3` (array, scale, shape,
  hierarchy, force, sankey, chord, delaunay, quadtree, geo, color, format, timeformat) —
  no runtime dependencies beyond templ.
- Comprehensive test suite: golden-snapshot SVG tests, benchmarks, and runnable
  example tests across all chart families.
- Runnable demo application under `examples/app` showcasing every chart with theme and
  palette switching.

### Changed

- Lowered the declared minimum Go version to the true floor (`go 1.25.0`, required by
  the templ dependency) so consumers are not forced onto a newer toolchain than needed.

[Unreleased]: https://github.com/geoffjay/templ-charts/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/geoffjay/templ-charts/releases/tag/v1.0.0
