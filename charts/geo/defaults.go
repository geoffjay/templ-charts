package geo

import "github.com/geoffjay/templ-charts/charts/core"

// commonDefaults mirrors @nivo/geo commonDefaultProps.
var commonDefaults = GeoBase{
	ProjectionType:        ProjectionMercator,
	ProjectionScale:       100,
	ProjectionTranslation: [2]float64{0.5, 0.5},
	ProjectionRotation:    [3]float64{0, 0, 0},

	EnableGraticule:    core.BoolPtr(false),
	GraticuleLineWidth: 0.5,
	GraticuleLineColor: "#999999",

	BorderWidth: 0,
	BorderColor: "#000000",

	Role: "img",
}

// GeoMapDefaults mirrors @nivo/geo GeoMapDefaultProps.
var GeoMapDefaults = GeoMapProps{
	FillColor: "#dddddd",
	GeoBase:   withLayers(commonDefaults, DefaultGeoMapLayers),
}

// ChoroplethDefaults mirrors @nivo/geo ChoroplethDefaultProps. nivo's default
// colors "PuBuGn" is this repo's "purple_blue_green" sequential scheme.
var ChoroplethDefaults = ChoroplethProps{
	Colors:       "purple_blue_green",
	Steps:        7,
	UnknownColor: "#999999",
	GeoBase:      withLayers(commonDefaults, DefaultChoroplethLayers),
}

func withLayers(b GeoBase, layers []GeoLayerId) GeoBase {
	b.Layers = layers
	return b
}
