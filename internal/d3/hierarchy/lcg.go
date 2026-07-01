// lcg.go — port of d3-hierarchy's built-in linear congruential generator.
// Used by the pack layout in place of Math.random so packing is deterministic
// and reproducible server-side. Constants match d3 (and glibc): a=1664525,
// c=1013904223, m=2^32. Values stay well within float64's exact-integer range,
// so results are bit-identical to d3's JS implementation.
package d3hierarchy

import "math"

type lcg struct{ s float64 }

func newLCG() *lcg { return &lcg{s: 1} }

func (l *lcg) next() float64 {
	const a = 1664525.0
	const c = 1013904223.0
	const m = 4294967296.0 // 2^32
	l.s = math.Mod(a*l.s+c, m)
	return l.s / m
}

// shuffle returns a Fisher–Yates shuffle of a copy of xs, driven by random.
// Mirrors d3-hierarchy's shuffle (i = random()*m | 0).
func shuffle(xs []*Node, random func() float64) []*Node {
	out := make([]*Node, len(xs))
	copy(out, xs)
	m := len(out)
	for m > 0 {
		i := int(random() * float64(m))
		m--
		out[m], out[i] = out[i], out[m]
	}
	return out
}
