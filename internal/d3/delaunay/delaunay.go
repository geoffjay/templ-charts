// Package delaunay is a port of the subset of d3-delaunay (and its underlying
// Delaunator) that templ-charts needs: a Delaunay triangulation of a 2-D point
// set, its Voronoi dual clipped to a rectangle, nearest-site location, and SVG
// path rendering. It powers the voronoi chart and the accurate voronoi-mesh
// hover retrofitted into line/scatterplot/bump/tree (docs/PLAN-v3.md §3.2, §4.3).
//
// The triangulation uses the classic Bowyer–Watson incremental algorithm with a
// bounding super-triangle: O(n²) worst case, which is ample for chart-scale
// point counts, and fully deterministic (no randomness, stable point ordering).
// The Voronoi cell for each site is computed as the intersection of the
// bounding rectangle with the half-planes defined by the perpendicular
// bisectors between the site and each of its Delaunay neighbours — a robust
// construction that clips unbounded (hull) cells to the rectangle automatically.
//
// Coordinates in rendered path strings are rounded to 3 decimal places, matching
// the d3-shape Path serializer convention used elsewhere in internal/d3.
package delaunay

import (
	"math"
	"sort"
	"strconv"
)

// Delaunay holds a point set and its Delaunay triangulation.
type Delaunay struct {
	// Points is the flat [x0,y0,x1,y1,...] coordinate array, mirroring
	// d3-delaunay's Delaunay.points.
	Points []float64
	// Triangles is the flat triangle index array: every three consecutive
	// entries are the point indices of one triangle, wound counter-clockwise.
	Triangles []int

	n         int     // number of input points
	neighbors [][]int // adjacency: neighbors[i] = sorted site indices sharing an edge with i
	hull      []int   // convex-hull site indices, counter-clockwise
}

// NewDelaunayFrom builds a Delaunay triangulation from a slice of [x,y] points.
// Duplicate points are tolerated (they simply never appear in a triangle).
func NewDelaunayFrom(points [][2]float64) *Delaunay {
	flat := make([]float64, 0, len(points)*2)
	for _, p := range points {
		flat = append(flat, p[0], p[1])
	}
	return NewDelaunay(flat)
}

// NewDelaunay builds a triangulation from a flat [x0,y0,x1,y1,...] array,
// matching d3-delaunay's primary constructor shape.
func NewDelaunay(points []float64) *Delaunay {
	d := &Delaunay{Points: points, n: len(points) / 2}
	d.triangulate()
	d.computeNeighbors()
	d.computeHull()
	return d
}

// point returns the i-th site coordinates.
func (d *Delaunay) point(i int) (float64, float64) { return d.Points[2*i], d.Points[2*i+1] }

// triangle is a working triangle referencing point indices a,b,c (indices into
// an augmented point list that appends the three super-triangle vertices).
type triangle struct{ a, b, c int }

