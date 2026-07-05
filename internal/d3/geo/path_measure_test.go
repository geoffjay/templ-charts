package geo

import (
	"math"
	"testing"
)

// Path measurement tests (bounds/area/centroid). Expected values were produced
// by d3-geo@3's path.bounds/path.area/path.centroid.

// diamond is a clockwise (small-interior) lon/lat polygon used for measurement
// and the conic path goldens.
func diamond() Geometry {
	return Geometry{Type: TypePolygon, Coordinates: [][][2]float64{{
		{-20, 0}, {0, 30}, {20, 0}, {0, -30}, {-20, 0},
	}}}
}

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
