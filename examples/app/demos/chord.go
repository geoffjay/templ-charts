package demos

import (
	"github.com/geoffjay/templ-charts/charts/chord"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/samples"
)

// ChordDemo is one chord tile on the /chord page (static SVG; the d3-chord
// layout is deterministic, so the render is stable).
type ChordDemo struct {
	ID          string
	Title       string
	Description string
	Props       chord.ChordProps
}

// chordSample is the shared chord flow matrix + keys, sourced from the public
// samples package.
func chordSample() ([][]float64, []string) {
	return samples.Chord()
}

// ChordDemos returns the chord demos for the /chord page: the default diagram,
// a padded/thicker-ring variant, and an interactive tile with a legend.
func ChordDemos() []ChordDemo {
	matrix, keys := chordSample()
	const w, h = 520.0, 520.0
	margin := core.Margin{Top: 60, Right: 60, Bottom: 60, Left: 60}

	base := func() chord.ChordProps {
		return chord.ChordProps{
			Width: w, Height: h, Margin: margin,
			Data: matrix, Keys: keys,
		}
	}

	basic := base()

	padded := base()
	padded.PadAngle = 0.04
	padded.InnerRadiusRatio = 0.86
	padded.InnerRadiusOffset = 0.04

	interactive := base()
	interactive.Interactive = true
	// A slightly thicker ring + inset ribbons make arcs easy to hover.
	interactive.InnerRadiusRatio = 0.8
	interactive.InnerRadiusOffset = 0.02
	interactive.Legends = []legends.LegendProps{{
		Anchor:     legends.LegendAnchorBottom,
		Direction:  legends.LegendDirectionRow,
		TranslateY: 56,
		ItemWidth:  90,
		ItemHeight: 16,
		SymbolSize: 12,
	}}
	interactive.Margin = core.Margin{Top: 60, Right: 60, Bottom: 90, Left: 60}

	return []ChordDemo{
		{
			ID:          "chord-basic",
			Title:       "Default chord diagram",
			Description: "The default layout: entities are arcs sized by total flow (internal/d3/chord, a faithful d3-chord port), with pairwise flows drawn as ribbons whose ends span each directed sub-flow.",
			Props:       basic,
		},
		{
			ID:          "chord-padded",
			Title:       "Padded arcs + inset ribbons",
			Description: "padAngle inserts gaps between the group arcs, innerRadiusRatio thickens the ring, and innerRadiusOffset pulls the ribbons inward so they don't touch the arcs.",
			Props:       padded,
		},
		{
			ID:          "chord-interactive",
			Title:       "Hover to highlight",
			Description: "Hover an entity's arc (or a ribbon) and it stays lit while every other arc and ribbon fades back — the hovered entity reads as the only selection. Pure CSS (:has hover) driven by the active/inactive opacity props, plus a client-side tooltip (charts/interact) and a box legend.",
			Props:       interactive,
		},
	}
}
