// Package text provides shared text primitives: SVG text-anchor /
// dominant-baseline derivation from nivo's textAlign/textBaseline values,
// and a tick-label truncation helper. Mirrors @nivo/text.
package text

import (
	"github.com/geoffjay/templ-charts/charts/theming"
)

// SvgTextAttrs is the SVG attribute pair derived from nivo's textAlign +
// textBaseline. textAnchor maps to the horizontal alignment; dominantBaseline
// maps to the vertical alignment.
type SvgTextAttrs struct {
	TextAnchor       string
	DominantBaseline string
}

// SvgTextAttrs maps nivo's textAlign/textBaseline to SVG attribute values via
// the theming bridge tables. Mirrors @nivo/text's textProps → svgAttrs
// conversion.
func SvgTextAttrsFrom(textAlign, textBaseline string) SvgTextAttrs {
	anchor := theming.ConvertStyleAttribute(theming.EngineSVG, "textAlign", textAlign)
	baseline := theming.ConvertStyleAttribute(theming.EngineSVG, "textBaseline", textBaseline)
	return SvgTextAttrs{TextAnchor: anchor, DominantBaseline: baseline}
}

// TruncateTickAt truncates a tick label to fit within `length` px, keeping an
// ellipsis on the side closest to the chart. Mirrors @nivo/text's
// truncateTickAt. rotation is "0" | "45" | "-45" | "90" | "-90".
//
// The algorithm: when the label is wider than length (in approximation: width
// ≈ len(label) * charWidth), keep the half nearest the axis and append/prepend
// "…". A length <= 0 disables truncation.
func TruncateTickAt(label string, length float64, rotation string, charWidth float64) string {
	if length <= 0 || label == "" {
		return label
	}
	if charWidth <= 0 {
		charWidth = 6 // approx average char width at 11px sans-serif
	}
	width := float64(len(label)) * charWidth
	if width <= length {
		return label
	}
	keep := int(length / charWidth)
	if keep < 1 {
		keep = 1
	}
	switch rotation {
	case "45", "-45", "90", "-90":
		// Rotated labels grow away from the axis on the leading side; keep the
		// first `keep` chars + ellipsis.
		if keep >= len(label) {
			return label
		}
		return label[:keep] + "…"
	default:
		// Horizontal labels are centered on the tick — split keep across both
		// sides with an ellipsis in the middle.
		if keep >= len(label) {
			return label
		}
		half := keep / 2
		if half < 1 {
			half = 1
		}
		head := label[:half]
		tail := label[len(label)-(keep-half):]
		return head + "…" + tail
	}
}
