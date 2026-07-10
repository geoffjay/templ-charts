package voronoi

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
)

// Defaults mirrors @nivo/voronoi defaultVoronoiProps. Fields left zero in a
// VoronoiProps fall back to these via applyDefaults.
var Defaults = VoronoiProps{
	XDomain: [2]float64{0, 1},
	YDomain: [2]float64{0, 1},

	Layers: DefaultLayers,

	EnableLinks:   core.BoolPtr(false),
	LinkLineWidth: 1,
	LinkLineColor: "#bbbbbb",

	EnableCells:     core.BoolPtr(true),
	CellLineWidth:   2,
	CellLineColor:   "#000000",
	EnableCellFill:  core.BoolPtr(false),
	CellFillOpacity: 0.85,
	Colors:          colors.OrdinalColorScaleConfig{Type: colors.OrdinalTypeScheme, Scheme: "nivo"},

	EnablePoints: core.BoolPtr(true),
	PointSize:    4,
	PointColor:   "#666666",

	Role: "img",
}
