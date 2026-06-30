package boxplot

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3array "github.com/geoffjay/templ-charts/internal/d3/array"
)

// UseBoxPlot mirrors @nivo/boxplot's stratify → summarize → generate pipeline:
// it groups the raw data, summarizes each distribution to quantiles, builds the
// band index scale + linear value scale, and lays out each box-plot glyph.
// props.Width/Height are the inner (margin-subtracted) dimensions.
func UseBoxPlot(props BoxPlotProps) BoxPlotResult {
	theme := resolveTheme(props.Theme)
	vertical := props.Layout != BoxPlotLayoutHorizontal

	groups := props.Groups
	if groups == nil {
		groups = uniqueFirstSeen(props.Data, func(d BoxPlotDatum) string { return d.Group })
	}
	subGroups := props.SubGroups
	if subGroups == nil {
		subGroups = uniqueFirstSeen(props.Data, func(d BoxPlotDatum) string { return d.SubGroup })
	}
	hasSubGroups := len(subGroups) > 0 && !(len(subGroups) == 1 && subGroups[0] == "")
	nSub := len(subGroups)
	if nSub == 0 {
		nSub = 1
		subGroups = []string{""}
	}

	groupIdx := indexMap(groups)
	subIdx := indexMap(subGroups)

	// Stratify into (group, subGroup) buckets.
	buckets := make(map[[2]int][]float64)
	for _, d := range props.Data {
		gi := groupIdx[d.Group]
		si := 0
		if hasSubGroups {
			si = subIdx[d.SubGroup]
		}
		buckets[[2]int{gi, si}] = append(buckets[[2]int{gi, si}], d.Value)
	}

	// Summarize each non-empty bucket.
	quantiles := props.Quantiles
	summaries := make([]BoxPlotSummary, 0, len(buckets))
	allValues := make([]float64, 0)
	for gi := range groups {
		for si := 0; si < nSub; si++ {
			vals := buckets[[2]int{gi, si}]
			if len(vals) == 0 {
				continue
			}
			sorted := append([]float64(nil), vals...)
			sort.Float64s(sorted)
			qVals := make([]float64, len(quantiles))
			for i, q := range quantiles {
				qVals[i] = d3array.QuantileSorted(sorted, q)
			}
			sub := ""
			if hasSubGroups {
				sub = subGroups[si]
			}
			s := BoxPlotSummary{
				Group:         groups[gi],
				SubGroup:      sub,
				GroupIndex:    gi,
				SubGroupIndex: si,
				N:             len(sorted),
				Extrema:       [2]float64{sorted[0], sorted[len(sorted)-1]},
				Quantiles:     quantiles,
				Values:        qVals,
				Mean:          mean(sorted),
			}
			summaries = append(summaries, s)
			allValues = append(allValues, qVals...)
		}
	}

	// Value domain.
	minV, maxV := math.Inf(1), math.Inf(-1)
	for _, v := range allValues {
		minV = math.Min(minV, v)
		maxV = math.Max(maxV, v)
	}
	if math.IsInf(minV, 1) {
		minV, maxV = 0, 0
	}
	if props.MinValue != nil {
		minV = *props.MinValue
	}
	if props.MaxValue != nil {
		maxV = *props.MaxValue
	}

	// Scales.
	var indexSize, valueSize float64
	var valueAxis scales.ScaleAxis
	if vertical {
		indexSize = props.Width
		valueSize = props.Height
		valueAxis = scales.ScaleAxisY
	} else {
		indexSize = props.Height
		valueSize = props.Width
		valueAxis = scales.ScaleAxisX
	}
	indexScale := scales.NewBandScaleWithRange(groups, 0, indexSize, props.Padding, true)
	valueScale := scales.ComputeScale(
		scales.ScaleLinearSpec{Min: scales.FloatVal(minV), Max: scales.FloatVal(maxV)},
		scales.ComputedSerieAxis{All: []any{minV, maxV}, Min: minV, Max: maxV},
		valueSize, valueAxis,
	)

	indexBandwidth := 0.0
	if bw, ok := indexScale.(scales.ScaleWithBandwidth); ok {
		indexBandwidth = bw.Bandwidth()
	}
	bandwidth := (indexBandwidth - props.InnerPadding*float64(nSub-1)) / float64(nSub)

	// Colors.
	colorBy := props.ColorBy
	getColor := colors.GetOrdinalColorScale[BoxPlotSummary](props.Colors, func(s BoxPlotSummary) string {
		if colorBy == "group" {
			return s.Group
		}
		return s.SubGroup
	})
	getBorder := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	getMedian := colors.GetInheritedColorGenerator(props.MedianColor, theme)
	getWhisker := colors.GetInheritedColorGenerator(props.WhiskerColor, theme)

	boxes := make([]ComputedBox, 0, len(summaries))
	for _, s := range summaries {
		// The box + whisker glyph needs the five-number summary
		// (q10/q25/median/q75/q90). Skip distributions with fewer quantiles
		// rather than index out of range — callers should pass 5 quantiles.
		if len(s.Values) < 5 {
			continue
		}
		indexCoord := indexScale.Call(s.Group) + bandwidth*float64(s.SubGroupIndex) + props.InnerPadding*float64(s.SubGroupIndex)
		coords := make([]float64, len(s.Values))
		for i, v := range s.Values {
			coords[i] = valueScale.Call(v)
		}
		median := coords[2]
		color := getColor(s)
		ctx := map[string]any{"color": color}
		var transform string
		if vertical {
			transform = fmt.Sprintf("translate(%s,%s)", fmtN(indexCoord+bandwidth/2), fmtN(median))
		} else {
			transform = fmt.Sprintf("translate(%s,%s) rotate(-90)", fmtN(median), fmtN(indexCoord+bandwidth/2))
		}
		rectY := coords[3] - coords[2]
		if !vertical {
			rectY = coords[1] - coords[2]
		}
		boxes = append(boxes, ComputedBox{
			Key:           strconv.Itoa(s.GroupIndex) + "." + strconv.Itoa(s.SubGroupIndex),
			Group:         s.Group,
			SubGroup:      s.SubGroup,
			Summary:       s,
			Transform:     transform,
			Bandwidth:     bandwidth,
			Color:         color,
			BorderColor:   getBorder(ctx),
			MedianColor:   getMedian(ctx),
			WhiskerColor:  getWhisker(ctx),
			Opacity:       props.Opacity,
			RectY:         rectY,
			ValueInterval: math.Abs(coords[3] - coords[1]),
			VD0:           coords[0] - coords[2],
			VD1:           coords[1] - coords[2],
			VD3:           coords[3] - coords[2],
			VD4:           coords[4] - coords[2],
			WhiskerEnd:    props.WhiskerEndSize * bandwidth / 2,
		})
	}

	var xScale, yScale scales.Scale
	if vertical {
		xScale, yScale = indexScale, valueScale
	} else {
		xScale, yScale = valueScale, indexScale
	}

	legendData := make([]legends.Datum, 0)
	if hasSubGroups {
		for _, sg := range subGroups {
			legendData = append(legendData, legends.Datum{ID: sg, Label: sg, Color: getColor(BoxPlotSummary{SubGroup: sg})})
		}
	}

	return BoxPlotResult{
		Boxes:      boxes,
		IndexScale: indexScale,
		ValueScale: valueScale,
		XScale:     xScale,
		YScale:     yScale,
		LegendData: legendData,
	}
}

func uniqueFirstSeen(data []BoxPlotDatum, key func(BoxPlotDatum) string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, d := range data {
		k := key(d)
		if !seen[k] {
			seen[k] = true
			out = append(out, k)
		}
	}
	return out
}

func indexMap(s []string) map[string]int {
	m := make(map[string]int, len(s))
	for i, v := range s {
		m[v] = i
	}
	return m
}

func mean(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
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
