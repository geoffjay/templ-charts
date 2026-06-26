// Package annotations mirrors @nivo/annotations: annotation matchers, the
// circle/dot/rect annotation specs, BindAnnotations / ComputeAnnotation /
// GetLinkAngle, and the Annotation (link + outline + symbol + note) templ
// component family.
package annotations

import "math"

// RelativeOrAbsolutePosition is a position that may be relative (0..1) or
// absolute (px). The discriminator is in the chart's bind step.
type RelativeOrAbsolutePosition struct {
	Relative   float64
	Absolute   float64
	IsRelative bool
}

// AnnotationMatcher selects chart data points that an annotation should
// attach to. Match returns true for matching datums.
type AnnotationMatcher[D any] func(d D) bool

// AnnotationType enumerates the annotation outline shapes.
type AnnotationType string

const (
	AnnotationTypeCircle AnnotationType = "circle"
	AnnotationTypeDot    AnnotationType = "dot"
	AnnotationTypeRect   AnnotationType = "rect"
)

// CircleAnnotationSpec annotates a data point with a circle outline.
type CircleAnnotationSpec[D any] struct {
	Match       AnnotationMatcher[D]
	Radius      float64
	OffsetX     float64
	OffsetY     float64
	Note        string
	NoteX       float64
	NoteY       float64
	NoteOffsetX float64
	NoteOffsetY float64
}

// DotAnnotationSpec annotates a data point with a filled dot.
type DotAnnotationSpec[D any] struct {
	Match       AnnotationMatcher[D]
	Size        float64
	OffsetX     float64
	OffsetY     float64
	Note        string
	NoteX       float64
	NoteY       float64
	NoteOffsetX float64
	NoteOffsetY float64
}

// RectAnnotationSpec annotates a rect region.
type RectAnnotationSpec[D any] struct {
	Match        AnnotationMatcher[D]
	X            float64
	Y            float64
	Width        float64
	Height       float64
	OffsetX      float64
	OffsetY      float64
	BorderRadius float64
	Note         string
	NoteX        float64
	NoteY        float64
	NoteOffsetX  float64
	NoteOffsetY  float64
}

// AnnotationSpec is the union of circle/dot/rect specs.
type AnnotationSpec[D any] struct {
	Type   AnnotationType
	Circle *CircleAnnotationSpec[D]
	Dot    *DotAnnotationSpec[D]
	Rect   *RectAnnotationSpec[D]
}

// BoundAnnotation is an annotation spec resolved to a concrete (x, y) +
// dimensions on the chart, ready to render.
type BoundAnnotation struct {
	Type             AnnotationType
	X, Y             float64 // symbol position
	Width, Height    float64 // for rect
	Radius           float64 // for circle
	Size             float64 // for dot
	OffsetX, OffsetY float64
	Note             string
	NoteX            float64
	NoteY            float64
	NoteOffsetX      float64
	NoteOffsetY      float64
}

// AnnotationInstructions is the resolved rendering plan for one annotation:
// the link line + the outline shape + the symbol + the note text.
type AnnotationInstructions struct {
	Bound                          BoundAnnotation
	LinkX1, LinkY1, LinkX2, LinkY2 float64
	OutlineX, OutlineY             float64
	OutlineW, OutlineH             float64
	OutlineR                       float64 // rect corner radius
	SymbolX, SymbolY               float64
	NoteX, NoteY                   float64
}

// GetPositionFunc returns the (x, y) position of a datum in chart coordinates.
type GetPositionFunc[D any] func(d D) (x, y float64)

// GetDimensionsFunc returns the (width, height) of a datum's bounding box.
type GetDimensionsFunc[D any] func(d D) (w, h float64)

