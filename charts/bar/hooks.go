package bar

import (
	"fmt"

	"github.com/geoffjay/templ-charts/charts/bar/compute"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// UseBar mirrors @nivo/bar useBar: resolves accessors, color scales, value
// formatter, generates bars (grouped/stacked), legend data, totals, and the
// shouldRenderLabel predicate. It merges `props` with Defaults so callers can
// pass a partial BarProps.
func UseBar(props BarProps) BarResult {
	// --- merge with defaults ---
	indexBy := props.IndexBy
	if indexBy == nil {
		indexBy = Defaults.IndexBy
	}
	keys := props.Keys
	if len(keys) == 0 {
		keys = Defaults.Keys
	}
	groupMode := props.GroupMode
	if groupMode == "" {
		groupMode = Defaults.GroupMode
	}
	layout := props.Layout
	if layout == "" {
		layout = Defaults.Layout
	}
	padding := props.Padding
	if padding == 0 {
		padding = Defaults.Padding
	}
	innerPadding := props.InnerPadding
	valueScale := props.ValueScale
	if valueScale == (scales.ScaleLinearSpec{}) {
		valueScale = Defaults.ValueScale
	}
	indexScale := props.IndexScale
	if indexScale == (scales.ScaleBandSpec{}) {
		indexScale = Defaults.IndexScale
	}
	enableLabel := props.EnableLabel || Defaults.EnableLabel
	labelSkipWidth := props.LabelSkipWidth
	labelSkipHeight := props.LabelSkipHeight
	labelPosition := props.LabelPosition
	if labelPosition == "" {
		labelPosition = Defaults.LabelPosition
	}
	colorBy := props.ColorBy
	if colorBy == "" {
		colorBy = Defaults.ColorBy
	}
	colorsCfg := props.Colors
	if colorsCfg.Type == 0 && colorsCfg.Static == "" && colorsCfg.Scheme == "" && len(colorsCfg.Colors) == 0 && colorsCfg.Func == nil && colorsCfg.DatumPath == "" {
		colorsCfg = Defaults.Colors
	}
	borderColor := props.BorderColor
	if borderColor.Type == 0 && borderColor.Static == "" && borderColor.ThemePath == "" && borderColor.FromPath == "" && borderColor.Func == nil {
		borderColor = Defaults.BorderColor
	}
	labelTextColor := props.LabelTextColor
	if labelTextColor.Type == 0 && labelTextColor.Static == "" && labelTextColor.ThemePath == "" && labelTextColor.FromPath == "" && labelTextColor.Func == nil {
		labelTextColor = Defaults.LabelTextColor
	}
	initialHidden := props.InitialHiddenIDs
	if initialHidden == nil {
		initialHidden = Defaults.InitialHiddenIDs
	}
	totalsOffset := props.TotalsOffset
	if totalsOffset == 0 {
		totalsOffset = Defaults.TotalsOffset
	}

	theme := props.Theme
	if theme == nil {
		theme = &theming.DefaultTheme
	}

	// --- accessors ---
	getIndex := core.GetPropertyAccessor[BarDatum, string](indexBy)
	getLabel := core.GetPropertyAccessor[ComputedDatum, string](props.Label)
	if props.Label == nil {
		getLabel = core.GetPropertyAccessor[ComputedDatum, string](Defaults.Label)
	}
	getTooltipLabel := core.GetPropertyAccessor[ComputedDatum, string](props.TooltipLabel)
	if props.TooltipLabel == nil {
		getTooltipLabel = func(d ComputedDatum) string {
			return fmt.Sprintf("%s - %s", d.ID, d.IndexValue)
		}
	}
	formatValue := core.GetValueFormatter[float64](props.ValueFormat)
	if props.ValueFormat == nil {
		formatValue = core.GetValueFormatter[float64](nil)
	}

	// --- color scales ---
	var getColor func(ComputedDatum) string
	{
		var identity any
		if colorBy == ColorByIndexValue {
			identity = func(d ComputedDatum) string { return d.IndexValue }
		} else {
			identity = func(d ComputedDatum) string { return d.ID }
		}
		getColor = colors.GetOrdinalColorScale[ComputedDatum](colorsCfg, identity)
	}
	getBorderColor := colors.GetInheritedColorGenerator(borderColor, theme)
	getLabelColor := colors.GetInheritedColorGenerator(labelTextColor, theme)

	// --- generate bars ---
	margin := compute.MarginLike{
		Top: props.Margin.Top, Right: props.Margin.Right,
		Bottom: props.Margin.Bottom, Left: props.Margin.Left,
	}
	layoutStr := string(layout)
	groupModeStr := string(groupMode)

	var bars []ComputedBarDatum
	var xScale, yScale scales.Scale
	if groupMode == GroupModeGrouped {
		res := compute.GenerateGroupedBars(compute.GroupedParams{
			Data: props.Data, Keys: keys, Layout: layoutStr,
			Width: props.Width, Height: props.Height,
			Padding: padding, InnerPadding: innerPadding,
			ValueScale: valueScale, IndexScale: indexScale,
			GetIndex: getIndex, GetColor: getColor,
			GetTooltipLabel: getTooltipLabel, FormatValue: formatValue,
			Margin: margin, HiddenIDs: initialHidden,
		})
		bars, xScale, yScale = res.Bars, res.XScale, res.YScale
	} else {
		res := compute.GenerateStackedBars(compute.StackedParams{
			Data: props.Data, Keys: keys, Layout: layoutStr,
			Width: props.Width, Height: props.Height,
			Padding: padding, InnerPadding: innerPadding,
			ValueScale: valueScale, IndexScale: indexScale,
			GetIndex: getIndex, GetColor: getColor,
			GetTooltipLabel: getTooltipLabel, FormatValue: formatValue,
			Margin: margin, HiddenIDs: initialHidden,
		})
		bars, xScale, yScale = res.Bars, res.XScale, res.YScale
	}

	// --- barsWithValue (filter nil values) ---
	barsWithValue := make([]ComputedBarDatum, 0, len(bars))
	for _, b := range bars {
		if b.Data.Value != nil {
			barsWithValue = append(barsWithValue, b)
		}
	}

	// --- shouldRenderLabel ---
	shouldRenderLabel := func(w, h float64) bool {
		if !enableLabel {
			return false
		}
		if labelSkipWidth > 0 && w < labelSkipWidth {
			return false
		}
		if labelSkipHeight > 0 && h < labelSkipHeight {
			return false
		}
		return true
	}

	// --- legend data ---
	// nivo builds legendData from the keys+bars, then per-legend calls
	// getLegendData with dataFrom/direction/groupMode/layout/reverse.
	legendData := buildLegendDataForKeys(keys, bars, initialHidden)
	reverse := valueScale.Reverse
	legendsWithData := make([]LegendWithData, 0, len(props.Legends))
	for _, legend := range props.Legends {
		from := legend.DataFrom
		if from == "" {
			from = "keys"
		}
		legendLabel := core.GetPropertyAccessor[map[string]any, string](props.LegendLabel)
		if props.LegendLabel == nil {
			if from == "indexes" {
				legendLabel = func(d map[string]any) string {
					if v, ok := d["indexValue"].(string); ok {
						return v
					}
					return fmt.Sprintf("%v", d["indexValue"])
				}
			} else {
				legendLabel = func(d map[string]any) string {
					if v, ok := d["id"].(string); ok {
						return v
					}
					return fmt.Sprintf("%v", d["id"])
				}
			}
		}
		// Build the source bars slice for getLegendData. For dataFrom=keys,
		// mirror nivo's legendData useMemo: one bar-like entry per key —
		// including hidden keys, which have no generated bars — so a toggled-
		// off series stays in the legend (dimmed) and can be toggled back on.
		srcBars := bars
		if from == "keys" {
			srcBars = legendBarsForKeys(keys, bars, legendData)
		}
		data := compute.GetLegendData(srcBars, from, string(legend.Direction), groupModeStr, layoutStr, reverse, legendLabel)
		legendsWithData = append(legendsWithData, LegendWithData{Props: legend, Data: data})
	}

	// --- totals ---
	var barTotals []BarTotalData
	if props.EnableTotals {
		barTotals = compute.ComputeBarTotals(bars, xScale, yScale, layoutStr, groupModeStr, totalsOffset, func(v float64) string {
			return formatValue(v)
		})
	}

	// --- label layout ---
	computeLabelLayout := compute.ComputeLabelLayout(layoutStr, valueScale.Reverse, string(labelPosition), props.LabelOffset)

	return BarResult{
		Bars:               bars,
		BarsWithValue:      barsWithValue,
		XScale:             xScale,
		YScale:             yScale,
		GetIndex:           getIndex,
		GetLabel:           getLabel,
		GetTooltipLabel:    getTooltipLabel,
		FormatValue:        formatValue,
		GetColor:           getColor,
		GetBorderColor:     getBorderColor,
		GetLabelColor:      getLabelColor,
		ShouldRenderLabel:  shouldRenderLabel,
		HiddenIDs:          initialHidden,
		LegendsWithData:    legendsWithData,
		BarTotals:          barTotals,
		ComputeLabelLayout: computeLabelLayout,
	}
}

// buildLegendDataForKeys mirrors nivo's `legendData` useMemo: per-key, find
// the first bar with data.id==key and emit {id, label, hidden, color}.
func buildLegendDataForKeys(keys []string, bars []ComputedBarDatum, hidden []string) []LegendData {
	out := make([]LegendData, 0, len(keys))
	for _, key := range keys {
		var color string
		for _, b := range bars {
			if b.Data.ID == key {
				color = b.Color
				break
			}
		}
		out = append(out, LegendData{
			ID:     key,
			Label:  key,
			Hidden: containsString(hidden, key),
			Color:  color,
		})
	}
	return out
}

// legendBarsForKeys mirrors nivo's legendData useMemo: one bar-like entry per
// key, in key order. Visible keys copy their first generated bar (so the raw
// datum stays available to custom LegendLabel accessors); hidden keys — which
// have no generated bars — get a synthetic entry carrying just the id, hidden
// flag, and legend color, keeping them present (and toggleable) in the legend.
func legendBarsForKeys(keys []string, bars []ComputedBarDatum, legendData []LegendData) []ComputedBarDatum {
	out := make([]ComputedBarDatum, 0, len(keys))
	for i, key := range keys {
		var entry ComputedBarDatum
		for j := range bars {
			if bars[j].Data.ID == key {
				entry = bars[j]
				break
			}
		}
		entry.Data.ID = key
		if i < len(legendData) {
			entry.Data.Hidden = legendData[i].Hidden
			entry.Color = legendData[i].Color
		}
		out = append(out, entry)
	}
	return out
}

func containsString(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
