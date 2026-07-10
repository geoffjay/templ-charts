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
// SVG only, static render. All projections render correctly, including the
// azimuthal family (clipCircle hides the far hemisphere). Per-feature hover
// tooltips are available via Interactive; HTMX hover-others dimming and
// interactive zoom/pan are deferred.
package geo

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3geo "github.com/geoffjay/templ-charts/internal/d3/geo"
)

// Public GeoJSON types owned by this package. These are concrete structs (not
// aliases to internal/d3/geo) so the v1.0 API is decoupled from the internal
// projection port: consumers construct these, and the package converts them to
// the internal representation via toD3 before rendering.
//
// Geometry.Coordinates holds one of these concrete shapes depending on Type,
// using the standard GeoJSON [lon,lat] ordering (degrees):
//
//	TypePoint                     → [2]float64
//	TypeMultiPoint / LineString   → [][2]float64
//	TypeMultiLineString / Polygon → [][][2]float64
//	TypeMultiPolygon              → [][][][2]float64
//	TypeSphere                    → nil

// Geometry is a GeoJSON geometry (or Sphere). See the package Type constants
// for the concrete shape held in Coordinates per Type.
type Geometry struct {
	// Type is the GeoJSON geometry type tag; use the package Type* constants.
	Type string
	// Coordinates holds the coordinate data whose concrete Go shape depends on
	// Type (see the Geometry doc), in [lon,lat] degree ordering.
	Coordinates any
	// Geometries holds child geometries, populated only for a
	// GeometryCollection.
	Geometries []Geometry // GeometryCollection only
}

// Feature is a GeoJSON Feature wrapping a single Geometry.
type Feature struct {
	// Type is the GeoJSON object type, normally "Feature".
	Type string
	// ID is the feature identifier; Choropleth binds data rows onto features by
	// this id.
	ID string
	// Properties holds arbitrary GeoJSON feature properties.
	Properties map[string]any
	// Geometry is the feature's geometry.
	Geometry Geometry
}

// FeatureCollection is a GeoJSON FeatureCollection.
type FeatureCollection struct {
	// Type is the GeoJSON object type, normally "FeatureCollection".
	Type string
	// Features are the member features.
	Features []Feature
}

// Geometry type tags for building Geometry values. See the Geometry doc for the
// concrete Coordinates shape each tag expects.
const (
	TypePoint           = d3geo.TypePoint
	TypeMultiPoint      = d3geo.TypeMultiPoint
	TypeLineString      = d3geo.TypeLineString
	TypeMultiLineString = d3geo.TypeMultiLineString
	TypePolygon         = d3geo.TypePolygon
	TypeMultiPolygon    = d3geo.TypeMultiPolygon
	TypeSphere          = d3geo.TypeSphere
)

// toD3 converts a public Geometry to the internal d3geo.Geometry. Coordinate
// shapes are identical, so they pass through unchanged; only nested
// GeometryCollection children are converted recursively.
func (g Geometry) toD3() d3geo.Geometry {
	out := d3geo.Geometry{Type: g.Type, Coordinates: g.Coordinates}
	if len(g.Geometries) > 0 {
		out.Geometries = make([]d3geo.Geometry, len(g.Geometries))
		for i, sub := range g.Geometries {
			out.Geometries[i] = sub.toD3()
		}
	}
	return out
}

// toD3 converts a public Feature to the internal d3geo.Feature.
func (f Feature) toD3() d3geo.Feature {
	return d3geo.Feature{
		Type:       f.Type,
		ID:         f.ID,
		Properties: f.Properties,
		Geometry:   f.Geometry.toD3(),
	}
}

// featuresToD3 converts a slice of public Features to internal features.
func featuresToD3(fs []Feature) []d3geo.Feature {
	out := make([]d3geo.Feature, len(fs))
	for i, f := range fs {
		out[i] = f.toD3()
	}
	return out
}

// ProjectionType selects the spherical projection. Mirrors @nivo/geo's
// projectionType ids; unknown values fall back to mercator.
type ProjectionType string

