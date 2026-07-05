package core

import (
	"fmt"
	"strconv"
)

// Shared SMIL enter-animation helpers used by the chart packages to render
// nivo-style "enter" transitions. Every helper emits nothing unless animation
// is enabled by the caller, so a chart rendered with Animate=false is
// byte-identical to the un-animated output (keeping goldens stable).
//
// The five geometry families reuse these as follows:
//   - rects / cells / paths → SMILFadeIn (opacity 0→1), optionally staggered
//   - circles               → SMILAnimate("r", "0", size)
//   - arcs                  → SMILAnimate("d", collapsedPath, finalPath)
//   - lines / areas         → SMILAnimate("d", collapsedPath, finalPath)
//
// The enter duration and easing match the bar/line/pie animations (600ms,
// fill="freeze").

const (
	// AnimateDuration is the enter-animation length (matches bar/line/pie).
	AnimateDuration = "0.6s"
	// AnimateBegin is the default start offset.
	AnimateBegin = "0s"
)

// SMILAnimate returns an <animate> element transitioning attributeName from→to
// over the enter duration, freezing on the final value. A blank begin defaults
// to AnimateBegin. Returns "" when from == to (nothing to animate).
func SMILAnimate(attr, from, to, begin string) string {
	if from == to {
		return ""
	}
	if begin == "" {
		begin = AnimateBegin
	}
	return fmt.Sprintf(
		`<animate attributeName=%q from=%q to=%q begin=%q dur=%q fill="freeze"></animate>`,
		attr, from, to, begin, AnimateDuration,
	)
}

// SMILFadeIn returns an <animate> fading opacity from 0 to 1 — the universal
// enter used by the rects/cells/paths families (works for any element type).
func SMILFadeIn(begin string) string {
	return SMILAnimate("opacity", "0", "1", begin)
}

// StaggerBegin returns the begin offset (e.g. "0.12s") for the i-th element
// given a per-element stagger in seconds. A non-positive stagger yields the
// shared AnimateBegin so all elements enter together.
func StaggerBegin(i int, stagger float64) string {
	if stagger <= 0 || i <= 0 {
		return AnimateBegin
	}
	return strconv.FormatFloat(float64(i)*stagger, 'f', -1, 64) + "s"
}
