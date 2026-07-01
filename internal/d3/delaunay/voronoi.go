// voronoi.go — the Voronoi dual of a Delaunay triangulation, clipped to a
// rectangle. Each cell is built by intersecting the clip rectangle with the
// half-planes defined by the perpendicular bisectors between a site and each of
// its Delaunay neighbours (Sutherland–Hodgman convex clipping). This construction
// is robust: it produces a correct finite convex polygon for both interior and
// unbounded (hull) sites, and needs no special-case ray extension.
package delaunay

// Voronoi is the clipped Voronoi diagram of a Delaunay triangulation.
type Voronoi struct {
	Delaunay *Delaunay
	// Xmin,Ymin,Xmax,Ymax is the clip rectangle (bounds = [xmin,ymin,xmax,ymax]).
	Xmin, Ymin, Xmax, Ymax float64

	cells [][][2]float64 // cells[i] = CCW polygon of site i (nil if empty)
}

// Voronoi returns the Voronoi diagram clipped to bounds = [xmin,ymin,xmax,ymax].
func (d *Delaunay) Voronoi(bounds [4]float64) *Voronoi {
	v := &Voronoi{
		Delaunay: d,
		Xmin:     bounds[0], Ymin: bounds[1],
		Xmax: bounds[2], Ymax: bounds[3],
	}
	if v.Xmax < v.Xmin {
		v.Xmin, v.Xmax = v.Xmax, v.Xmin
	}
	if v.Ymax < v.Ymin {
		v.Ymin, v.Ymax = v.Ymax, v.Ymin
	}
	v.compute()
	return v
}

func (v *Voronoi) compute() {
	d := v.Delaunay
	v.cells = make([][][2]float64, d.n)

	// Starting polygon: the clip rectangle, wound counter-clockwise.
	rect := [][2]float64{
		{v.Xmin, v.Ymin},
		{v.Xmax, v.Ymin},
		{v.Xmax, v.Ymax},
		{v.Xmin, v.Ymax},
	}

	for i := 0; i < d.n; i++ {
		xi, yi := d.point(i)
		poly := make([][2]float64, len(rect))
		copy(poly, rect)

		neighbors := d.neighbors[i]
		if len(neighbors) == 0 && d.n > 1 {
			// Isolated / duplicate site: fall back to bisectors against every
			// other site so its cell is still well defined.
			for j := 0; j < d.n; j++ {
				if j == i {
					continue
				}
				xj, yj := d.point(j)
				poly = clipHalfPlane(poly, xi, yi, xj, yj)
				if len(poly) == 0 {
					break
				}
			}
		} else {
			for _, j := range neighbors {
				xj, yj := d.point(j)
				poly = clipHalfPlane(poly, xi, yi, xj, yj)
				if len(poly) == 0 {
					break
				}
			}
		}
		if len(poly) >= 3 {
			v.cells[i] = poly
		}
	}
}

// clipHalfPlane clips convex polygon `poly` to the half-plane of points at least
// as close to site (sx,sy) as to neighbour (nx,ny): the perpendicular-bisector
// half-plane. Sutherland–Hodgman against the bisector line.
func clipHalfPlane(poly [][2]float64, sx, sy, nx, ny float64) [][2]float64 {
	// Bisector: points p with dot(p - mid, n) <= 0 are on the site's side,
	// where n = (neighbour - site) and mid is the segment midpoint.
	dx, dy := nx-sx, ny-sy
	mx, my := (sx+nx)/2, (sy+ny)/2
	// Signed distance of p from the bisector (negative = keep).
	side := func(px, py float64) float64 { return (px-mx)*dx + (py-my)*dy }

	if len(poly) == 0 {
		return poly
	}
	out := make([][2]float64, 0, len(poly)+1)
	for i := 0; i < len(poly); i++ {
		cur := poly[i]
		prev := poly[(i-1+len(poly))%len(poly)]
		dCur := side(cur[0], cur[1])
		dPrev := side(prev[0], prev[1])
		curIn := dCur <= 1e-9
		prevIn := dPrev <= 1e-9
		if curIn != prevIn {
			// Edge crosses the bisector: add the intersection point.
			t := dPrev / (dPrev - dCur)
			ix := prev[0] + t*(cur[0]-prev[0])
			iy := prev[1] + t*(cur[1]-prev[1])
			out = append(out, [2]float64{ix, iy})
		}
		if curIn {
			out = append(out, cur)
		}
	}
	return out
}

// CellPolygon returns the counter-clockwise vertex list of site i's clipped
// Voronoi cell, or nil if the cell is empty. The polygon is not explicitly
// closed (the first vertex is not repeated).
func (v *Voronoi) CellPolygon(i int) [][2]float64 {
	if i < 0 || i >= len(v.cells) {
		return nil
	}
	return v.cells[i]
}

// RenderCell returns the SVG path string for site i's cell (a closed polygon),
// or "" if the cell is empty.
func (v *Voronoi) RenderCell(i int) string {
	poly := v.CellPolygon(i)
	if len(poly) < 3 {
		return ""
	}
	var b pathBuf
	appendCell(&b, poly)
	return b.String()
}

// Render returns the SVG path string for the whole diagram: every cell as a
// closed sub-path.
func (v *Voronoi) Render() string {
	var b pathBuf
	for i := range v.cells {
		if len(v.cells[i]) >= 3 {
			appendCell(&b, v.cells[i])
		}
	}
	return b.String()
}

// RenderBounds returns the SVG path string of the clip rectangle.
func (v *Voronoi) RenderBounds() string {
	var b pathBuf
	b.s("M")
	b.n(v.Xmin)
	b.c()
	b.n(v.Ymin)
	b.s("L")
	b.n(v.Xmax)
	b.c()
	b.n(v.Ymin)
	b.s("L")
	b.n(v.Xmax)
	b.c()
	b.n(v.Ymax)
	b.s("L")
	b.n(v.Xmin)
	b.c()
	b.n(v.Ymax)
	b.s("Z")
	return b.String()
}

func appendCell(b *pathBuf, poly [][2]float64) {
	b.s("M")
	b.n(poly[0][0])
	b.c()
	b.n(poly[0][1])
	for k := 1; k < len(poly); k++ {
		b.s("L")
		b.n(poly[k][0])
		b.c()
		b.n(poly[k][1])
	}
	b.s("Z")
}

// Contains reports whether point (x,y) falls in site i's cell — i.e. i is the
// nearest site. A thin convenience over Find used by tests.
func (v *Voronoi) Contains(i int, x, y float64) bool {
	return v.Delaunay.Find(x, y) == i
}
