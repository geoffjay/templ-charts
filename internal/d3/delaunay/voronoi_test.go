package delaunay

import (
	"math"
	"testing"
)

var sampleBounds = [4]float64{0, 0, 100, 100}

// TestVoronoiCellsContainSites asserts each site lies inside its own cell.
func TestVoronoiCellsContainSites(t *testing.T) {
	d := newSample()
	v := d.Voronoi(sampleBounds)
	for i, p := range samplePoints {
		poly := v.CellPolygon(i)
		if len(poly) < 3 {
			t.Errorf("site %d has degenerate cell (%d verts)", i, len(poly))
			continue
		}
		if !pointInPolygon(p[0], p[1], poly) {
			t.Errorf("site %d not inside its own cell", i)
		}
	}
}

// TestVoronoiCellMembership asserts a cell's polygon agrees with Find: a point
// strictly inside cell i must have i as its nearest site.
func TestVoronoiCellMembership(t *testing.T) {
	d := newSample()
	v := d.Voronoi(sampleBounds)
	for i := range samplePoints {
		poly := v.CellPolygon(i)
		if len(poly) < 3 {
			continue
		}
		// Sample the polygon centroid — guaranteed interior for convex cells.
		cx, cy := centroid(poly)
		if got := d.Find(cx, cy); got != i {
			t.Errorf("centroid of cell %d nearest to site %d", i, got)
		}
	}
}

// TestVoronoiCellsWithinBounds asserts every cell vertex lies within the clip
// rectangle (allowing a tiny epsilon).
func TestVoronoiCellsWithinBounds(t *testing.T) {
	d := newSample()
	v := d.Voronoi(sampleBounds)
	const eps = 1e-6
	for i := range samplePoints {
		for _, pt := range v.CellPolygon(i) {
			if pt[0] < -eps || pt[0] > 100+eps || pt[1] < -eps || pt[1] > 100+eps {
				t.Errorf("cell %d vertex %v outside bounds", i, pt)
			}
		}
	}
}

// TestVoronoiPartitionCoversArea asserts the cells tile the clip rectangle: a
// dense grid of probe points each falls inside exactly one cell.
func TestVoronoiPartitionCoversArea(t *testing.T) {
	d := newSample()
	v := d.Voronoi(sampleBounds)
	for gx := 2; gx < 100; gx += 5 {
		for gy := 2; gy < 100; gy += 5 {
			x, y := float64(gx), float64(gy)
			inside := 0
			for i := range samplePoints {
				if poly := v.CellPolygon(i); len(poly) >= 3 && pointInPolygon(x, y, poly) {
					inside++
				}
			}
			if inside != 1 {
				t.Errorf("probe (%g,%g) inside %d cells, want exactly 1", x, y, inside)
			}
		}
	}
}

// TestVoronoiRenderNonEmpty checks the render helpers produce path strings.
func TestVoronoiRenderNonEmpty(t *testing.T) {
	d := newSample()
	v := d.Voronoi(sampleBounds)
	if v.Render() == "" {
		t.Error("Render() empty")
	}
	if v.RenderCell(0) == "" {
		t.Error("RenderCell(0) empty")
	}
	if got := v.RenderBounds(); got != "M0,0L100,0L100,100L0,100Z" {
		t.Errorf("RenderBounds = %q", got)
	}
}

// TestVoronoiDeterministic asserts stable rendering.
func TestVoronoiDeterministic(t *testing.T) {
	a := newSample().Voronoi(sampleBounds)
	b := newSample().Voronoi(sampleBounds)
	if a.Render() != b.Render() {
		t.Error("voronoi render not deterministic")
	}
}

// --- helpers ---------------------------------------------------------------

func pointInPolygon(x, y float64, poly [][2]float64) bool {
	in := false
	n := len(poly)
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		xi, yi := poly[i][0], poly[i][1]
		xj, yj := poly[j][0], poly[j][1]
		if (yi > y) != (yj > y) {
			xint := (xj-xi)*(y-yi)/(yj-yi) + xi
			if x < xint {
				in = !in
			}
		}
	}
	return in
}

func centroid(poly [][2]float64) (float64, float64) {
	var sx, sy float64
	for _, p := range poly {
		sx += p[0]
		sy += p[1]
	}
	n := float64(len(poly))
	return sx / n, sy / n
}

func TestCircumcenter(t *testing.T) {
	// Right triangle with hypotenuse endpoints (0,0)-(2,2); circumcenter is the
	// midpoint of the hypotenuse (1,1).
	cx, cy := circumcenter(0, 0, 2, 0, 0, 2)
	if math.Abs(cx-1) > 1e-9 || math.Abs(cy-1) > 1e-9 {
		t.Errorf("circumcenter = (%g,%g), want (1,1)", cx, cy)
	}
}
