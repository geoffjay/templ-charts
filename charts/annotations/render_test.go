package annotations_test

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
	"github.com/geoffjay/templ-charts/charts/annotations"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return b.String()
}

func TestAnnotation_Circle(t *testing.T) {
	inst := annotations.ComputeAnnotation(annotations.BoundAnnotation{
		Type: annotations.AnnotationTypeCircle,
		X:    50, Y: 60,
		Radius:      12,
		Note:        "peak",
		NoteOffsetX: 30, NoteOffsetY: -20,
	})
	out := render(t, annotations.Annotation(inst))
	for _, want := range []string{
		`style="pointer-events: none"`,
		// Link from the symbol (50,60) to the note (80,40).
		`<line x1="50" y1="60" x2="80" y2="40" stroke="#000" stroke-width="1"`,
		`<circle cx="50" cy="60" r="12" fill="none" stroke="#000" stroke-width="2"`,
		`<text x="80" y="40" font-size="13">peak</text>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestAnnotation_Dot(t *testing.T) {
	inst := annotations.ComputeAnnotation(annotations.BoundAnnotation{
		Type: annotations.AnnotationTypeDot,
		X:    10, Y: 20,
		Size: 5,
		Note: "dot note",
	})
	out := render(t, annotations.Annotation(inst))
	if !strings.Contains(out, `<circle cx="10" cy="20" r="5" fill="#000"`) {
		t.Errorf("dot outline missing: %s", out)
	}
	if !strings.Contains(out, ">dot note</text>") {
		t.Errorf("note missing: %s", out)
	}
}

func TestAnnotation_DotDefaultSize(t *testing.T) {
	inst := annotations.ComputeAnnotation(annotations.BoundAnnotation{
		Type: annotations.AnnotationTypeDot,
		X:    10, Y: 20,
	})
	out := render(t, annotations.Annotation(inst))
	// Zero Size falls back to the default dot size of 4.
	if !strings.Contains(out, `r="4"`) {
		t.Errorf("default dot size missing: %s", out)
	}
}

func TestAnnotation_Rect(t *testing.T) {
	inst := annotations.ComputeAnnotation(annotations.BoundAnnotation{
		Type: annotations.AnnotationTypeRect,
		X:    100, Y: 40,
		Width: 60, Height: 30,
		Note: "region",
	})
	out := render(t, annotations.Annotation(inst))
	// Rect outline with the default border radius of 6, symbol at its center.
	if !strings.Contains(out, `<rect x="100" y="40" width="60" height="30" rx="6" ry="6" fill="none" stroke="#000" stroke-width="2"`) {
		t.Errorf("rect outline missing: %s", out)
	}
	if !strings.Contains(out, `<line x1="130" y1="55"`) {
		t.Errorf("link should start at rect center: %s", out)
	}
	if !strings.Contains(out, ">region</text>") {
		t.Errorf("note missing: %s", out)
	}
}

func TestAnnotation_ExplicitNotePosition(t *testing.T) {
	inst := annotations.ComputeAnnotation(annotations.BoundAnnotation{
		Type: annotations.AnnotationTypeCircle,
		X:    10, Y: 10,
		Radius: 5,
		Note:   "n",
		NoteX:  99, NoteY: 88,
	})
	out := render(t, annotations.Annotation(inst))
	if !strings.Contains(out, `<text x="99" y="88"`) {
		t.Errorf("explicit note position not used: %s", out)
	}
	if !strings.Contains(out, `x2="99" y2="88"`) {
		t.Errorf("link should end at note position: %s", out)
	}
}

func TestAnnotation_NoNote(t *testing.T) {
	inst := annotations.ComputeAnnotation(annotations.BoundAnnotation{
		Type: annotations.AnnotationTypeCircle,
		X:    10, Y: 10,
		Radius: 5,
	})
	out := render(t, annotations.Annotation(inst))
	if strings.Contains(out, "<text") {
		t.Errorf("note text rendered without a Note: %s", out)
	}
	if !strings.Contains(out, "<circle") {
		t.Errorf("outline missing: %s", out)
	}
}

func TestAnnotation_UnknownTypeRendersNoOutline(t *testing.T) {
	inst := annotations.ComputeAnnotation(annotations.BoundAnnotation{
		Type: annotations.AnnotationType("bogus"),
		X:    10, Y: 10,
	})
	out := render(t, annotations.Annotation(inst))
	if strings.Contains(out, "<circle") || strings.Contains(out, "<rect") {
		t.Errorf("unexpected outline for unknown type: %s", out)
	}
	// The link line is always rendered.
	if got := strings.Count(out, "<line"); got != 1 {
		t.Errorf("line count = %d, want 1", got)
	}
}

func TestAnnotation_Offsets(t *testing.T) {
	inst := annotations.ComputeAnnotation(annotations.BoundAnnotation{
		Type: annotations.AnnotationTypeCircle,
		X:    10, Y: 20,
		OffsetX: 5, OffsetY: -5,
		Radius: 3,
	})
	out := render(t, annotations.Annotation(inst))
	if !strings.Contains(out, `<circle cx="15" cy="15" r="3"`) {
		t.Errorf("offsets not applied to outline: %s", out)
	}
}

func TestCircleAnnotationOutline(t *testing.T) {
	out := render(t, annotations.CircleAnnotationOutline(annotations.AnnotationInstructions{
		OutlineX: 1.25, OutlineY: 2.5, OutlineR: 7,
	}))
	if !strings.Contains(out, `<circle cx="1.25" cy="2.5" r="7" fill="none" stroke="#000" stroke-width="2"`) {
		t.Errorf("circle outline wrong: %s", out)
	}
}

func TestDotAnnotationOutline(t *testing.T) {
	out := render(t, annotations.DotAnnotationOutline(annotations.AnnotationInstructions{
		SymbolX: 3, SymbolY: 4,
		Bound: annotations.BoundAnnotation{Size: 6},
	}))
	if !strings.Contains(out, `<circle cx="3" cy="4" r="6" fill="#000"`) {
		t.Errorf("dot outline wrong: %s", out)
	}
}

func TestRectAnnotationOutline(t *testing.T) {
	out := render(t, annotations.RectAnnotationOutline(annotations.AnnotationInstructions{
		OutlineX: 5, OutlineY: 6, OutlineW: 70, OutlineH: 35, OutlineR: 4,
	}))
	if !strings.Contains(out, `<rect x="5" y="6" width="70" height="35" rx="4" ry="4" fill="none" stroke="#000" stroke-width="2"`) {
		t.Errorf("rect outline wrong: %s", out)
	}
}

func TestAnnotations_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var b strings.Builder
	inst := annotations.AnnotationInstructions{}
	if err := annotations.Annotation(inst).Render(ctx, &b); err == nil {
		t.Error("Annotation: expected error with canceled context")
	}
	if err := annotations.CircleAnnotationOutline(inst).Render(ctx, &b); err == nil {
		t.Error("CircleAnnotationOutline: expected error with canceled context")
	}
	if err := annotations.DotAnnotationOutline(inst).Render(ctx, &b); err == nil {
		t.Error("DotAnnotationOutline: expected error with canceled context")
	}
	if err := annotations.RectAnnotationOutline(inst).Render(ctx, &b); err == nil {
		t.Error("RectAnnotationOutline: expected error with canceled context")
	}
}

func TestAnnotation_NonFiniteValuesRenderAsZero(t *testing.T) {
	// NaN/Inf coordinates render as 0; a tiny negative normalizes -0 to 0.
	out := render(t, annotations.CircleAnnotationOutline(annotations.AnnotationInstructions{
		OutlineX: math.NaN(), OutlineY: math.Inf(1), OutlineR: -1e-12,
	}))
	if !strings.Contains(out, `<circle cx="0" cy="0" r="0"`) {
		t.Errorf("non-finite values should render as 0: %s", out)
	}
}

func TestBindAnnotations_DotAndNilSpecs(t *testing.T) {
	type point struct{ X, Y float64 }
	data := []point{{X: 1, Y: 2}}
	specs := []annotations.AnnotationSpec[point]{
		// Dot spec that matches.
		{Type: annotations.AnnotationTypeDot, Dot: &annotations.DotAnnotationSpec[point]{
			Match: func(d point) bool { return true },
			Size:  7, Note: "d",
		}},
		// Specs whose payload pointer is nil are skipped.
		{Type: annotations.AnnotationTypeCircle},
		{Type: annotations.AnnotationTypeDot},
		{Type: annotations.AnnotationTypeRect},
		// Unknown type is skipped.
		{Type: annotations.AnnotationType("bogus")},
	}
	bound := annotations.BindAnnotations(data, specs,
		func(d point) (float64, float64) { return d.X, d.Y },
		func(d point) (float64, float64) { return 0, 0 },
	)
	if len(bound) != 1 {
		t.Fatalf("got %d bound, want 1", len(bound))
	}
	if bound[0].Type != annotations.AnnotationTypeDot || bound[0].Size != 7 || bound[0].X != 1 || bound[0].Y != 2 {
		t.Fatalf("bound = %+v, want dot size 7 at (1,2)", bound[0])
	}
}

// failingWriter fails every Write after the first `remaining` calls succeed.
type failingWriter struct{ remaining int }

func (w *failingWriter) Write(p []byte) (int, error) {
	if w.remaining <= 0 {
		return 0, errors.New("write failed")
	}
	w.remaining--
	return len(p), nil
}

// TestAnnotations_WriteErrorsPropagate verifies that a failing writer surfaces
// an error from Render no matter which write fails (walking the failure point
// through every write in the component).
func TestAnnotations_WriteErrorsPropagate(t *testing.T) {
	old := templruntime.DefaultBufferSize
	templruntime.DefaultBufferSize = 1
	defer func() { templruntime.DefaultBufferSize = old }()

	components := map[string]templ.Component{
		"circle": annotations.Annotation(annotations.ComputeAnnotation(annotations.BoundAnnotation{
			Type: annotations.AnnotationTypeCircle, X: 1, Y: 2, Radius: 3, Note: "n",
		})),
		"dot": annotations.Annotation(annotations.ComputeAnnotation(annotations.BoundAnnotation{
			Type: annotations.AnnotationTypeDot, X: 1, Y: 2, Note: "n",
		})),
		"rect": annotations.Annotation(annotations.ComputeAnnotation(annotations.BoundAnnotation{
			Type: annotations.AnnotationTypeRect, X: 1, Y: 2, Width: 10, Height: 5, Note: "n",
		})),
	}
	for name, c := range components {
		t.Run(name, func(t *testing.T) {
			for n := 0; n < 10000; n++ {
				buf := new(templruntime.Buffer)
				buf.Reset(&failingWriter{remaining: n})
				err := c.Render(context.Background(), buf)
				if err == nil {
					err = buf.Flush()
				}
				if err == nil {
					if n == 0 {
						t.Fatal("expected an error when every write fails")
					}
					return // failure point walked past the last write
				}
			}
			t.Fatal("render never succeeded")
		})
	}
}
