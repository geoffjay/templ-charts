// lcg.go — port of d3-force's built-in linear congruential generator, used for
// the jiggle applied to coincident bodies so runs are deterministic without
// Math.random. Constants match d3 (and glibc): a=1664525, c=1013904223,
// m=2^32; the seed is 1. Values stay within float64's exact-integer range, so
// results are bit-identical to d3's JS. (Identical to internal/d3/hierarchy's
// LCG — each package keeps its own copy to stay self-contained.)
package force

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
