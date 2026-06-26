// curve.go — port of d3-shape's curve factories.
//
// A Curve receives a stream of point/lineStart/lineEnd/areaStart/areaEnd
// commands from a generator (Line/Area) and emits drawing commands to a
// Path (the "context"). Each curve factory is a function CurveFactory that
// takes a *Path and returns a Curve.
//
// Ports: linear, linearClosed, step/stepBefore/stepAfter, monotoneX/Y,
// basis, cardinal (tension 0), catmullRom (alpha 0.5), natural. The closed
// /open variants of basis/cardinal/catmullRom are not needed by nivo's line
// chart (which uses non-closed curves) and are omitted; they can be added
// later if a future chart type requires them.
package d3shape

import "math"

// Curve is the interface d3-shape generators drive. Methods mirror d3-shape's
// curve output protocol: areaStart/areaEnd bracket an area's top+bottom
// lines; lineStart/lineEnd bracket a single polyline; point adds a vertex.
type Curve interface {
	AreaStart()
	AreaEnd()
	LineStart()
	LineEnd()
	Point(x, y float64)
}

// CurveFactory constructs a Curve bound to a *Path context.
type CurveFactory func(*Path) Curve

// jsTruthy returns true if v would be truthy in JavaScript (0 and NaN are
// falsy; all other floats including negative and >0 are truthy).
func jsTruthy(v float64) bool {
	return !math.IsNaN(v) && v != 0
}

// --- linear ---------------------------------------------------------------

type linearCurve struct {
	ctx   *Path
	line  float64 // 0, 1, or NaN; d3 uses NaN as sentinel
	point int
}

func (c *linearCurve) AreaStart() { c.line = 0 }
func (c *linearCurve) AreaEnd()   { c.line = math.NaN() }
func (c *linearCurve) LineStart() { c.point = 0 }
func (c *linearCurve) LineEnd() {
	// d3: if (this._line || (this._line !== 0 && this._point === 1)) closePath()
	if jsTruthy(c.line) || (!math.IsNaN(c.line) && c.line != 0 && c.point == 1) {
		c.ctx.closePath()
	}
	c.line = 1 - c.line
}

func (c *linearCurve) Point(x, y float64) {
	switch c.point {
	case 0:
		c.point = 1
		if jsTruthy(c.line) {
			c.ctx.lineTo(x, y)
		} else {
			c.ctx.moveTo(x, y)
		}
	case 1:
		c.point = 2
		fallthrough
	default:
		c.ctx.lineTo(x, y)
	}
}

// CurveLinear is the linear curve factory (d3-shape curve/linear.js).
func CurveLinear(ctx *Path) Curve { return &linearCurve{ctx: ctx, line: math.NaN()} }

// --- linearClosed ---------------------------------------------------------

type linearClosedCurve struct {
	ctx   *Path
	point int
}

func (c *linearClosedCurve) AreaStart() {} // noop
func (c *linearClosedCurve) AreaEnd()   {} // noop
func (c *linearClosedCurve) LineStart() { c.point = 0 }

func (c *linearClosedCurve) LineEnd() {
	if c.point != 0 {
		c.ctx.closePath()
	}
}

func (c *linearClosedCurve) Point(x, y float64) {
	if c.point != 0 {
		c.ctx.lineTo(x, y)
	} else {
		c.point = 1
		c.ctx.moveTo(x, y)
	}
}

// CurveLinearClosed is the linearClosed curve factory.
func CurveLinearClosed(ctx *Path) Curve { return &linearClosedCurve{ctx: ctx} }

// --- step -----------------------------------------------------------------

type stepCurve struct {
	ctx   *Path
	line  float64
	point int
	x, y  float64
	t     float64
}

func (c *stepCurve) AreaStart() { c.line = 0 }
func (c *stepCurve) AreaEnd()   { c.line = math.NaN() }
func (c *stepCurve) LineStart() {
	c.point = 0
	c.x = math.NaN()
	c.y = math.NaN()
}

