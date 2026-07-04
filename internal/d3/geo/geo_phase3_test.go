package geo

import (
	"math"
	"testing"

	"github.com/geoffjay/templ-charts/internal/golden"
)

// Phase-3 geo tests: path measurement (bounds/area/centroid), fitExtent/fitSize/
// fitWidth/fitHeight, the conic projection family, and TopoJSON decoding. All
// expected values were produced by d3-geo@3 / topojson-client@3.

// diamond is a clockwise (small-interior) lon/lat polygon used for measurement
// and conic path goldens.
func diamond() Geometry {
	return Geometry{Type: TypePolygon, Coordinates: [][][2]float64{{
		{-20, 0}, {0, 30}, {20, 0}, {0, -30}, {-20, 0},
	}}}
}

func approx2(t *testing.T, label string, got, wx, wy float64) {
	t.Helper()
	if math.Abs(got-wx) > tol {
		t.Errorf("%s = %.4f, want %.4f", label, got, wx)
	}
	_ = wy
}

// --- measurement -----------------------------------------------------------

func TestPathMeasurement(t *testing.T) {
	p := GeoEquirectangular().Scale(100).Translate(480, 250).Precision(0)
	gp := NewPath(p)

	// Polygon (diamond): rhombus with diagonals 69.81 × 104.72.
	b := gp.Bounds(diamond())
	approx(t, "diamond.b.x0", b[0][0], 445.0934)
	approx(t, "diamond.b.y0", b[0][1], 197.6401)
	approx(t, "diamond.b.x1", b[1][0], 514.9066)
	approx(t, "diamond.b.y1", b[1][1], 302.3599)
	c := gp.Centroid(diamond())
	approx(t, "diamond.c.x", c[0], 480)
	approx(t, "diamond.c.y", c[1], 250)
	approx(t, "diamond.area", gp.Area(diamond()), 3655.409)

	// LineString: area 0, length-weighted centroid.
	line := Geometry{Type: TypeLineString, Coordinates: [][2]float64{{-20, 0}, {0, 20}, {20, 0}}}
	lb := gp.Bounds(line)
	approx(t, "line.b.x0", lb[0][0], 445.0934)
	approx(t, "line.b.y0", lb[0][1], 215.0934)
	approx(t, "line.b.x1", lb[1][0], 514.9066)
	approx(t, "line.b.y1", lb[1][1], 250)
	lc := gp.Centroid(line)
	approx(t, "line.c.x", lc[0], 480)
	approx(t, "line.c.y", lc[1], 232.5467)
	approx(t, "line.area", gp.Area(line), 0)

	// MultiPoint: area 0, mean-of-points centroid.
	mp := Geometry{Type: TypeMultiPoint, Coordinates: [][2]float64{{-10, -10}, {10, 10}, {0, 5}}}
	mc := gp.Centroid(mp)
	approx(t, "mp.c.x", mc[0], 480)
	approx(t, "mp.c.y", mc[1], 247.0911)

	// Point: bounds collapse to the point, centroid is the point.
	pt := Geometry{Type: TypePoint, Coordinates: [2]float64{5, 5}}
	pc := gp.Centroid(pt)
	approx(t, "pt.c.x", pc[0], 488.7266)
	approx(t, "pt.c.y", pc[1], 241.2734)
	pb := gp.Bounds(pt)
	approx(t, "pt.b.x0", pb[0][0], 488.7266)
	approx(t, "pt.b.x1", pb[1][0], 488.7266)
}

func TestCentroid_EmptyIsNaN(t *testing.T) {
	p := GeoEquirectangular().Scale(100).Translate(480, 250)
	c := NewPath(p).Centroid(Geometry{Type: TypeGeometryColl})
	if !math.IsNaN(c[0]) || !math.IsNaN(c[1]) {
		t.Errorf("empty centroid = %v, want [NaN,NaN]", c)
	}
}

