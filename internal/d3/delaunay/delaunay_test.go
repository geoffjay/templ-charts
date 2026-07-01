package delaunay

import (
	"math"
	"testing"
)

// sample point set used across the tests: an irregular scatter with no special
// symmetry, so the triangulation and Voronoi cells are non-trivial.
var samplePoints = [][2]float64{
	{20, 20},
	{80, 30},
	{50, 70},
	{90, 90},
	{15, 85},
	{55, 25},
	{35, 50},
	{75, 60},
}

func newSample() *Delaunay { return NewDelaunayFrom(samplePoints) }

// TestTriangulationEmptyCircumcircle asserts the defining Delaunay property:
// no site lies strictly inside the circumcircle of any triangle.
func TestTriangulationEmptyCircumcircle(t *testing.T) {
	d := newSample()
	if len(d.Triangles) == 0 {
		t.Fatal("no triangles produced")
	}
	if len(d.Triangles)%3 != 0 {
		t.Fatalf("triangle index count %d not a multiple of 3", len(d.Triangles))
	}
	for tr := 0; tr < len(d.Triangles); tr += 3 {
		a, b, c := d.Triangles[tr], d.Triangles[tr+1], d.Triangles[tr+2]
		ax, ay := d.point(a)
		bx, by := d.point(b)
		cx, cy := d.point(c)
		for i := 0; i < d.n; i++ {
			if i == a || i == b || i == c {
				continue
			}
			px, py := d.point(i)
			if inCircumcircle(ax, ay, bx, by, cx, cy, px, py) {
				t.Errorf("point %d inside circumcircle of triangle (%d,%d,%d)", i, a, b, c)
			}
		}
	}
}

// TestTriangleWindingCCW asserts every emitted triangle is counter-clockwise.
func TestTriangleWindingCCW(t *testing.T) {
	d := newSample()
	for tr := 0; tr < len(d.Triangles); tr += 3 {
		a, b, c := d.Triangles[tr], d.Triangles[tr+1], d.Triangles[tr+2]
		ax, ay := d.point(a)
		bx, by := d.point(b)
		cx, cy := d.point(c)
		if orient2d(ax, ay, bx, by, cx, cy) <= 0 {
			t.Errorf("triangle (%d,%d,%d) is not counter-clockwise", a, b, c)
		}
	}
}

// TestEulerTriangleCount checks the triangle count matches Euler's relation for
// a triangulation of points in general position: T = 2n - 2 - h, where h is the
// number of hull vertices.
func TestEulerTriangleCount(t *testing.T) {
	d := newSample()
	nTris := len(d.Triangles) / 3
	h := len(d.Hull())
	want := 2*d.n - 2 - h
	if nTris != want {
		t.Errorf("triangle count = %d, want 2n-2-h = %d (n=%d, h=%d)", nTris, want, d.n, h)
	}
}

// TestFind checks nearest-site queries, including exact hits on the sites.
func TestFind(t *testing.T) {
	d := newSample()
	for i, p := range samplePoints {
		if got := d.Find(p[0], p[1]); got != i {
			t.Errorf("Find at site %d returned %d", i, got)
		}
	}
	// A query point clearly closest to site 0 (20,20).
	if got := d.Find(18, 22); got != 0 {
		t.Errorf("Find(18,22) = %d, want 0", got)
	}
	// Empty set.
	if got := NewDelaunayFrom(nil).Find(0, 0); got != -1 {
		t.Errorf("Find on empty = %d, want -1", got)
	}
}

// TestFindMatchesLinearNearest fuzzes Find against a brute-force nearest search
// over a grid, confirming they always agree (Find defines cell membership).
func TestFindMatchesLinearNearest(t *testing.T) {
	d := newSample()
	for gx := 0; gx <= 100; gx += 7 {
		for gy := 0; gy <= 100; gy += 7 {
			x, y := float64(gx), float64(gy)
			got := d.Find(x, y)
			want, bestD := -1, math.Inf(1)
			for i := 0; i < d.n; i++ {
				px, py := d.point(i)
				dd := (px-x)*(px-x) + (py-y)*(py-y)
				if dd < bestD {
					bestD, want = dd, i
				}
			}
			if got != want {
				t.Fatalf("Find(%g,%g) = %d, brute-force = %d", x, y, got, want)
			}
		}
	}
}

// TestNeighborsSymmetric asserts the adjacency is symmetric and self-free.
func TestNeighborsSymmetric(t *testing.T) {
	d := newSample()
	for i := 0; i < d.n; i++ {
		for _, j := range d.Neighbors(i) {
			if j == i {
				t.Errorf("site %d lists itself as a neighbor", i)
			}
			found := false
			for _, k := range d.Neighbors(j) {
				if k == i {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("neighbor relation not symmetric: %d->%d but not back", i, j)
			}
		}
	}
}

// TestDeterministic asserts repeated builds produce identical output.
func TestDeterministic(t *testing.T) {
	a := newSample()
	b := newSample()
	if len(a.Triangles) != len(b.Triangles) {
		t.Fatalf("triangle counts differ: %d vs %d", len(a.Triangles), len(b.Triangles))
	}
	for i := range a.Triangles {
		if a.Triangles[i] != b.Triangles[i] {
			t.Fatalf("triangulation not deterministic at %d", i)
		}
	}
	if a.Render() != b.Render() || a.RenderPoints(4) != b.RenderPoints(4) {
		t.Error("render output not deterministic")
	}
}

func TestRenderPointsShape(t *testing.T) {
	d := NewDelaunayFrom([][2]float64{{10, 10}})
	got := d.RenderPoints(4)
	want := "M14,10A4,4,0,1,1,6,10A4,4,0,1,1,14,10"
	if got != want {
		t.Errorf("RenderPoints = %q, want %q", got, want)
	}
}
