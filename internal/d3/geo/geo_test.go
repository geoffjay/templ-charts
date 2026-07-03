package geo

import (
	"math"
	"strconv"
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
	for _, proj := range []string{
		"mercator", "equirectangular", "naturalEarth1", "orthographic",
		"gnomonic", "stereographic", "azimuthalEqualArea", "azimuthalEquidistant",
	} {
		p := ProjectionByType(proj).Scale(150).Translate(480, 250)
		golden.Assert(t, "path-"+proj, NewPath(p).Geometry(feature))
	}
	golden.Assert(t, "graticule-equirectangular",
		NewPath(GeoEquirectangular().Scale(100).Translate(480, 250)).Geometry(Graticule()))
}

// farSideBand is a small counter-clockwise box on the far hemisphere (around
// lon 167°) of an orthographic/azimuthal projection centered at (0°,0°). It does
// not cross the antimeridian, so it unambiguously encloses a small far-side area.
func farSideBand() Geometry {
	return Geometry{
		Type: TypePolygon,
		Coordinates: [][][2]float64{{
			{160, -10}, {160, 10}, {175, 10}, {175, -10}, {160, -10},
		}},
	}
}

// nearSideBand is the same shape (same winding) on the front hemisphere.
func nearSideBand() Geometry {
	return Geometry{
		Type: TypePolygon,
		Coordinates: [][][2]float64{{
			{-10, -10}, {-10, 10}, {10, 10}, {10, -10}, {-10, -10},
		}},
	}
}

// TestClipCircle_HidesFarSide is the correctness assertion for the fix: with the
// clipCircle preclip installed, orthographic must drop geometry on the far
// hemisphere entirely (empty path) while still rendering the near side. Before
// the fix the far-side band rendered over the near side.
func TestClipCircle_HidesFarSide(t *testing.T) {
	p := GeoOrthographic().Scale(150).Translate(480, 250).Precision(0)

	if got := NewPath(p).Geometry(farSideBand()); got != "" {
		t.Errorf("far-side band should be clipped to empty, got %q", got)
	}
	near := NewPath(p).Geometry(nearSideBand())
	if near == "" || !strings.HasPrefix(near, "M") || !strings.HasSuffix(near, "Z") {
		t.Errorf("near-side band should render a closed path, got %q", near)
	}
	if strings.Contains(near, "NaN") {
		t.Errorf("near-side path contains NaN: %q", near)
	}
}

// TestClipCircle_CutsStraddlingRing checks that a ring spanning the visible limb
// is cut and rejoined along the clip circle into a valid closed path (not
// dropped, not NaN).
func TestClipCircle_CutsStraddlingRing(t *testing.T) {
	// A wide band from the front hemisphere across the limb into the far side.
	straddle := Geometry{
		Type: TypePolygon,
		Coordinates: [][][2]float64{{
			{-30, -20}, {120, -20}, {120, 20}, {-30, 20}, {-30, -20},
		}},
	}
	p := GeoOrthographic().Scale(150).Translate(480, 250).Precision(0)
	got := NewPath(p).Geometry(straddle)
	if got == "" {
		t.Fatalf("straddling ring should not be fully clipped")
	}
	if !strings.HasPrefix(got, "M") || !strings.HasSuffix(got, "Z") {
		t.Errorf("clipped ring should be a closed path, got %q", got)
	}
	if strings.Contains(got, "NaN") {
		t.Errorf("clipped ring contains NaN: %q", got)
	}
	golden.Assert(t, "orthographic-clipped-ring", got)
}

// TestClipExtent_CropsToBox checks the rectangular postclip: a large polygon is
// cropped to the viewport box, and geometry entirely outside the box is dropped.
func TestClipExtent_CropsToBox(t *testing.T) {
	// equirectangular scale 100 @ (480,250): lon/lat map linearly, so a world-ish
	// polygon spans well beyond a small box centered on the origin.
	// Wound so the interior (not the complement) is filled — see the winding
	// note in farSideBand; the reversed winding would make this a hole.
	big := Geometry{
		Type: TypePolygon,
		Coordinates: [][][2]float64{{
			{-80, -60}, {-80, 60}, {80, 60}, {80, -60}, {-80, -60},
		}},
	}
	box := [4]float64{400, 200, 560, 300} // x0,y0,x1,y1 around the center
	p := GeoEquirectangular().Scale(100).Translate(480, 250).Precision(0).
		ClipExtent(box[0], box[1], box[2], box[3])

	got := NewPath(p).Geometry(big)
	if got == "" {
		t.Fatalf("clipped polygon should not be empty")
	}
	if !strings.HasSuffix(got, "Z") {
		t.Errorf("clipped polygon should be closed, got %q", got)
	}
	if strings.Contains(got, "NaN") {
		t.Errorf("clipExtent path contains NaN: %q", got)
	}
	// Every coordinate must fall within the clip box (± rounding tolerance).
	assertWithinBox(t, got, box)

	// A feature entirely outside the box is dropped.
	away := Geometry{Type: TypePolygon, Coordinates: [][][2]float64{{
		{-179, 80}, {-179, 85}, {-170, 85}, {-170, 80}, {-179, 80},
	}}}
	if out := NewPath(p).Geometry(away); out != "" {
		t.Errorf("feature outside clip box should be dropped, got %q", out)
	}

	golden.Assert(t, "equirectangular-clipextent", got)
}

// assertWithinBox parses the numeric pairs out of an SVG path and asserts each
// lies inside the box (with a small tolerance for path rounding).
func assertWithinBox(t *testing.T, path string, box [4]float64) {
	t.Helper()
	const eps = 0.5
	repl := strings.NewReplacer("M", " ", "L", " ", "Z", " ", ",", " ")
	fields := strings.Fields(repl.Replace(path))
	for i := 0; i+1 < len(fields); i += 2 {
		x, errX := strconv.ParseFloat(fields[i], 64)
		y, errY := strconv.ParseFloat(fields[i+1], 64)
		if errX != nil || errY != nil {
			continue
		}
		if x < box[0]-eps || x > box[2]+eps || y < box[1]-eps || y > box[3]+eps {
			t.Errorf("point (%.3f,%.3f) outside clip box %v", x, y, box)
		}
	}
}

// TestClipAngle_Toggle confirms clipAngle(0) restores the antimeridian preclip
// (so the far-side band is no longer hidden — it renders through, as it did
// before the fix).
func TestClipAngle_Toggle(t *testing.T) {
	p := GeoOrthographic().Scale(150).Translate(480, 250).Precision(0).ClipAngle(0)
	if p.ClipAngleValue() != 0 {
		t.Errorf("ClipAngleValue after ClipAngle(0) = %v, want 0", p.ClipAngleValue())
	}
	if got := NewPath(p).Geometry(farSideBand()); got == "" {
		t.Errorf("with antimeridian preclip the far-side band should render, got empty")
	}
}