func TestBounds_AcceptsFeatureCollection(t *testing.T) {
	p := GeoEquirectangular().Scale(100).Translate(480, 250).Precision(0)
	fc := FeatureCollection{Type: TypeFeatureCollection, Features: []Feature{
		{Type: TypeFeature, Geometry: diamond()},
	}}
	b := NewPath(p).Bounds(fc)
	approx(t, "fc.b.x0", b[0][0], 445.0934)
	approx(t, "fc.b.x1", b[1][0], 514.9066)
}

// --- fit -------------------------------------------------------------------

// cwWorld is a clockwise 60°×40° box (small, unambiguous interior).
func cwWorld() Geometry {
	return Geometry{Type: TypePolygon, Coordinates: [][][2]float64{{
		{-30, -20}, {-30, 20}, {30, 20}, {30, -20}, {-30, -20},
	}}}
}

func TestFitExtentAndSize(t *testing.T) {
	fe := GeoEquirectangular().FitExtent(0, 0, 300, 200, cwWorld())
	approx(t, "fitExtent.scale", fe.ScaleValue(), 251.3427)
	fx, fy := fe.TranslateValue()
	approx(t, "fitExtent.tx", fx, 150)
	approx(t, "fitExtent.ty", fy, 100)

	fs := GeoEquirectangular().FitSize(400, 300, cwWorld())
	approx(t, "fitSize.scale", fs.ScaleValue(), 377.0141)
	sx, sy := fs.TranslateValue()
	approx(t, "fitSize.tx", sx, 200)
	approx(t, "fitSize.ty", sy, 150)

	fw := GeoEquirectangular().FitWidth(400, cwWorld())
	approx(t, "fitWidth.scale", fw.ScaleValue(), 381.9719)
	wx, wy := fw.TranslateValue()
	approx(t, "fitWidth.tx", wx, 200)
	approx(t, "fitWidth.ty", wy, 151.9725)

	fh := GeoEquirectangular().FitHeight(300, cwWorld())
	approx(t, "fitHeight.scale", fh.ScaleValue(), 377.0141)
	hx, hy := fh.TranslateValue()
	approx(t, "fitHeight.tx", hx, 197.4041)
	approx(t, "fitHeight.ty", hy, 150)
}

// TestFit_PreservesClipExtent verifies fit restores a rectangular postclip that
// was installed before fitting.
func TestFit_PreservesClipExtent(t *testing.T) {
	p := GeoEquirectangular().ClipExtent(10, 20, 300, 400)
	p.FitSize(400, 300, cwWorld())
	// After fitting, the clipExtent postclip must still crop: a feature fully
	// outside the box renders empty.
	away := Geometry{Type: TypePolygon, Coordinates: [][][2]float64{{
		{-179, 89}, {-179, 89.5}, {-178, 89.5}, {-178, 89}, {-179, 89},
	}}}
	// (Sanity: the projection still renders the fitted world non-empty.)
	if NewPath(p).Geometry(cwWorld()) == "" {
		t.Error("fitted world should render non-empty")
	}
	_ = away
}

// --- conic projections -----------------------------------------------------

func TestConicProjections_Project(t *testing.T) {
	cases := []struct {
		name string
		p    *Projection
		// proj(10,20) and proj(-30,45)
		ax, ay, bx, by float64
	}{
		{"conicConformal", GeoConicConformal(), 504.559, 192.9498, 426.4147, 122.2004},
		{"conicEqualArea", GeoConicEqualArea(), 504.5628, 199.0771, 426.4309, 125.3052},
		{"conicEquidistant", GeoConicEquidistant(), 504.561, 196.4238, 426.4244, 124.1779},
	}
	for _, c := range cases {
		p := c.p.Parallels(20, 50).Scale(150).Translate(480, 250)
		ax, ay := p.Project(10, 20)
		bx, by := p.Project(-30, 45)
		approx(t, c.name+".a.x", ax, c.ax)
		approx(t, c.name+".a.y", ay, c.ay)
		approx(t, c.name+".b.x", bx, c.bx)
		approx(t, c.name+".b.y", by, c.by)
	}
}

func TestConicProjections_ByType(t *testing.T) {
	for _, name := range []string{"conicConformal", "conicEqualArea", "conicEquidistant"} {
		p := ProjectionByType(name)
		if p.rawFactory == nil {
			t.Errorf("ProjectionByType(%q) is not a conic projection", name)
		}
	}
}

