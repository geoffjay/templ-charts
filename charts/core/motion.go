package core

// MotionProps mirrors nivo's motion config: the `animate` boolean gates SMIL
// `<animate>` enter animations, and MotionConfig names a preset ("gentle",
// "wobbly", "stiff" in nivo). templ-charts uses fixed 600ms ease enter
// animations and treats MotionConfig as an opaque label.
type MotionProps struct {
	Animate      bool
	MotionConfig string
}

// DefaultAnimate is nivo's defaultAnimate = true.
var DefaultAnimate = true

// DefaultMotionConfig is nivo's default motion preset label.
const DefaultMotionConfig = "gentle"
