// Package chord ports d3-chord (v3): a deterministic, single-pass layout that
// turns a square flow matrix into a set of group arcs around a circle plus the
// ribbons connecting them. Given an n×n matrix where matrix[i][j] is the flow
// from group i to group j, the layout assigns each group an angular span
// (startAngle,endAngle) proportional to its total flow and each subgroup
// (i→j flow) a sub-span, then pairs subgroups into ribbons.
//
// The port mirrors d3-chord's algorithm verbatim, including the subgroup index
// arithmetic and the source/target ordering (the larger of the two reciprocal
// flows becomes the ribbon source). It is fully deterministic — no randomness —
// and optional group/subgroup/chord comparators use stable sorts so goldens are
// byte-stable.
//
// Charts consume this via charts/chord, which runs the layout, renders the
// group arcs with charts/arcs (d3-shape Arc), and renders the ribbons with the
// Ribbon generator in this package.
package chord

import (
	"math"
	"sort"
)

const tau = 2 * math.Pi

// Subgroup is one directed sub-arc: the slice of group Index's span attributed
// to the flow toward SubIndex. Value is matrix[Index][SubIndex].
type Subgroup struct {
	Index      int
	SubIndex   int
	StartAngle float64
	EndAngle   float64
	Value      float64
}

// Group is one node's full arc around the circle. Value is the row sum.
type Group struct {
	Index      int
	StartAngle float64
	EndAngle   float64
	Value      float64
}

// Ribbon connects a Source subgroup to a Target subgroup. By convention Source
// is the subgroup with the larger value (d3-chord orders them so).
type Ribbon struct {
	Source Subgroup
	Target Subgroup
}

// Result is the computed layout: the group arcs (indexed by group) and the
// ribbons (in d3-chord emission order).
type Result struct {
	Groups  []Group
	Ribbons []Ribbon
}

// CompareFunc orders two values (returns <0, 0, >0). Used by the optional
// SortGroups/SortSubgroups/SortChords comparators.
type CompareFunc func(a, b float64) int

// Chord is a configurable chord layout. Create via New; configure with the
// chainable setters; run with Compute.
type Chord struct {
	padAngle      float64
	sortGroups    CompareFunc
	sortSubgroups CompareFunc
	sortChords    CompareFunc
}

// New returns a chord layout with padAngle 0 and no sorting (d3-chord defaults).
func New() *Chord { return &Chord{} }

// PadAngle sets the angular padding between adjacent groups (radians, clamped to
// >= 0). Mirrors d3-chord chord.padAngle.
func (c *Chord) PadAngle(v float64) *Chord {
	c.padAngle = math.Max(0, v)
	return c
}

// SortGroups sets the group comparator (by group value). Mirrors
// chord.sortGroups.
func (c *Chord) SortGroups(cmp CompareFunc) *Chord { c.sortGroups = cmp; return c }

// SortSubgroups sets the subgroup comparator (by cell value). Mirrors
// chord.sortSubgroups.
func (c *Chord) SortSubgroups(cmp CompareFunc) *Chord { c.sortSubgroups = cmp; return c }

// SortChords sets the ribbon comparator (by source.value+target.value). Mirrors
// chord.sortChords.
func (c *Chord) SortChords(cmp CompareFunc) *Chord { c.sortChords = cmp; return c }

// Compute runs the layout over a square matrix and returns the group arcs and
// ribbons. matrix must be n×n; ragged rows are treated as zero-padded.
func (c *Chord) Compute(matrix [][]float64) *Result {
	n := len(matrix)

	groupSums := make([]float64, 0, n)
	groupIndex := rangeN(n)
	subgroupIndex := make([][]int, 0, n)

	// Compute row sums and the running total.
	var k float64
	for i := 0; i < n; i++ {
		var x float64
		for j := 0; j < n; j++ {
			x += cell(matrix, i, j)
		}
		groupSums = append(groupSums, x)
		subgroupIndex = append(subgroupIndex, rangeN(n))
		k += x
	}

	// Sort groups by descending/ascending value if requested (stable).
	if c.sortGroups != nil {
		sort.SliceStable(groupIndex, func(a, b int) bool {
			return c.sortGroups(groupSums[groupIndex[a]], groupSums[groupIndex[b]]) < 0
		})
	}

	// Sort subgroups within each row if requested (stable).
	if c.sortSubgroups != nil {
		for i, d := range subgroupIndex {
			row := i
			sort.SliceStable(d, func(a, b int) bool {
				return c.sortSubgroups(cell(matrix, row, d[a]), cell(matrix, row, d[b])) < 0
			})
		}
	}

	// Convert the sum to a scaling factor for [0, tau], accounting for padding.
	if k > 0 {
		k = math.Max(0, tau-c.padAngle*float64(n)) / k
	}
	dx := 0.0
	if k != 0 {
		dx = c.padAngle
	} else if n > 0 {
		dx = tau / float64(n)
	}

	groups := make([]Group, n)
	subgroups := make([]Subgroup, n*n)

	var x float64
	for i := 0; i < n; i++ {
		di := groupIndex[i]
		x0 := x
		for j := 0; j < n; j++ {
			dj := subgroupIndex[di][j]
			v := cell(matrix, di, dj)
			a0 := x
			x += v * k
			a1 := x
			subgroups[dj*n+di] = Subgroup{
				Index:      di,
				SubIndex:   dj,
				StartAngle: a0,
				EndAngle:   a1,
				Value:      v,
			}
		}
		groups[di] = Group{
			Index:      di,
			StartAngle: x0,
			EndAngle:   x,
			Value:      groupSums[di],
		}
		x += dx
	}

	// Generate ribbons for each non-empty subgroup-subgroup pair.
	ribbons := make([]Ribbon, 0)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			source := subgroups[j*n+i]
			target := subgroups[i*n+j]
			if source.Value != 0 || target.Value != 0 {
				if source.Value < target.Value {
					ribbons = append(ribbons, Ribbon{Source: target, Target: source})
				} else {
					ribbons = append(ribbons, Ribbon{Source: source, Target: target})
				}
			}
		}
	}

	if c.sortChords != nil {
		sort.SliceStable(ribbons, func(a, b int) bool {
			va := ribbons[a].Source.Value + ribbons[a].Target.Value
			vb := ribbons[b].Source.Value + ribbons[b].Target.Value
			return c.sortChords(va, vb) < 0
		})
	}

	return &Result{Groups: groups, Ribbons: ribbons}
}

// rangeN returns [0, 1, …, n-1].
func rangeN(n int) []int {
	r := make([]int, n)
	for i := range r {
		r[i] = i
	}
	return r
}

// cell reads matrix[i][j], treating out-of-range/ragged entries as 0.
func cell(matrix [][]float64, i, j int) float64 {
	if i < 0 || i >= len(matrix) {
		return 0
	}
	row := matrix[i]
	if j < 0 || j >= len(row) {
		return 0
	}
	return row[j]
}
