package geo

import "math"

// clipRectangle is a screen-space rectangular postclip: after projection, it
// crops geometry to the axis-aligned box [x0,y0,x1,y1], cutting lines at the
// edges and rejoining clipped polygon rings along the boundary. Ported from
// d3-geo src/clip/rectangle.js (with the parametric line clip from
// src/clip/line.js). Unlike the preclips it operates in plane coordinates, so it
// is installed as a postclip wrapping the sink.

const (
	clipMax = 1e9
	clipMin = -clipMax
)

// clipLineRect clips segment a→b to the rectangle, mutating a and b to the
// clipped endpoints. Returns false if the segment lies entirely outside. Ported
// from d3-geo src/clip/line.js.
func clipLineRect(a, b *[2]float64, x0, y0, x1, y1 float64) bool {
	ax, ay := a[0], a[1]
	bx, by := b[0], b[1]
	t0, t1 := 0.0, 1.0
	dx := bx - ax
	dy := by - ay

	r := x0 - ax
	if dx == 0 && r > 0 {
		return false
	}
	r /= dx
	if dx < 0 {
		if r < t0 {
			return false
		}
		if r < t1 {
			t1 = r
		}
	} else if dx > 0 {
		if r > t1 {
			return false
		}
		if r > t0 {
			t0 = r
		}
	}

	r = x1 - ax
	if dx == 0 && r < 0 {
		return false
	}
	r /= dx
	if dx < 0 {
		if r > t1 {
			return false
		}
		if r > t0 {
			t0 = r
		}
	} else if dx > 0 {
		if r < t0 {
			return false
		}
		if r < t1 {
			t1 = r
		}
	}

	r = y0 - ay
	if dy == 0 && r > 0 {
		return false
	}
	r /= dy
	if dy < 0 {
		if r < t0 {
			return false
		}
		if r < t1 {
			t1 = r
		}
	} else if dy > 0 {
		if r > t1 {
			return false
		}
		if r > t0 {
			t0 = r
		}
	}

	r = y1 - ay
	if dy == 0 && r < 0 {
		return false
	}
	r /= dy
	if dy < 0 {
		if r > t1 {
			return false
		}
		if r > t0 {
			t0 = r
		}
	} else if dy > 0 {
		if r < t0 {
			return false
		}
		if r < t1 {
			t1 = r
		}
	}

	if t0 > 0 {
		a[0] = ax + t0*dx
		a[1] = ay + t0*dy
	}
	if t1 < 1 {
		b[0] = ax + t1*dx
		b[1] = ay + t1*dy
	}
	return true
}

// clipRectangle returns a postclip stream factory for the given box.
func clipRectangle(x0, y0, x1, y1 float64) func(Sink) Sink {
	visible := func(x, y float64) bool {
		return x0 <= x && x <= x1 && y0 <= y && y <= y1
	}

	// corner returns which edge/corner index (0..3) a boundary point sits on,
	// oriented by the walk direction. Ported from rectangle.js corner().
	corner := func(p [3]float64, direction float64) int {
		switch {
		case math.Abs(p[0]-x0) < epsilon:
			if direction > 0 {
				return 0
			}
			return 3
		case math.Abs(p[0]-x1) < epsilon:
			if direction > 0 {
				return 2
			}
			return 1
		case math.Abs(p[1]-y0) < epsilon:
			if direction > 0 {
				return 1
			}
			return 0
		default: // abs(p[1]-y1) < epsilon
			if direction > 0 {
				return 3
			}
			return 2
		}
	}

	comparePoint := func(a, b [3]float64) float64 {
		ca := corner(a, 1)
		cb := corner(b, 1)
		switch {
		case ca != cb:
			return float64(ca - cb)
		case ca == 0:
			return b[1] - a[1]
		case ca == 1:
			return a[0] - b[0]
		case ca == 2:
			return a[1] - b[1]
		default:
			return b[0] - a[0]
		}
	}

	compareIntersection := func(a, b *intersection) float64 { return comparePoint(a.x, b.x) }

	interpolate := func(from, to *[3]float64, direction float64, stream Sink) {
		a, a1 := 0, 0
		loop := from == nil
		if !loop {
			a = corner(*from, direction)
			a1 = corner(*to, direction)
			loop = a != a1 || boolXor(comparePoint(*from, *to) < 0, direction > 0)
		}
		if !loop {
			stream.Point(to[0], to[1])
			return
		}
		for { // do…while ((a = (a+direction+4)%4) !== a1)
			var px, py float64
			if a == 0 || a == 3 {
				px = x0
			} else {
				px = x1
			}
			if a > 1 {
				py = y1
			} else {
				py = y0
			}
			stream.Point(px, py)
			a = (a + int(direction) + 4) % 4
			if a == a1 {
				break
			}
		}
	}

	return func(stream Sink) Sink {
		return &rectClip{
			stream:       stream,
			activeStream: stream,
			bufferStream: &clipBuffer{},
			visible:      visible,
			interpolate:  interpolate,
			compare:      compareIntersection,
			x0:           x0,
			y0:           y0,
			x1:           x1,
			y1:           y1,
		}
	}
}

