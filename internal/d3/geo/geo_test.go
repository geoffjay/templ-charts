package geo

import (
	"math"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/internal/golden"
)

const tol = 1e-3

func approx(t *testing.T, label string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("%s = %.6f, want %.6f", label, got, want)
	}
}

// projectPoint tests: values verified against d3-geo v3 output.
func TestMercator_Project(t *testing.T) {
	p := GeoMercator().Scale(100).Translate(480, 250)
	x, y := p.Project(0, 0)
	approx(t, "mercator(0,0).x", x, 480)
	approx(t, "mercator(0,0).y", y, 250)

	x, y = p.Project(90, 0)
	approx(t, "mercator(90,0).x", x, 637.0796)
	approx(t, "mercator(90,0).y", y, 250)

	_, y = p.Project(0, 45)
	approx(t, "mercator(0,45).y", y, 161.8626)
}

func TestEquirectangular_Project(t *testing.T) {
	p := GeoEquirectangular().Scale(100).Translate(480, 250)
	x, y := p.Project(0, 0)
	approx(t, "equirect(0,0).x", x, 480)
	approx(t, "equirect(0,0).y", y, 250)

	x, _ = p.Project(90, 0)
	approx(t, "equirect(90,0).x", x, 637.0796)

	_, y = p.Project(0, 90)
	approx(t, "equirect(0,90).y", y, 92.9204)
}

func TestRotate_ShiftsLongitude(t *testing.T) {
	p := GeoEquirectangular().Scale(100).Translate(480, 250).Rotate(90, 0, 0)
	x, y := p.Project(0, 0)
	approx(t, "rotated(0,0).x", x, 637.0796)
	approx(t, "rotated(0,0).y", y, 250)
}

func TestInvert_RoundTrip(t *testing.T) {
	p := GeoMercator().Scale(120).Translate(400, 300)
	for _, ll := range [][2]float64{{0, 0}, {45, 30}, {-120, -20}, {170, 60}} {
		x, y := p.Project(ll[0], ll[1])
		lon, lat := p.Invert(x, y)
		approx(t, "invert.lon", lon, ll[0])
		approx(t, "invert.lat", lat, ll[1])
	}
}

func TestProjectionByType_FallbackMercator(t *testing.T) {
	p := ProjectionByType("nope").Scale(100).Translate(480, 250)
	x, y := p.Project(90, 0)
	approx(t, "fallback.x", x, 637.0796)
	approx(t, "fallback.y", y, 250)
}

// square is a small closed polygon ring (lon/lat) used for path tests.
func square() Geometry {
	return Geometry{
		Type: TypePolygon,
		Coordinates: [][][2]float64{{
			{-10, -10}, {10, -10}, {10, 10}, {-10, 10}, {-10, -10},
		}},
	}
}

func TestPath_PolygonClosed(t *testing.T) {
	p := GeoEquirectangular().Scale(100).Translate(480, 250).Precision(0)
	path := NewPath(p).Geometry(square())
	if !strings.HasPrefix(path, "M") {
		t.Errorf("path should start with M, got %q", path)
	}
	if !strings.HasSuffix(path, "Z") {
		t.Errorf("polygon path should end with Z, got %q", path)
	}
	if strings.Contains(path, "NaN") {
		t.Errorf("path contains NaN: %q", path)
	}
}

func TestPath_PointRadiusCircle(t *testing.T) {
	p := GeoEquirectangular().Scale(100).Translate(480, 250)
	path := NewPath(p).PointRadius(5).Geometry(Geometry{Type: TypePoint, Coordinates: [2]float64{0, 0}})
	if !strings.Contains(path, "a5,5") {
		t.Errorf("point path should contain a circle of radius 5, got %q", path)
	}
}

func TestGraticule_Structure(t *testing.T) {
	g := Graticule()
	if g.Type != TypeMultiLineString {
		t.Fatalf("graticule type = %q, want MultiLineString", g.Type)
	}
	lines, ok := g.Coordinates.([][][2]float64)
	if !ok {
		t.Fatalf("graticule coordinates wrong type: %T", g.Coordinates)
	}
	// 4 major meridians + 1 major parallel (equator) + minor meridians/parallels.
	if len(lines) < 40 {
		t.Errorf("graticule line count = %d, want >= 40", len(lines))
	}
}

func TestAntimeridian_SplitsCrossingLine(t *testing.T) {
	// A line from +170° to -170° crosses the antimeridian; the clipper should
	// split it into two subpaths (two M commands).
	p := GeoEquirectangular().Scale(100).Translate(480, 250).Precision(0)
	line := Geometry{Type: TypeLineString, Coordinates: [][2]float64{{170, 0}, {-170, 0}}}
	path := NewPath(p).Geometry(line)
	if got := strings.Count(path, "M"); got < 2 {
		t.Errorf("antimeridian-crossing line should split into >=2 subpaths, got %d in %q", got, path)
	}
}

func TestPath_Deterministic(t *testing.T) {
	p := GeoNaturalEarth1().Scale(120).Translate(400, 300)
	a := NewPath(p).Geometry(square())
	b := NewPath(p).Geometry(square())
	if a != b {
		t.Errorf("path render not deterministic")
	}
}

// Golden path strings for a fixed feature under several projections, so
// unintended geometry drift is caught.
func TestPath_Golden(t *testing.T) {
	feature := Geometry{
		Type: TypePolygon,
		Coordinates: [][][2]float64{{
			{-20, 0}, {0, 30}, {20, 0}, {0, -30}, {-20, 0},
		}},
	}
	for _, proj := range []string{"mercator", "equirectangular", "naturalEarth1", "orthographic"} {
		p := ProjectionByType(proj).Scale(150).Translate(480, 250)
		golden.Assert(t, "path-"+proj, NewPath(p).Geometry(feature))
	}
	golden.Assert(t, "graticule-equirectangular",
		NewPath(GeoEquirectangular().Scale(100).Translate(480, 250)).Geometry(Graticule()))
}
