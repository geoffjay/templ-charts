package geo

import (
	"math"
	"strconv"
	"strings"
)

// Path renders GeoJSON geometry to SVG path strings through a projection.
// Mirrors d3-geo's geoPath with the string PathContext (bounds/area/centroid
// are not ported — GeoMap/Choropleth only need the path string).
type Path struct {
	projection *Projection
	radius     float64
}

// NewPath returns a path generator for the given projection with the default
// point radius (4.5, as in d3-geo).
func NewPath(projection *Projection) *Path {
	return &Path{projection: projection, radius: 4.5}
}

// PointRadius sets the radius used to render Point/MultiPoint geometries.
func (gp *Path) PointRadius(r float64) *Path { gp.radius = r; return gp }

// Geometry returns the SVG path string for a geometry (empty if it projects to
// nothing).
func (gp *Path) Geometry(g Geometry) string {
	ctx := newPathContext(gp.radius)
	StreamGeometry(g, gp.projection.Stream(ctx))
	return ctx.result()
}

// Feature returns the SVG path string for a feature's geometry.
func (gp *Path) Feature(f Feature) string { return gp.Geometry(f.Geometry) }

// pathContext accumulates SVG path commands from a projected stream. Ported
// from d3-geo src/path/string.js PathString.
type pathContext struct {
	b          strings.Builder
	radius     float64
	circle     string
	inPolygon  bool // true between polygonStart/polygonEnd (rings get a Z)
	pointState int  // -1 isolated point, 0 first vertex, 1 subsequent
}

func newPathContext(radius float64) *pathContext {
	return &pathContext{radius: radius, pointState: -1}
}

func (c *pathContext) PolygonStart() { c.inPolygon = true }
func (c *pathContext) PolygonEnd()   { c.inPolygon = false }
func (c *pathContext) Sphere()       {}

func (c *pathContext) LineStart() { c.pointState = 0 }

func (c *pathContext) LineEnd() {
	if c.inPolygon {
		c.b.WriteByte('Z')
	}
	c.pointState = -1
}

func (c *pathContext) Point(x, y float64) {
	switch c.pointState {
	case 0:
		c.b.WriteByte('M')
		c.writeXY(x, y)
		c.pointState = 1
	case 1:
		c.b.WriteByte('L')
		c.writeXY(x, y)
	default: // isolated point → circle marker
		if c.circle == "" {
			c.circle = circlePath(c.radius)
		}
		c.b.WriteByte('M')
		c.writeXY(x, y)
		c.b.WriteString(c.circle)
	}
}

func (c *pathContext) writeXY(x, y float64) {
	c.b.WriteString(geoNum(x))
	c.b.WriteByte(',')
	c.b.WriteString(geoNum(y))
}

func (c *pathContext) result() string { return c.b.String() }

// circlePath is the two-arc circle used for point markers (d3-geo string.js).
func circlePath(radius float64) string {
	r := geoNum(radius)
	return "m0," + r +
		"a" + r + "," + r + " 0 1,1 0," + geoNum(-2*radius) +
		"a" + r + "," + r + " 0 1,1 0," + geoNum(2*radius) +
		"z"
}

// geoNum formats a coordinate rounded to 3 decimals, matching the precision the
// rest of templ-charts uses for path strings (charts/chord fmtF etc.).
func geoNum(v float64) string {
	r := math.Round(v*1000) / 1000
	if r == 0 {
		r = 0 // normalize -0 to +0 for cross-platform-stable output
	}
	return strconv.FormatFloat(r, 'g', -1, 64)
}