type rectClip struct {
	stream         Sink
	activeStream   Sink
	bufferStream   *clipBuffer
	visible        func(x, y float64) bool
	interpolate    interpolateFunc
	compare        func(a, b *intersection) float64
	x0, y0, x1, y1 float64

	buffering bool             // inside a polygon (segments are being buffered)
	segments  [][][][3]float64 // one buffer-result (list of lines) per ring
	polygon   [][][2]float64
	ring      [][2]float64

	xFirst, yFirst float64
	vFirst         bool
	xPrev, yPrev   float64
	vPrev          bool
	first          bool
	clean          bool

	pointFn func(x, y float64)
}

func (c *rectClip) Sphere() { c.stream.Sphere() }

func (c *rectClip) Point(x, y float64) {
	if c.pointFn != nil {
		c.pointFn(x, y)
		return
	}
	if c.visible(x, y) {
		c.activeStream.Point(x, y)
	}
}

// polygonInside returns whether the top-left clip corner is enclosed by the
// buffered polygon (nonzero winding). Ported from rectangle.js polygonInside.
func (c *rectClip) polygonInside() bool {
	winding := 0
	for _, ring := range c.polygon {
		m := len(ring)
		if m == 0 {
			continue
		}
		b0, b1 := ring[0][0], ring[0][1]
		for j := 1; j < m; j++ {
			a0, a1 := b0, b1
			b0, b1 = ring[j][0], ring[j][1]
			if a1 <= c.y1 {
				if b1 > c.y1 && (b0-a0)*(c.y1-a1) > (b1-a1)*(c.x0-a0) {
					winding++
				}
			} else {
				if b1 <= c.y1 && (b0-a0)*(c.y1-a1) < (b1-a1)*(c.x0-a0) {
					winding--
				}
			}
		}
	}
	return winding != 0
}

func (c *rectClip) PolygonStart() {
	c.activeStream = c.bufferStream
	c.buffering = true
	c.segments = nil
	c.polygon = [][][2]float64{}
	c.clean = true
}

func (c *rectClip) PolygonEnd() {
	startInside := c.polygonInside()
	cleanInside := c.clean && startInside
	merged := mergeSegments(c.segments)
	visible := len(merged) > 0
	if cleanInside || visible {
		c.stream.PolygonStart()
		if cleanInside {
			c.stream.LineStart()
			c.interpolate(nil, nil, 1, c.stream)
			c.stream.LineEnd()
		}
		if visible {
			clipRejoin(merged, c.compare, startInside, c.interpolate, c.stream)
		}
		c.stream.PolygonEnd()
	}
	c.activeStream = c.stream
	c.buffering = false
	c.segments = nil
	c.polygon = nil
	c.ring = nil
}

func (c *rectClip) LineStart() {
	c.pointFn = c.linePoint
	c.ring = nil
	c.first = true
	c.vPrev = false
	c.xPrev = math.NaN()
	c.yPrev = math.NaN()
}

func (c *rectClip) LineEnd() {
	if c.buffering {
		c.linePoint(c.xFirst, c.yFirst)
		if c.vFirst && c.vPrev {
			c.bufferStream.rejoin()
		}
		c.segments = append(c.segments, c.bufferStream.result())
		c.polygon = append(c.polygon, c.ring)
	}
	c.pointFn = nil
	if c.vPrev {
		c.activeStream.LineEnd()
	}
}

func (c *rectClip) linePoint(x, y float64) {
	v := c.visible(x, y)
	if c.buffering {
		c.ring = append(c.ring, [2]float64{x, y})
	}
	if c.first {
		c.xFirst, c.yFirst, c.vFirst = x, y, v
		c.first = false
		if v {
			c.activeStream.LineStart()
			c.activeStream.Point(x, y)
		}
	} else {
		if v && c.vPrev {
			c.activeStream.Point(x, y)
		} else {
			a := [2]float64{clamp(c.xPrev), clamp(c.yPrev)}
			b := [2]float64{clamp(x), clamp(y)}
			if clipLineRect(&a, &b, c.x0, c.y0, c.x1, c.y1) {
				if !c.vPrev {
					c.activeStream.LineStart()
					c.activeStream.Point(a[0], a[1])
				}
				c.activeStream.Point(b[0], b[1])
				if !v {
					c.activeStream.LineEnd()
				}
				c.clean = false
			} else if v {
				c.activeStream.LineStart()
				c.activeStream.Point(x, y)
				c.clean = false
			}
		}
	}
	c.xPrev, c.yPrev, c.vPrev = x, y, v
}

func clamp(v float64) float64 { return math.Max(clipMin, math.Min(clipMax, v)) }

// mergeSegments flattens the per-ring buffer results (each a list of lines) into
// a single list of lines, mirroring d3-array merge() as used by rectangle.js.
func mergeSegments(segs [][][][3]float64) [][][3]float64 {
	var out [][][3]float64
	for _, s := range segs {
		out = append(out, s...)
	}
	return out
}
