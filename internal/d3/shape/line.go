// line.go — port of d3-shape's line generator (src/line.js).
//
// Line() returns a generator that walks a slice of data, calls x/y accessors
// for each defined point, and feeds the resulting coordinates to a Curve,
// which emits drawing commands to a *Path. The generated SVG path-data
// string is returned via the path's String() method.
//
// The API mirrors d3-shape's fluent builders: X(), Y(), Defined(), Curve()
// configure accessors and the curve factory; Call(data) produces the string.
package d3shape

import "math"

// Point2D is the default data type: a [x, y] pair. Generators accept any
// data type via accessor functions.
type Point2D [2]float64

// accessor default for x: p[0]
func defaultX(p Point2D) float64 { return p[0] }
func defaultY(p Point2D) float64 { return p[1] }

// FloatAccessor extracts a float64 from a datum of type D.
type FloatAccessor[D any] func(d D, i int, data []D) float64

// BoolAccessor reports whether a datum is defined.
type BoolAccessor[D any] func(d D, i int, data []D) bool

// alwaysDefined is the default Defined accessor.
func alwaysDefined[D any](d D, i int, data []D) bool { return true }

// Line is a line generator parameterized by the data type D.
type Line[D any] struct {
	x       FloatAccessor[D]
	y       FloatAccessor[D]
	defined BoolAccessor[D]
	curve   CurveFactory
	digits  int
}

// NewLine constructs a Line with default x/y accessors (expecting Point2D
// data), Defined=true everywhere, and CurveLinear.
func NewLine() *Line[Point2D] {
	return &Line[Point2D]{
		x:       func(d Point2D, i int, data []Point2D) float64 { return d[0] },
		y:       func(d Point2D, i int, data []Point2D) float64 { return d[1] },
		defined: alwaysDefined[Point2D],
		curve:   CurveLinear,
		digits:  3,
	}
}

// NewLineTyped constructs a Line for a custom data type D with the given
// x/y accessors.
func NewLineTyped[D any](x FloatAccessor[D], y FloatAccessor[D]) *Line[D] {
	return &Line[D]{
		x:       x,
		y:       y,
		defined: alwaysDefined[D],
		curve:   CurveLinear,
		digits:  3,
	}
}

// X sets the x accessor.
func (l *Line[D]) X(fn FloatAccessor[D]) *Line[D] { l.x = fn; return l }

// Y sets the y accessor.
func (l *Line[D]) Y(fn FloatAccessor[D]) *Line[D] { l.y = fn; return l }

// Defined sets the defined predicate.
func (l *Line[D]) Defined(fn BoolAccessor[D]) *Line[D] { l.defined = fn; return l }

// Curve sets the curve factory.
func (l *Line[D]) Curve(c CurveFactory) *Line[D] { l.curve = c; return l }

// Digits sets the path rounding precision (default 3).
func (l *Line[D]) Digits(d int) *Line[D] { l.digits = d; return l }

// Call generates the SVG path-data string for the given data. Returns "" for
// empty input (d3 returns null; empty string is the Go equivalent for "no path").
func (l *Line[D]) Call(data []D) string {
	n := len(data)
	if n == 0 {
		return ""
	}
	path := NewPathDigits(l.digits)
	curve := l.curve(path)
	defined0 := false
	for i := 0; i <= n; i++ {
		// d3: !(i < n && defined(d, i, data)) === defined0
		// When i == n, the condition is !(false) === defined0, i.e. true === defined0,
		// which toggles when defined0 is false (closing any open segment).
		var isDefined bool
		if i < n {
			isDefined = l.defined(data[i], i, data)
		}
		within := i < n && isDefined
		if !within == defined0 {
			defined0 = !defined0
			if defined0 {
				curve.LineStart()
			} else {
				curve.LineEnd()
			}
		}
		if defined0 {
			curve.Point(l.x(data[i], i, data), l.y(data[i], i, data))
		}
	}
	return path.String()
}

// Centroid computes the centroid of the points (d3-shape line.centroid).
// Not used by nivo but included for completeness.
func (l *Line[D]) Centroid(data []D) [2]float64 {
	var sx, sy float64
	var count int
	for i, d := range data {
		if l.defined(d, i, data) {
			sx += l.x(d, i, data)
			sy += l.y(d, i, data)
			count++
		}
	}
	if count == 0 {
		return [2]float64{math.NaN(), math.NaN()}
	}
	return [2]float64{sx / float64(count), sy / float64(count)}
}
