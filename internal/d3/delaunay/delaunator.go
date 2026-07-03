// delaunator.go — a port of Mapbox's Delaunator (the incremental sweep-hull
// algorithm d3-delaunay is built on): an O(n log n) Delaunay triangulation that
// replaces the O(n²) Bowyer–Watson core this package used to carry. It produces
// the triangle index array and its parallel half-edge array; the public
// Delaunay wrapper re-winds the triangles to this package's counter-clockwise
// convention and derives neighbours/hull/voronoi from them exactly as before,
// so the Delaunay edge set — and therefore every rendered voronoi golden — is
// unchanged (the Delaunay triangulation of points in general position is
// unique; only the triangle *ordering* differs from Bowyer–Watson, and nothing
// in the public surface depends on that order).
//
// The predicates here (orientD, inCircleD, circumradiusD, pseudoAngle) use
// Delaunator's own sign conventions and are internally consistent with the
// sweep logic; they are deliberately kept separate from the package's
// orient2d/inCircumcircle helpers, which use the opposite winding convention.
package delaunay

import (
	"math"
	"sort"
)

// delEpsilon matches delaunator-js's EPSILON (2⁻⁵²): the threshold below which
// two consecutive (distance-sorted) points are treated as coincident and the
// later one skipped, so duplicate sites never enter a triangle.
const delEpsilon = 2.220446049250313e-16

// delaunator holds the working state of one sweep-hull triangulation. triangles
// and halfedges are trimmed to trianglesLen on completion; hull is the
// counter-clockwise hull site list (Delaunator's convention).
type delaunator struct {
	coords       []float64
	n            int
	triangles    []int
	halfedges    []int
	trianglesLen int

	hashSize  int
	hullPrev  []int // hullPrev[i] = site before i on the hull
	hullNext  []int // hullNext[i] = site after i on the hull
	hullTri   []int // hullTri[i]  = half-edge of the triangle at hull site i
	hullHash  []int // angular hash bucket -> hull site
	hullStart int
	hull      []int

	ids   []int
	dists []float64
	cx    float64
	cy    float64

	edgeStack [512]int
}

// newDelaunator triangulates the flat [x0,y0,x1,y1,...] coordinate array.
func newDelaunator(coords []float64) *delaunator {
	n := len(coords) / 2
	maxTriangles := 2*n - 5
	if maxTriangles < 0 {
		maxTriangles = 0
	}
	hashSize := int(math.Ceil(math.Sqrt(float64(n))))
	d := &delaunator{
		coords:    coords,
		n:         n,
		triangles: make([]int, maxTriangles*3),
		halfedges: make([]int, maxTriangles*3),
		hashSize:  hashSize,
		hullPrev:  make([]int, n),
		hullNext:  make([]int, n),
		hullTri:   make([]int, n),
		hullHash:  make([]int, hashSize),
		ids:       make([]int, n),
		dists:     make([]float64, n),
	}
	d.update()
	return d
}

