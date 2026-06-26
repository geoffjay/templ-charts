package legends

import "testing"

func TestComputePositionFromAnchor(t *testing.T) {
	x, y := ComputePositionFromAnchor(LegendAnchorTopLeft, 200, 100, 50, 30, 10, 5)
	if x != 10 || y != 5 {
		t.Fatalf("top-left got (%v,%v), want (10,5)", x, y)
	}
	x, y = ComputePositionFromAnchor(LegendAnchorBottomRight, 200, 100, 50, 30, 0, 0)
	if x != 150 || y != 70 {
		t.Fatalf("bottom-right got (%v,%v), want (150,70)", x, y)
	}
	x, y = ComputePositionFromAnchor(LegendAnchorTop, 200, 100, 50, 30, 0, 0)
	if x != 75 || y != 0 {
		t.Fatalf("top got (%v,%v), want (75,0)", x, y)
	}
}

func TestComputeDimensions_Column(t *testing.T) {
	props := LegendProps{Direction: LegendDirectionColumn, Items: []Datum{{}, {}, {}}, ItemWidth: 100, ItemHeight: 20}
	w, h := ComputeDimensions(props)
	if w != 100 || h != 60 {
		t.Fatalf("column dims got (%v,%v), want (100,60)", w, h)
	}
}

func TestComputeDimensions_Row(t *testing.T) {
	props := LegendProps{Direction: LegendDirectionRow, Items: []Datum{{}, {}, {}}, ItemWidth: 100, ItemHeight: 20}
	w, h := ComputeDimensions(props)
	if w != 300 || h != 20 {
		t.Fatalf("row dims got (%v,%v), want (300,20)", w, h)
	}
}

func TestComputeItemLayout_Column(t *testing.T) {
	props := LegendProps{Direction: LegendDirectionColumn, Items: []Datum{{}, {}, {}}, ItemWidth: 100, ItemHeight: 20}
	layouts := ComputeItemLayout(props)
	if layouts[0].Y != 0 || layouts[1].Y != 20 || layouts[2].Y != 40 {
		t.Fatalf("column layouts Y = %v, want [0,20,40]", []float64{layouts[0].Y, layouts[1].Y, layouts[2].Y})
	}
}

func TestComputeItemLayout_Row(t *testing.T) {
	props := LegendProps{Direction: LegendDirectionRow, Items: []Datum{{}, {}, {}}, ItemWidth: 100, ItemHeight: 20}
	layouts := ComputeItemLayout(props)
	if layouts[0].X != 0 || layouts[1].X != 100 || layouts[2].X != 200 {
		t.Fatalf("row layouts X = %v, want [0,100,200]", []float64{layouts[0].X, layouts[1].X, layouts[2].X})
	}
}

func TestComputeContinuousColorsLegend(t *testing.T) {
	scale := func(v float64) string {
		if v < 5 {
			return "blue"
		}
		return "red"
	}
	stops := ComputeContinuousColorsLegend(scale, 0, 10, 4)
	if len(stops) != 4 {
		t.Fatalf("got %d stops, want 4", len(stops))
	}
	if stops[0].Offset != 0 || stops[3].Offset != 1 {
		t.Fatalf("offsets got %v, %v, want 0 and 1", stops[0].Offset, stops[3].Offset)
	}
	if stops[0].Color != "blue" || stops[3].Color != "red" {
		t.Fatalf("colors got %q, %q, want blue, red", stops[0].Color, stops[3].Color)
	}
}

func TestSymbolGeometry(t *testing.T) {
	g := symbolGeometry(LegendProps{})
	if g.Size != LegendDefaults.SymbolSize || g.Shape != LegendDefaults.SymbolShape {
		t.Fatalf("defaults not applied: %+v", g)
	}
}

func TestSymbolLabelX(t *testing.T) {
	p := LegendProps{SymbolSize: 12, SymbolSpacing: 4, ItemDirection: LegendItemDirectionLeftToRight}
	if got := symbolLabelX(p); got != 16 {
		t.Fatalf("ltr label x = %v, want 16", got)
	}
	p.ItemDirection = LegendItemDirectionRightToLeft
	if got := symbolLabelX(p); got != -4 {
		t.Fatalf("rtl label x = %v, want -4", got)
	}
}
