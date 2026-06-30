package scatterplot

import (
	"fmt"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
)

// UseScatterPlot mirrors @nivo/scatterplot's useScatterPlot: it builds the x/y
// scales via ComputeXYScalesForSeries, flattens the series into positioned
// nodes, assigns a per-serie ordinal color and a constant node size, and
// derives the legend data. props.Width/Height are the inner dimensions.
func UseScatterPlot(props ScatterPlotProps) ScatterPlotResult {
	seriesInput := make([]scales.Serie, len(props.Data))
	for i, s := range props.Data {
		data := make([]scales.SerieDatum, len(s.Data))
		for j, d := range s.Data {
			data[j] = scales.SerieDatum{X: d.X, Y: d.Y}
		}
		seriesInput[i] = scales.Serie{Data: data, Extra: s.ID}
	}

	xy := scales.ComputeXYScalesForSeries(seriesInput, props.XScale, props.YScale, props.Width, props.Height)

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	formatX := valueFormatter(props.XFormat)
	formatY := valueFormatter(props.YFormat)

	nodes := make([]ComputedNode, 0)
	for _, s := range props.Data {
		// Find the matching computed serie (by Extra == serie id).
		var rs *scales.ComputedSerie
		for i := range xy.Series {
			if id, ok := xy.Series[i].Extra.(string); ok && id == s.ID {
				rs = &xy.Series[i]
				break
			}
		}
		if rs == nil {
			continue
		}
		color := getColor(s.ID)
		for j, d := range rs.Data {
			nodes = append(nodes, ComputedNode{
				ID:         s.ID + "." + strconv.Itoa(j),
				SerieID:    s.ID,
				Index:      len(nodes),
				SerieIndex: j,
				X:          ptrToFloat(d.Position.X),
				Y:          ptrToFloat(d.Position.Y),
				XValue:     s.Data[j].X,
				YValue:     s.Data[j].Y,
				FormattedX: formatX(s.Data[j].X),
				FormattedY: formatY(s.Data[j].Y),
				Size:       props.NodeSize,
				Color:      color,
			})
		}
	}

	legendData := make([]legends.Datum, 0, len(props.Data))
	for _, s := range props.Data {
		legendData = append(legendData, legends.Datum{ID: s.ID, Label: s.ID, Color: getColor(s.ID)})
	}

	return ScatterPlotResult{
		Nodes:      nodes,
		XScale:     xy.XScale,
		YScale:     xy.YScale,
		LegendData: legendData,
	}
}

// valueFormatter returns an any→string formatter. Empty spec falls back to a
// generic %v / %g formatter; a non-empty spec is a d3-format spec applied to
// numeric values.
func valueFormatter(spec string) func(any) string {
	if spec == "" {
		return func(v any) string {
			switch x := v.(type) {
			case float64:
				return strconv.FormatFloat(x, 'g', -1, 64)
			case int:
				return strconv.Itoa(x)
			case nil:
				return ""
			default:
				return fmt.Sprintf("%v", x)
			}
		}
	}
	return func(v any) string {
		if f, ok := toFloatOK(v); ok {
			return d3format.FormatString(spec, f)
		}
		return fmt.Sprintf("%v", v)
	}
}

func toFloatOK(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	}
	return 0, false
}

func ptrToFloat(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}