// triangulate runs Bowyer–Watson over the input points and fills d.Triangles
// with the point indices of the final (super-triangle-free) triangles.
func (d *Delaunay) triangulate() {
	d.Triangles = nil
	if d.n < 3 {
		return
	}

	// Augmented coordinate accessor: indices [0,n) are input points; indices
	// [n,n+3) are the three super-triangle vertices stored in `super`.
	var super [3][2]float64
	coord := func(i int) (float64, float64) {
		if i < d.n {
			return d.Points[2*i], d.Points[2*i+1]
		}
		v := super[i-d.n]
		return v[0], v[1]
	}

	// Super-triangle: large enough to enclose every input point.
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for i := 0; i < d.n; i++ {
		x, y := d.point(i)
		minX, minY = math.Min(minX, x), math.Min(minY, y)
		maxX, maxY = math.Max(maxX, x), math.Max(maxY, y)
	}
	dx, dy := maxX-minX, maxY-minY
	dmax := math.Max(dx, dy)
	if dmax == 0 {
		dmax = 1
	}
	midX, midY := (minX+maxX)/2, (minY+maxY)/2
	// A generous margin keeps the super-triangle well clear of the data so its
	// removal leaves the true triangulation of the input points.
	const m = 20
	super[0] = [2]float64{midX - m*dmax, midY - dmax}
	super[1] = [2]float64{midX, midY + m*dmax}
	super[2] = [2]float64{midX + m*dmax, midY - dmax}

	tris := []triangle{{d.n, d.n + 1, d.n + 2}}

	// Insert input points one at a time in input order (deterministic).
	for i := 0; i < d.n; i++ {
		px, py := d.point(i)

		// Find triangles whose circumcircle contains the point.
		bad := make([]bool, len(tris))
		for ti := range tris {
			t := tris[ti]
			ax, ay := coord(t.a)
			bx, by := coord(t.b)
			cx, cy := coord(t.c)
			if inCircumcircle(ax, ay, bx, by, cx, cy, px, py) {
				bad[ti] = true
			}
		}

		// The boundary of the cavity is the set of edges belonging to exactly
		// one bad triangle.
		type edge struct{ u, v int }
		count := map[edge]int{}
		norm := func(u, v int) edge {
			if u < v {
				return edge{u, v}
			}
			return edge{v, u}
		}
		// Preserve directed edges to re-wind new triangles consistently.
		var dirEdges []edge
		for ti := range tris {
			if !bad[ti] {
				continue
			}
			t := tris[ti]
			for _, e := range [3]edge{{t.a, t.b}, {t.b, t.c}, {t.c, t.a}} {
				count[norm(e.u, e.v)]++
				dirEdges = append(dirEdges, e)
			}
		}

		// Remove bad triangles.
		kept := tris[:0]
		for ti := range tris {
			if !bad[ti] {
				kept = append(kept, tris[ti])
			}
		}
		tris = kept

		// Re-triangulate the cavity: connect the new point to each boundary
		// (directed) edge, keeping the winding of the removed triangles.
		for _, e := range dirEdges {
			if count[norm(e.u, e.v)] == 1 {
				tris = append(tris, triangle{e.u, e.v, i})
			}
		}
	}

	// Drop triangles that touch a super-triangle vertex, then emit CCW-wound
	// point-index triples.
	for _, t := range tris {
		if t.a >= d.n || t.b >= d.n || t.c >= d.n {
			continue
		}
		ax, ay := d.point(t.a)
		bx, by := d.point(t.b)
		cx, cy := d.point(t.c)
		// orient2d > 0 ⇒ counter-clockwise.
		if orient2d(ax, ay, bx, by, cx, cy) < 0 {
			t.b, t.c = t.c, t.b
		}
		d.Triangles = append(d.Triangles, t.a, t.b, t.c)
	}
}

// computeNeighbors derives the per-site Delaunay adjacency from the triangle
// list. Two sites are neighbours iff they share a triangle edge.
func (d *Delaunay) computeNeighbors() {
	sets := make([]map[int]struct{}, d.n)
	for i := range sets {
		sets[i] = map[int]struct{}{}
	}
	add := func(a, b int) {
		sets[a][b] = struct{}{}
		sets[b][a] = struct{}{}
	}
	for t := 0; t < len(d.Triangles); t += 3 {
		a, b, c := d.Triangles[t], d.Triangles[t+1], d.Triangles[t+2]
		add(a, b)
		add(b, c)
		add(c, a)
	}
	d.neighbors = make([][]int, d.n)
	for i := range sets {
		ns := make([]int, 0, len(sets[i]))
		for j := range sets[i] {
			ns = append(ns, j)
		}
		sort.Ints(ns)
		d.neighbors[i] = ns
	}
}

// Neighbors returns the Delaunay-neighbour site indices of site i (sorted).
// These are the endpoints of the Delaunay edges incident to i — the voronoi
// "links".
func (d *Delaunay) Neighbors(i int) []int {
	if i < 0 || i >= len(d.neighbors) {
		return nil
	}
	return d.neighbors[i]
}

// Find returns the index of the site nearest to (x,y). It performs a linear
// scan, which is exact and deterministic; site counts in charts are small
// enough that the edge-walk optimisation is unnecessary. Returns -1 when there
// are no points.
func (d *Delaunay) Find(x, y float64) int {
	best, bestD := -1, math.Inf(1)
	for i := 0; i < d.n; i++ {
		px, py := d.point(i)
		dx, dy := px-x, py-y
		dd := dx*dx + dy*dy
		if dd < bestD {
			bestD, best = dd, i
		}
	}
	return best
}

// Hull returns the convex-hull site indices in counter-clockwise order.
func (d *Delaunay) Hull() []int { return d.hull }

// computeHull extracts the convex hull as the boundary edges of the
// triangulation (edges belonging to exactly one triangle), chained into a loop.
func (d *Delaunay) computeHull() {
	if d.n < 3 {
		// Degenerate: hull is just the points in order.
		d.hull = nil
		for i := 0; i < d.n; i++ {
			d.hull = append(d.hull, i)
		}
		return
	}
	type edge struct{ u, v int }
	// Directed boundary edges: an edge (u->v) is on the hull if its reverse
	// (v->u) is not also present.
	dir := map[edge]bool{}
	for t := 0; t < len(d.Triangles); t += 3 {
		a, b, c := d.Triangles[t], d.Triangles[t+1], d.Triangles[t+2]
		for _, e := range [3]edge{{a, b}, {b, c}, {c, a}} {
			dir[e] = true
		}
	}
	next := map[int]int{}
	for e := range dir {
		if !dir[edge{e.v, e.u}] {
			next[e.u] = e.v
		}
	}
	if len(next) == 0 {
		return
	}
	// Chain the boundary edges into a loop starting from the smallest index.
	start := d.n
	for u := range next {
		if u < start {
			start = u
		}
	}
	d.hull = []int{start}
	for cur := next[start]; cur != start && len(d.hull) <= d.n; cur = next[cur] {
		d.hull = append(d.hull, cur)
	}
}