func (c *stepCurve) LineEnd() {
	if c.t > 0 && c.t < 1 && c.point == 2 {
		c.ctx.lineTo(c.x, c.y)
	}
	if jsTruthy(c.line) || (!math.IsNaN(c.line) && c.line != 0 && c.point == 1) {
		c.ctx.closePath()
	}
	if !math.IsNaN(c.line) && c.line >= 0 {
		c.t = 1 - c.t
		c.line = 1 - c.line
	}
}

func (c *stepCurve) Point(x, y float64) {
	switch c.point {
	case 0:
		c.point = 1
		if jsTruthy(c.line) {
			c.ctx.lineTo(x, y)
		} else {
			c.ctx.moveTo(x, y)
		}
	case 1:
		c.point = 2
		fallthrough
	default:
		if c.t <= 0 {
			c.ctx.lineTo(c.x, y)
			c.ctx.lineTo(x, y)
		} else {
			x1 := c.x*(1-c.t) + x*c.t
			c.ctx.lineTo(x1, c.y)
			c.ctx.lineTo(x1, y)
		}
	}
	c.x = x
	c.y = y
}

// CurveStep is the step curve factory (t=0.5, d3-shape curve/step.js default).
func CurveStep(ctx *Path) Curve { return &stepCurve{ctx: ctx, line: math.NaN(), t: 0.5} }

// CurveStepBefore is the stepBefore curve factory (t=0).
func CurveStepBefore(ctx *Path) Curve { return &stepCurve{ctx: ctx, line: math.NaN(), t: 0} }

// CurveStepAfter is the stepAfter curve factory (t=1).
func CurveStepAfter(ctx *Path) Curve { return &stepCurve{ctx: ctx, line: math.NaN(), t: 1} }

// --- monotone (X and Y) ---------------------------------------------------

// monotoneCurve implements the Steffen 1990 monotonic cubic interpolation
// from d3-shape's curve/monotone.js. MonotoneY wraps MonotoneX with a
// reflecting context that swaps x/y.
type monotoneCurve struct {
	ctx            *Path
	line           float64
	point          int
	x0, y0, x1, y1 float64
	t0             float64
	reflect        bool // true for MonotoneY
}

func (c *monotoneCurve) AreaStart() { c.line = 0 }
func (c *monotoneCurve) AreaEnd()   { c.line = math.NaN() }
func (c *monotoneCurve) LineStart() {
	c.point = 0
	c.x0 = math.NaN()
	c.y0 = math.NaN()
	c.x1 = math.NaN()
	c.y1 = math.NaN()
	c.t0 = math.NaN()
}

func (c *monotoneCurve) LineEnd() {
	switch c.point {
	case 2:
		c.lineToCtx(c.x1, c.y1)
	case 3:
		c.monotonePoint(c.t0, monotoneSlope2(c, c.t0))
	}
	if jsTruthy(c.line) || (!math.IsNaN(c.line) && c.line != 0 && c.point == 1) {
		c.ctx.closePath()
	}
	c.line = 1 - c.line
}

func (c *monotoneCurve) Point(x, y float64) {
	t1 := math.NaN()
	if x == c.x1 && y == c.y1 {
		return // ignore coincident points
	}
	switch c.point {
	case 0:
		c.point = 1
		if jsTruthy(c.line) {
			c.lineToCtx(x, y)
		} else {
			c.moveToCtx(x, y)
		}
	case 1:
		c.point = 2
	case 2:
		c.point = 3
		t1 = monotoneSlope3OrZero(c, x, y)
		c.monotonePoint(monotoneSlope2(c, t1), t1)
	default:
		t1 = monotoneSlope3OrZero(c, x, y)
		c.monotonePoint(c.t0, t1)
	}
	c.x0 = c.x1
	c.x1 = x
	c.y0 = c.y1
	c.y1 = y
	c.t0 = t1
}

// moveToCtx/lineToCtx/bezierCtx account for the reflect flag (MonotoneY).
func (c *monotoneCurve) moveToCtx(x, y float64) {
	if c.reflect {
		c.ctx.moveTo(y, x)
	} else {
		c.ctx.moveTo(x, y)
	}
}

