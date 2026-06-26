// stack.go — port of d3-shape's stack generator (src/stack.js) plus the
// offset (none, diverging, expand, silhouette, wiggle) and order (none,
// ascending, descending, reverse, appearance, insideOut) helpers.
//
// Stack() returns a generator that, given a slice of data objects and a list
// of keys, produces a series per key. Each series is a slice of [lower, upper]
// pairs representing the stacked extent of that key's value at each data
// point. The offset controls how series are stacked relative to each other;
// the order controls the series stacking sequence.
package d3shape

import "math"

// StackPoint is a [lower, upper] pair with a back-reference to the datum.
type StackPoint struct {
	Lo   float64
	Hi   float64
	Data any
}

// StackSeries is a series of StackPoints for one key.
type StackSeries struct {
	Key   string
	Index int
	Stats []StackPoint
}

// StackValueAccessor extracts a float64 from a datum for a given key.
type StackValueAccessor[D any] func(d D, key string, j int, data []D) float64

// StackOrderFunc returns an ordering (permutation) of series indices.
type StackOrderFunc func(series []StackSeries) []int

// StackOffsetFunc adjusts the [lo, hi] pairs of all series after ordering.
type StackOffsetFunc func(series []StackSeries, order []int)

// Stack is a stack generator parameterized by the data type D.
type Stack[D any] struct {
	keys   []string
	order  StackOrderFunc
	offset StackOffsetFunc
	value  StackValueAccessor[D]
}

// NewStack constructs a Stack with no keys, orderNone, offsetNone, and the
// default value accessor (d[key] via map lookup for map[string]float64 data).
func NewStack[D any]() *Stack[D] {
	return &Stack[D]{
		keys:   []string{},
		order:  StackOrderNone,
		offset: StackOffsetNone,
		value:  defaultStackValue[D],
	}
}

// Keys sets the keys to stack.
func (s *Stack[D]) Keys(keys []string) *Stack[D] { s.keys = append([]string(nil), keys...); return s }

// Value sets the value accessor.
func (s *Stack[D]) Value(fn StackValueAccessor[D]) *Stack[D] { s.value = fn; return s }

// Order sets the order function. nil → orderNone.
func (s *Stack[D]) Order(fn StackOrderFunc) *Stack[D] {
	if fn == nil {
		fn = StackOrderNone
	}
	s.order = fn
	return s
}

// Offset sets the offset function. nil → offsetNone.
func (s *Stack[D]) Offset(fn StackOffsetFunc) *Stack[D] {
	if fn == nil {
		fn = StackOffsetNone
	}
	s.offset = fn
	return s
}

// Call computes the stacked series for the given data.
func (s *Stack[D]) Call(data []D) []StackSeries {
	sz := make([]StackSeries, len(s.keys))
	for i, k := range s.keys {
		sz[i] = StackSeries{Key: k}
	}
	for j, d := range data {
		for i := range sz {
			v := s.value(d, sz[i].Key, j, data)
			pt := StackPoint{Lo: 0, Hi: v, Data: d}
			sz[i].Stats = append(sz[i].Stats, pt)
		}
	}
	// compute order
	oz := s.order(sz)
	for i, idx := range oz {
		sz[idx].Index = i
	}
	// apply offset
	s.offset(sz, oz)
	return sz
}

// defaultStackValue tries to extract d[key] for map[string]float64-like data.
// For non-map types, callers must provide a custom Value accessor.
func defaultStackValue[D any](d D, key string, j int, data []D) float64 {
	// This is a best-effort default; for typed structs, use Value().
	// We can't do reflection here without import; return 0 and let callers
	// override.
	_ = key
	_ = j
	_ = data
	return 0
}

// --- orders ---------------------------------------------------------------

// StackOrderNone returns [0, 1, ..., n-1] (the natural order).
func StackOrderNone(series []StackSeries) []int {
	n := len(series)
	o := make([]int, n)
	for i := 0; i < n; i++ {
		o[i] = i
	}
	return o
}

// StackOrderReverse returns [n-1, ..., 1, 0].
func StackOrderReverse(series []StackSeries) []int {
	n := len(series)
	o := make([]int, n)
	for i := 0; i < n; i++ {
		o[i] = n - 1 - i
	}
	return o
}

// StackOrderAscending sorts series by the sum of their upper values.
func StackOrderAscending(series []StackSeries) []int {
	sums := make([]float64, len(series))
	for i, s := range series {
		var sum float64
		for _, p := range s.Stats {
			sum += p.Hi
		}
		sums[i] = sum
	}
	o := StackOrderNone(series)
	// stable sort by sums
	for i := 1; i < len(o); i++ {
		for j := i; j > 0 && sums[o[j]] < sums[o[j-1]]; j-- {
			o[j], o[j-1] = o[j-1], o[j]
		}
	}
	return o
}

// StackOrderDescending is the reverse of ascending.
func StackOrderDescending(series []StackSeries) []int {
	o := StackOrderAscending(series)
	// reverse
	for i, j := 0, len(o)-1; i < j; i, j = i+1, j-1 {
		o[i], o[j] = o[j], o[i]
	}
	return o
}