// BindAnnotations matches annotation specs against data and produces a bound
// annotation per match. Mirrors @nivo/annotations bindAnnotations.
func BindAnnotations[D any](
	data []D,
	matchers []AnnotationSpec[D],
	getPosition GetPositionFunc[D],
	getDimensions GetDimensionsFunc[D],
) []BoundAnnotation {
	out := []BoundAnnotation{}
	for _, m := range matchers {
		for _, d := range data {
			switch m.Type {
			case AnnotationTypeCircle:
				if m.Circle == nil || !m.Circle.Match(d) {
					continue
				}
				x, y := getPosition(d)
				out = append(out, BoundAnnotation{
					Type: AnnotationTypeCircle,
					X:    x, Y: y,
					Radius:  m.Circle.Radius,
					OffsetX: m.Circle.OffsetX, OffsetY: m.Circle.OffsetY,
					Note:  m.Circle.Note,
					NoteX: m.Circle.NoteX, NoteY: m.Circle.NoteY,
					NoteOffsetX: m.Circle.NoteOffsetX, NoteOffsetY: m.Circle.NoteOffsetY,
				})
			case AnnotationTypeDot:
				if m.Dot == nil || !m.Dot.Match(d) {
					continue
				}
				x, y := getPosition(d)
				out = append(out, BoundAnnotation{
					Type: AnnotationTypeDot,
					X:    x, Y: y,
					Size:    m.Dot.Size,
					OffsetX: m.Dot.OffsetX, OffsetY: m.Dot.OffsetY,
					Note:  m.Dot.Note,
					NoteX: m.Dot.NoteX, NoteY: m.Dot.NoteY,
					NoteOffsetX: m.Dot.NoteOffsetX, NoteOffsetY: m.Dot.NoteOffsetY,
				})
			case AnnotationTypeRect:
				if m.Rect == nil || !m.Rect.Match(d) {
					continue
				}
				x, y := getPosition(d)
				w, h := getDimensions(d)
				out = append(out, BoundAnnotation{
					Type: AnnotationTypeRect,
					X:    x, Y: y,
					Width: w, Height: h,
					OffsetX: m.Rect.OffsetX, OffsetY: m.Rect.OffsetY,
					Note:  m.Rect.Note,
					NoteX: m.Rect.NoteX, NoteY: m.Rect.NoteY,
					NoteOffsetX: m.Rect.NoteOffsetX, NoteOffsetY: m.Rect.NoteOffsetY,
				})
			}
		}
	}
	return out
}

// Defaults mirror @nivo/annotations default note/symbol sizing.
const (
	defaultNoteWidth        = 120
	defaultNoteTextOffset   = 8
	defaultDotSize          = 4
	defaultRectBorderRadius = 6
)

// ComputeAnnotation turns a BoundAnnotation into rendering instructions:
// link line from the symbol to the note, the outline shape, the symbol, and
// the note text position. Mirrors @nivo/annotations computeAnnotation.
func ComputeAnnotation(b BoundAnnotation) AnnotationInstructions {
	inst := InstructionsWithDefaults(b)
	switch b.Type {
	case AnnotationTypeCircle:
		inst.OutlineX = b.X + b.OffsetX
		inst.OutlineY = b.Y + b.OffsetY
		inst.OutlineR = b.Radius
		inst.SymbolX = b.X + b.OffsetX
		inst.SymbolY = b.Y + b.OffsetY
	case AnnotationTypeDot:
		inst.SymbolX = b.X + b.OffsetX
		inst.SymbolY = b.Y + b.OffsetY
		inst.OutlineX = inst.SymbolX
		inst.OutlineY = inst.SymbolY
	case AnnotationTypeRect:
		inst.OutlineX = b.X + b.OffsetX
		inst.OutlineY = b.Y + b.OffsetY
		inst.OutlineW = b.Width
		inst.OutlineH = b.Height
		inst.OutlineR = defaultRectBorderRadius
		inst.SymbolX = inst.OutlineX + inst.OutlineW/2
		inst.SymbolY = inst.OutlineY + inst.OutlineH/2
	}
	// Note position: default to the symbol + an offset when not specified.
	if b.NoteX == 0 && b.NoteY == 0 {
		inst.NoteX = inst.SymbolX + b.NoteOffsetX
		inst.NoteY = inst.SymbolY + b.NoteOffsetY
	} else {
		inst.NoteX = b.NoteX
		inst.NoteY = b.NoteY
	}
	// Link from the symbol to the note.
	inst.LinkX1 = inst.SymbolX
	inst.LinkY1 = inst.SymbolY
	inst.LinkX2 = inst.NoteX
	inst.LinkY2 = inst.NoteY
	return inst
}

// InstructionsWithDefaults initializes an AnnotationInstructions from a bound
// annotation, applying the package defaults.
func InstructionsWithDefaults(b BoundAnnotation) AnnotationInstructions {
	inst := AnnotationInstructions{Bound: b}
	if b.Size == 0 && b.Type == AnnotationTypeDot {
		inst.Bound.Size = defaultDotSize
	}
	return inst
}

// GetLinkAngle returns the angle (radians) of the link line from the symbol to
// the note. Mirrors @nivo/annotations getLinkAngle.
func GetLinkAngle(inst AnnotationInstructions) float64 {
	dx := inst.LinkX2 - inst.LinkX1
	dy := inst.LinkY2 - inst.LinkY1
	if dx == 0 && dy == 0 {
		return 0
	}
	return math.Atan2(dy, dx)
}