func (c *monotoneCurve) lineToCtx(x, y float64) {
	if c.reflect {
		c.ctx.lineTo(y, x)
	} else {
		c.ctx.lineTo(x, y)
	}
}

func (c *monotoneCurve) bezierCtx(x1, y1, x2, y2, x, y float64) {
	if c.reflect {
		c.ctx.bezierCurveTo(y1, x1, y2, x2, y, x)
	} else {
		c.ctx.bezierCurveTo(x1, y1, x2, y2, x, y)
	}
}

// monotonePoint emits a cubic bezier segment per the Hermite representation
// (d3-shape monotone.js point()).
func (c *monotoneCurve) monotonePoint(t0, t1 float64) {
	x0 := c.x0
	y0 := c.y0
	x1 := c.x1
	y1 := c.y1
	dx := (x1 - x0) / 3
	c.bezierCtx(x0+dx, y0+dx*t0, x1-dx, y1-dx*t1, x1, y1)
}

func sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	return 1
}

// monotoneSlope3 computes the two-sided slope (d3-shape slope3).
// d3: s0 = (y1-y0) / (h0 || h1<0 && -0). In JS, h0||X evaluates X when h0
// is falsy (0 or NaN).
func monotoneSlope3(c *monotoneCurve, x2, y2 float64) float64 {
	h0 := c.x1 - c.x0
	h1 := x2 - c.x1
	// d3: s0 = (y1-y0) / (h0 || (h1 < 0 ? -0 : false))
	var s0Denom float64
	if jsTruthy(h0) {
		s0Denom = h0
	} else if h1 < 0 {
		s0Denom = math.Copysign(0, -1) // -0
	} else {
		s0Denom = 0 // NaN would result; use 0 to produce Inf/NaN as d3 does
	}
	s0 := (c.y1 - c.y0) / s0Denom

	var s1Denom float64
	if jsTruthy(h1) {
		s1Denom = h1
	} else if h0 < 0 {
		s1Denom = math.Copysign(0, -1)
	} else {
		s1Denom = 0
	}
	s1 := (y2 - c.y1) / s1Denom

	var p float64
	if !math.IsNaN(h0+h1) && h0+h1 != 0 {
		p = (s0*h1 + s1*h0) / (h0 + h1)
	}
	return (sign(s0) + sign(s1)) * math.Min(math.Min(math.Abs(s0), math.Abs(s1)), 0.5*math.Abs(p))
}

// monotoneSlope3OrZero wraps monotoneSlope3 to return 0 when the result is
// NaN (d3's `|| 0` fallback).
func monotoneSlope3OrZero(c *monotoneCurve, x2, y2 float64) float64 {
	v := monotoneSlope3(c, x2, y2)
	if math.IsNaN(v) {
		return 0
	}
	return v
}

// monotoneSlope2 computes the one-sided slope (d3-shape slope2).
func monotoneSlope2(c *monotoneCurve, t float64) float64 {
	h := c.x1 - c.x0
	if !math.IsNaN(h) && h != 0 {
		return (3*(c.y1-c.y0)/h - t) / 2
	}
	return t
}

// CurveMonotoneX is the monotoneX curve factory.
func CurveMonotoneX(ctx *Path) Curve {
	return &monotoneCurve{ctx: ctx, line: math.NaN(), reflect: false}
}

// CurveMonotoneY is the monotoneY curve factory (reflects x/y).
func CurveMonotoneY(ctx *Path) Curve {
	return &monotoneCurve{ctx: ctx, line: math.NaN(), reflect: true}
}

// --- basis ----------------------------------------------------------------

type basisCurve struct {
	ctx            *Path
	line           float64
	point          int
	x0, y0, x1, y1 float64
}

func (c *basisCurve) AreaStart() { c.line = 0 }
func (c *basisCurve) AreaEnd()   { c.line = math.NaN() }
func (c *basisCurve) LineStart() {
	c.point = 0
	c.x0 = math.NaN()
	c.y0 = math.NaN()
	c.x1 = math.NaN()
	c.y1 = math.NaN()
}

