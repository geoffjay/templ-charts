// pie.go — port of d3-shape's pie generator (src/pie.js).
//
// Pie() returns a generator that computes arc angles for a dataset of values.
// Each value gets a [startAngle, endAngle] slice of the total, optionally
// sorted by value (descending, the d3 default) or by data, with optional
// padding between arcs and a configurable start/end angle for the whole pie.
//
// Angles are in radians, starting from the +x axis (12 o'clock after the arc
// generator's -π/2 offset). The default span is [0, 2π].
package d3shape

import "math"

// PieArc is one computed arc from the pie generator.
type PieArc[D any] struct {
	Data       D
	Index      int
	Value      float64
	StartAngle float64
	EndAngle   float64
	PadAngle   float64
}

// PieValueAccessor extracts a float64 from a datum.
type PieValueAccessor[D any] func(d D, i int, data []D) float64

// PieComparator compares two values or data, returning -1/0/1 (d3 ascending).
type (
	PieComparator[D any] func(a, b D) int
	PieValueComparator   func(a, b float64) int
)

// Pie is a pie generator parameterized by the data type D.
type Pie[D any] struct {
	value      PieValueAccessor[D]
	sortValues PieValueComparator
	sort       PieComparator[D]
	startAngle func() float64
	endAngle   func() float64
	padAngle   func() float64
}

// identityValue returns d as a float64 (for numeric data).
func identityValue[D any](d D, _ int, _ []D) float64 {
	// Can't convert generic D to float64; this is a placeholder.
	// For numeric slices, use NewPieNumeric or set a Value accessor.
	return 0
}

// NewPie constructs a Pie with value=identity (0 for non-float data — set
// Value() for real usage), sortValues=descending, startAngle=0, endAngle=2π,
// padAngle=0.
func NewPie[D any]() *Pie[D] {
	return &Pie[D]{
		value:      identityValue[D],
		sortValues: descendingFloat,
		sort:       nil,
		startAngle: func() float64 { return 0 },
		endAngle:   func() float64 { return tau },
		padAngle:   func() float64 { return 0 },
	}
}

// NewPieNumeric constructs a Pie for []float64 data with value=identity.
func NewPieNumeric() *Pie[float64] {
	return &Pie[float64]{
		value:      func(d float64, _ int, _ []float64) float64 { return d },
		sortValues: descendingFloat,
		startAngle: func() float64 { return 0 },
		endAngle:   func() float64 { return tau },
		padAngle:   func() float64 { return 0 },
	}
}

// Value sets the value accessor.
func (p *Pie[D]) Value(fn PieValueAccessor[D]) *Pie[D] {
	if fn != nil {
		p.value = fn
	}
	return p
}

// SortValues sets the value comparator. nil disables value sorting.
func (p *Pie[D]) SortValues(fn PieValueComparator) *Pie[D] {
	p.sortValues = fn
	p.sort = nil
	return p
}

// Sort sets the data comparator. nil disables data sorting.
func (p *Pie[D]) Sort(fn PieComparator[D]) *Pie[D] {
	p.sort = fn
	p.sortValues = nil
	return p
}

// StartAngle sets the start angle (radians).
func (p *Pie[D]) StartAngle(a float64) *Pie[D] {
	p.startAngle = func() float64 { return a }
	return p
}

// EndAngle sets the end angle (radians).
func (p *Pie[D]) EndAngle(a float64) *Pie[D] {
	p.endAngle = func() float64 { return a }
	return p
}

// PadAngle sets the padding between arcs (radians).
func (p *Pie[D]) PadAngle(a float64) *Pie[D] {
	p.padAngle = func() float64 { return a }
	return p
}

// Call computes the pie arcs for the given data.
func (p *Pie[D]) Call(data []D) []PieArc[D] {
	n := len(data)
	if n == 0 {
		return nil
	}
	arcs := make([]PieArc[D], n)
	index := make([]int, n)
	sum := 0.0
	a0 := p.startAngle()
	da := math.Min(tau, math.Max(-tau, p.endAngle()-a0))
	pad := math.Min(math.Abs(da)/float64(n), p.padAngle())
	pa := pad
	if da < 0 {
		pa = -pad
	}

	for i := 0; i < n; i++ {
		index[i] = i
		v := p.value(data[i], i, data)
		arcs[i].Value = v
		arcs[i].Data = data[i]
		if v > 0 {
			sum += v
		}
	}

	// Sort
	if p.sortValues != nil {
		// stable sort index by arcs[index].Value
		for i := 1; i < n; i++ {
			for j := i; j > 0 && p.sortValues(arcs[index[j]].Value, arcs[index[j-1]].Value) < 0; j-- {
				index[j], index[j-1] = index[j-1], index[j]
			}
		}
	} else if p.sort != nil {
		for i := 1; i < n; i++ {
			for j := i; j > 0 && p.sort(data[index[j]], data[index[j-1]]) < 0; j-- {
				index[j], index[j-1] = index[j-1], index[j]
			}
		}
	}

	// Compute arcs
	k := 0.0
	if sum != 0 {
		k = (da - float64(n)*pa) / sum
	}
	a0cur := a0
	for i := 0; i < n; i++ {
		j := index[i]
		v := arcs[j].Value
		a1 := a0cur
		if v > 0 {
			a1 += v * k
		}
		a1 += pa
		arcs[j].Index = i
		arcs[j].StartAngle = a0cur
		arcs[j].EndAngle = a1
		arcs[j].PadAngle = pad
		a0cur = a1
	}
	return arcs
}

// descendingFloat is d3's default sortValues (descending by value).
func descendingFloat(a, b float64) int {
	if a < b {
		return 1
	}
	if a > b {
		return -1
	}
	return 0
}
