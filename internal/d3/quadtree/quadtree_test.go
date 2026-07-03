package quadtree_test

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/internal/d3/quadtree"
	"github.com/geoffjay/templ-charts/internal/golden"
)

// pt is a test datum positioned by X/Y.
type pt struct {
	X, Y float64
}

func ax(d any) float64 { return d.(pt).X }
func ay(d any) float64 { return d.(pt).Y }

func build(pts []pt) *quadtree.Quadtree {
	data := make([]any, len(pts))
	for i, p := range pts {
		data[i] = p
	}
	return quadtree.From(data, ax, ay)
}

// fixedPoints is a deterministic, tie-free point set (no two share a coordinate
// pair; distances to the query points below are distinct).
func fixedPoints() []pt {
	var pts []pt
	for i := 0; i < 40; i++ {
		x := math.Mod(float64(i)*7.3, 97)
		y := math.Mod(float64(i)*13.1+3, 89)
		pts = append(pts, pt{X: x, Y: y})
	}
	return pts
}

func TestFrom_SizeAndData(t *testing.T) {
	pts := fixedPoints()
	q := build(pts)
	if got := q.Size(); got != len(pts) {
		t.Fatalf("Size = %d, want %d", got, len(pts))
	}
	if got := len(q.Data()); got != len(pts) {
		t.Fatalf("len(Data) = %d, want %d", got, len(pts))
	}
}

// TestLeafPointsWithinCell asserts the core quadtree invariant: every leaf's
// point lies inside that leaf's cell bounds. This is what makes the spatial
// pruning in Find (and Barnes–Hut) correct.
func TestLeafPointsWithinCell(t *testing.T) {
	q := build(fixedPoints())
	q.Visit(func(n *quadtree.Node, x0, y0, x1, y1 float64) bool {
		if n.IsLeaf() {
			for node := n; node != nil; node = node.Next() {
				p := node.Data().(pt)
				if p.X < x0 || p.X > x1 || p.Y < y0 || p.Y > y1 {
					t.Errorf("leaf point (%g,%g) outside cell [%g,%g]-[%g,%g]", p.X, p.Y, x0, y0, x1, y1)
				}
			}
		}
		return false
	})
}

// TestFindMatchesBruteForce is the strong correctness check: Find must return
// the same nearest point a linear scan does, for a range of queries.
func TestFindMatchesBruteForce(t *testing.T) {
	pts := fixedPoints()
	q := build(pts)
	queries := []pt{{5, 5}, {50, 40}, {96, 88}, {0, 0}, {33.3, 71.7}, {-10, -10}, {120, 120}}
	for _, query := range queries {
		want := brute(pts, query.X, query.Y)
		got, ok := q.Find(query.X, query.Y)
		if !ok {
			t.Fatalf("Find(%g,%g) found nothing", query.X, query.Y)
		}
		if got.(pt) != want {
			t.Errorf("Find(%g,%g) = %+v, brute-force nearest = %+v", query.X, query.Y, got.(pt), want)
		}
	}
}

func TestFind_RadiusBounded(t *testing.T) {
	pts := fixedPoints()
	q := build(pts)
	// A query far from every point, with a radius smaller than the nearest
	// distance, must find nothing.
	if _, ok := q.Find(-100, -100, 5); ok {
		t.Errorf("Find with radius 5 at (-100,-100) should find nothing")
	}
	// With a huge radius it finds the brute-force nearest.
	want := brute(pts, -100, -100)
	got, ok := q.Find(-100, -100, 1e6)
	if !ok || got.(pt) != want {
		t.Errorf("Find(-100,-100, 1e6) = %v ok=%v, want %+v", got, ok, want)
	}
}

func TestFind_Empty(t *testing.T) {
	q := quadtree.New(ax, ay)
	if _, ok := q.Find(0, 0); ok {
		t.Errorf("Find on empty tree should return ok=false")
	}
}

func brute(pts []pt, x, y float64) pt {
	best := pts[0]
	bestD := math.Inf(1)
	for _, p := range pts {
		dx, dy := p.X-x, p.Y-y
		d := dx*dx + dy*dy
		if d < bestD {
			bestD = d
			best = p
		}
	}
	return best
}