func (c *basisCurve) LineEnd() {
	switch c.point {
	case 3:
		basisPoint(c, c.x1, c.y1)
		fallthrough
	case 2:
		c.ctx.lineTo(c.x1, c.y1)
	}
	if jsTruthy(c.line) || (!math.IsNaN(c.line) && c.line != 0 && c.point == 1) {
		c.ctx.closePath()
	}
	c.line = 1 - c.line
}

func (c *basisCurve) Point(x, y float64) {
	switch c.point {
	case 0:
		c.point = 1
		if jsTruthy(c.line) {
			c.ctx.lineTo(x, y)
		} else {
			c.ctx.moveTo(x, y)
		}
	case 1:
		c.point = 2
	case 2:
		c.point = 3
		c.ctx.lineTo((5*c.x0+c.x1)/6, (5*c.y0+c.y1)/6)
		fallthrough
	default:
		basisPoint(c, x, y)
	}
	c.x0 = c.x1
	c.x1 = x
	c.y0 = c.y1
	c.y1 = y
}

// basisPoint emits the cubic bezier for the basis spline (d3-shape basis.js point()).
func basisPoint(c *basisCurve, x, y float64) {
	c.ctx.bezierCurveTo(
		(2*c.x0+c.x1)/3, (2*c.y0+c.y1)/3,
		(c.x0+2*c.x1)/3, (c.y0+2*c.y1)/3,
		(c.x0+4*c.x1+x)/6, (c.y0+4*c.y1+y)/6,
	)
}

// CurveBasis is the basis curve factory.
func CurveBasis(ctx *Path) Curve { return &basisCurve{ctx: ctx, line: math.NaN()} }

// --- cardinal -------------------------------------------------------------

type cardinalCurve struct {
	ctx                    *Path
	line                   float64
	point                  int
	x0, y0, x1, y1, x2, y2 float64
	k                      float64 // (1 - tension) / 6
}

func (c *cardinalCurve) AreaStart() { c.line = 0 }
func (c *cardinalCurve) AreaEnd()   { c.line = math.NaN() }
func (c *cardinalCurve) LineStart() {
	c.point = 0
	c.x0 = math.NaN()
	c.y0 = math.NaN()
	c.x1 = math.NaN()
	c.y1 = math.NaN()
	c.x2 = math.NaN()
	c.y2 = math.NaN()
}

func (c *cardinalCurve) LineEnd() {
	switch c.point {
	case 2:
		c.ctx.lineTo(c.x2, c.y2)
	case 3:
		cardinalPoint(c, c.x1, c.y1)
	}
	if jsTruthy(c.line) || (!math.IsNaN(c.line) && c.line != 0 && c.point == 1) {
		c.ctx.closePath()
	}
	c.line = 1 - c.line
}

func (c *cardinalCurve) Point(x, y float64) {
	switch c.point {
	case 0:
		c.point = 1
		if jsTruthy(c.line) {
			c.ctx.lineTo(x, y)
		} else {
			c.ctx.moveTo(x, y)
		}
	case 1:
		c.point = 2
		c.x1 = x
		c.y1 = y
	case 2:
		c.point = 3
		fallthrough
	default:
		cardinalPoint(c, x, y)
	}
	c.x0 = c.x1
	c.x1 = c.x2
	c.x2 = x
	c.y0 = c.y1
	c.y1 = c.y2
	c.y2 = y
}

// cardinalPoint emits the cubic bezier for the cardinal spline.
func cardinalPoint(c *cardinalCurve, x, y float64) {
	c.ctx.bezierCurveTo(
		c.x1+c.k*(c.x2-c.x0), c.y1+c.k*(c.y2-c.y0),
		c.x2+c.k*(c.x1-x), c.y2+c.k*(c.y1-y),
		c.x2, c.y2,
	)
}

// CurveCardinal is the cardinal curve factory with tension 0 (d3 default).
func CurveCardinal(ctx *Path) Curve {
	return &cardinalCurve{ctx: ctx, line: math.NaN(), k: (1.0 - 0) / 6}
}

