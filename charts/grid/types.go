// Package grid provides the grid-cell layout primitives shared by heatmap,
// waffle, and similar cell-based chart types. It is scaffold-only in v1 —
// not exercised by bar/line/pie — but the layout math (computeCellDimensions,
// generateGrid, bounding-box overlap, perpendicular polygon merge) is ported
// verbatim from @nivo/grid so future chart types plug in without rework.
//
// Mirrors @nivo/grid (types.ts, grid.ts, boundingBoxes.ts, polygon.ts).
package grid

// Vertex is an [x, y] coordinate pair.
type Vertex [2]float64

// BoundingBox is an axis-aligned rectangle described by its four edges.
type BoundingBox struct {
	Top, Right, Bottom, Left float64
}

// GridFillDirection affects the order cells are indexed when generating a
// grid. Mirrors @nivo/grid GridFillDirection.
//
//	│   top             │   right           │   bottom          │   left            │
//	│                   │                   │   ↓               │                   │
//	│   8 ─── 7 ─── 6   │ → 0  ╭─ 3  ╭─ 6   │   0 ─── 1 ─── 2   │   8  ╭─ 5  ╭─ 2   │
//	│   ╭───────────╯   │   │  │  │  │  │   │   ╭───────────╯   │   │  │  │  │  │   │
//	│   5 ─── 4 ─── 3   │   1  │  4  │  7   │   3 ─── 4 ─── 5   │   5  │  4  │  1   │
//	│   ╭───────────╯   │   │  │  │  │  │   │   ╭───────────╯   │   │  │  │  │  │   │
//	│   2 ─── 1 ─── 0   │   2 ─╯  5 ─╯  8   │   6 ─── 7 ─── 8   │   6 ─╯  3 ─╯  0 ← │
//	│               ↑   │                   │                   │                   │
type GridFillDirection string

const (
	GridFillTop    GridFillDirection = "top"
	GridFillRight  GridFillDirection = "right"
	GridFillBottom GridFillDirection = "bottom"
	GridFillLeft   GridFillDirection = "left"
)

// GridCell is one cell of a generated grid. Key is "row.column".
type GridCell struct {
	Key    string
	Index  int
	Column int
	Row    int
	X, Y   float64
	Width  float64
	Height float64
}

// CellExtender extends a base GridCell with chart-specific fields. origin is
// the grid's top-left offset (so cells can be positioned in absolute svg
// coordinates). Used by heatmap/waffle to attach values/colors.
type CellExtender[C any] func(cell GridCell, origin Vertex) C

// GeneratedGrid is the output of GenerateGrid: the grid's origin, the
// resolved cell dimensions, and the (possibly extended) cells in fill order.
type GeneratedGrid[C any] struct {
	X, Y       float64
	CellWidth  float64
	CellHeight float64
	Cells      []C
}