func TestCoincidentPoints(t *testing.T) {
	// Three datums at the same coordinate must chain in one leaf, and Size/Data
	// must count all three.
	data := []any{pt{2, 2}, pt{2, 2}, pt{2, 2}, pt{9, 9}}
	q := quadtree.From(data, ax, ay)
	if got := q.Size(); got != 4 {
		t.Fatalf("Size = %d, want 4", got)
	}
	chain := 0
	q.Visit(func(n *quadtree.Node, _, _, _, _ float64) bool {
		if n.IsLeaf() {
			c := 0
			for node := n; node != nil; node = node.Next() {
				c++
			}
			if c > chain {
				chain = c
			}
		}
		return false
	})
	if chain != 3 {
		t.Errorf("longest coincident chain = %d, want 3", chain)
	}
}

func TestRemove(t *testing.T) {
	pts := fixedPoints()
	q := build(pts)
	// Remove half the points; the rest must survive and Find must still work.
	removed := map[pt]bool{}
	for i := 0; i < len(pts); i += 2 {
		q.Remove(pts[i])
		removed[pts[i]] = true
	}
	if got, want := q.Size(), len(pts)-len(removed); got != want {
		t.Fatalf("Size after removals = %d, want %d", got, want)
	}
	for _, d := range q.Data() {
		if removed[d.(pt)] {
			t.Errorf("removed point %+v still present", d.(pt))
		}
	}
	// The nearest of the survivors, per brute force, must match Find.
	var survivors []pt
	for _, p := range pts {
		if !removed[p] {
			survivors = append(survivors, p)
		}
	}
	got, ok := q.Find(50, 40)
	if !ok || got.(pt) != brute(survivors, 50, 40) {
		t.Errorf("Find after removals = %v, want %+v", got, brute(survivors, 50, 40))
	}
}

func TestRemove_CoincidentChain(t *testing.T) {
	a, b, c := pt{2, 2}, pt{2, 2}, pt{2, 2}
	q := quadtree.From([]any{a, b, c, pt{9, 9}}, ax, ay)
	q.Remove(b)
	if got := q.Size(); got != 3 {
		t.Fatalf("Size after removing one coincident = %d, want 3", got)
	}
	// Removing the remaining two coincident points, then the last, leaves 0.
	q.Remove(a)
	q.Remove(c)
	if got := q.Size(); got != 1 {
		t.Fatalf("Size = %d, want 1 (only (9,9) left)", got)
	}
}

func TestVisitPreOrderSkip(t *testing.T) {
	q := build(fixedPoints())
	// Returning true (skip) at the root visits only the root.
	visits := 0
	q.Visit(func(n *quadtree.Node, _, _, _, _ float64) bool {
		visits++
		return true
	})
	if visits != 1 {
		t.Errorf("skipping at root should visit exactly 1 node, got %d", visits)
	}
}

// TestVisitAfterAggregation exercises the Barnes–Hut use case: post-order,
// accumulate each internal node's total weight and centre of mass from its
// children, and verify against a direct computation over all leaves.
func TestVisitAfterAggregation(t *testing.T) {
	pts := fixedPoints()
	q := build(pts)
	q.VisitAfter(func(n *quadtree.Node, _, _, _, _ float64) {
		if n.IsLeaf() {
			// A leaf's weight is its coincident count; centre is its point.
			w := 0.0
			p := n.Data().(pt)
			for node := n; node != nil; node = node.Next() {
				w++
			}
			n.Value = w
			n.X, n.Y = p.X, p.Y
			return
		}
		var w, sx, sy float64
		for i := 0; i < 4; i++ {
			if c := n.Child(i); c != nil {
				w += c.Value
				sx += c.Value * c.X
				sy += c.Value * c.Y
			}
		}
		n.Value = w
		if w > 0 {
			n.X, n.Y = sx/w, sy/w
		}
	})
	// Direct totals.
	var wantW, wantSX, wantSY float64
	for _, p := range pts {
		wantW++
		wantSX += p.X
		wantSY += p.Y
	}
	root := q.Root()
	if root.Value != wantW {
		t.Errorf("root weight = %g, want %g", root.Value, wantW)
	}
	if !approx(root.X, wantSX/wantW) || !approx(root.Y, wantSY/wantW) {
		t.Errorf("root centre = (%g,%g), want (%g,%g)", root.X, root.Y, wantSX/wantW, wantSY/wantW)
	}
}

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestCopyIndependent(t *testing.T) {
	pts := fixedPoints()
	q := build(pts)
	c := q.Copy()
	if c.Size() != q.Size() {
		t.Fatalf("copy size %d != original %d", c.Size(), q.Size())
	}
	// Mutating the copy (remove) must not change the original.
	c.Remove(pts[0])
	if q.Size() != len(pts) {
		t.Errorf("original changed after mutating copy: size %d, want %d", q.Size(), len(pts))
	}
	if c.Size() != len(pts)-1 {
		t.Errorf("copy size %d, want %d", c.Size(), len(pts)-1)
	}
}

