package demos

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/geo"
	"github.com/geoffjay/templ-charts/charts/legends"
)

// GeoMapDemo / ChoroplethDemo are one tile each on the /geo page. Both are
// static SVG (the projection layout is deterministic).
type GeoMapDemo struct {
	ID          string
	Title       string
	Description string
	Props       geo.GeoMapProps
}

type ChoroplethDemo struct {
	ID          string
	Title       string
	Description string
	Props       geo.ChoroplethProps
}

const (
	geoW      = 700.0
	geoH      = 440.0
	geoMargin = 10.0
)

// fitScale returns the projection scale that spans the full 360° of longitude
// across the inner width, for a projection whose planar half-width at ±180° is
// `halfWidthRad` radians-scaled units (π for cylindrical, ~2.73 for
// naturalEarth1). This stands in for the fitSize helper (not yet ported).
func fitScale(innerWidth, halfWidthRad float64) float64 {
	return innerWidth / (2 * halfWidthRad)
}

// GeoMapDemos returns the GeoMap tiles rendered from the bundled world-countries
// GeoJSON: a mercator fill map, a graticule variant, a Natural Earth
// projection, and an interactive tile.
func GeoMapDemos() []GeoMapDemo {
	features, _ := worldCountries()
	margin := core.Margin{Top: geoMargin, Right: geoMargin, Bottom: geoMargin, Left: geoMargin}
	innerW := geoW - 2*geoMargin

	base := func() geo.GeoMapProps {
		return geo.GeoMapProps{
			Features: features,
			GeoBase: geo.GeoBase{
				Width: geoW, Height: geoH, Margin: margin,
				ProjectionScale:    fitScale(innerW, math.Pi),
				BorderWidth:        0.4,
				BorderColor:        "#152238",
				GraticuleLineColor: "#dbe1ea",
			},
			FillColor: "#a7c6da",
		}
	}

	basic := base()

	// Graticule on equirectangular: a clean rectangular lat/lon grid. (Mercator
	// diverges toward the poles, where its graticule would shoot off-canvas —
	// d3 avoids that with geoMercator's built-in reclip, deferred here.)
	graticule := base()
	graticule.ProjectionType = "equirectangular"
	graticule.EnableGraticule = true

	natural := base()
	natural.ProjectionType = "naturalEarth1"
	natural.ProjectionScale = fitScale(innerW, 2.73)
	natural.EnableGraticule = true

	// Orthographic globe: an azimuthal projection whose clipCircle preclip hides
	// the far hemisphere (the v5 correctness fix). Rotated to center on Africa;
	// the graticule and borders end cleanly at the visible limb instead of
	// wrapping the whole sphere.
	globe := base()
	globe.ProjectionType = "orthographic"
	globe.ProjectionScale = innerW / 2 // radius = scale for the unit-sphere raw
	globe.ProjectionRotation = [3]float64{-10, -25, 0}
	globe.EnableGraticule = true

	interactive := base()
	interactive.Interactive = true

	return []GeoMapDemo{
		{
			ID:          "geomap-basic",
			Title:       "GeoMap (mercator)",
			Description: "The bundled Natural Earth world-countries GeoJSON (176 features) projected with internal/d3/geo (a faithful d3-geo subset) and filled a single color. Antimeridian clipping keeps features that wrap ±180° (Russia, Fiji, …) correct; adaptive resampling curves the projected borders.",
			Props:       basic,
		},
		{
			ID:          "geomap-graticule",
			Title:       "Graticule (equirectangular)",
			Description: "enableGraticule overlays the 10° meridian/parallel mesh (internal/d3/geo Graticule) on an equirectangular projection, projected through the same path generator.",
			Props:       graticule,
		},
		{
			ID:          "geomap-natural",
			Title:       "Natural Earth projection",
			Description: "projectionType:'naturalEarth1' — a pseudocylindrical projection with curved meridians. The graticule and every country border are adaptively resampled so straight lat/lon segments become the projection's true curves.",
			Props:       natural,
		},
		{
			ID:          "geomap-globe",
			Title:       "Orthographic globe",
			Description: "projectionType:'orthographic' — an azimuthal projection. internal/d3/geo's clipCircle preclip hides the far hemisphere, so only the visible cap renders and the graticule ends cleanly at the limb (before v5 the whole sphere drew, overlaying far-side geometry).",
			Props:       globe,
		},
		{
			ID:          "geomap-interactive",
			Title:       "Hover for country",
			Description: "Interactive tiles emit a per-feature client-side tooltip (charts/interact) showing the feature id on hover.",
			Props:       interactive,
		},
	}
}

// ChoroplethDemos returns the Choropleth tiles: a deterministic pseudo-value per
// country (stable across runs) colored by a quantize scale, with the
// continuous-color legend.
func ChoroplethDemos() []ChoroplethDemo {
	features, _ := worldCountries()
	margin := core.Margin{Top: geoMargin, Right: geoMargin, Bottom: 44, Left: geoMargin}
	innerW := geoW - 2*geoMargin
	data := syntheticChoroplethData(features)

	base := func() geo.ChoroplethProps {
		return geo.ChoroplethProps{
			Features: features,
			Data:     data,
			GeoBase: geo.GeoBase{
				Width: geoW, Height: geoH, Margin: margin,
				ProjectionType:  "naturalEarth1",
				ProjectionScale: fitScale(innerW, 2.73),
				BorderWidth:     0.4,
				BorderColor:     "#152238",
			},
			ValueFormat: ",.0f",
		}
	}

	basic := base()
	basic.Legends = []legends.LegendProps{{
		Anchor:     legends.LegendAnchorBottomLeft,
		Direction:  legends.LegendDirectionRow,
		ItemWidth:  300,
		ItemHeight: 14,
	}}

	interactive := base()
	interactive.Interactive = true
	interactive.Legends = basic.Legends

	return []ChoroplethDemo{
		{
			ID:          "choropleth-basic",
			Title:       "Choropleth (Natural Earth)",
			Description: "A deterministic pseudo-value per country (id-hashed, stable across runs — like nivo's random demo) bound to features by id and colored with a quantize color scale (nivo's default 'PuBuGn' — this repo's purple_blue_green). A couple of ids are omitted (Antarctica, Greenland) so unknownColor shows. The continuous-color legend spans the value domain.",
			Props:       basic,
		},
		{
			ID:          "choropleth-interactive",
			Title:       "Hover for value",
			Description: "The interactive tile adds a per-feature tooltip showing the country id and its formatted value.",
			Props:       interactive,
		},
	}
}