// CurveCardinalTension returns a cardinal curve factory with custom tension.
func CurveCardinalTension(tension float64) CurveFactory {
	k := (1.0 - tension) / 6
	return func(ctx *Path) Curve {
		return &cardinalCurve{ctx: ctx, line: math.NaN(), k: k}
	}
}

// --- catmullRom -----------------------------------------------------------

// catmullRomCurve is the centripetal Catmull-Rom spline (d3-shape
// curve/catmullRom.js). alpha defaults to 0.5; alpha=0 degenerates to
// cardinal (tension 0).
type catmullRomCurve struct {
	ctx                    *Path
	line                   float64
	point                  int
	x0, y0, x1, y1, x2, y2 float64
	l01_a, l12_a, l23_a    float64
	l01_2a, l12_2a, l23_2a float64
	alpha                  float64
	// when alpha == 0 we delegate to the cardinal implementation
	cardinal *cardinalCurve
}

func (c *catmullRomCurve) AreaStart() {
	if c.cardinal != nil {
		c.cardinal.AreaStart()
		return
	}
	c.line = 0
}

func (c *catmullRomCurve) AreaEnd() {
	if c.cardinal != nil {
		c.cardinal.AreaEnd()
		return
	}
	c.line = math.NaN()
}

func (c *catmullRomCurve) LineStart() {
	if c.cardinal != nil {
		c.cardinal.LineStart()
		return
	}
	c.point = 0
	c.x0 = math.NaN()
	c.y0 = math.NaN()
	c.x1 = math.NaN()
	c.y1 = math.NaN()
	c.x2 = math.NaN()
	c.y2 = math.NaN()
	c.l01_a = 0
	c.l12_a = 0
	c.l23_a = 0
	c.l01_2a = 0
	c.l12_2a = 0
	c.l23_2a = 0
}

func (c *catmullRomCurve) LineEnd() {
	if c.cardinal != nil {
		c.cardinal.LineEnd()
		return
	}
	switch c.point {
	case 2:
		c.ctx.lineTo(c.x2, c.y2)
	case 3:
		c.catmullPoint(c.x2, c.y2)
	}
	if jsTruthy(c.line) || (!math.IsNaN(c.line) && c.line != 0 && c.point == 1) {
		c.ctx.closePath()
	}
	c.line = 1 - c.line
}

func (c *catmullRomCurve) Point(x, y float64) {
	if c.cardinal != nil {
		c.cardinal.Point(x, y)
		return
	}
	if c.point != 0 {
		x23 := c.x2 - x
		y23 := c.y2 - y
		c.l23_2a = math.Pow(x23*x23+y23*y23, c.alpha)
		c.l23_a = math.Sqrt(c.l23_2a)
	}
	switch c.point {
	case 0:
		c.point = 1
		if jsTruthy(c.line) {
			c.ctx.lineTo(x, y)
		} else {
			c.ctx.moveTo(x, y)
		}
	case 1:
		c.point = 2
	case 2:
		c.point = 3
		fallthrough
	default:
		c.catmullPoint(x, y)
	}
	c.l01_a = c.l12_a
	c.l12_a = c.l23_a
	c.l01_2a = c.l12_2a
	c.l12_2a = c.l23_2a
	c.x0 = c.x1
	c.x1 = c.x2
	c.x2 = x
	c.y0 = c.y1
	c.y1 = c.y2
	c.y2 = y
}

// catmullPoint emits the cubic bezier for the catmull-rom spline.
func (c *catmullRomCurve) catmullPoint(x, y float64) {
	x1 := c.x1
	y1 := c.y1
	x2 := c.x2
	y2 := c.y2
	if c.l01_a > epsilon {
		a := 2*c.l01_2a + 3*c.l01_a*c.l12_a + c.l12_2a
		n := 3 * c.l01_a * (c.l01_a + c.l12_a)
		x1 = (x1*a - c.x0*c.l12_2a + c.x2*c.l01_2a) / n
		y1 = (y1*a - c.y0*c.l12_2a + c.y2*c.l01_2a) / n
	}
	if c.l23_a > epsilon {
		b := 2*c.l23_2a + 3*c.l23_a*c.l12_a + c.l12_2a
		m := 3 * c.l23_a * (c.l23_a + c.l12_a)
		x2 = (x2*b + c.x1*c.l23_2a - x*c.l12_2a) / m
		y2 = (y2*b + c.y1*c.l23_2a - y*c.l12_2a) / m
	}
	c.ctx.bezierCurveTo(x1, y1, x2, y2, c.x2, c.y2)
}