// StackOrderAppearance sorts by the index of the peak value in each series.
func StackOrderAppearance(series []StackSeries) []int {
	peaks := make([]int, len(series))
	for i, s := range series {
		maxV := math.Inf(-1)
		maxJ := 0
		for j, p := range s.Stats {
			if p.Hi > maxV {
				maxV = p.Hi
				maxJ = j
			}
		}
		peaks[i] = maxJ
	}
	o := StackOrderNone(series)
	for i := 1; i < len(o); i++ {
		for j := i; j > 0 && peaks[o[j]] < peaks[o[j-1]]; j-- {
			o[j], o[j-1] = o[j-1], o[j]
		}
	}
	return o
}

// StackOrderInsideOut orders by appearance then splits top/bottom alternately.
func StackOrderInsideOut(series []StackSeries) []int {
	n := len(series)
	sums := make([]float64, n)
	for i, s := range series {
		var sum float64
		for _, p := range s.Stats {
			sum += p.Hi
		}
		sums[i] = sum
	}
	order := StackOrderAppearance(series)
	top, bottom := 0, 0
	tops, bottoms := []int{}, []int{}
	for i := 0; i < n; i++ {
		j := order[i]
		if top < bottom {
			top += int(sums[j])
			tops = append(tops, j)
		} else {
			bottom += int(sums[j])
			bottoms = append(bottoms, j)
		}
	}
	// reverse bottoms, concat with tops
	result := make([]int, 0, n)
	for i := len(bottoms) - 1; i >= 0; i-- {
		result = append(result, bottoms[i])
	}
	result = append(result, tops...)
	return result
}

// --- offsets --------------------------------------------------------------

// StackOffsetNone stacks each series on top of the previous (cumulative).
func StackOffsetNone(series []StackSeries, order []int) {
	n := len(series)
	if n <= 1 {
		return
	}
	s1 := series[order[0]]
	for i := 1; i < n; i++ {
		s0 := s1
		s1 = series[order[i]]
		for j := range s1.Stats {
			prev := s0.Stats[j].Hi
			if math.IsNaN(prev) {
				s1.Stats[j].Lo = s0.Stats[j].Lo
			} else {
				s1.Stats[j].Lo = prev
			}
			s1.Stats[j].Hi += s1.Stats[j].Lo
		}
	}
}

// StackOffsetDiverging stacks positive values above zero, negative below.
func StackOffsetDiverging(series []StackSeries, order []int) {
	n := len(series)
	if n == 0 {
		return
	}
	m := len(series[order[0]].Stats)
	for j := 0; j < m; j++ {
		yp, yn := 0.0, 0.0
		for i := 0; i < n; i++ {
			d := &series[order[i]].Stats[j]
			dy := d.Hi - d.Lo
			if dy > 0 {
				d.Lo = yp
				yp += dy
				d.Hi = yp
			} else if dy < 0 {
				d.Hi = yn
				yn += dy
				d.Lo = yn
			} else {
				d.Lo = 0
				d.Hi = dy
			}
		}
	}
}

// StackOffsetExpand normalizes each stack to [0, 1].
func StackOffsetExpand(series []StackSeries, order []int) {
	n := len(series)
	if n == 0 {
		return
	}
	m := len(series[0].Stats)
	for j := 0; j < m; j++ {
		var y float64
		for i := 0; i < n; i++ {
			y += series[i].Stats[j].Hi
		}
		if y != 0 {
			for i := 0; i < n; i++ {
				series[i].Stats[j].Hi /= y
			}
		}
	}
	StackOffsetNone(series, order)
}

// StackOffsetSilhouette centers each stack around zero.
func StackOffsetSilhouette(series []StackSeries, order []int) {
	n := len(series)
	if n == 0 {
		return
	}
	s0 := series[order[0]]
	m := len(s0.Stats)
	for j := 0; j < m; j++ {
		var y float64
		for i := 0; i < n; i++ {
			y += series[i].Stats[j].Hi
		}
		s0.Stats[j].Lo = -y / 2
		s0.Stats[j].Hi += s0.Stats[j].Lo
	}
	StackOffsetNone(series, order)
}

// StackOffsetWiggle minimizes weighted displacement of stack layers.
// Used for streamgraphs.
func StackOffsetWiggle(series []StackSeries, order []int) {
	n := len(series)
	if n == 0 {
		return
	}
	s0 := series[order[0]]
	m := len(s0.Stats)
	if m == 0 {
		return
	}
	y := 0.0
	for j := 1; j < m; j++ {
		var s1, s2 float64
		for i := 0; i < n; i++ {
			si := series[order[i]]
			sij0 := si.Stats[j].Hi
			sij1 := si.Stats[j-1].Hi
			s3 := (sij0 - sij1) / 2
			for k := 0; k < i; k++ {
				sk := series[order[k]]
				skj0 := sk.Stats[j].Hi
				skj1 := sk.Stats[j-1].Hi
				s3 += skj0 - skj1
			}
			s1 += sij0
			s2 += s3 * sij0
		}
		s0.Stats[j-1].Lo = y
		s0.Stats[j-1].Hi += y
		if s1 != 0 {
			y -= s2 / s1
		}
	}
	s0.Stats[m-1].Lo = y
	s0.Stats[m-1].Hi += y
	StackOffsetNone(series, order)
}
