package chord

import (
	"math"

	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

const halfPi = math.Pi / 2

// ribbonEpsilon matches d3's epsilon used to decide whether a pad angle is
// significant.
const ribbonEpsilon = 1e-12

// RibbonEndpoint is one end (source or target) of a ribbon: an angular span at
// the ribbon radius.
type RibbonEndpoint struct {
	StartAngle float64
	EndAngle   float64
}

// RibbonArg is the input to the ribbon generator — a source and target angular
// span. Mirrors the object @nivo/chord passes to d3.ribbon.
type RibbonArg struct {
	Source RibbonEndpoint
	Target RibbonEndpoint
}

// RibbonGenerator renders ribbon paths at a fixed radius. Mirrors d3-chord's
// ribbon() with a constant .radius and no headRadius/padAngle (the @nivo/chord
// configuration). Create via NewRibbon.
type RibbonGenerator struct {
	radius   float64
	padAngle float64
	digits   int
}

// NewRibbon returns a ribbon generator with the given radius, no pad angle, and
// 3-decimal coordinate rounding (matching d3-shape's default path precision).
func NewRibbon(radius float64) *RibbonGenerator {
	return &RibbonGenerator{radius: radius, digits: 3}
}

// Radius sets the (constant) source and target radius. Mirrors ribbon.radius.
func (r *RibbonGenerator) Radius(v float64) *RibbonGenerator { r.radius = v; return r }

// PadAngle sets the ribbon pad angle (radians). Mirrors ribbon.padAngle.
func (r *RibbonGenerator) PadAngle(v float64) *RibbonGenerator { r.padAngle = v; return r }

// Digits sets the coordinate rounding precision (decimals). digits < 0 disables
// rounding.
func (r *RibbonGenerator) Digits(d int) *RibbonGenerator { r.digits = d; return r }

// Call renders the ribbon SVG path-data string, exactly mirroring d3-chord's
// ribbon generator (moveTo → arc → quadraticCurveTo → arc → quadraticCurveTo →
// closePath). Coordinates are relative to the circle center at (0,0).
func (r *RibbonGenerator) Call(d RibbonArg) string {
	ap := r.padAngle / 2
	sr := r.radius
	tr := r.radius
	sa0 := d.Source.StartAngle - halfPi
	sa1 := d.Source.EndAngle - halfPi
	ta0 := d.Target.StartAngle - halfPi
	ta1 := d.Target.EndAngle - halfPi

	if ap > ribbonEpsilon {
		if math.Abs(sa1-sa0) > ap*2+ribbonEpsilon {
			if sa1 > sa0 {
				sa0 += ap
				sa1 -= ap
			} else {
				sa0 -= ap
				sa1 += ap
			}
		} else {
			sa0 = (sa0 + sa1) / 2
			sa1 = sa0
		}
		if math.Abs(ta1-ta0) > ap*2+ribbonEpsilon {
			if ta1 > ta0 {
				ta0 += ap
				ta1 -= ap
			} else {
				ta0 -= ap
				ta1 += ap
			}
		} else {
			ta0 = (ta0 + ta1) / 2
			ta1 = ta0
		}
	}

	p := d3shape.NewPathDigits(r.digits)
	p.MoveTo(sr*math.Cos(sa0), sr*math.Sin(sa0))
	p.Arc(0, 0, sr, sa0, sa1, false)
	if sa0 != ta0 || sa1 != ta1 {
		p.QuadraticCurveTo(0, 0, tr*math.Cos(ta0), tr*math.Sin(ta0))
		p.Arc(0, 0, tr, ta0, ta1, false)
	}
	p.QuadraticCurveTo(0, 0, sr*math.Cos(sa0), sr*math.Sin(sa0))
	p.ClosePath()
	return p.String()
}