func (d *delaunator) update() {
	coords := d.coords
	n := d.n
	if n < 3 {
		d.triangles = d.triangles[:0]
		d.halfedges = d.halfedges[:0]
		return
	}

	// Bounding box + id list.
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for i := 0; i < n; i++ {
		x, y := coords[2*i], coords[2*i+1]
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		d.ids[i] = i
	}
	cx := (minX + maxX) / 2
	cy := (minY + maxY) / 2

	// Seed i0: the point closest to the bounding-box centre.
	i0 := 0
	minDist := math.Inf(1)
	for i := 0; i < n; i++ {
		dd := distSq(cx, cy, coords[2*i], coords[2*i+1])
		if dd < minDist {
			i0, minDist = i, dd
		}
	}
	i0x, i0y := coords[2*i0], coords[2*i0+1]

	// i1: the point closest to i0.
	i1 := 0
	minDist = math.Inf(1)
	for i := 0; i < n; i++ {
		if i == i0 {
			continue
		}
		dd := distSq(i0x, i0y, coords[2*i], coords[2*i+1])
		if dd < minDist && dd > 0 {
			i1, minDist = i, dd
		}
	}
	i1x, i1y := coords[2*i1], coords[2*i1+1]

	// i2: the point giving the smallest circumcircle with i0,i1.
	i2 := 0
	minRadius := math.Inf(1)
	for i := 0; i < n; i++ {
		if i == i0 || i == i1 {
			continue
		}
		r := circumradiusD(i0x, i0y, i1x, i1y, coords[2*i], coords[2*i+1])
		if r < minRadius {
			i2, minRadius = i, r
		}
	}
	i2x, i2y := coords[2*i2], coords[2*i2+1]

	if math.IsInf(minRadius, 1) {
		// All points are collinear: order them along the dominant axis and
		// expose the ordering as the hull. No triangles exist.
		for i := 0; i < n; i++ {
			dv := coords[2*i] - coords[0]
			if dv == 0 {
				dv = coords[2*i+1] - coords[1]
			}
			d.dists[i] = dv
		}
		d.sortIDs()
		hull := make([]int, 0, n)
		d0 := math.Inf(-1)
		for _, id := range d.ids {
			if d.dists[id] > d0 {
				hull = append(hull, id)
				d0 = d.dists[id]
			}
		}
		d.hull = hull
		d.triangles = d.triangles[:0]
		d.halfedges = d.halfedges[:0]
		return
	}

	// Orient the seed triangle to Delaunator's canonical winding.
	if orientD(i0x, i0y, i1x, i1y, i2x, i2y) < 0 {
		i1, i2 = i2, i1
		i1x, i2x = i2x, i1x
		i1y, i2y = i2y, i1y
	}

	// Sort the remaining points by distance from the seed circumcentre.
	d.cx, d.cy = circumcenter(i0x, i0y, i1x, i1y, i2x, i2y)
	for i := 0; i < n; i++ {
		d.dists[i] = distSq(coords[2*i], coords[2*i+1], d.cx, d.cy)
	}
	d.sortIDs()

	// Seed the hull as the initial triangle.
	d.hullStart = i0
	hullSize := 3
	d.hullNext[i0], d.hullPrev[i2] = i1, i1
	d.hullNext[i1], d.hullPrev[i0] = i2, i2
	d.hullNext[i2], d.hullPrev[i1] = i0, i0
	d.hullTri[i0], d.hullTri[i1], d.hullTri[i2] = 0, 1, 2
	for i := range d.hullHash {
		d.hullHash[i] = -1
	}
	d.hullHash[d.hashKey(i0x, i0y)] = i0
	d.hullHash[d.hashKey(i1x, i1y)] = i1
	d.hullHash[d.hashKey(i2x, i2y)] = i2

	d.trianglesLen = 0
	d.addTriangle(i0, i1, i2, -1, -1, -1)

	var xp, yp float64
	for k := 0; k < len(d.ids); k++ {
		i := d.ids[k]
		x, y := coords[2*i], coords[2*i+1]

		// Skip near-duplicate points and the seed vertices.
		if k > 0 && math.Abs(x-xp) <= delEpsilon && math.Abs(y-yp) <= delEpsilon {
			continue
		}
		xp, yp = x, y
		if i == i0 || i == i1 || i == i2 {
			continue
		}

		// Find a visible edge on the convex hull using the angular hash.
		start := 0
		key := d.hashKey(x, y)
		for j := 0; j < d.hashSize; j++ {
			start = d.hullHash[(key+j)%d.hashSize]
			if start != -1 && start != d.hullNext[start] {
				break
			}
		}
		start = d.hullPrev[start]
		e := start
		var q int
		for {
			q = d.hullNext[e]
			if orientD(x, y, coords[2*e], coords[2*e+1], coords[2*q], coords[2*q+1]) < 0 {
				break
			}
			e = q
			if e == start {
				e = -1
				break
			}
		}
		if e == -1 {
			continue // likely a near-duplicate; skip
		}

		// Add the first triangle from this point.
		t := d.addTriangle(e, i, d.hullNext[e], -1, -1, d.hullTri[e])
		d.hullTri[i] = d.legalize(t + 2)
		d.hullTri[e] = t
		hullSize++

		// Walk forward, adding triangles while the hull turns the right way.
		next := d.hullNext[e]
		for {
			q = d.hullNext[next]
			if orientD(x, y, coords[2*next], coords[2*next+1], coords[2*q], coords[2*q+1]) >= 0 {
				break
			}
			t = d.addTriangle(next, i, q, d.hullTri[i], -1, d.hullTri[next])
			d.hullTri[i] = d.legalize(t + 2)
			d.hullNext[next] = next // mark removed
			hullSize--
			next = q
		}

		// Walk backward when the point was inserted at the hull start.
		if e == start {
			for {
				q = d.hullPrev[e]
				if orientD(x, y, coords[2*q], coords[2*q+1], coords[2*e], coords[2*e+1]) >= 0 {
					break
				}
				t = d.addTriangle(q, i, e, -1, d.hullTri[e], d.hullTri[q])
				d.legalize(t + 2)
				d.hullTri[q] = t
				d.hullNext[e] = e // mark removed
				hullSize--
				e = q
			}
		}

		// Splice the new point into the hull linked list + hash.
		d.hullStart = e
		d.hullPrev[i] = e
		d.hullNext[e] = i
		d.hullPrev[next] = i
		d.hullNext[i] = next
		d.hullHash[d.hashKey(x, y)] = i
		d.hullHash[d.hashKey(coords[2*e], coords[2*e+1])] = e
	}

	// Materialise the hull site list from the linked list.
	d.hull = make([]int, hullSize)
	e := d.hullStart
	for i := 0; i < hullSize; i++ {
		d.hull[i] = e
		e = d.hullNext[e]
	}

	d.triangles = d.triangles[:d.trianglesLen]
	d.halfedges = d.halfedges[:d.trianglesLen]
}