const (
	ProjectionMercator             ProjectionType = "mercator"
	ProjectionEquirectangular      ProjectionType = "equirectangular"
	ProjectionTransverseMercator   ProjectionType = "transverseMercator"
	ProjectionNaturalEarth1        ProjectionType = "naturalEarth1"
	ProjectionEqualEarth           ProjectionType = "equalEarth"
	ProjectionAzimuthalEqualArea   ProjectionType = "azimuthalEqualArea"
	ProjectionAzimuthalEquidistant ProjectionType = "azimuthalEquidistant"
	ProjectionGnomonic             ProjectionType = "gnomonic"
	ProjectionOrthographic         ProjectionType = "orthographic"
	ProjectionStereographic        ProjectionType = "stereographic"
	ProjectionConicConformal       ProjectionType = "conicConformal"
	ProjectionConicEqualArea       ProjectionType = "conicEqualArea"
	ProjectionConicEquidistant     ProjectionType = "conicEquidistant"
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
	// Width is the total chart width in pixels (including Margin).
	Width float64
	// Height is the total chart height in pixels (including Margin).
	Height float64
	// Margin reserves space around the projected map.
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	ProjectionType        ProjectionType // see the Projection* constants
	ProjectionScale       float64        // d3 projection scale
	ProjectionTranslation [2]float64     // fractions of inner width/height (default [0.5,0.5])
	ProjectionRotation    [3]float64     // [λ,φ,γ] degrees

	// Fit, when true, auto-scales and centers the projection so the features
	// fill the inner frame (d3-geo fitExtent), ignoring ProjectionScale /
	// ProjectionTranslation. ProjectionRotation is still applied first. This
	// removes the manual scale/center fiddling a map otherwise needs.
	Fit bool

	EnableGraticule *bool // nil → false (nivo default)
	// GraticuleLineWidth is the stroke width of the graticule lines, in pixels.
	// Default 0.5.
	GraticuleLineWidth float64
	// GraticuleLineColor is the CSS color of the graticule lines. Default
	// "#999999".
	GraticuleLineColor string

	// BorderWidth is the stroke width of each feature's border, in pixels.
	// Default 0 (no border drawn).
	BorderWidth float64
	// BorderColor is the CSS color of the feature borders. Default "#000000".
	BorderColor string

	// Interactive enables per-feature client-side hover tooltips (charts/interact).
	Interactive bool

	// Theme overrides the styling theme (colors, fonts, label styles). Nil uses
	// the default theme.
	Theme *theming.Theme
	// Layers selects which render layers to draw and their order (graticule,
	// features, legends). Defaults to the per-component layer set.
	Layers []GeoLayerId

	// Role is the ARIA role of the root svg element. Default "img".
	Role string
	// AriaLabel sets aria-label on the root svg element. Empty by default.
	AriaLabel string
	// AriaLabelledBy sets aria-labelledby on the root svg element. Empty by default.
	AriaLabelledBy string
	// AriaDescribedBy sets aria-describedby on the root svg element. Empty by default.
	AriaDescribedBy string
	// Title sets the svg <title> element. Empty by default.
	Title string
	// Desc sets the svg <desc> element. Empty by default.
	Desc string
	// IsFocusable makes the root svg keyboard-focusable. Defaults to false.
	IsFocusable bool

	// Animate, when true, emits a SMIL opacity fade-in enter animation on each
	// feature path (600ms). MotionStagger delays successive features by that
	// many seconds (0 ⇒ all enter together). Defaults off, so static output is
	// unchanged. Shared by GeoMap and Choropleth via GeoBase. Mirrors @nivo/geo's
	// enter transition (opacity component).
	Animate       bool
	MotionStagger float64
}

// GeoMapProps mirrors @nivo/geo GeoMapSvgProps (supported subset).
type GeoMapProps struct {
	// Features are the GeoJSON features to project and draw.
	Features []Feature
	// FillColor is the CSS color painted on every feature. Default "#dddddd".
	FillColor string
	GeoBase
}

// ChoroplethDatum is one {id, value} row bound onto a feature by id.
type ChoroplethDatum struct {
	// ID matches a Feature.ID to bind this row's Value onto that feature.
	ID string
	// Value is the datum value driving the feature's color via the scale.
	Value float64
}

// ChoroplethProps mirrors @nivo/geo ChoroplethSvgProps (supported subset).
// Features are matched to Data by feature id (nivo's match:'id' default).
type ChoroplethProps struct {
	// Features are the GeoJSON features to project and draw.
	Features []Feature
	// Data are the {id, value} rows matched onto Features by feature id.
	Data []ChoroplethDatum

	// Colors is the quantize color scheme id (nivo "PuBuGn" ↔ this repo's
	// "purple_blue_green"); Steps is the number of quantization bands.
	Colors string
	Steps  int
	// Domain is the value domain; when zero it falls back to the data min/max.
	Domain [2]float64
	// UnknownColor is the CSS color used for features without a matching datum.
	// Default "#999999".
	UnknownColor string
	// ValueFormat is a d3-format spec for datum values; empty → %g.
	ValueFormat string

	// Legends configures the continuous color legends to draw.
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

// GraticuleEnabled resolves EnableGraticule (nil → false, nivo default).
func (b GeoBase) GraticuleEnabled() bool { return b.EnableGraticule != nil && *b.EnableGraticule }
