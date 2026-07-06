---
name: templ-charts-composition
description: Build custom charts from templ-charts' Use* layout hooks — compute scales/series/generators with UseLine (or UseBar, UsePie, …), reuse the exported grid/axes/lines sub-components, and hand-write SVG marks (custom point symbols, direct labels, sparklines). Use when a stock templ-charts chart component isn't enough and you need a bespoke chart or custom layer.
---

# templ-charts composition (build-your-own charts)

Every chart component is a thin templ layer over an exported layout hook.
When the stock component doesn't fit, call the hook yourself: it returns the
computed scales, positioned series, and path generators, and you assemble the
SVG — reusing the exported sub-components for the standard layers and writing
custom marks for the rest.

Hooks exist per family — `line.UseLine(line.LineProps) line.LineResult`,
`bar.UseBar`, `pie.UsePie`, `heatmap.UseHeatMap`, `swarmplot.UseSwarmPlot`,
etc. — plus `core.UseDimensions(width, height, margin)` for margin math.

## The pattern (from the demo app's /composition page)

```go
import (
    "fmt"
    "strings"

    "github.com/geoffjay/templ-charts/charts/axes"
    "github.com/geoffjay/templ-charts/charts/core"
    "github.com/geoffjay/templ-charts/charts/line"
    "github.com/geoffjay/templ-charts/charts/render"
    "github.com/geoffjay/templ-charts/charts/theming"
)

// 1. Dimensions: outer size + margin → inner plot area.
dims := core.UseDimensions(700, 400, core.Margin{Top: 40, Right: 40, Bottom: 60, Left: 60})
theme := &theming.DefaultTheme

// 2. The hook computes scales, positioned series, and generators.
result := line.UseLine(line.LineProps{
    Width: dims.InnerWidth, Height: dims.InnerHeight,
    Curve: core.CurveMonotoneX,
    Data:  seriesData,
})

var inner strings.Builder

// 3. Standard layers via exported sub-components.
grid, _ := render.String(axes.Grid(axes.GridProps{
    Axis: "y", Scale: result.YScale,
    Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
}))
inner.WriteString(grid)

xa := axes.DefaultAxisProps
xa.Axis, xa.Scale, xa.Length = "x", result.XScale, dims.InnerWidth
xa.Y = dims.InnerHeight
xa.TicksPosition = "after"
ya := axes.DefaultAxisProps
ya.Axis, ya.Scale, ya.Length = "y", result.YScale, dims.InnerHeight
ya.TicksPosition = "after"
axesSVG, _ := render.String(axes.Axes(axes.AxesProps{
    XAxis: &xa, YAxis: &ya,
    Width: dims.InnerWidth, Height: dims.InnerHeight, Theme: theme,
}))
inner.WriteString(axesSVG)

lines, _ := render.String(line.Lines(line.LinesProps{
    Series: result.Series, LineGenerator: result.LineGenerator, LineWidth: 2.5,
}))
inner.WriteString(lines)

// 4. Custom marks: positioned data is on result.Series[i].Data[j].Position.
for _, s := range result.Series {
    for _, d := range s.Data {
        fmt.Fprintf(&inner,
            `<rect transform="translate(%.2f,%.2f) rotate(45)" x="-3.5" y="-3.5" width="7" height="7" fill="%s" stroke="#fff" stroke-width="1.5"></rect>`,
            d.Position.X, d.Position.Y, s.Color)
    }
}

// 5. Wrap in the accessible <svg> shell (applies margin translate, role, defs).
svg, err := render.String(core.SvgWrapper(core.SvgWrapperProps{
    Width: dims.OuterWidth, Height: dims.OuterHeight, Margin: dims.Margin,
    Role: "img", AriaLabel: "custom line chart",
}, inner.String()))
```

## What the hook gives you

`line.LineResult` (other hooks are analogous):

- `Series []ComputedSeries` — `{ID string; Color string; Data []ComputedDatum}`
  where each datum has the raw `Data` (X/Y values) and `Position` (pixel
  `X, Y float64` in inner-chart coordinates).
- `XScale`, `YScale scales.Scale` — pass to `axes.Grid`/`axes.Axes`, or call
  directly to place your own marks.
- `LineGenerator`, `AreaGenerator` — `func([]PointXY) string` producing SVG
  path `d` strings (curve-aware).
- `Points`, `Slices`, `LegendData`, `GetColor`, `HiddenIDs` — for tooltips,
  legends, and series filtering.

## Recipes proven in the demo

- **Custom point symbols** — render `line.Lines` normally, then draw your own
  markers at each `Position` (diamonds, rings on maxima, value labels).
- **Direct labels instead of a legend** — after `line.Lines`, take
  `s.Data[len(s.Data)-1].Position` per series and draw an end dot + a
  `<text>` label in `s.Color`; widen `Margin.Right` to make room.
- **Sparkline KPI cards** — run `UseLine` at tiny sizes (e.g. 190×48) with no
  axes/grid; draw gradient area + line + end dot from the returned
  generators; overlay `<text>` for the KPI value. Several cards can share one
  `<svg>` with `<g transform="translate(…)">` per card.

Full working source: `examples/app/demos/composition.go` (and
`dashboard.go`) in the repo — view at `/composition` via `make run-demo`.
