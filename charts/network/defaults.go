package network

import "github.com/geoffjay/templ-charts/charts/core"

// Defaults mirrors @nivo/network svgDefaultProps. Fields left zero in a
// NetworkProps fall back to these via applyDefaults.
var Defaults = NetworkProps{
	LinkDistance:      30,
	CenteringStrength: 1,
	Repulsivity:       10,
	DistanceMin:       1,
	DistanceMax:       0, // 0 is treated as +Inf by the charge force
	Iterations:        120,

	NodeSize:        12,
	NodeColor:       "#000000",
	NodeBorderWidth: 0,

	LinkThickness: 1,

	Layers:      DefaultLayers,
	Role:        "img",
	MotionProps: core.MotionProps{Animate: true, MotionConfig: "gentle"},
}
