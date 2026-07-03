// Package geo mirrors @nivo/geo's two SVG components — GeoMap and Choropleth —
// rendering GeoJSON features through a spherical projection to server-side SVG.
// The projection math, GeoJSON→path generation, and graticule come from
// internal/d3/geo (a faithful d3-geo subset); colors/legends/theming/core are
// reused from the existing chart foundation, exactly as nivo composes them.
//
// GeoMap paints every feature a single fill color. Choropleth binds a
// [{id,value}] dataset onto the features by id and colors each by a quantize
// color scale (nivo's default scheme "PuBuGn" maps to this repo's
// "purple_blue_green" sequential scheme).
//
// v3 scope (see docs/PLAN-v3.md §3.6 and NOTES.md): SVG only, static render.
// Cylindrical/pseudocylindrical projections (mercator, equirectangular,
// transverseMercator, naturalEarth1, equalEarth) are fully correct;
// azimuthal-family projections render the whole sphere (clipCircle is
// deferred). Per-feature hover tooltips are available via Interactive; HTMX
// hover-others dimming and interactive zoom/pan are deferred.
package geo

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3geo "github.com/geoffjay/templ-charts/internal/d3/geo"
)

// Re-exported GeoJSON types. internal/d3/geo is not importable by library
// consumers, so callers construct features via these aliases.
type (
	Feature           = d3geo.Feature
	Geometry          = d3geo.Geometry
	FeatureCollection = d3geo.FeatureCollection
)

// Geometry type tags (re-exported for building Geometry values). See
// internal/d3/geo for the concrete Coordinates shape each tag expects.
const (
	TypePoint           = d3geo.TypePoint
	TypeMultiPoint      = d3geo.TypeMultiPoint
	TypeLineString      = d3geo.TypeLineString
	TypeMultiLineString = d3geo.TypeMultiLineString
	TypePolygon         = d3geo.TypePolygon
	TypeMultiPolygon    = d3geo.TypeMultiPolygon
	TypeSphere          = d3geo.TypeSphere
)

// GeoLayerId enumerates the render layers. Mirrors @nivo/geo layers.
type GeoLayerId string

const (
	GeoLayerGraticule GeoLayerId = "graticule"
	GeoLayerFeatures  GeoLayerId = "features"
	GeoLayerLegends   GeoLayerId = "legends"
)

// DefaultGeoMapLayers / DefaultChoroplethLayers mirror nivo's defaults.
var (
	DefaultGeoMapLayers     = []GeoLayerId{GeoLayerGraticule, GeoLayerFeatures}
	DefaultChoroplethLayers = []GeoLayerId{GeoLayerGraticule, GeoLayerFeatures, GeoLayerLegends}
)

// GeoBase holds the projection, border, graticule, dimensions, and
// accessibility props shared by GeoMap and Choropleth.
type GeoBase struct {
	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	ProjectionType        string     // mercator | equirectangular | … (see internal/d3/geo)
	ProjectionScale       float64    // d3 projection scale
	ProjectionTranslation [2]float64 // fractions of inner width/height (default [0.5,0.5])
	ProjectionRotation    [3]float64 // [λ,φ,γ] degrees

	EnableGraticule    bool
	GraticuleLineWidth float64
	GraticuleLineColor string

	BorderWidth float64
	BorderColor string

	// Interactive enables per-feature client-side hover tooltips (charts/interact).
	Interactive bool

	Theme  *theming.Theme
	Layers []GeoLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool
}

// GeoMapProps mirrors @nivo/geo GeoMapSvgProps (supported subset).
type GeoMapProps struct {
	Features  []Feature
	FillColor string
	GeoBase
}

// ChoroplethDatum is one {id, value} row bound onto a feature by id.
type ChoroplethDatum struct {
	ID    string
	Value float64
}

// ChoroplethProps mirrors @nivo/geo ChoroplethSvgProps (supported subset).
// Features are matched to Data by feature id (nivo's match:'id' default).
type ChoroplethProps struct {
	Features []Feature
	Data     []ChoroplethDatum

	// Colors is the quantize color scheme id (nivo "PuBuGn" ↔ this repo's
	// "purple_blue_green"); Steps is the number of quantization bands.
	Colors string
	Steps  int
	// Domain is the value domain; when zero it falls back to the data min/max.
	Domain       [2]float64
	UnknownColor string
	ValueFormat  string

	Legends []legends.LegendProps

	GeoBase
}

// ComputedFeature is one projected, colored feature ready to render.
type ComputedFeature struct {
	ID          string
	Path        string
	FillColor   string
	BorderWidth float64
	BorderColor string

	// Choropleth-only: whether the feature matched a datum, and its value.
	HasValue       bool
	Value          float64
	FormattedValue string
	Label          string
}

// GeoResult is the computed model produced by UseGeoMap / UseChoropleth.
type GeoResult struct {
	Features      []ComputedFeature
	GraticulePath string

	// Choropleth-only:
	ColorScale func(float64) string
	ValueMin   float64
	ValueMax   float64
}
