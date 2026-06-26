// area.go — port of d3-shape's area generator (src/area.js).
//
// Area() returns a generator that walks a slice of data, calls x/y0/y1
// accessors for each defined point, and feeds the resulting coordinates to a
// Curve. The area consists of a "top" line (y1 values) followed by a reversed
// "bottom" line (y0 values), closed into a single subpath. The generated SVG
// path-data string is returned via the path's String() method.
//
// The API mirrors d3-shape's fluent builders: X(), X0(), X1(), Y(), Y0(), Y1(),
// Defined(), Curve() configure accessors and the curve factory; Call(data)
// produces the string.
package d3shape

// Area is an area generator parameterized by the data type D.
type Area[D any] struct {
	x0, x1  FloatAccessor[D]
	y0, y1  FloatAccessor[D]
	defined BoolAccessor[D]
	curve   CurveFactory
	digits  int
}

// NewArea constructs an Area with default accessors (expecting Point2D data):
// x = p[0], y0 = constant 0, y1 = p[1]. Defined=true everywhere, CurveLinear.
func NewArea() *Area[Point2D] {
	return &Area[Point2D]{
		x0:      func(d Point2D, _ int, _ []Point2D) float64 { return d[0] },
		x1:      nil, // nil means use x0
		y0:      func(_ Point2D, _ int, _ []Point2D) float64 { return 0 },
		y1:      func(d Point2D, _ int, _ []Point2D) float64 { return d[1] },
		defined: alwaysDefined[Point2D],
		curve:   CurveLinear,
		digits:  3,
	}
}

// NewAreaTyped constructs an Area for a custom data type D with the given
// x0, y0, y1 accessors.
func NewAreaTyped[D any](x0, y0, y1 FloatAccessor[D]) *Area[D] {
	return &Area[D]{
		x0:      x0,
		x1:      nil,
		y0:      y0,
		y1:      y1,
		defined: alwaysDefined[D],
		curve:   CurveLinear,
		digits:  3,
	}
}

// X sets the x accessor (clears x1, so x0=x).
func (a *Area[D]) X(fn FloatAccessor[D]) *Area[D] { a.x0 = fn; a.x1 = nil; return a }

// X0 sets the x0 accessor.
func (a *Area[D]) X0(fn FloatAccessor[D]) *Area[D] { a.x0 = fn; return a }

// X1 sets the x1 accessor. nil means use x0.
func (a *Area[D]) X1(fn FloatAccessor[D]) *Area[D] { a.x1 = fn; return a }

// Y sets the y accessor (clears y1, so y0=y).
func (a *Area[D]) Y(fn FloatAccessor[D]) *Area[D] { a.y0 = fn; a.y1 = nil; return a }

// Y0 sets the y0 accessor (the area's baseline).
func (a *Area[D]) Y0(fn FloatAccessor[D]) *Area[D] { a.y0 = fn; return a }

// Y1 sets the y1 accessor (the area's top line).
func (a *Area[D]) Y1(fn FloatAccessor[D]) *Area[D] { a.y1 = fn; return a }

// Defined sets the defined predicate.
func (a *Area[D]) Defined(fn BoolAccessor[D]) *Area[D] { a.defined = fn; return a }

// Curve sets the curve factory.
func (a *Area[D]) Curve(c CurveFactory) *Area[D] { a.curve = c; return a }

// Digits sets the path rounding precision (default 3).
func (a *Area[D]) Digits(d int) *Area[D] { a.digits = d; return a }

// Call generates the SVG path-data string for the given data. Returns "" for
// empty input.
func (a *Area[D]) Call(data []D) string {
	n := len(data)
	if n == 0 {
		return ""
	}
	path := NewPathDigits(a.digits)
	curve := a.curve(path)
	defined0 := false
	// x0z/y0z buffer the bottom-line coordinates for replay in reverse.
	x0z := make([]float64, n)
	y0z := make([]float64, n)
	j := 0 // start index of the current defined segment

	for i := 0; i <= n; i++ {
		var isDefined bool
		if i < n {
			isDefined = a.defined(data[i], i, data)
		}
		within := i < n && isDefined
		if !within == defined0 {
			defined0 = !defined0
			if defined0 {
				// entering a defined segment
				j = i
				curve.AreaStart()
				curve.LineStart()
			} else {
				// leaving a defined segment: emit the bottom line in reverse
				curve.LineEnd()
				curve.LineStart()
				for k := i - 1; k >= j; k-- {
					curve.Point(x0z[k], y0z[k])
				}
				curve.LineEnd()
				curve.AreaEnd()
			}
		}
		if defined0 {
			x0z[i] = a.x0(data[i], i, data)
			y0z[i] = a.y0(data[i], i, data)
			var px, py float64
			if a.x1 != nil {
				px = a.x1(data[i], i, data)
			} else {
				px = x0z[i]
			}
			if a.y1 != nil {
				py = a.y1(data[i], i, data)
			} else {
				py = y0z[i]
			}
			curve.Point(px, py)
		}
	}
	return path.String()
}
