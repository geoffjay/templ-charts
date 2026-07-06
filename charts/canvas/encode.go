// encode.go — the JSON transport/golden form of a draw-list, plus the HTML
// markup that pairs a <canvas> with its ops for the replay script.
//
// The draw-list encodes as a JSON array of `[opcode, args…]` arrays, one op per
// line: compact enough to embed inline, fully deterministic (array order, no
// map-key ordering), and diffable as a golden. Coordinates are rounded to 3
// decimals with the shortest round-trip representation, matching the SVG path
// serializer used across internal/d3. The same JSON is what the JS replay
// (Script) decodes, so goldens are exactly what ships to the browser.
package canvas

import (
	"math"
	"strconv"
	"strings"
)

// opcode maps each OpKind to the short string tag shared by the encoder and the
// JS replay switch.
var opcode = [...]string{
	OpFillStyle:    "fs",
	OpStrokeStyle:  "ss",
	OpLineWidth:    "lw",
	OpGlobalAlpha:  "ga",
	OpFont:         "ft",
	OpTextAlign:    "ta",
	OpTextBaseline: "tb",
	OpSave:         "sv",
	OpRestore:      "rs",
	OpTranslate:    "tr",
	OpFillRect:     "fr",
	OpStrokeRect:   "sr",
	OpFillCircle:   "fc",
	OpStrokeCircle: "sc",
	OpFillText:     "fx",
	OpLine:         "ln",
	OpFillPath:     "pf",
	OpStrokePath:   "sp",
}

// EncodeJSON serialises a draw-list to its canonical JSON form: a `[...]` array
// with one `[opcode, args…]` element per line. Deterministic and 3-decimal
// rounded — suitable both as the browser transport and as a golden snapshot.
func EncodeJSON(ops []Op) string {
	var b strings.Builder
	b.WriteByte('[')
	for i := range ops {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('\n')
		encodeOp(&b, ops[i])
	}
	if len(ops) > 0 {
		b.WriteByte('\n')
	}
	b.WriteByte(']')
	return b.String()
}

func encodeOp(b *strings.Builder, op Op) {
	b.WriteByte('[')
	writeJSONString(b, opcode[op.Kind])
	switch op.Kind {
	case OpFillStyle, OpStrokeStyle, OpFont, OpTextAlign, OpTextBaseline:
		b.WriteByte(',')
		writeJSONString(b, op.Str)
	case OpLineWidth, OpGlobalAlpha:
		writeNums(b, op.A)
	case OpTranslate:
		writeNums(b, op.A, op.B)
	case OpFillRect, OpStrokeRect, OpLine:
		writeNums(b, op.A, op.B, op.C, op.D)
	case OpFillCircle, OpStrokeCircle:
		writeNums(b, op.A, op.B, op.C)
	case OpFillText:
		b.WriteByte(',')
		writeJSONString(b, op.Str)
		writeNums(b, op.A, op.B)
	case OpFillPath, OpStrokePath:
		b.WriteByte(',')
		writeJSONString(b, op.Str)
	case OpSave, OpRestore:
		// no args
	}
	b.WriteByte(']')
}

func writeNums(b *strings.Builder, nums ...float64) {
	for _, n := range nums {
		b.WriteByte(',')
		b.WriteString(formatNum(n))
	}
}

// formatNum rounds to 3 decimals and formats with the shortest round-trip
// representation (matching internal/d3/shape's floatToShortest convention).
func formatNum(v float64) string {
	const k = 1000
	v = math.Round(v*k) / k
	if v == 0 {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// writeJSONString writes s as a quoted JSON string, escaping the characters
// encoding/json escapes by default (including <, > and & for safe embedding in
// an HTML <script> element).
func writeJSONString(b *strings.Builder, s string) {
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case '<':
			b.WriteString("\\u003c")
		case '>':
			b.WriteString("\\u003e")
		case '&':
			b.WriteString("\\u0026")
		default:
			if r < 0x20 {
				b.WriteString(`\u00`)
				const hex = "0123456789abcdef"
				b.WriteByte(hex[r>>4])
				b.WriteByte(hex[r&0xf])
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
}

// Markup returns the HTML for a Canvas-rendered chart: a <canvas> filling its
// wrapper (which carries the chart's intrinsic width×height + aspect-ratio)
// paired (by id) with a <script type="application/json"> carrying its
// draw-list. data-tc-w/data-tc-h record the chart coordinate space; the
// replay script (CanvasScriptTag) finds each canvas, reads its ops, sizes the
// backing store to displayed-size × devicePixelRatio, scales chart space onto
// it, and renders. The width/height attributes double as a no-JS fallback
// size, and the ops script keeps the (potentially large) draw-list out of an
// HTML attribute so it needs no attribute-escaping.
func Markup(id string, width, height float64, ops []Op) string {
	w := formatNum(width)
	h := formatNum(height)
	opsID := id + "-ops"
	var b strings.Builder
	b.WriteString(`<canvas class="tc-canvas" id="`)
	b.WriteString(attrEscape(id))
	b.WriteString(`" data-tc-canvas="`)
	b.WriteString(attrEscape(opsID))
	b.WriteString(`" data-tc-w="`)
	b.WriteString(w)
	b.WriteString(`" data-tc-h="`)
	b.WriteString(h)
	b.WriteString(`" width="`)
	b.WriteString(w)
	b.WriteString(`" height="`)
	b.WriteString(h)
	b.WriteString(`" style="display:block;width:100%;height:100%"></canvas><script type="application/json" class="tc-canvas-ops" id="`)
	b.WriteString(attrEscape(opsID))
	b.WriteString(`">`)
	b.WriteString(EncodeJSON(ops))
	b.WriteString(`</script>`)
	return b.String()
}

// attrEscape escapes a string for a double-quoted HTML attribute value.
func attrEscape(s string) string {
	if !strings.ContainsAny(s, `"&<>`) {
		return s
	}
	r := strings.NewReplacer(`&`, "&amp;", `"`, "&#34;", `<`, "&lt;", `>`, "&gt;")
	return r.Replace(s)
}