// sortIDs orders d.ids by ascending dists, tie-breaking by id for determinism.
func (d *delaunator) sortIDs() {
	sort.Slice(d.ids, func(a, b int) bool {
		ia, ib := d.ids[a], d.ids[b]
		if d.dists[ia] != d.dists[ib] {
			return d.dists[ia] < d.dists[ib]
		}
		return ia < ib
	})
}

func (d *delaunator) hashKey(x, y float64) int {
	k := int(math.Floor(pseudoAngle(x-d.cx, y-d.cy) * float64(d.hashSize)))
	return ((k % d.hashSize) + d.hashSize) % d.hashSize
}

func (d *delaunator) addTriangle(i0, i1, i2, a, b, c int) int {
	t := d.trianglesLen
	d.triangles[t] = i0
	d.triangles[t+1] = i1
	d.triangles[t+2] = i2
	d.link(t, a)
	d.link(t+1, b)
	d.link(t+2, c)
	d.trianglesLen += 3
	return t
}

func (d *delaunator) link(a, b int) {
	d.halfedges[a] = b
	if b != -1 {
		d.halfedges[b] = a
	}
}

// legalize restores the Delaunay property around half-edge a by flipping
// illegal edges, iterating with an explicit stack (Delaunator's approach). It
// returns the half-edge opposite the apex, used to update hullTri.
func (d *delaunator) legalize(a int) int {
	triangles := d.triangles
	halfedges := d.halfedges
	coords := d.coords

	i := 0
	var ar int
	for {
		b := halfedges[a]
		a0 := a - a%3
		ar = a0 + (a+2)%3

		if b == -1 { // convex hull edge — nothing opposite to flip
			if i == 0 {
				break
			}
			i--
			a = d.edgeStack[i]
			continue
		}

		b0 := b - b%3
		al := a0 + (a+1)%3
		bl := b0 + (b+2)%3

		p0 := triangles[ar]
		pr := triangles[a]
		pl := triangles[al]
		p1 := triangles[bl]

		illegal := inCircleD(
			coords[2*p0], coords[2*p0+1],
			coords[2*pr], coords[2*pr+1],
			coords[2*pl], coords[2*pl+1],
			coords[2*p1], coords[2*p1+1],
		)

		if illegal {
			triangles[a] = p1
			triangles[b] = p0

			hbl := halfedges[bl]
			// The flipped edge may lie on the hull; retarget its hullTri.
			if hbl == -1 {
				e := d.hullStart
				for {
					if d.hullTri[e] == bl {
						d.hullTri[e] = a
						break
					}
					e = d.hullPrev[e]
					if e == d.hullStart {
						break
					}
				}
			}
			d.link(a, hbl)
			d.link(b, halfedges[ar])
			d.link(ar, bl)

			br := b0 + (b+1)%3
			if i < len(d.edgeStack) {
				d.edgeStack[i] = br
				i++
			}
		} else {
			if i == 0 {
				break
			}
			i--
			a = d.edgeStack[i]
		}
	}
	return ar
}

// --- Delaunator predicates (its own winding/sign convention) --------------

func distSq(ax, ay, bx, by float64) float64 {
	dx, dy := ax-bx, ay-by
	return dx*dx + dy*dy
}

// orientD is Delaunator's fast orientation test; it is the negation of this
// package's orient2d, so the sweep logic's >=0 / <0 comparisons carry
// Delaunator's meaning directly.
func orientD(ax, ay, bx, by, cx, cy float64) float64 {
	return (ay-cy)*(bx-cx) - (ax-cx)*(by-cy)
}

// circumradiusD returns the squared circumradius of triangle (a,b,c), or +Inf
// when the points are collinear.
func circumradiusD(ax, ay, bx, by, cx, cy float64) float64 {
	dx, dy := bx-ax, by-ay
	ex, ey := cx-ax, cy-ay
	bl := dx*dx + dy*dy
	cl := ex*ex + ey*ey
	det := dx*ey - dy*ex
	if det == 0 {
		return math.Inf(1)
	}
	d := 0.5 / det
	x := (ey*bl - dy*cl) * d
	y := (dx*cl - ex*bl) * d
	return x*x + y*y
}

// inCircleD reports whether p lies inside the circumcircle of (a,b,c) under
// Delaunator's winding convention (used to decide edge flips).
func inCircleD(ax, ay, bx, by, cx, cy, px, py float64) bool {
	dx, dy := ax-px, ay-py
	ex, ey := bx-px, by-py
	fx, fy := cx-px, cy-py
	ap := dx*dx + dy*dy
	bp := ex*ex + ey*ey
	cp := fx*fx + fy*fy
	return dx*(ey*cp-bp*fy)-dy*(ex*cp-bp*fx)+ap*(ex*fy-ey*fx) < 0
}

// pseudoAngle maps a direction to a monotonic value in [0,1) without a trig
// call — the ordering key for Delaunator's angular hull hash.
func pseudoAngle(dx, dy float64) float64 {
	p := dx / (math.Abs(dx) + math.Abs(dy))
	if dy > 0 {
		return (3 - p) / 4
	}
	return (1 + p) / 4
}
