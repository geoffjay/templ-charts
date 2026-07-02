package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/sankey"
)

// SankeyDemo is one sankey tile on the /sankey page (static SVG; the layout is
// deterministic, so the render is stable).
type SankeyDemo struct {
	ID          string
	Title       string
	Description string
	Props       sankey.SankeyProps
}

// sankeySample builds a small energy-flow graph: two sources feed two
// intermediates, which converge on a single sink.
func sankeySample() ([]sankey.SankeyInputNode, []sankey.SankeyInputLink) {
	nodes := []sankey.SankeyInputNode{
		{ID: "Coal"},
		{ID: "Gas"},
		{ID: "Grid"},
		{ID: "Solar"},
		{ID: "Homes"},
		{ID: "Industry"},
	}
	links := []sankey.SankeyInputLink{
		{Source: "Coal", Target: "Grid", Value: 12},
		{Source: "Gas", Target: "Grid", Value: 8},
		{Source: "Solar", Target: "Grid", Value: 5},
		{Source: "Solar", Target: "Homes", Value: 3},
		{Source: "Grid", Target: "Homes", Value: 15},
		{Source: "Grid", Target: "Industry", Value: 10},
	}
	return nodes, links
}

// SankeyDemos returns the sankey demos for the /sankey page: horizontal +
// vertical layouts, an alignment variant, and an interactive tile with a
// legend.
func SankeyDemos() []SankeyDemo {
	nodes, links := sankeySample()
	margin := core.Margin{Top: 20, Right: 100, Bottom: 20, Left: 100}
	const w, h = 700.0, 420.0

	base := func() sankey.SankeyProps {
		return sankey.SankeyProps{
			Width: w, Height: h, Margin: margin,
			Nodes: nodes, Links: links,
		}
	}

	horizontal := base()

	vertical := base()
	vertical.Layout = sankey.SankeyLayoutVertical
	vertical.Margin = core.Margin{Top: 40, Right: 40, Bottom: 40, Left: 40}

	justify := base()
	justify.Align = sankey.SankeyAlignJustify

	interactive := base()
	interactive.Interactive = true
	interactive.Legends = []legends.LegendProps{{
		Anchor:     legends.LegendAnchorBottom,
		Direction:  legends.LegendDirectionRow,
		TranslateY: 50,
		ItemWidth:  90,
		ItemHeight: 16,
		SymbolSize: 12,
	}}
	interactive.Margin = core.Margin{Top: 20, Right: 100, Bottom: 60, Left: 100}

	return []SankeyDemo{
		{
			ID:          "sankey-horizontal",
			Title:       "Horizontal (center align)",
			Description: "The default layout: flow runs left→right, nodes centered by depth, ribbons drawn with a monotone curve (internal/d3/shape) over the d3-sankey layout (internal/d3/sankey).",
			Props:       horizontal,
		},
		{
			ID:          "sankey-vertical",
			Title:       "Vertical layout",
			Description: "The same graph transposed: flow runs top→bottom. The d3-sankey coordinates are swapped exactly as nivo does, and ribbons use curveMonotoneY.",
			Props:       vertical,
		},
		{
			ID:          "sankey-justify",
			Title:       "Justify alignment",
			Description: "align: justify pins pure sinks to the last column instead of centering them, stretching the terminal ribbons.",
			Props:       justify,
		},
		{
			ID:          "sankey-interactive",
			Title:       "Interactive + legend",
			Description: "Hover a node or ribbon for a client-side tooltip (charts/interact), with a box legend of the node ids beneath.",
			Props:       interactive,
		},
	}
}
