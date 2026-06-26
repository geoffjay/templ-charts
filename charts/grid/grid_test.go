package grid

import (
	"math"
	"testing"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestComputeCellDimensions_NonSquare(t *testing.T) {
	w, h := ComputeCellDimensions(ComputeCellDimensionsArgs{
		Width: 300, Height: 200, Rows: 2, Columns: 3, Padding: 0, Square: false,
	})
	// (300 - 2*0)/3 = 100 ; (200 - 1*0)/2 = 100
	if !approx(w, 100) || !approx(h, 100) {
		t.Fatalf("got %v,%v want 100,100", w, h)
	}
}

func TestComputeCellDimensions_Square(t *testing.T) {
	w, h := ComputeCellDimensions(ComputeCellDimensionsArgs{
		Width: 300, Height: 150, Rows: 2, Columns: 3, Padding: 0, Square: true,
	})
	// cellWidth=100, cellHeight=75 → min=75 for both
	if !approx(w, 75) || !approx(h, 75) {
		t.Fatalf("got %v,%v want 75,75", w, h)
	}
}

func TestComputeCellDimensions_Padding(t *testing.T) {
	w, h := ComputeCellDimensions(ComputeCellDimensionsArgs{
		Width: 320, Height: 210, Rows: 2, Columns: 3, Padding: 10, Square: false,
	})
	// (320 - 2*10)/3 = 100 ; (210 - 1*10)/2 = 100
	if !approx(w, 100) || !approx(h, 100) {
		t.Fatalf("got %v,%v want 100,100", w, h)
	}
}

func TestGenerateGrid_BottomFill(t *testing.T) {
	g := GenerateGrid(GenerateGridArgs{
		Width: 300, Height: 200, Columns: 3, Rows: 3, Square: false,
	})
	if len(g.Cells) != 9 {
		t.Fatalf("got %d cells want 9", len(g.Cells))
	}
	// bottom fill: row-major, index 0 = (row0,col0), 8 = (row2,col2)
	if g.Cells[0].Index != 0 || g.Cells[0].Row != 0 || g.Cells[0].Column != 0 {
		t.Fatalf("cell0 = %+v", g.Cells[0])
	}
	last := g.Cells[8]
	if last.Index != 8 || last.Row != 2 || last.Column != 2 {
		t.Fatalf("cell8 = %+v", last)
	}
}

func TestGenerateGrid_TopFillReversesIndex(t *testing.T) {
	g := GenerateGrid(GenerateGridArgs{
		Width: 300, Height: 200, Columns: 3, Rows: 3, Square: false,
		FillDirection: GridFillTop,
	})
	// top fill: index 0 should be the last row-major cell (row2,col2).
	if g.Cells[0].Row != 2 || g.Cells[0].Column != 2 {
		t.Fatalf("top-fill cell0 = %+v want row2,col2", g.Cells[0])
	}
	if g.Cells[0].Index != 0 {
		t.Fatalf("cell0 index = %d want 0", g.Cells[0].Index)
	}
}

func TestGenerateGrid_KeyFormat(t *testing.T) {
	g := GenerateGrid(GenerateGridArgs{
		Width: 100, Height: 100, Columns: 2, Rows: 2, Square: false,
	})
	keys := map[string]bool{}
	for _, c := range g.Cells {
		keys[c.Key] = true
	}
	for _, want := range []string{"0.0", "0.1", "1.0", "1.1"} {
		if !keys[want] {
			t.Fatalf("missing key %q in %v", want, keys)
		}
	}
}

func TestGenerateGrid_OriginCentered(t *testing.T) {
	// 2x2 in 100x100, square=false → cellWidth=50, cellHeight=50, origin (0,0).
	g := GenerateGrid(GenerateGridArgs{
		Width: 100, Height: 100, Columns: 2, Rows: 2, Square: false,
	})
	if !approx(g.X, 0) || !approx(g.Y, 0) {
		t.Fatalf("origin = %v,%v want 0,0", g.X, g.Y)
	}
}

func TestGenerateGridWith_Extender(t *testing.T) {
	type cell struct {
		GridCell
		Value int
	}
	g := GenerateGridWith[cell](GenerateGridWithArgs[cell]{
		Width: 100, Height: 100, Columns: 2, Rows: 2, Square: false,
		Extend: func(c GridCell, origin Vertex) cell {
			return cell{GridCell: c, Value: c.Index * 10}
		},
	})
	if len(g.Cells) != 4 {
		t.Fatalf("got %d cells want 4", len(g.Cells))
	}
	if g.Cells[3].Value != 30 {
		t.Fatalf("cell3 value = %d want 30", g.Cells[3].Value)
	}
}

func TestAreBoundingBoxTouching(t *testing.T) {
	a := BoundingBox{Top: 0, Right: 10, Bottom: 10, Left: 0}
	b := BoundingBox{Top: 10, Right: 10, Bottom: 20, Left: 0} // touches on bottom edge
	if !AreBoundingBoxTouching(a, b) {
		t.Fatal("touching boxes reported as not touching")
	}
	c := BoundingBox{Top: 11, Right: 10, Bottom: 21, Left: 0} // gap
	if AreBoundingBoxTouching(a, c) {
		t.Fatal("non-touching boxes reported as touching")
	}
}

func TestGetCellsPolygons_SingleRow(t *testing.T) {
	cells := []GridCell{
		{X: 0, Y: 0, Width: 10, Height: 10},
		{X: 10, Y: 0, Width: 10, Height: 10},
	}
	polys := GetCellsPolygons(cells)
	if len(polys) != 1 {
		t.Fatalf("got %d polygons want 1", len(polys))
	}
	// 2 top-right + 2 bottom-left = 4 vertices.
	if len(polys[0]) != 4 {
		t.Fatalf("got %d vertices want 4", len(polys[0]))
	}
}
