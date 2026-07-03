// Package chord mirrors @nivo/chord: a radial diagram of the flows between a
// set of entities. Given a square flow matrix and a list of keys, entities are
// laid out as arcs around a circle (sized by total flow) and the pairwise flows
// are drawn as ribbons connecting them. The layout is produced server-side by
// internal/d3/chord (a faithful d3-chord port); arcs are rendered with
// charts/arcs (internal/d3/shape's Arc) and ribbons with internal/d3/chord's
// Ribbon generator, exactly as nivo does.
//
// It reuses internal/d3/chord (layout + ribbon), charts/arcs (group arc paths),
// charts/colors (ordinal arc colors + inherited border/label colors),
// charts/core (SvgWrapper + dimensions), charts/legends, charts/theming and
// charts/interact (optional client-side hover tooltips).
//
// v3 scope: SVG only, static render. The layout is deterministic, so goldens
// are byte-stable. HTMX active-arc/ribbon dimming and animated transitions are
// deferred; per-arc/ribbon hover tooltips are available via Interactive.
package chord

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// ComputedArc is one positioned, colored group arc ready to render. Angles are
// in radians (d3 convention, arc generator applies the -π/2 offset).
type ComputedArc struct {
	ID             string
	Index          int
	Label          string
	Value          float64
	FormattedValue string
	StartAngle     float64
	EndAngle       float64
	Color          string
	Path           string
}

// ComputedRibbon is one positioned ribbon connecting two arcs. Path is the SVG
// path string (relative to the circle center); Color is the source arc color
// (the ribbon border color is derived from it at render time).
type ComputedRibbon struct {
	ID     string
	Source ComputedArc
	Target ComputedArc
	Path   string
	Color  string
}

// ChordLayerId enumerates the render layers. Mirrors @nivo/chord LayerId.
type ChordLayerId string

const (
	ChordLayerRibbons ChordLayerId = "ribbons"
	ChordLayerArcs    ChordLayerId = "arcs"
	ChordLayerLabels  ChordLayerId = "labels"
	ChordLayerLegends ChordLayerId = "legends"
)

// DefaultLayers mirrors @nivo/chord svgDefaultProps.layers.
var DefaultLayers = []ChordLayerId{ChordLayerRibbons, ChordLayerArcs, ChordLayerLabels, ChordLayerLegends}

// ChordProps mirrors @nivo/chord ChordSvgProps (the supported subset). Fields
// left zero fall back to Defaults via applyDefaults.
type ChordProps struct {
	// Data is the square flow matrix; Data[i][j] is the flow from Keys[i] to
	// Keys[j]. Keys names each row/column and is used as the arc identity.
	Data [][]float64
	Keys []string

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	PadAngle          float64 // radians between adjacent group arcs
	InnerRadiusRatio  float64
	InnerRadiusOffset float64

	Colors colors.OrdinalColorScaleConfig

	ArcOpacity         float64
	ActiveArcOpacity   float64
	InactiveArcOpacity float64
	ArcBorderWidth     float64
	ArcBorderColor     colors.InheritedColorConfig

	RibbonOpacity         float64
	ActiveRibbonOpacity   float64
	InactiveRibbonOpacity float64
	RibbonBorderWidth     float64
	RibbonBorderColor     colors.InheritedColorConfig
	RibbonBlendMode       string

	EnableLabel    *bool // nil → true
	Label          string
	LabelOffset    float64
	LabelRotation  float64
	LabelTextColor colors.InheritedColorConfig

	ValueFormat string // d3-format spec; empty → %g

	// Interactive enables per-arc/ribbon client-side hover tooltips (charts/interact).
	Interactive bool

	Legends []legends.LegendProps

	Theme  *theming.Theme
	Layers []ChordLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool
}

// LabelEnabled resolves EnableLabel (nil → true, matching nivo's default).
func (p ChordProps) LabelEnabled() bool { return p.EnableLabel == nil || *p.EnableLabel }

// ChordResult is the computed model produced by UseChord.
type ChordResult struct {
	Arcs       []ComputedArc
	Ribbons    []ComputedRibbon
	Center     [2]float64
	Radius     float64
	LegendData []legends.Datum
}
