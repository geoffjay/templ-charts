package heatmap

import (
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
)

// UseHeatMap mirrors @nivo/heatmap's compute pipeline: it derives the x/y band
// scales, the value domain + continuous color scale, and the positioned,
// colored cells. props.Width/Height are the inner (margin-subtracted)
// dimensions, matching the HeatMap templ component's contract.
func UseHeatMap(props HeatMapProps) HeatMapResult {
	theme := resolveTheme(props.Theme)

	// Ordered, unique x categories (first-seen) and serie ids.
	xValues := make([]any, 0)
	seenX := map[string]bool{}
	serieIDs := make([]any, 0, len(props.Data))
	min := math.Inf(1)
	max := math.Inf(-1)
	hasValue := false
	for _, serie := range props.Data {
		serieIDs = append(serieIDs, serie.ID)
		for _, d := range serie.Data {
			if !seenX[d.X] {
				seenX[d.X] = true
				xValues = append(xValues, d.X)
			}
			if d.Y != nil {
				hasValue = true
				if *d.Y < min {
					min = *d.Y
				}
				if *d.Y > max {
					max = *d.Y
				}
			}
		}
	}
	if !hasValue {
		min, max = 0, 0
	}

	offsetX, offsetY, layoutW, layoutH := computeLayout(props.Width, props.Height, len(xValues), len(serieIDs), props.ForceSquare)

	// Band scales: built with the X-axis range orientation ([0,size], not
	// reversed) so serie ids run top→bottom like nivo's heatmap.
	xScale := scales.ComputeScale(scales.ScaleBandSpec{}, scales.ComputedSerieAxis{All: xValues}, layoutW, scales.ScaleAxisX)
	yScale := scales.ComputeScale(scales.ScaleBandSpec{}, scales.ComputedSerieAxis{All: serieIDs}, layoutH, scales.ScaleAxisX)

	cellWidth, cellHeight := 0.0, 0.0
	if bw, ok := xScale.(scales.ScaleWithBandwidth); ok {
		cellWidth = bw.Bandwidth()
	}
	if bh, ok := yScale.(scales.ScaleWithBandwidth); ok {
		cellHeight = bh.Bandwidth()
	}

	colorScale := buildColorScale(props.Colors, min, max)
	borderGen := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	labelGen := colors.GetInheritedColorGenerator(props.LabelTextColor, theme)
	format := valueFormatter(props.ValueFormat)
	labelsEnabled := props.LabelsEnabled()

	hovering := props.HoveredKey != ""

	cells := make([]ComputedCell, 0, len(xValues)*len(serieIDs))
	for _, serie := range props.Data {
		yCenter := yScale.Call(serie.ID) + cellHeight/2 + offsetY
		for _, d := range serie.Data {
			color := props.EmptyColor
			formatted := ""
			label := ""
			id := serie.ID + "." + d.X
			if d.Y != nil {
				color = colorScale(*d.Y)
				formatted = format(*d.Y)
				if labelsEnabled {
					label = formatted
				}
			}
			// Opacity: dim non-hovered cells when a cell is hovered.
			opacity := props.Opacity
			if hovering {
				if id == props.HoveredKey {
					opacity = props.ActiveOpacity
				} else {
					opacity = props.InactiveOpacity
				}
			}
			ctx := map[string]any{"color": color}
			cells = append(cells, ComputedCell{
				ID:             id,
				SerieID:        serie.ID,
				X:              d.X,
				Value:          d.Y,
				FormattedValue: formatted,
				XPos:           xScale.Call(d.X) + cellWidth/2 + offsetX,
				YPos:           yCenter,
				Width:          cellWidth,
				Height:         cellHeight,
				Color:          color,
				Opacity:        opacity,
				BorderColor:    borderGen(ctx),
				Label:          label,
				LabelTextColor: labelGen(ctx),
			})
		}
	}

	return HeatMapResult{
		Cells:      cells,
		XScale:     xScale,
		YScale:     yScale,
		XValues:    xValues,
		SerieIDs:   serieIDs,
		OffsetX:    offsetX,
		OffsetY:    offsetY,
		Width:      layoutW,
		Height:     layoutH,
		MinValue:   min,
		MaxValue:   max,
		ColorScale: colorScale,
	}
}

// computeLayout mirrors @nivo/heatmap computeLayout: when forceSquare, cells
// are squared and the grid is centered, returning the offset + laid-out size.
func computeLayout(w, h float64, cols, rows int, forceSquare bool) (offsetX, offsetY, width, height float64) {
	width, height = w, h
	if !forceSquare || cols == 0 || rows == 0 {
		return 0, 0, width, height
	}
	cellWidth := math.Max(w/float64(cols), 0)
	cellHeight := math.Max(h/float64(rows), 0)
	cellSize := math.Min(cellWidth, cellHeight)
	width = cellSize * float64(cols)
	height = cellSize * float64(rows)
	offsetX = (w - width) / 2
	offsetY = (h - height) / 2
	return offsetX, offsetY, width, height
}

// buildColorScale resolves the continuous color scale (sequential or
// diverging) for the value domain [min, max].
func buildColorScale(cfg HeatMapColorConfig, min, max float64) func(float64) string {
	vals := colors.SequentialColorScaleValues{Min: min, Max: max}
	if cfg.Type == "diverging" {
		return colors.GetDivergingColorScale(colors.DivergingColorScaleConfig{
			Type:      "diverging",
			Scheme:    cfg.Scheme,
			MinValue:  cfg.MinValue,
			MaxValue:  cfg.MaxValue,
			DivergeAt: cfg.DivergeAt,
		}, vals)
	}
	return colors.GetSequentialColorScale(colors.SequentialColorScaleConfig{
		Type:     "sequential",
		Scheme:   cfg.Scheme,
		MinValue: cfg.MinValue,
		MaxValue: cfg.MaxValue,
	}, vals)
}

// valueFormatter returns a number→string formatter. An empty spec uses %g; a
// non-empty spec is treated as a d3-format spec.
func valueFormatter(spec string) func(float64) string {
	if spec == "" {
		return func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	}
	return func(v float64) string { return d3format.FormatString(spec, v) }
}

// resolveTheme returns props.Theme or the default theme when nil.
func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

// themeBackground returns the theme background color (or "" for transparent).
func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}
