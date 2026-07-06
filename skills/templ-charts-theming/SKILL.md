---
name: templ-charts-theming
description: Style templ-charts charts — color palettes and the five Colors config shapes, custom/dark themes, legends, and gradient/pattern fills with match rules. Use when setting chart colors, building a dark or branded theme, adding a legend, or applying gradients/patterns to templ-charts marks.
---

# templ-charts theming & styling

## Colors

Charts color marks through a `Colors` field (default: the `nivo` scheme).
There are **five config shapes** depending on the chart family:

1. **Ordinal/categorical** (most charts: bar, line, pie, radar, treemap,
   sankey, …) — field type `colors.OrdinalColorScaleConfig`:

   ```go
   import "github.com/geoffjay/templ-charts/charts/colors"

   Colors: colors.Scheme(colors.PaletteTableau10),      // named palette
   Colors: colors.PaletteColors("#4269d0", "#efb118"),  // explicit list
   ```

2. **Continuous** (heatmap) — `heatmap.HeatMapColorConfig` (sequential or
   diverging value→color scale) + `EmptyColor string`.
3. **Quantize palette** (calendar) — `Colors []string` buckets + `EmptyColor`.
4. **Quantize scheme id** (geo Choropleth) — `Colors string` (e.g.
   `"purple_blue_green"`) + `Steps int` + `UnknownColor`. GeoMap takes a
   single `FillColor string`.
5. **Plain color strings** (network) — `NodeColor`, `NodeBorderColor`,
   `LinkColor string`.

**One entry point over all five:** `colors.Set(v)` accepts a `PaletteID`, a
`[]string`, a color string, a `func(any) string`, or a pre-built config, and
resolves to whatever shape a field wants:

```go
Colors:      colors.Set(colors.PaletteTableau10).Ordinal(),
BorderColor: colors.Set("#333").Inherited(),
// also: .Sequential(), .Diverging(), .Static()
// variants: colors.SetFromDatum("color"), colors.SetFunc(fn)
// modifiers: colors.Set(...).WithModifiers(colors.ColorModifier{"darker", 0.6})
```

**Perceptual ramps (opt-in):** sequential/diverging scales interpolate in RGB
by default; pass `.InSpace(colors.SpaceLab)` or `colors.SpaceLch` for a
perceptually-uniform ramp.

### Palettes

Categorical: `PaletteNivo`, `PaletteCategory10`, `PaletteAccent`,
`PaletteDark2`, `PalettePaired`, `PalettePastel1/2`, `PaletteSet1/2/3`,
`PaletteTableau10`, `PaletteTableau20`, `PaletteObservable10`,
`PaletteMaterial`, `PaletteOkabeIto`, `PaletteTolVibrant`, `PaletteTolMuted`.
Sequential: `PaletteViridis`, `PaletteMagma`, `PaletteInferno`,
`PalettePlasma`, `PaletteTurbo`, `PaletteCividis`, `PaletteWarm`,
`PaletteCool`, `PaletteSunset`, `PaletteTeal`, `PaletteEmerald`.
Diverging: `PaletteRedBlue`, `PaletteRedYellowBlue`, `PaletteBrownBlueGreen`,
`PaletteSpectral`, `PaletteTemps`.

Enumerate at runtime: `colors.Palettes()`, `colors.PalettesByKind(kind)`,
`colors.ColorblindSafePalettes()`, `colors.LookupPalette(id)`. Full catalog:
`docs/PALETTES.md` in the repo; visual gallery on the demo app's `/palettes`
page.

## Themes

Charts take `Theme *theming.Theme` (nil → `theming.DefaultTheme`). Build a
custom theme by extending the default with a `PartialTheme` (deep-merged,
text-style inheritance resolved):

```go
import "github.com/geoffjay/templ-charts/charts/theming"

bg := "#1a1a1a"
dark := theming.ExtendDefaultTheme(theming.DefaultTheme, &theming.PartialTheme{
    Background: &bg,
    // Text, Axis, Grid, Crosshair, Legends, Labels, Markers, Dots,
    // Tooltip, Annotations — all optional pointer overrides
})

bar.Bar(bar.BarProps{ /* … */ Theme: &dark })
```

`PartialTheme` fields are pointers — set only what you override. The demo
app's `/themes` page shows default/dark/custom side by side.

## Legends

Charts accept a `Legends` slice. Most take `[]legends.LegendProps`; bar wraps
it with a data source selector:

```go
import "github.com/geoffjay/templ-charts/charts/legends"

// line: Legends []legends.LegendProps
Legends: []legends.LegendProps{{
    Anchor:      legends.LegendAnchorBottomRight,
    Direction:   legends.LegendDirectionColumn,
    ItemWidth:   100, ItemHeight: 20,
    SymbolShape: legends.SymbolShapeCircle,  // Circle|Square|Diamond|Triangle
    SymbolSize:  12,
    TranslateX:  110,                        // offset from anchor
}},

// bar: Legends []bar.BarLegendProps — embeds legends.LegendProps
Legends: []bar.BarLegendProps{{
    LegendProps: legends.LegendProps{ /* as above */ },
    DataFrom:    "keys",                     // "keys" | "indexes"
}},
```

Anchors: `top-left`, `top`, `top-right`, `right`, `bottom-right`, `bottom`,
`bottom-left`, `left`. Leave room in `Margin` for the legend (e.g.
`Margin: core.Margin{Right: 120}` for a right-anchored column). With the
HTMX registry wired (see templ-charts-interactivity), legend items become
clickable series toggles.

## Gradients & patterns (defs + match rules)

Fill marks with SVG defs from `charts/core`, matched to data by rules:

```go
import "github.com/geoffjay/templ-charts/charts/core"

// Gradient that inherits each series' color:
Defs: []core.Def{
    core.LinearGradientDef("areaGrad", []core.GradientStop{
        {Offset: 0, Color: "inherit", Opacity: 0.55},
        {Offset: 100, Color: "inherit", Opacity: 0.02},
    }, nil),
},
Fill: []core.DefRule{{ID: "areaGrad", Match: "*"}},   // apply to all

// Pattern matched to one datum by property equality:
Defs: []core.Def{{
    ID: "dots", Type: core.DefTypePatternDots,
    Background: "inherit", Color: "rgba(255,255,255,0.55)",
    Size: 4, Padding: 3, Stagger: true,
}},
Fill: []core.DefRule{{ID: "dots", Match: map[string]any{"id": "fries"}}},
```

`Match` accepts `"*"` (all), a `map[string]any` (field equality against the
datum), or a `func(any) bool`. Def types: `core.DefTypeLinearGradient`,
`DefTypePatternDots`, `DefTypePatternLines`, `DefTypePatternSquares`
(constructors: `LinearGradientDef`, `PatternDotsDef`, `PatternLinesDef`,
`PatternSquaresDef`). `Color: "inherit"` picks up the mark's ordinal color.
The demo app's `/styling` page exercises all of these.
