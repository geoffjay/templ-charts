package interact

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/geoffjay/templ-charts/internal/d3/delaunay"
)

// MeshRadiusAttrName is the optional attribute carrying the detection radius (in
// inner-chart user units) for a voronoi mesh: the client ignores points farther
// than this from the cursor. Absent/≤0 means unlimited (nivo's default).
const MeshRadiusAttrName = "data-tc-mesh-radius"

// MeshPoint is one nearest-point hover target in inner-chart coordinates, with
// the tooltip HTML to show when it is the closest point to the cursor.
type MeshPoint struct {
	X    float64
	Y    float64
	HTML string
}

// meshPointJSON is the wire shape the client script reads from data-tc-mesh.
type meshPointJSON struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	HTML string  `json:"html"`
}

// MeshData serializes points into the data-tc-mesh JSON array the client
// interactivity layer consumes for nearest-point (voronoi-cell) hover.
func MeshData(points []MeshPoint) string {
	pts := make([]meshPointJSON, len(points))
	for i, p := range points {
		pts[i] = meshPointJSON{X: p.X, Y: p.Y, HTML: p.HTML}
	}
	b, err := json.Marshal(pts)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// MeshCellsPath returns the SVG path string of the Voronoi cells for the given
// mesh points, clipped to the inner chart rectangle [0,0,width,height]. This is
// the exact detection partition the client hit-tests against (nearest site =
// enclosing cell), so drawing it visualises the mesh precisely — used for the
// debug-mesh overlay. Returns "" when there are too few points to triangulate.
func MeshCellsPath(points []MeshPoint, width, height float64) string {
	if len(points) == 0 {
		return ""
	}
	pts := make([][2]float64, len(points))
	for i, p := range points {
		pts[i] = [2]float64{p.X, p.Y}
	}
	d := delaunay.NewDelaunayFrom(pts)
	v := d.Voronoi([4]float64{0, 0, width, height})
	return v.Render()
}

// MeshOverlay builds the SVG for a voronoi-mesh hover layer over the inner chart
// area: an optional faint rendering of the actual Voronoi cells (when debug),
// then a transparent capture <rect> carrying the nearest-point data and — when
// detectionRadius > 0 — the detection-radius attribute. The client script
// resolves the nearest point on mousemove, honours the detection radius, shows
// its tooltip and draws a crosshair, entirely client-side.
//
// When len(points) == 0 the overlay is empty (nothing to detect).
func MeshOverlay(points []MeshPoint, width, height float64, debug bool, detectionRadius float64) string {
	if len(points) == 0 {
		return ""
	}
	var b strings.Builder
	if debug {
		if cells := MeshCellsPath(points, width, height); cells != "" {
			b.WriteString(`<path d="`)
			b.WriteString(cells)
			b.WriteString(`" fill="none" stroke="red" stroke-width="1" stroke-opacity="0.35"></path>`)
		}
	}
	b.WriteString(`<rect x="0" y="0" width="`)
	b.WriteString(fmtDim(width))
	b.WriteString(`" height="`)
	b.WriteString(fmtDim(height))
	b.WriteString(`" fill="black" fill-opacity="0" data-tc-mesh="`)
	b.WriteString(attrEscape(MeshData(points)))
	b.WriteString(`"`)
	if detectionRadius > 0 {
		b.WriteString(` `)
		b.WriteString(MeshRadiusAttrName)
		b.WriteString(`="`)
		b.WriteString(fmtDim(detectionRadius))
		b.WriteString(`"`)
	}
	b.WriteString(`></rect>`)
	return b.String()
}

// fmtDim formats a dimension rounded to 3 decimals.
func fmtDim(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
