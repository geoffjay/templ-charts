package compute

// GetLegendData mirrors @nivo/bar getLegendData: dispatches on `from`
// ("keys"|"indexes") and applies the order-reversal rules nivo uses so legend
// item order matches the visual stacking/grouping order.
func GetLegendData(
	bars []ComputedBarDatum,
	from string,
	direction string,
	groupMode string,
	layout string,
	reverse bool,
	legendLabel func(map[string]any) string,
) []LegendData {
	if from == "indexes" {
		return getLegendDataForIndexes(bars, layout, legendLabel)
	}
	return getLegendDataForKeys(bars, layout, direction, groupMode, reverse, legendLabel)
}

// getLegendDataForKeys mirrors @nivo/bar getLegendDataForKeys.
func getLegendDataForKeys(
	bars []ComputedBarDatum,
	layout string,
	direction string,
	groupMode string,
	reverse bool,
	getLabel func(map[string]any) string,
) []LegendData {
	// Build unique-by-id list preserving first-seen order.
	seen := map[string]bool{}
	data := []LegendData{}
	for _, b := range bars {
		if seen[b.Data.ID] {
			continue
		}
		seen[b.Data.ID] = true
		color := b.Color
		if color == "" {
			color = "#000"
		}
		data = append(data, LegendData{
			ID:     b.Data.ID,
			Label:  getLabel(legendLabelDatum(b)),
			Hidden: b.Data.Hidden,
			Color:  color,
		})
	}
	// Reversal rules verbatim from nivo:
	if (layout == "vertical" && groupMode == "stacked" && direction == "column" && !reverse) ||
		(layout == "horizontal" && groupMode == "stacked" && reverse) {
		reverseLegendData(data)
	}
	return data
}

// getLegendDataForIndexes mirrors @nivo/bar getLegendDataForIndexes.
func getLegendDataForIndexes(
	bars []ComputedBarDatum,
	layout string,
	getLabel func(map[string]any) string,
) []LegendData {
	seen := map[string]bool{}
	data := []LegendData{}
	for _, b := range bars {
		id := b.Data.IndexValue
		if seen[id] {
			continue
		}
		seen[id] = true
		color := b.Color
		if color == "" {
			color = "#000"
		}
		data = append(data, LegendData{
			ID:     id,
			Label:  getLabel(legendLabelDatum(b)),
			Hidden: b.Data.Hidden,
			Color:  color,
		})
	}
	if layout == "horizontal" {
		reverseLegendData(data)
	}
	return data
}

// legendLabelDatum builds the map handed to the legend-label accessor,
// mirroring the fields nivo exposes on ComputedDatum (bar.data): the raw
// datum stays nested under "data", so the default accessors ("id",
// "indexValue") and custom paths like "data.<field>" both resolve.
func legendLabelDatum(b ComputedBarDatum) map[string]any {
	return map[string]any{
		"id":             b.Data.ID,
		"value":          b.Data.Value,
		"formattedValue": b.Data.FormattedValue,
		"hidden":         b.Data.Hidden,
		"indexValue":     b.Data.IndexValue,
		"data":           b.Data.Data,
	}
}

func reverseLegendData(d []LegendData) {
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
}
