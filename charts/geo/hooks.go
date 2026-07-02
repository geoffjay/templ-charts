package geo

import (
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3geo "github.com/geoffjay/templ-charts/internal/d3/geo"
)

// buildProjection mirrors @nivo/geo useGeoMap: projectionById(type) with scale,
// translate (inner width/height × translation fraction), and rotation. width /
// height are the inner dimensions.
func buildProjection(b GeoBase, width, height float64) *d3geo.Projection {
	return d3geo.ProjectionByType(b.ProjectionType).
		Scale(b.ProjectionScale).
		Translate(width*b.ProjectionTranslation[0], height*b.ProjectionTranslation[1]).
		Rotate(b.ProjectionRotation[0], b.ProjectionRotation[1], b.ProjectionRotation[2])
}

// UseGeoMap computes the projected feature paths and graticule for a GeoMap.
// width/height are the inner dimensions.
func UseGeoMap(props GeoMapProps, width, height float64) GeoResult {
	proj := buildProjection(props.GeoBase, width, height)
	path := d3geo.NewPath(proj)

	features := make([]ComputedFeature, 0, len(props.Features))
	for _, f := range props.Features {
		features = append(features, ComputedFeature{
			ID:          f.ID,
			Path:        path.Feature(f),
			FillColor:   props.FillColor,
			BorderWidth: props.BorderWidth,
			BorderColor: props.BorderColor,
		})
	}

	return GeoResult{
		Features:      features,
		GraticulePath: graticulePath(props.EnableGraticule, path),
	}
}

// UseChoropleth mirrors @nivo/geo useChoropleth: it binds Data onto Features by
// id, builds a quantize color scale over the value domain, and colors each
// matched feature (unmatched → unknownColor). width/height are inner dims.
func UseChoropleth(props ChoroplethProps, width, height float64) GeoResult {
	proj := buildProjection(props.GeoBase, width, height)
	path := d3geo.NewPath(proj)

	byID := make(map[string]float64, len(props.Data))
	min, max := math.Inf(1), math.Inf(-1)
	for _, d := range props.Data {
		byID[d.ID] = d.Value
		if d.Value < min {
			min = d.Value
		}
		if d.Value > max {
			max = d.Value
		}
	}
	if math.IsInf(min, 0) {
		min, max = 0, 0
	}

	scale := colors.GetQuantizeColorScale(
		colors.QuantizeColorScaleConfig{
			Type:   "quantize",
			Domain: props.Domain,
			Scheme: props.Colors,
			Steps:  props.Steps,
		},
		colors.SequentialColorScaleValues{Min: min, Max: max},
	)
	format := valueFormatter(props.ValueFormat)

	features := make([]ComputedFeature, 0, len(props.Features))
	for _, f := range props.Features {
		cf := ComputedFeature{
			ID:          f.ID,
			Path:        path.Feature(f),
			BorderWidth: props.BorderWidth,
			BorderColor: props.BorderColor,
			Label:       f.ID,
		}
		if v, ok := byID[f.ID]; ok {
			cf.HasValue = true
			cf.Value = v
			cf.FormattedValue = format(v)
			cf.FillColor = scale(v)
		} else {
			cf.FillColor = props.UnknownColor
		}
		features = append(features, cf)
	}

	// Legend domain: explicit Domain if set, else data min/max.
	lmin, lmax := props.Domain[0], props.Domain[1]
	if lmin == 0 && lmax == 0 {
		lmin, lmax = min, max
	}

	return GeoResult{
		Features:      features,
		GraticulePath: graticulePath(props.EnableGraticule, path),
		ColorScale:    scale,
		ValueMin:      lmin,
		ValueMax:      lmax,
	}
}

func graticulePath(enabled bool, path *d3geo.Path) string {
	if !enabled {
		return ""
	}
	return path.Geometry(d3geo.Graticule())
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
