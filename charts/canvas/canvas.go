// Package canvas is templ-charts' Canvas rendering backend (docs/PLAN-v6.md §4).
//
// A Go library cannot draw pixels in the browser, but the chart hooks already
// *compute* every mark's geometry server-side. This package captures that
// geometry as a compact, ordered draw-list — a sequence of primitive ops
// (fillStyle, fillRect, fillCircle, fillText, line, fillPath, …) recorded by a
// Recorder, the Canvas analogue of the SVG string-builder — and ships a small
// dependency-free JS renderer (Script / CanvasScriptTag, siblings of
// charts/interact's) that replays the list into a <canvas> 2D context.
//
// The draw-list is plain data: it golden-tests as canonical JSON (EncodeJSON,
// 3-decimal rounding like the SVG path) with no pixels in CI, and it is
// backend-agnostic — the JS replay is one consumer; a future pure-Go PNG
// rasterizer (§4.4) is another, over the same []Op. Canvas is opt-in per chart
// via a RenderEngine prop (SVG stays the default), so existing SVG goldens are
// byte-stable; this package adds the backend, and the chart wiring lands in the
// Canvas-variant tranches (§5).
package canvas

// OpKind identifies a draw-list primitive. The numeric args (A,B,C,D) and the
// string arg (Str) that each kind uses are documented per Recorder method; the
// JSON encoder and the JS replay agree on the same mapping.
type OpKind uint8

const (
	OpFillStyle    OpKind = iota // Str = color
	OpStrokeStyle                // Str = color
	OpLineWidth                  // A = width
	OpGlobalAlpha                // A = alpha (0..1)
	OpFont                       // Str = CSS font shorthand
	OpTextAlign                  // Str = start|end|left|right|center
	OpTextBaseline               // Str = top|middle|bottom|alphabetic
	OpSave                       // (no args)
	OpRestore                    // (no args)
	OpTranslate                  // A,B = dx,dy
	OpFillRect                   // A,B,C,D = x,y,w,h
	OpStrokeRect                 // A,B,C,D = x,y,w,h
	OpFillCircle                 // A,B,C = cx,cy,r
	OpStrokeCircle               // A,B,C = cx,cy,r
	OpFillText                   // Str = text; A,B = x,y
	OpLine                       // A,B,C,D = x1,y1,x2,y2
	OpFillPath                   // Str = SVG path data
	OpStrokePath                 // Str = SVG path data
)

// Op is a single draw-list primitive. A flat struct (kind + up to four floats +
// one string) keeps the list a plain, allocation-light slice that both the JSON
// encoder and any future rasterizer can switch over without per-op boxing.
type Op struct {
	Kind       OpKind
	A, B, C, D float64
	Str        string
}

// Recorder accumulates draw ops in order — the Canvas analogue of the SVG
// string-builder each chart's render layer writes into. It is backend-agnostic:
// call Ops (or EncodeJSON) to hand the recorded list to a renderer.
type Recorder struct {
	ops []Op
}

// NewRecorder returns an empty Recorder.
func NewRecorder() *Recorder { return &Recorder{} }

// Ops returns the recorded draw-list. The slice aliases the Recorder's storage;
// treat it as read-only.
func (r *Recorder) Ops() []Op { return r.ops }

// Len reports the number of recorded ops.
func (r *Recorder) Len() int { return len(r.ops) }

func (r *Recorder) push(op Op) { r.ops = append(r.ops, op) }

// --- state ops -------------------------------------------------------------

// FillStyle sets the fill color for subsequent fill ops (any CSS color string).
func (r *Recorder) FillStyle(color string) { r.push(Op{Kind: OpFillStyle, Str: color}) }

// StrokeStyle sets the stroke color for subsequent stroke ops.
func (r *Recorder) StrokeStyle(color string) { r.push(Op{Kind: OpStrokeStyle, Str: color}) }

// LineWidth sets the stroke width in CSS pixels.
func (r *Recorder) LineWidth(w float64) { r.push(Op{Kind: OpLineWidth, A: w}) }

// GlobalAlpha sets the global opacity (0..1) applied to subsequent draws.
func (r *Recorder) GlobalAlpha(a float64) { r.push(Op{Kind: OpGlobalAlpha, A: a}) }

// Font sets the text font (a CSS font shorthand, e.g. "11px sans-serif").
func (r *Recorder) Font(font string) { r.push(Op{Kind: OpFont, Str: font}) }

// TextAlign sets the horizontal text alignment (canvas textAlign values).
func (r *Recorder) TextAlign(a string) { r.push(Op{Kind: OpTextAlign, Str: a}) }

// TextBaseline sets the vertical text baseline (canvas textBaseline values).
func (r *Recorder) TextBaseline(b string) { r.push(Op{Kind: OpTextBaseline, Str: b}) }

// Save pushes the current drawing state (ctx.save).
func (r *Recorder) Save() { r.push(Op{Kind: OpSave}) }

// Restore pops the last saved drawing state (ctx.restore).
func (r *Recorder) Restore() { r.push(Op{Kind: OpRestore}) }

// Translate offsets the coordinate origin by (dx,dy) (ctx.translate).
func (r *Recorder) Translate(dx, dy float64) { r.push(Op{Kind: OpTranslate, A: dx, B: dy}) }

// --- drawing ops -----------------------------------------------------------

// FillRect fills an axis-aligned rectangle — the heatmap/bar primitive.
func (r *Recorder) FillRect(x, y, w, h float64) {
	r.push(Op{Kind: OpFillRect, A: x, B: y, C: w, D: h})
}

// StrokeRect strokes an axis-aligned rectangle outline.
func (r *Recorder) StrokeRect(x, y, w, h float64) {
	r.push(Op{Kind: OpStrokeRect, A: x, B: y, C: w, D: h})
}

// FillCircle fills a circle of radius r centred at (cx,cy) — the scatterplot
// primitive. Kept a dedicated op (rather than a path) so thousands of points
// stay compact in the draw-list.
func (r *Recorder) FillCircle(cx, cy, radius float64) {
	r.push(Op{Kind: OpFillCircle, A: cx, B: cy, C: radius})
}

// StrokeCircle strokes a circle outline of radius r centred at (cx,cy).
func (r *Recorder) StrokeCircle(cx, cy, radius float64) {
	r.push(Op{Kind: OpStrokeCircle, A: cx, B: cy, C: radius})
}

// FillText draws text with the current font/align/baseline at (x,y).
func (r *Recorder) FillText(text string, x, y float64) {
	r.push(Op{Kind: OpFillText, Str: text, A: x, B: y})
}

// Line strokes a single segment from (x1,y1) to (x2,y2) — axis rules/gridlines.
func (r *Recorder) Line(x1, y1, x2, y2 float64) {
	r.push(Op{Kind: OpLine, A: x1, B: y1, C: x2, D: y2})
}

// FillPath fills an SVG path (the same "d" string the SVG builders produce) via
// the browser's Path2D — so charts that already build arc/line/area paths can
// emit them straight into the Canvas backend with no geometry duplication.
func (r *Recorder) FillPath(d string) { r.push(Op{Kind: OpFillPath, Str: d}) }

// StrokePath strokes an SVG path (via Path2D), the stroke counterpart of FillPath.
func (r *Recorder) StrokePath(d string) { r.push(Op{Kind: OpStrokePath, Str: d}) }
