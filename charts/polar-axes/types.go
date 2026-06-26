// Package polaraxes provides the polar-axis scaffolding (CircularAxis,
// RadialAxis, PolarGrid + their grid components) ported from @nivo/polar-axes.
// It is scaffold-only in v1 — pie does its own layout and no v1 chart uses
// these — but the components are fully implemented so future radar/chord/
// sunburst chart types plug in without rework.
//
// Angle convention (matching @nivo/polar-axes): user-facing angles are in
// degrees, 0 = top, clockwise. Internally the -90 offset (so angle 0 maps to
// the +x axis / 3 o'clock) is applied before passing to PositionFromAngle.
// All animation is SMIL (<animate> on transform/opacity), gated by Animate,
// since there is no react-spring in the server-rendered SVG world.
package polaraxes

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// CircularAxisType is whether the circular axis is drawn at the inner or outer
// edge of the circle (affects tick direction — inward vs outward).
type CircularAxisType string

const (
	CircularAxisInner CircularAxisType = "inner"
	CircularAxisOuter CircularAxisType = "outer"
)

// CircularAxisConfig mirrors @nivo/polar-axes CircularAxisConfig.
type CircularAxisConfig struct {
	TickSize      float64
	TickPadding   float64
	TickRotation  float64
	Format        core.ValueFormat[any] // label formatter
	TickComponent *CircularTickRenderer // override; nil = default
	Style         *theming.PartialAxisTheme
	AriaHidden    bool
}

// RadialAxisConfig mirrors @nivo/polar-axes RadialAxisConfig.
type RadialAxisConfig struct {
	Ticks         scales.TicksSpec
	TickSize      float64
	TickPadding   float64
	TickRotation  float64
	Format        core.ValueFormat[any]
	TickComponent *RadialTickRenderer
	Style         *theming.PartialAxisTheme
	AriaHidden    bool
}

// TicksPosition is whether radial-axis ticks are drawn before (toward center)
// or after (away from center) the axis line.
type TicksPosition string

const (
	TicksBefore TicksPosition = "before"
	TicksAfter  TicksPosition = "after"
)

// CircularTickRenderer renders one circular-axis tick. Override via
// CircularAxisConfig.TickComponent; nil uses the default CircularAxisTick.
type CircularTickRenderer = func(props CircularAxisTickProps) string

// RadialTickRenderer renders one radial-axis tick. Override via
// RadialAxisConfig.TickComponent; nil uses the default RadialAxisTick.
type RadialTickRenderer = func(props RadialAxisTickProps) string

// CircularAxisTickProps mirrors @nivo/polar-axes CircularAxisTickProps (minus
// the react-spring animated wrapper; animation is handled via SMIL).
type CircularAxisTickProps struct {
	Label      string
	TextAnchor string
	Theme      theming.AxisTheme
	// Geometry (resolved from angle + radii by the parent).
	X1, Y1, X2, Y2 float64 // tick line endpoints
	TextX, TextY   float64 // label position
	// Animate, when true, signals the renderer to emit a fade-in <animate>.
	Animate bool
}

// RadialAxisTickProps mirrors @nivo/polar-axes RadialAxisTickProps.
type RadialAxisTickProps struct {
	Label      string
	TextAnchor string
	Theme      theming.AxisTheme
	// Geometry (in the axis's local rotated frame): y is the position along
	// the axis line, length is the tick line length, textX the label offset.
	Y        float64
	Length   float64
	TextX    float64
	Rotation float64
	Animate  bool
}

// PolarGridProps mirrors @nivo/polar-axes PolarGridProps. nivo derives the
// inner/outer radius from radiusScale.range(); our scales.Scale interface
// doesn't expose the range, so v1 callers pass InnerRadius/OuterRadius
// explicitly (typically 0 and the chart's outer radius).
type PolarGridProps struct {
	Center             [2]float64
	EnableRadialGrid   bool
	AngleScale         scales.Scale
	StartAngle         float64 // degrees
	EndAngle           float64 // degrees
	EnableCircularGrid bool
	RadiusScale        scales.Scale
	InnerRadius        float64
	OuterRadius        float64
	CircularGridTicks  scales.TicksSpec
	Theme              *theming.Theme
	Animate            bool
}

// RadialGridProps mirrors @nivo/polar-axes RadialGridProps.
type RadialGridProps struct {
	Scale       scales.Scale
	Ticks       scales.TicksSpec
	InnerRadius float64
	OuterRadius float64
	Theme       *theming.Theme
	Animate     bool
}

// CircularGridProps mirrors @nivo/polar-axes CircularGridProps. Angles are in
// degrees.
type CircularGridProps struct {
	Scale      scales.Scale
	Ticks      scales.TicksSpec
	StartAngle float64 // degrees
	EndAngle   float64 // degrees
	Theme      *theming.Theme
	Animate    bool
}

// CircularAxisProps mirrors @nivo/polar-axes CircularAxisProps (the component
// props, not the config).
type CircularAxisProps struct {
	Type       CircularAxisType
	Center     [2]float64
	Radius     float64
	StartAngle float64 // degrees
	EndAngle   float64 // degrees
	Scale      scales.Scale
	Theme      *theming.Theme
	Animate    bool
	CircularAxisConfig
}

// RadialAxisProps mirrors @nivo/polar-axes RadialAxisProps.
type RadialAxisProps struct {
	Center        [2]float64
	Angle         float64 // degrees
	Scale         scales.Scale
	TicksPosition TicksPosition
	Theme         *theming.Theme
	Animate       bool
	RadialAxisConfig
}
