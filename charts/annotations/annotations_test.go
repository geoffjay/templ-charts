package annotations

import (
	"math"
	"testing"
)

type point struct{ X, Y float64 }

func TestBindAnnotations_Circle(t *testing.T) {
	data := []point{{X: 10, Y: 20}, {X: 30, Y: 40}}
	specs := []AnnotationSpec[point]{
		{Type: AnnotationTypeCircle, Circle: &CircleAnnotationSpec[point]{
			Match:  func(d point) bool { return d.X == 30 },
			Radius: 8, Note: "hi", NoteOffsetX: 20, NoteOffsetY: 0,
		}},
	}
	bound := BindAnnotations(data, specs,
		func(d point) (float64, float64) { return d.X, d.Y },
		func(d point) (float64, float64) { return 0, 0 },
	)
	if len(bound) != 1 {
		t.Fatalf("got %d bound, want 1", len(bound))
	}
	if bound[0].X != 30 || bound[0].Y != 40 || bound[0].Radius != 8 {
		t.Fatalf("bound = %+v, want X=30 Y=40 R=8", bound[0])
	}
	if bound[0].Note != "hi" {
		t.Fatalf("note = %q", bound[0].Note)
	}
}

func TestBindAnnotations_Rect(t *testing.T) {
	data := []point{{X: 5, Y: 5}}
	specs := []AnnotationSpec[point]{
		{Type: AnnotationTypeRect, Rect: &RectAnnotationSpec[point]{
			Match: func(d point) bool { return true },
		}},
	}
	bound := BindAnnotations(data, specs,
		func(d point) (float64, float64) { return d.X, d.Y },
		func(d point) (float64, float64) { return 20, 10 },
	)
	if len(bound) != 1 {
		t.Fatalf("got %d, want 1", len(bound))
	}
	if bound[0].Width != 20 || bound[0].Height != 10 {
		t.Fatalf("rect dims = %v,%v, want 20,10", bound[0].Width, bound[0].Height)
	}
}

func TestBindAnnotations_NoMatch(t *testing.T) {
	data := []point{{X: 1, Y: 1}}
	specs := []AnnotationSpec[point]{
		{Type: AnnotationTypeDot, Dot: &DotAnnotationSpec[point]{
			Match: func(d point) bool { return d.X > 100 },
		}},
	}
	bound := BindAnnotations(data, specs,
		func(d point) (float64, float64) { return d.X, d.Y },
		func(d point) (float64, float64) { return 0, 0 },
	)
	if len(bound) != 0 {
		t.Fatalf("got %d bound, want 0", len(bound))
	}
}

func TestComputeAnnotation_Circle(t *testing.T) {
	b := BoundAnnotation{Type: AnnotationTypeCircle, X: 10, Y: 20, Radius: 8, Note: "n", NoteOffsetX: 30, NoteOffsetY: 0}
	inst := ComputeAnnotation(b)
	if inst.SymbolX != 10 || inst.SymbolY != 20 {
		t.Fatalf("symbol = (%v,%v), want (10,20)", inst.SymbolX, inst.SymbolY)
	}
	if inst.OutlineR != 8 {
		t.Fatalf("outline r = %v, want 8", inst.OutlineR)
	}
	if inst.NoteX != 40 || inst.NoteY != 20 {
		t.Fatalf("note = (%v,%v), want (40,20)", inst.NoteX, inst.NoteY)
	}
	if inst.LinkX1 != 10 || inst.LinkX2 != 40 {
		t.Fatalf("link x1=%v x2=%v, want 10,40", inst.LinkX1, inst.LinkX2)
	}
}

func TestComputeAnnotation_DotDefaults(t *testing.T) {
	b := BoundAnnotation{Type: AnnotationTypeDot, X: 5, Y: 5}
	inst := ComputeAnnotation(b)
	// Dot size defaults to defaultDotSize (4).
	if inst.Bound.Size != defaultDotSize {
		t.Fatalf("dot size = %v, want %d", inst.Bound.Size, defaultDotSize)
	}
}

func TestComputeAnnotation_Rect(t *testing.T) {
	b := BoundAnnotation{Type: AnnotationTypeRect, X: 0, Y: 0, Width: 40, Height: 20, Note: "r"}
	inst := ComputeAnnotation(b)
	if inst.OutlineW != 40 || inst.OutlineH != 20 {
		t.Fatalf("outline dims = %v,%v, want 40,20", inst.OutlineW, inst.OutlineH)
	}
	if inst.OutlineR != defaultRectBorderRadius {
		t.Fatalf("rect radius = %v, want %d", inst.OutlineR, defaultRectBorderRadius)
	}
	if inst.SymbolX != 20 || inst.SymbolY != 10 {
		t.Fatalf("rect symbol center = (%v,%v), want (20,10)", inst.SymbolX, inst.SymbolY)
	}
}

func TestGetLinkAngle(t *testing.T) {
	inst := AnnotationInstructions{LinkX1: 0, LinkY1: 0, LinkX2: 10, LinkY2: 0}
	if got := GetLinkAngle(inst); got != 0 {
		t.Fatalf("horizontal link angle = %v, want 0", got)
	}
	inst = AnnotationInstructions{LinkX1: 0, LinkY1: 0, LinkX2: 0, LinkY2: 10}
	if got := GetLinkAngle(inst); math.Abs(got-math.Pi/2) > 1e-9 {
		t.Fatalf("vertical link angle = %v, want π/2", got)
	}
	inst = AnnotationInstructions{}
	if got := GetLinkAngle(inst); got != 0 {
		t.Fatalf("degenerate link angle = %v, want 0", got)
	}
}
