// Package interact provides the client-side interactivity layer for
// templ-charts v2: a small vanilla-JS hover/tooltip/crosshair module (Script /
// ScriptTag) plus helpers for charts to emit the data-tc-* attributes it reads.
//
// The design (docs/PLAN-v2.md §5) is hybrid: ephemeral cursor-tracking
// interactions (tooltip show/position, crosshair, nearest-point hit-testing)
// run entirely client-side off data-* attributes with no server round-trip,
// while state-changing interactions (series toggle, active-arc) stay on the
// server via charts/htmx. Charts emit interaction attributes only when asked
// (an Interactive flag), so the default static render — and zero-JS embedding —
// is unchanged.
package interact

import (
	"fmt"
	"strings"
)

// TooltipAttrName is the attribute charts set on a hoverable element to carry
// its tooltip HTML; the client script shows it on hover.
const TooltipAttrName = "data-tc-tooltip"

// MeshAttrName is the attribute carrying a JSON array of {x,y,html} points for
// nearest-point (mesh) hover.
const MeshAttrName = "data-tc-mesh"

// ObserveAttrName is the opt-in attribute that turns on the ResizeObserver
// re-fetch: its value is a GET URL the client re-requests (with the new
// content-box ?w=&h= appended) whenever the element resizes, swapping the
// response in for a pixel-accurate re-render of an axes-heavy chart. An optional
// data-tc-observe-target CSS selector (ObserveTargetAttrName) redirects the swap
// to another element; absent, the observed element is swapped. Costs nothing on
// pages that never set it.
const ObserveAttrName = "data-tc-observe"

// ObserveTargetAttrName optionally overrides which element receives the
// ResizeObserver re-fetch swap (a CSS selector); default is the observed one.
const ObserveTargetAttrName = "data-tc-observe-target"

// TooltipHTML builds the standard one-line tooltip markup: a small color chip
// followed by "label: value" (or just value/label when the other is empty).
// All styling is inline so the tooltip renders correctly without any extra CSS
// beyond the .tc-chart-tooltip container box. The returned string is raw HTML
// (label/value are HTML-escaped); emit it into the data-tc-tooltip attribute
// via a templ attribute (auto-escaped) or EscapeAttr for hand-built SVG.
func TooltipHTML(color, label, value string) string {
	var b strings.Builder
	b.WriteString(`<span style="display:flex;align-items:center;white-space:pre">`)
	if color != "" {
		fmt.Fprintf(&b, `<span style="display:inline-block;width:12px;height:12px;margin-right:7px;background:%s"></span>`, color)
	}
	b.WriteString(`<span>`)
	b.WriteString(htmlEscape(tooltipText(label, value)))
	b.WriteString(`</span></span>`)
	return b.String()
}

// EscapeAttr escapes a string (typically TooltipHTML output) for a
// double-quoted attribute value, for charts that build SVG via string-builders
// rather than templ (which auto-escapes). The browser decodes it back to the
// original HTML when the client script reads the attribute.
func EscapeAttr(s string) string { return attrEscape(s) }

// tooltipText returns "label: value", or whichever side is non-empty.
func tooltipText(label, value string) string {
	switch {
	case label != "" && value != "":
		return label + ": " + value
	case value != "":
		return value
	default:
		return label
	}
}

// htmlEscape escapes text for HTML body context.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// attrEscape escapes a string for a double-quoted attribute value. The
// generated templ code and Go string-builders both emit data-tc-tooltip inside
// double quotes, so escaping quotes (and the HTML metacharacters that could
// break out of the attribute) is sufficient.
func attrEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
