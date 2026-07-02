package swarmplot

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/internal/d3/force"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
)

// UseSwarmPlot mirrors @nivo/swarmplot's useSwarmPlot: it builds a value scale
// (along the layout axis) and a band group scale (across it), seeds a force
// simulation with ForceX/ForceY toward those targets plus a ForceCollide to
// avoid overlap, runs a fixed iteration count, then reads back the node
// positions and assigns per-group (or per-id) ordinal colors. props.Width/Height
// are the inner dimensions.
func UseSwarmPlot(props SwarmPlotProps) SwarmPlotResult {
	horizontal := props.Layout == "horizontal"

	groups := props.Groups
	if len(groups) == 0 {
		groups = deriveGroups(props.Data)
	}

	// Value scale runs along the layout axis; the group (band) scale across it.
	valueAxisSize := props.Height
	groupAxisSize := props.Width
	scaleAxis := scales.ScaleAxisY
	if horizontal {
		valueAxisSize = props.Width
		groupAxisSize = props.Height
		scaleAxis = scales.ScaleAxisX
	}

	minV, maxV := valueExtent(props.Data)
	all := make([]any, len(props.Data))
	for i, d := range props.Data {
		all[i] = d.Value
	}
	valueScale := scales.ComputeScale(
		props.ValueScale,
		scales.ComputedSerieAxis{All: all, Min: minV, Max: maxV},
		valueAxisSize, scaleAxis,
	)

	bandScale := scales.NewBandScaleWithRange(groups, 0, groupAxisSize, 0, false)
	bandwidth := 0.0
	if bw, ok := bandScale.(scales.ScaleWithBandwidth); ok {
		bandwidth = bw.Bandwidth()
	}
	ordinalCenter := func(group string) float64 { return bandScale.Call(group) + bandwidth/2 }

	// Build simulation nodes carrying each datum.
	fnodes := make([]*force.Node, len(props.Data))
	for i, d := range props.Data {
		fnodes[i] = force.NewNode(d)
	}

	valueTarget := func(n *force.Node) float64 { return valueScale.Call(n.Data.(SwarmPlotDatum).Value) }
	groupTarget := func(n *force.Node) float64 { return ordinalCenter(n.Data.(SwarmPlotDatum).Group) }
	radius := func(*force.Node) float64 { return props.Size/2 + props.Spacing/2 }

	var xForce *force.XForce
	var yForce *force.YForce
	if horizontal {
		xForce = force.ForceX(valueTarget).Strength(props.ForceStrength)
		yForce = force.ForceY(groupTarget)
	} else {
		xForce = force.ForceX(groupTarget)
		yForce = force.ForceY(valueTarget).Strength(props.ForceStrength)
	}

	sim := force.NewSimulation(fnodes).
		Force("x", xForce).
		Force("y", yForce).
		Force("collide", force.ForceCollide(radius)).
		Stop()
	sim.Tick(props.SimulationIterations)

	getColorID := func(d SwarmPlotDatum) string {
		if props.ColorBy == "id" {
			return d.ID
		}
		return d.Group
	}
	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(k string) string { return k })
	format := valueFormatter(props.ValueFormat)

	nodes := make([]ComputedNode, len(fnodes))
	for i, fn := range sim.Nodes() {
		d := fn.Data.(SwarmPlotDatum)
		nodes[i] = ComputedNode{
			ID:             d.ID,
			Group:          d.Group,
			Value:          d.Value,
			FormattedValue: format(d.Value),
			X:              fn.X,
			Y:              fn.Y,
			Size:           props.Size,
			Color:          getColor(getColorID(d)),
		}
	}

	// For vertical layout the value scale is y and the group scale is x; for
	// horizontal it's the reverse (mirrors nivo's config map).
	xScale, yScale := bandScale, valueScale
	if horizontal {
		xScale, yScale = valueScale, bandScale
	}

	return SwarmPlotResult{Nodes: nodes, XScale: xScale, YScale: yScale}
}

func deriveGroups(data []SwarmPlotDatum) []string {
	seen := make(map[string]bool)
	groups := make([]string, 0)
	for _, d := range data {
		if !seen[d.Group] {
			seen[d.Group] = true
			groups = append(groups, d.Group)
		}
	}
	return groups
}

func valueExtent(data []SwarmPlotDatum) (min, max float64) {
	if len(data) == 0 {
		return 0, 0
	}
	min, max = data[0].Value, data[0].Value
	for _, d := range data[1:] {
		if d.Value < min {
			min = d.Value
		}
		if d.Value > max {
			max = d.Value
		}
	}
	return min, max
}

func valueFormatter(spec string) func(float64) string {
	if spec == "" {
		return func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	}
	return func(v float64) string { return d3format.FormatString(spec, v) }
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
