// enclose.go — port of d3-hierarchy's smallest-enclosing-circle (Welzl's
// algorithm with a randomized basis, matching d3's packEncloseRandom).
package d3hierarchy

import "math"

// circle is a plain (x, y, r) used by the enclose math (distinct from *Node so
// basis circles computed by encloseBasis2/3 need no backing node).
type circle struct{ x, y, r float64 }

func nodeCircle(n *Node) circle { return circle{x: n.X, y: n.Y, r: n.R} }

// packEncloseRandom returns the smallest circle enclosing all the given nodes'
// circles, using the provided random source to shuffle (Welzl). Mirrors
// d3-hierarchy packEncloseRandom.
func packEncloseRandom(nodes []*Node, random func() float64) circle {
	shuffled := shuffle(nodes, random)
	n := len(shuffled)
	var b []circle
	var e circle
	haveE := false
	i := 0
	for i < n {
		p := nodeCircle(shuffled[i])
		if haveE && enclosesWeak(e, p) {
			i++
		} else {
			b = extendBasis(b, p)
			e = encloseBasis(b)
			haveE = true
			i = 0
		}
	}
	return e
}

func extendBasis(b []circle, p circle) []circle {
	if enclosesWeakAll(p, b) {
		return []circle{p}
	}
	// b must have at least one element.
	for i := 0; i < len(b); i++ {
		if enclosesNot(p, b[i]) && enclosesWeakAll(encloseBasis2(b[i], p), b) {
			return []circle{b[i], p}
		}
	}
	// b must have at least two elements.
	for i := 0; i < len(b)-1; i++ {
		for j := i + 1; j < len(b); j++ {
			if enclosesNot(encloseBasis2(b[i], b[j]), p) &&
				enclosesNot(encloseBasis2(b[i], p), b[j]) &&
				enclosesNot(encloseBasis2(b[j], p), b[i]) &&
				enclosesWeakAll(encloseBasis3(b[i], b[j], p), b) {
				return []circle{b[i], b[j], p}
			}
		}
	}
	// Should be unreachable for valid input.
	return []circle{p}
}

func enclosesNot(a, b circle) bool {
	dr := a.r - b.r
	dx := b.x - a.x
	dy := b.y - a.y
	return dr < 0 || dr*dr < dx*dx+dy*dy
}

func enclosesWeak(a, b circle) bool {
	dr := a.r - b.r + math.Max(math.Max(a.r, b.r), 1)*1e-9
	dx := b.x - a.x
	dy := b.y - a.y
	return dr > 0 && dr*dr > dx*dx+dy*dy
}

func enclosesWeakAll(a circle, b []circle) bool {
	for i := 0; i < len(b); i++ {
		if !enclosesWeak(a, b[i]) {
			return false
		}
	}
	return true
}

func encloseBasis(b []circle) circle {
	switch len(b) {
	case 1:
		return b[0]
	case 2:
		return encloseBasis2(b[0], b[1])
	case 3:
		return encloseBasis3(b[0], b[1], b[2])
	}
	return circle{}
}

func encloseBasis2(a, b circle) circle {
	x1, y1, r1 := a.x, a.y, a.r
	x2, y2, r2 := b.x, b.y, b.r
	x21 := x2 - x1
	y21 := y2 - y1
	r21 := r2 - r1
	l := math.Sqrt(x21*x21 + y21*y21)
	return circle{
		x: (x1 + x2 + x21/l*r21) / 2,
		y: (y1 + y2 + y21/l*r21) / 2,
		r: (l + r1 + r2) / 2,
	}
}

func encloseBasis3(a, b, c circle) circle {
	x1, y1, r1 := a.x, a.y, a.r
	x2, y2, r2 := b.x, b.y, b.r
	x3, y3, r3 := c.x, c.y, c.r
	a2 := x1 - x2
	a3 := x1 - x3
	b2 := y1 - y2
	b3 := y1 - y3
	c2 := r2 - r1
	c3 := r3 - r1
	d1 := x1*x1 + y1*y1 - r1*r1
	d2 := d1 - x2*x2 - y2*y2 + r2*r2
	d3 := d1 - x3*x3 - y3*y3 + r3*r3
	ab := a3*b2 - a2*b3
	xa := (b2*d3-b3*d2)/(ab*2) - x1
	xb := (b3*c2 - b2*c3) / ab
	ya := (a3*d2-a2*d3)/(ab*2) - y1
	yb := (a2*c3 - a3*c2) / ab
	A := xb*xb + yb*yb - 1
	B := 2 * (r1 + xa*xb + ya*yb)
	C := xa*xa + ya*ya - r1*r1
	var r float64
	if math.Abs(A) > 1e-6 {
		r = -((B + math.Sqrt(B*B-4*A*C)) / (2 * A))
	} else {
		r = -(C / B)
	}
	return circle{
		x: x1 + xa + xb*r,
		y: y1 + ya + yb*r,
		r: r,
	}
}