// --- rendering ------------------------------------------------------------

// RenderPoints returns an SVG path string drawing every site as a circle of the
// given radius, matching d3-delaunay's Delaunay.renderPoints output shape.
func (d *Delaunay) RenderPoints(radius float64) string {
	if radius <= 0 {
		radius = 2
	}
	var b pathBuf
	for i := 0; i < d.n; i++ {
		x, y := d.point(i)
		b.s("M")
		b.n(x + radius)
		b.c()
		b.n(y)
		b.s("A")
		b.n(radius)
		b.c()
		b.n(radius)
		b.s(",0,1,1,")
		b.n(x - radius)
		b.c()
		b.n(y)
		b.s("A")
		b.n(radius)
		b.c()
		b.n(radius)
		b.s(",0,1,1,")
		b.n(x + radius)
		b.c()
		b.n(y)
	}
	return b.String()
}

// Render returns an SVG path string drawing every Delaunay edge once (the
// triangulation mesh / voronoi links).
func (d *Delaunay) Render() string {
	var b pathBuf
	for i := 0; i < d.n; i++ {
		xi, yi := d.point(i)
		for _, j := range d.neighbors[i] {
			if j <= i {
				continue
			}
			xj, yj := d.point(j)
			b.s("M")
			b.n(xi)
			b.c()
			b.n(yi)
			b.s("L")
			b.n(xj)
			b.c()
			b.n(yj)
		}
	}
	return b.String()
}

// pathBuf is a tiny SVG path-string builder with 3-decimal coordinate rounding,
// matching the d3-shape Path serializer convention.
type pathBuf struct{ b []byte }

func (p *pathBuf) s(str string)   { p.b = append(p.b, str...) }
func (p *pathBuf) c()             { p.b = append(p.b, ',') }
func (p *pathBuf) n(v float64)    { p.b = append(p.b, formatCoord(v)...) }
func (p *pathBuf) String() string { return string(p.b) }

// formatCoord rounds to 3 decimals and formats with the shortest round-trip
// representation (matching internal/d3/shape's floatToShortest).
func formatCoord(v float64) string {
	const k = 1000
	v = math.Round(v*k) / k
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// --- geometry primitives --------------------------------------------------

// orient2d returns >0 if (a,b,c) is counter-clockwise, <0 if clockwise, 0 if
// collinear.
func orient2d(ax, ay, bx, by, cx, cy float64) float64 {
	return (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
}

// inCircumcircle reports whether point (px,py) lies strictly inside the
// circumcircle of triangle (a,b,c). Uses the standard 3×3 in-circle
// determinant; the triangle is normalised to CCW so the sign test is stable.
func inCircumcircle(ax, ay, bx, by, cx, cy, px, py float64) bool {
	// Ensure CCW winding so a positive determinant means "inside".
	if orient2d(ax, ay, bx, by, cx, cy) < 0 {
		bx, by, cx, cy = cx, cy, bx, by
	}
	adx, ady := ax-px, ay-py
	bdx, bdy := bx-px, by-py
	cdx, cdy := cx-px, cy-py
	ad := adx*adx + ady*ady
	bd := bdx*bdx + bdy*bdy
	cd := cdx*cdx + cdy*cdy
	det := adx*(bdy*cd-bd*cdy) -
		ady*(bdx*cd-bd*cdx) +
		ad*(bdx*cdy-bdy*cdx)
	return det > 0
}

// circumcenter returns the circumcenter of triangle (a,b,c). Used by the
// Voronoi dual for interior vertices.
func circumcenter(ax, ay, bx, by, cx, cy float64) (float64, float64) {
	dax, day := bx-ax, by-ay
	dbx, dby := cx-ax, cy-ay
	dd := 2 * (dax*dby - day*dbx)
	if dd == 0 {
		return (ax + bx + cx) / 3, (ay + by + cy) / 3
	}
	al := dax*dax + day*day
	bl := dbx*dbx + dby*dby
	ux := (dby*al - day*bl) / dd
	uy := (dax*bl - dbx*al) / dd
	return ax + ux, ay + uy
}