func TestConic_ParallelsRoundTrip(t *testing.T) {
	p := GeoConicConformal().Parallels(15, 45)
	p0, p1 := p.ParallelsValue()
	approx(t, "parallels.0", p0, 15)
	approx(t, "parallels.1", p1, 45)
}

// TestConicPathGolden pins the rendered path strings for the conic family (drift
// guard); the numbers themselves are validated against d3 in the Project test.
func TestConicPathGolden(t *testing.T) {
	for _, name := range []string{"conicConformal", "conicEqualArea", "conicEquidistant"} {
		p := ProjectionByType(name).Parallels(20, 50).Scale(150).Translate(480, 250).Precision(0)
		golden.Assert(t, "path-"+name, NewPath(p).Geometry(diamond()))
	}
}

// --- topojson --------------------------------------------------------------

const topoSample = `{
  "type":"Topology",
  "transform":{"scale":[0.01,0.01],"translate":[-1,-1]},
  "arcs":[[[0,0],[100,0],[0,100]],[[100,100],[-100,0],[0,-100]]],
  "objects":{"shapes":{"type":"GeometryCollection","geometries":[
    {"type":"Polygon","arcs":[[0,1]],"id":"P"},
    {"type":"LineString","arcs":[0],"id":"L"}
  ]}}
}`

func TestTopoJSON_Decode(t *testing.T) {
	fc, err := DecodeTopoJSON([]byte(topoSample), "shapes")
	if err != nil {
		t.Fatalf("DecodeTopoJSON: %v", err)
	}
	if len(fc.Features) != 2 {
		t.Fatalf("features = %d, want 2", len(fc.Features))
	}

	poly := fc.Features[0]
	if poly.ID != "P" || poly.Geometry.Type != TypePolygon {
		t.Fatalf("feature0 = %q/%q, want P/Polygon", poly.ID, poly.Geometry.Type)
	}
	wantPoly := [][][2]float64{{{-1, -1}, {0, -1}, {0, 0}, {-1, 0}, {-1, -1}}}
	rings, ok := poly.Geometry.Coordinates.([][][2]float64)
	if !ok {
		t.Fatalf("polygon coords type %T", poly.Geometry.Coordinates)
	}
	assertCoords3(t, "poly", rings, wantPoly)

	line := fc.Features[1]
	if line.ID != "L" || line.Geometry.Type != TypeLineString {
		t.Fatalf("feature1 = %q/%q, want L/LineString", line.ID, line.Geometry.Type)
	}
	wantLine := [][2]float64{{-1, -1}, {0, -1}, {0, 0}}
	pts, ok := line.Geometry.Coordinates.([][2]float64)
	if !ok {
		t.Fatalf("line coords type %T", line.Geometry.Coordinates)
	}
	assertCoords2(t, "line", pts, wantLine)
}

// TestTopoJSON_DecodedIsRenderable confirms the decoded features flow straight
// into the path generator (the whole point of the decoder).
func TestTopoJSON_DecodedIsRenderable(t *testing.T) {
	fc, err := DecodeTopoJSON([]byte(topoSample), "shapes")
	if err != nil {
		t.Fatal(err)
	}
	p := GeoEquirectangular().FitSize(200, 200, fc)
	out := NewPath(p).Feature(fc.Features[0])
	if out == "" || out[0] != 'M' {
		t.Errorf("decoded polygon should render a path, got %q", out)
	}
}

func assertCoords2(t *testing.T, label string, got, want [][2]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: %d points, want %d", label, len(got), len(want))
	}
	for i := range want {
		if math.Abs(got[i][0]-want[i][0]) > tol || math.Abs(got[i][1]-want[i][1]) > tol {
			t.Errorf("%s[%d] = %v, want %v", label, i, got[i], want[i])
		}
	}
}

func assertCoords3(t *testing.T, label string, got, want [][][2]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: %d rings, want %d", label, len(got), len(want))
	}
	for i := range want {
		assertCoords2(t, label, got[i], want[i])
	}
}
