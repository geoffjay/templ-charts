// compose.go — the Canvas *chart* layout: a wrapper that layers the data marks
// (drawn into a <canvas>) between two SVG panes — one below for content that
// must sit behind the marks (grid), one above for content that sits in front
// and/or handles pointer events (axes, legends, and the charts/interact mesh
// hit-surface). Everything shares one outer coordinate
// space: the canvas draw-list and both SVG panes translate their inner content
// by the chart margin, so a Canvas chart lines up exactly with its SVG twin.
//
// Reusing the existing SVG axis/grid/legend renderers for the panes means no
// layout code is duplicated for Canvas — only the high-cardinality marks move
// to the draw-list, which is where the large-N win is.
package canvas

import "strings"

// Compose returns the full Canvas-chart markup: a positioned wrapper div sized
// to outerW×outerH (carrying the chart background) containing, bottom to top,
// the underlay SVG pane (e.g. grid), the <canvas> of marks, and the overlay SVG
// pane (e.g. axes/legends/mesh). Each SVG pane wraps its inner content in a
// <g transform="translate(marginLeft,marginTop)"> so it aligns with the
// margin-translated draw-list; empty panes are omitted. The draw-list itself is
// responsible for the same margin translation (record a Translate first).
func Compose(id string, outerW, outerH, marginLeft, marginTop float64, background string, ops []Op, underlaySVG, overlaySVG string) string {
	var b strings.Builder
	b.WriteString(`<div class="tc-canvas-chart" style="position:relative;width:`)
	b.WriteString(formatNum(outerW))
	b.WriteString(`px;height:`)
	b.WriteString(formatNum(outerH))
	b.WriteString(`px`)
	if background != "" {
		b.WriteString(`;background:`)
		b.WriteString(background)
	}
	b.WriteString(`">`)
	if underlaySVG != "" {
		writePane(&b, outerW, outerH, marginLeft, marginTop, underlaySVG)
	}
	b.WriteString(`<div class="tc-canvas-marks" style="position:absolute;left:0;top:0">`)
	b.WriteString(Markup(id, outerW, outerH, ops))
	b.WriteString(`</div>`)
	if overlaySVG != "" {
		writePane(&b, outerW, outerH, marginLeft, marginTop, overlaySVG)
	}
	b.WriteString(`</div>`)
	return b.String()
}

// writePane emits one absolutely-positioned SVG pane wrapping inner in a
// margin-translate group. panes are transparent (no background), so lower
// layers show through.
func writePane(b *strings.Builder, w, h, mx, my float64, inner string) {
	b.WriteString(`<svg class="tc-canvas-pane" xmlns="http://www.w3.org/2000/svg" width="`)
	b.WriteString(formatNum(w))
	b.WriteString(`" height="`)
	b.WriteString(formatNum(h))
	b.WriteString(`" viewBox="0 0 `)
	b.WriteString(formatNum(w))
	b.WriteByte(' ')
	b.WriteString(formatNum(h))
	b.WriteString(`" style="position:absolute;left:0;top:0;overflow:visible"><g transform="translate(`)
	b.WriteString(formatNum(mx))
	b.WriteByte(',')
	b.WriteString(formatNum(my))
	b.WriteString(`)">`)
	b.WriteString(inner)
	b.WriteString(`</g></svg>`)
}
