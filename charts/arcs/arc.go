// Package arcs provides the arc primitives shared by pie and future polar
// chart types: the Arc type, an ArcGenerator wrapping internal/d3/shape.Arc,
// bounding-box / center / link-label geometry helpers, and the ArcShape /
// ArcLabel / ArcLinkLabel templ components.
//
// Angle convention (matches nivo/d3): internal arc angles are in radians,
// measured from the +x axis. d3-shape's arc generator offsets by -π/2 so that
// angle 0 ends up at the top (12 o'clock) and angles increase clockwise. The
// helpers in this package take degrees (0 = top, clockwise) and convert.
package arcs

import (
	"fmt"
	"math"
	"strings"

	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// Arc is an annular sector in radians. StartAngle/EndAngle are measured from
// the +x axis (d3 convention); the arc generator applies the -π/2 offset.
type Arc struct {
	StartAngle  float64 // radians
	EndAngle    float64 // radians
	InnerRadius float64
	OuterRadius float64
	PadAngle    float64 // radians
}

// DatumWithArc is implemented by computed pie/bar data that carries an Arc.
type DatumWithArc interface {
	GetArc() Arc
}

// ArcGenerator wraps a configured d3-shape Arc generator. The zero value is
// not usable — create via CreateArcGenerator.
type ArcGenerator struct {
	inner *d3shape.Arc
}

// CreateArcGenerator builds an ArcGenerator with the given corner radius
// (constant across arcs) and pad angle (radians). Mirrors nivo's
// createArcGenerator.
func CreateArcGenerator(cornerRadius, padAngle float64) *ArcGenerator {
	a := d3shape.NewArc().
		CornerRadius(cornerRadius).
		PadAngle(padAngle)
	return &ArcGenerator{inner: a}
}

// GenerateSvgArc produces the SVG path-data string for one arc using the
// generator's configured cornerRadius/padAngle and the arc's own radii. The
// arc's padAngle is taken from the Arc datum (overriding the generator's pad
// for per-arc padding). Mirrors nivo's generateArcPath.
func (g *ArcGenerator) GenerateSvgArc(arc Arc) string {
	return g.inner.Call(d3shape.ArcDatum{
		StartAngle:  arc.StartAngle,
		EndAngle:    arc.EndAngle,
		InnerRadius: arc.InnerRadius,
		OuterRadius: arc.OuterRadius,
		PadAngle:    arc.PadAngle,
	})
}

// Centroid returns the [x, y] centroid of an arc (the midpoint of the angular
// span at the midpoint of the two radii), in the arc's local coordinate
// system (centered at the origin). Mirrors d3-shape arc.centroid.
func (g *ArcGenerator) Centroid(arc Arc) [2]float64 {
	return g.inner.Centroid(d3shape.ArcDatum{
		StartAngle:  arc.StartAngle,
		EndAngle:    arc.EndAngle,
		InnerRadius: arc.InnerRadius,
		OuterRadius: arc.OuterRadius,
		PadAngle:    arc.PadAngle,
	})
}

// DegToRad converts degrees to radians.
func DegToRad(deg float64) float64 { return deg * math.Pi / 180 }

// RadToDeg converts radians to degrees.
func RadToDeg(rad float64) float64 { return rad * 180 / math.Pi }

// Point is an (x, y) pair in SVG units.
type Point struct{ X, Y float64 }

// PositionFromAngle returns the cartesian point at the given angle (radians,
// measured from the +x axis) and distance from the origin. Mirrors nivo's
// core/lib/polar/utils.positionFromAngle.
func PositionFromAngle(angle, distance float64) Point {
	return Point{X: math.Cos(angle) * distance, Y: math.Sin(angle) * distance}
}

// NormalizeAngleDegrees maps an arbitrary angle (degrees) into [0, 360).
// Mirrors nivo's core/lib/polar/utils.normalizeAngleDegrees.
func NormalizeAngleDegrees(angle float64) float64 {
	a := math.Mod(angle, 360)
	if a < 0 {
		a += 360
	}
	return a
}

// GetNormalizedAngle maps an arbitrary angle (degrees, 0 = top, clockwise)
// into [0, 360).
func GetNormalizedAngle(angleDeg float64) float64 {
	a := math.Mod(angleDeg, 360)
	if a < 0 {
		a += 360
	}
	return a
}

// FilterDataBySkipAngle returns the subset of `data` whose arc angular span
// (degrees) is >= skipAngle. Used by pie's arc-label/arc-link-label layers to
// drop labels for tiny slices. data must be a slice of DatumWithArc.
func FilterDataBySkipAngle(data []DatumWithArc, skipAngleDeg float64) []DatumWithArc {
	if skipAngleDeg <= 0 {
		out := make([]DatumWithArc, len(data))
		copy(out, data)
		return out
	}
	out := make([]DatumWithArc, 0, len(data))
	for _, d := range data {
		span := RadToDeg(d.GetArc().EndAngle - d.GetArc().StartAngle)
		if span >= skipAngleDeg {
			out = append(out, d)
		}
	}
	return out
}

// ComputeArcCenter returns the (x, y) position of an arc's label anchor point,
// at radiusOffset * (outerRadius - innerRadius) + innerRadius from the center,
// offset by (centerX, centerY). Mirrors nivo's computeArcCenter
// (angle = midAngle - π/2).
func ComputeArcCenter(arc Arc, radiusOffset float64) (float64, float64) {
	r := arc.InnerRadius + (arc.OuterRadius-arc.InnerRadius)*radiusOffset
	midAngle := (arc.StartAngle + arc.EndAngle) / 2
	// d3-shape offsets by -π/2 so angle 0 is at the top.
	a := midAngle - math.Pi/2
	return math.Cos(a) * r, math.Sin(a) * r
}

// ComputeArcBoundingBox returns the {x, y, width, height} of the bounding box
// of the arc sector, in the arc's local coordinate system (centered at 0,0),
// given the start/end angles in DEGREES (0 = top, clockwise) and the outer
// radius. Mirrors @nivo/arcs computeArcBoundingBox.
//
// Algorithm: sample the two arc endpoints, plus any axis crossings (multiples
// of 90°) that fall within the arc's angular span; the bounding box is the
// extent of those points (and the origin if includeCenter is true).
func ComputeArcBoundingBox(centerX, centerY, radius, startAngleDeg, endAngleDeg float64, includeCenter bool) (float64, float64, float64, float64) {
	// Convert to d3 radians (angle 0 at +x axis, after subtracting 90°).
	toRad := func(deg float64) float64 { return DegToRad(deg - 90) }
	a0 := toRad(startAngleDeg)
	a1 := toRad(endAngleDeg)

	// Sample points: the two endpoints at radius r, plus axis crossings.
	xs := []float64{centerX + radius*math.Cos(a0), centerX + radius*math.Cos(a1)}
	ys := []float64{centerY + radius*math.Sin(a0), centerY + radius*math.Sin(a1)}

	if includeCenter {
		xs = append(xs, centerX)
		ys = append(ys, centerY)
	}

	// Collect the axis-crossing angles (in d3 radians) within [a0, a1].
	// The +x axis is at 0 d3-rad, +y at π/2, -x at π, -y at 3π/2 (or -π/2).
	crossings := []float64{0, math.Pi / 2, math.Pi, 3 * math.Pi / 2, -math.Pi / 2, -math.Pi, -3 * math.Pi / 2}
	lo, hi := a0, a1
	if lo > hi {
		lo, hi = hi, lo
	}
	for _, c := range crossings {
		// Normalize crossing into [lo, hi] by adding/subtracting 2π.
		for k := -2; k <= 2; k++ {
			ck := c + float64(k)*2*math.Pi
			if ck >= lo-epsilon && ck <= hi+epsilon {
				xs = append(xs, centerX+radius*math.Cos(ck))
				ys = append(ys, centerY+radius*math.Sin(ck))
			}
		}
	}

	minX, maxX := extent(xs)
	minY, maxY := extent(ys)
	return minX, minY, maxX - minX, maxY - minY
}

// ArcLink is the geometry of an arc's link label: a two-segment polyline from
// the arc edge (p0) through a diagonal bend (p1) to a horizontal straight
// segment ending at p2, plus the text-anchor side.
type ArcLink struct {
	P0x, P0y   float64 // start (on the arc's outer edge)
	P1x, P1y   float64 // diagonal bend end
	P2x, P2y   float64 // straight-segment end (where the text begins)
	TextAnchor string  // "start" | "end"
	Side       string  // "right" | "left"
}

// ComputeArcLink computes the link-label geometry for an arc. Mirrors nivo's
// computeArcLink. offset shifts the start point inward/outward along the arc's
// mid-angle; diagLength is the diagonal segment length; straightLength is the
// horizontal tail length.
func ComputeArcLink(arc Arc, offset, diagLength, straightLength float64) ArcLink {
	// Mid angle (d3 radians, with -π/2 offset applied).
	mid := (arc.StartAngle+arc.EndAngle)/2 - math.Pi/2
	r := arc.OuterRadius + offset
	p0x := math.Cos(mid) * r
	p0y := math.Sin(mid) * r

	// p1 is diagLength away from p0, continuing along the mid-angle direction
	// for the diagonal, then p2 extends horizontally by straightLength in the
	// same left/right direction.
	p1x := math.Cos(mid) * (r + diagLength)
	p1y := math.Sin(mid) * (r + diagLength)

	side := "right"
	textAnchor := "start"
	if p1x < 0 {
		side = "left"
		textAnchor = "end"
	}
	p2x := p1x + straightLength
	if side == "left" {
		p2x = p1x - straightLength
	}
	p2y := p1y
	return ArcLink{
		P0x: p0x, P0y: p0y,
		P1x: p1x, P1y: p1y,
		P2x: p2x, P2y: p2y,
		TextAnchor: textAnchor, Side: side,
	}
}

// ComputeArcLinkTextAnchor returns the text-anchor for an arc's link label,
// based on which side of the pie the arc falls on. Mirrors
// @nivo/arcs getArcLinkLabelAnchor.
func ComputeArcLinkTextAnchor(arc Arc) string {
	mid := (arc.StartAngle+arc.EndAngle)/2 - math.Pi/2
	if math.Cos(mid) < 0 {
		return "end"
	}
	return "start"
}

// FindArcUnderCursor returns the index of the arc under the (x, y) cursor
// position relative to the pie center, or -1 if none. Used by the htmx hover
// handler. Mirrors @nivo/arcs findArcUnderCursor.
func FindArcUnderCursor(arcs []Arc, cx, cy, x, y float64) int {
	dx := x - cx
	dy := y - cy
	dist := math.Sqrt(dx*dx + dy*dy)
	for i, a := range arcs {
		if dist < a.InnerRadius || dist > a.OuterRadius+epsilon {
			continue
		}
		// Angle of the cursor from the center, d3 convention (0 at +x axis).
		cursorAngle := math.Atan2(dy, dx) + math.Pi/2
		if cursorAngle < 0 {
			cursorAngle += 2 * math.Pi
		}
		lo := a.StartAngle
		hi := a.EndAngle
		// Normalize so lo <= hi.
		if lo > hi {
			lo, hi = hi, lo
		}
		// Handle wrap-around arcs (start > end after normalization).
		angle := cursorAngle
		if angle < lo {
			angle += 2 * math.Pi
		}
		if angle >= lo-epsilon && angle <= hi+epsilon {
			return i
		}
		// Also try the un-wrapped position.
		angle = cursorAngle
		if angle > hi {
			angle -= 2 * math.Pi
		}
		if angle >= lo-epsilon && angle <= hi+epsilon {
			return i
		}
	}
	return -1
}

const epsilon = 1e-9

func extent(xs []float64) (min, max float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	min, max = xs[0], xs[0]
	for _, x := range xs[1:] {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	return min, max
}

// fmtF formats a float for SVG output (3 dp, trimmed).
func fmtF(v float64) string {
	return strings.TrimSuffix(strings.TrimRight(fmt.Sprintf("%.3f", v), "0"), ".")
}