func TestCoverExtentGrows(t *testing.T) {
	q := quadtree.New(ax, ay)
	if _, _, _, _, ok := q.Extent(); ok {
		t.Errorf("empty tree should be uncovered")
	}
	q.Cover(0, 0)
	x0, y0, x1, y1, ok := q.Extent()
	if !ok || x0 != 0 || y0 != 0 || x1 != 1 || y1 != 1 {
		t.Errorf("cover(0,0) extent = [%g,%g,%g,%g] ok=%v, want [0,0,1,1]", x0, y0, x1, y1, ok)
	}
	q.Cover(3, 3)
	x0, y0, x1, y1, _ = q.Extent()
	// Must now contain (3,3) as a square power-of-two extent from the origin.
	if !(x0 <= 3 && 3 < x1 && y0 <= 3 && 3 < y1) {
		t.Errorf("cover(3,3) extent [%g,%g,%g,%g] does not contain (3,3)", x0, y0, x1, y1)
	}
	if (x1 - x0) != (y1 - y0) {
		t.Errorf("extent not square: %g x %g", x1-x0, y1-y0)
	}
}

func TestDeterministic(t *testing.T) {
	if dumpStructure(build(fixedPoints())) != dumpStructure(build(fixedPoints())) {
		t.Errorf("quadtree structure is not deterministic across builds")
	}
}

// TestGolden locks the exact tree structure for a fixed input. The structure is
// a pure function of the d3-quadtree insertion algorithm, so this snapshot
// encodes d3-equivalent output (and catches any regression in the port).
func TestGolden(t *testing.T) {
	golden.Assert(t, "quadtree-structure", dumpStructure(build(fixedPoints())))
}

// dumpStructure renders the tree as a deterministic pre-order outline: one line
// per node with its depth, kind, cell bounds, and (for leaves) the point(s).
func dumpStructure(q *quadtree.Quadtree) string {
	var b strings.Builder
	depth := map[*quadtree.Node]int{}
	depth[q.Root()] = 0
	q.Visit(func(n *quadtree.Node, x0, y0, x1, y1 float64) bool {
		d := depth[n]
		b.WriteString(strings.Repeat("  ", d))
		fmt.Fprintf(&b, "[%s,%s %s,%s] ", f(x0), f(y0), f(x1), f(y1))
		if n.IsLeaf() {
			b.WriteString("leaf")
			for node := n; node != nil; node = node.Next() {
				p := node.Data().(pt)
				fmt.Fprintf(&b, " (%s,%s)", f(p.X), f(p.Y))
			}
		} else {
			b.WriteString("node")
			xm, ym := (x0+x1)/2, (y0+y1)/2
			for i := 0; i < 4; i++ {
				if c := n.Child(i); c != nil {
					cx0, cy0, cx1, cy1 := x0, y0, xm, ym
					if i&1 == 1 {
						cx0, cx1 = xm, x1
					}
					if i&2 == 2 {
						cy0, cy1 = ym, y1
					}
					_ = cx0
					_ = cy0
					_ = cx1
					_ = cy1
					depth[c] = d + 1
				}
			}
		}
		b.WriteByte('\n')
		return false
	})
	return b.String()
}

func f(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