// CurveCatmullRom is the catmullRom curve factory with alpha 0.5 (d3 default).
// alpha=0 delegates to Cardinal(tension 0).
func CurveCatmullRom(ctx *Path) Curve {
	return &catmullRomCurve{ctx: ctx, line: math.NaN(), alpha: 0.5}
}

// CurveCatmullRomAlpha returns a catmullRom curve factory with custom alpha.
func CurveCatmullRomAlpha(alpha float64) CurveFactory {
	if alpha == 0 {
		return CurveCardinal
	}
	return func(ctx *Path) Curve {
		return &catmullRomCurve{ctx: ctx, line: math.NaN(), alpha: alpha}
	}
}

// --- natural --------------------------------------------------------------

type naturalCurve struct {
	ctx  *Path
	line float64
	xs   []float64
	ys   []float64
}

func (c *naturalCurve) AreaStart() { c.line = 0 }
func (c *naturalCurve) AreaEnd()   { c.line = math.NaN() }
func (c *naturalCurve) LineStart() {
	c.xs = nil
	c.ys = nil
}

func (c *naturalCurve) LineEnd() {
	x := c.xs
	y := c.ys
	n := len(x)
	if n > 0 {
		if jsTruthy(c.line) {
			c.ctx.lineTo(x[0], y[0])
		} else {
			c.ctx.moveTo(x[0], y[0])
		}
		if n == 2 {
			c.ctx.lineTo(x[1], y[1])
		} else if n > 2 {
			px := naturalControlPoints(x)
			py := naturalControlPoints(y)
			for i0, i1 := 0, 1; i1 < n; i0, i1 = i0+1, i1+1 {
				c.ctx.bezierCurveTo(px[0][i0], py[0][i0], px[1][i0], py[1][i0], x[i1], y[i1])
			}
		}
	}
	if jsTruthy(c.line) || (!math.IsNaN(c.line) && c.line != 0 && n == 1) {
		c.ctx.closePath()
	}
	c.line = 1 - c.line
	c.xs = nil
	c.ys = nil
}

func (c *naturalCurve) Point(x, y float64) {
	c.xs = append(c.xs, x)
	c.ys = append(c.ys, y)
}

// naturalControlPoints computes the bezier control points for a natural
// cubic spline through the given points (d3-shape curve/natural.js
// controlPoints). Returns [a, b] where a[i] and b[i] are the first and
// second control points for the segment from x[i] to x[i+1].
func naturalControlPoints(x []float64) [2][]float64 {
	n := len(x) - 1
	if n < 1 {
		return [2][]float64{{}, {}}
	}
	a := make([]float64, n)
	b := make([]float64, n)
	r := make([]float64, n)
	a[0] = 0
	b[0] = 2
	r[0] = x[0] + 2*x[1]
	for i := 1; i < n-1; i++ {
		a[i] = 1
		b[i] = 4
		r[i] = 4*x[i] + 2*x[i+1]
	}
	if n >= 2 {
		a[n-1] = 2
		b[n-1] = 7
		r[n-1] = 8*x[n-1] + x[n]
	}
	for i := 1; i < n; i++ {
		m := a[i] / b[i-1]
		b[i] -= m
		r[i] -= m * r[i-1]
	}
	a[n-1] = r[n-1] / b[n-1]
	for i := n - 2; i >= 0; i-- {
		a[i] = (r[i] - a[i+1]) / b[i]
	}
	b[n-1] = (x[n] + a[n-1]) / 2
	for i := 0; i < n-1; i++ {
		b[i] = 2*x[i+1] - a[i+1]
	}
	return [2][]float64{a, b}
}

// CurveNatural is the natural curve factory.
func CurveNatural(ctx *Path) Curve { return &naturalCurve{ctx: ctx, line: math.NaN()} }
