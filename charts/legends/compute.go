package legends

// ComputePositionFromAnchor returns the (x, y) translate offset that pins a
// legend box of (width, height) to the given anchor of a chart of
// (chartWidth, chartHeight). Mirrors @nivo/legends computePositionFromAnchor.
func ComputePositionFromAnchor(anchor LegendAnchor, chartWidth, chartHeight, width, height, translateX, translateY float64) (float64, float64) {
	var x, y float64
	switch anchor {
	case LegendAnchorTopLeft:
		x, y = 0, 0
	case LegendAnchorTop:
		x = (chartWidth - width) / 2
		y = 0
	case LegendAnchorTopRight:
		x = chartWidth - width
		y = 0
	case LegendAnchorRight:
		x = chartWidth - width
		y = (chartHeight - height) / 2
	case LegendAnchorBottomRight:
		x = chartWidth - width
		y = chartHeight - height
	case LegendAnchorBottom:
		x = (chartWidth - width) / 2
		y = chartHeight - height
	case LegendAnchorBottomLeft:
		x, y = 0, chartHeight-height
	case LegendAnchorLeft:
		x = 0
		y = (chartHeight - height) / 2
	}
	return x + translateX, y + translateY
}

// ComputeDimensions returns the (width, height) of a legend box given its
// direction, item count, and item dimensions. Mirrors @nivo/legends
// computeDimensions.
func ComputeDimensions(props LegendProps) (float64, float64) {
	n := len(props.Items)
	if n == 0 {
		return 0, 0
	}
	itemW := props.ItemWidth
	itemH := props.ItemHeight
	if itemW == 0 {
		itemW = LegendDefaults.ItemWidth
	}
	if itemH == 0 {
		itemH = LegendDefaults.ItemHeight
	}
	switch props.Direction {
	case LegendDirectionRow:
		return float64(n) * itemW, itemH
	default: // column
		return itemW, float64(n) * itemH
	}
}

// ItemLayout is the (x, y) position of one legend item inside the legend box.
type ItemLayout struct {
	X, Y float64
}

// ComputeItemLayout returns the positions of each legend item within the box.
// Mirrors @nivo/legends computeItemLayout.
func ComputeItemLayout(props LegendProps) []ItemLayout {
	n := len(props.Items)
	out := make([]ItemLayout, n)
	itemW := props.ItemWidth
	itemH := props.ItemHeight
	if itemW == 0 {
		itemW = LegendDefaults.ItemWidth
	}
	if itemH == 0 {
		itemH = LegendDefaults.ItemHeight
	}
	switch props.Direction {
	case LegendDirectionRow:
		for i := range props.Items {
			out[i] = ItemLayout{X: float64(i) * itemW, Y: 0}
		}
	default:
		for i := range props.Items {
			out[i] = ItemLayout{X: 0, Y: float64(i) * itemH}
		}
	}
	return out
}

// ComputeContinuousColorsLegend produces the discrete color samples for a
// continuous legend gradient bar. Returns a slice of (offsetPct, color)
// pairs for the bar's stops. Mirrors @nivo/legends
// computeContinuousColorsLegend.
func ComputeContinuousColorsLegend(scale func(float64) string, min, max float64, samples int) []ContinuousColorStop {
	if samples <= 0 {
		samples = ContinuousColorsLegendDefaults.Samples
	}
	if samples < 2 {
		samples = 2
	}
	span := max - min
	if span == 0 {
		span = 1
	}
	out := make([]ContinuousColorStop, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(samples-1)
		value := min + t*span
		out[i] = ContinuousColorStop{Offset: t, Color: scale(value)}
	}
	return out
}

// ContinuousColorStop is one sample along a continuous color legend bar.
type ContinuousColorStop struct {
	Offset float64 // 0..1
	Color  string
}
