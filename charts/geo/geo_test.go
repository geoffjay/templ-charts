package geo_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/geo"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/internal/golden"
)

// sampleFeatures returns three rectangular "countries" as GeoJSON polygons.
func sampleFeatures() []geo.Feature {
	rect := func(id string, lon0, lat0, lon1, lat1 float64) geo.Feature {
		return geo.Feature{
			Type: "Feature",
			ID:   id,
			Geometry: geo.Geometry{
				Type: geo.TypePolygon,
				Coordinates: [][][2]float64{{
					{lon0, lat0}, {lon1, lat0}, {lon1, lat1}, {lon0, lat1}, {lon0, lat0},
				}},
			},
		}
	}
	return []geo.Feature{
		rect("AAA", -40, 0, -10, 30),
		rect("BBB", 0, -20, 30, 10),
		rect("CCC", 40, 20, 70, 50),
	}
}

func renderGeoMap(t *testing.T, props geo.GeoMapProps) string {
	t.Helper()
	var b strings.Builder
	if err := geo.GeoMap(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("GeoMap.Render: %v", err)
	}
	return b.String()
}

func renderChoropleth(t *testing.T, props geo.ChoroplethProps) string {
	t.Helper()
	var b strings.Builder
	if err := geo.Choropleth(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Choropleth.Render: %v", err)
	}
	return b.String()
}

func baseMap() geo.GeoMapProps {
	return geo.GeoMapProps{
		Features: sampleFeatures(),
		GeoBase: geo.GeoBase{
			Width: 600, Height: 400,
			Margin:          core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20},
			ProjectionScale: 120,
		},
	}
}

func baseChoropleth() geo.ChoroplethProps {
	return geo.ChoroplethProps{
		Features: sampleFeatures(),
		Data: []geo.ChoroplethDatum{
			{ID: "AAA", Value: 10},
			{ID: "BBB", Value: 50},
			// CCC intentionally missing → unknownColor.
		},
		GeoBase: geo.GeoBase{
			Width: 600, Height: 400,
			Margin:          core.Margin{Top: 20, Right: 20, Bottom: 20, Left: 20},
			ProjectionScale: 120,
		},
	}
}

func TestGeoMap_RendersSVG(t *testing.T) {
	out := renderGeoMap(t, baseMap())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("not a well-formed svg")
	}
	if got := strings.Count(out, "<path"); got != 3 {
		t.Errorf("expected 3 feature paths, got %d", got)
	}
	if strings.Contains(out, "NaN") {
		t.Errorf("svg contains NaN")
	}
	if !strings.Contains(out, `fill="#dddddd"`) {
		t.Errorf("expected default fill #dddddd")
	}
}

func TestGeoMap_Graticule(t *testing.T) {
	p := baseMap()
	p.EnableGraticule = true
	out := renderGeoMap(t, p)
	if !strings.Contains(out, `fill="none"`) || !strings.Contains(out, `stroke="#999999"`) {
		t.Errorf("expected graticule path with default line color")
	}
	// graticule + 3 features
	if got := strings.Count(out, "<path"); got != 4 {
		t.Errorf("expected 4 paths (graticule + 3 features), got %d", got)
	}
}

func TestGeoMap_Projections(t *testing.T) {
	for _, proj := range []string{"mercator", "equirectangular", "naturalEarth1", "equalEarth", "orthographic"} {
		p := baseMap()
		p.ProjectionType = proj
		out := renderGeoMap(t, p)
		if strings.Contains(out, "NaN") {
			t.Errorf("projection %s produced NaN", proj)
		}
		if strings.Count(out, "<path") != 3 {
			t.Errorf("projection %s: expected 3 paths", proj)
		}
	}
}

func TestChoropleth_ColorsAndUnknown(t *testing.T) {
	out := renderChoropleth(t, baseChoropleth())
	if strings.Count(out, "<path") != 3 {
		t.Errorf("expected 3 feature paths")
	}
	// The unmatched feature (CCC) uses the unknown color.
	if !strings.Contains(out, `fill="#999999"`) {
		t.Errorf("expected unknownColor #999999 for unmatched feature")
	}
}

func TestChoropleth_Legend(t *testing.T) {
	p := baseChoropleth()
	p.Legends = []legends.LegendProps{{
		Anchor:     legends.LegendAnchorBottomLeft,
		Direction:  legends.LegendDirectionColumn,
		TranslateY: -20,
		ItemWidth:  200,
		ItemHeight: 16,
	}}
	out := renderChoropleth(t, p)
	// The continuous legend emits a gradient of <rect> samples.
	if strings.Count(out, "<rect") < 2 {
		t.Errorf("expected continuous-color legend rects, got %d", strings.Count(out, "<rect"))
	}
}

func TestGeoMap_Interactive(t *testing.T) {
	p := baseMap()
	p.Interactive = true
	if !strings.Contains(renderGeoMap(t, p), "data-tc-tooltip") {
		t.Errorf("interactive GeoMap should emit data-tc-tooltip")
	}
	if strings.Contains(renderGeoMap(t, baseMap()), "data-tc-tooltip") {
		t.Errorf("non-interactive GeoMap must not emit data-tc-tooltip")
	}
}

func TestGeoMap_A11yTitleDesc(t *testing.T) {
	p := baseMap()
	p.Title = "World"
	p.Desc = "A small map."
	out := renderGeoMap(t, p)
	if !strings.Contains(out, "<title>World</title>") || !strings.Contains(out, "<desc>A small map.</desc>") {
		t.Errorf("expected title/desc threaded to SvgWrapper")
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected default role=img")
	}
}

func TestGeo_Deterministic(t *testing.T) {
	if renderGeoMap(t, baseMap()) != renderGeoMap(t, baseMap()) {
		t.Errorf("GeoMap render not deterministic")
	}
	if renderChoropleth(t, baseChoropleth()) != renderChoropleth(t, baseChoropleth()) {
		t.Errorf("Choropleth render not deterministic")
	}
}

func TestGeo_Animate(t *testing.T) {
	const fade = `<animate attributeName="opacity" from="0" to="1"`

	m := baseMap()
	m.Animate = true
	m.MotionStagger = 0.01
	if !strings.Contains(renderGeoMap(t, m), fade) {
		t.Errorf("animated GeoMap should emit an opacity fade-in <animate>")
	}
	if strings.Contains(renderGeoMap(t, baseMap()), "<animate") {
		t.Errorf("non-animated GeoMap must not emit <animate>")
	}

	c := baseChoropleth()
	c.Animate = true
	c.MotionStagger = 0.01
	if !strings.Contains(renderChoropleth(t, c), fade) {
		t.Errorf("animated Choropleth should emit an opacity fade-in <animate>")
	}
	if strings.Contains(renderChoropleth(t, baseChoropleth()), "<animate") {
		t.Errorf("non-animated Choropleth must not emit <animate>")
	}
}

func TestGeoMap_Golden(t *testing.T) {
	p := baseMap()
	p.EnableGraticule = true
	golden.Assert(t, "geomap-basic", renderGeoMap(t, p))
}

func TestChoropleth_Golden(t *testing.T) {
	p := baseChoropleth()
	p.Legends = []legends.LegendProps{{
		Anchor:     legends.LegendAnchorBottomLeft,
		Direction:  legends.LegendDirectionColumn,
		ItemWidth:  200,
		ItemHeight: 16,
	}}
	golden.Assert(t, "choropleth-basic", renderChoropleth(t, p))
}
