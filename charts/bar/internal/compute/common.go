// Package compute contains the bar layout algorithms: index/value scale
// construction, data normalization, label-layout, grouped/stacked bar
// generation, legend-data derivation, and totals computation. Mirrors
// @nivo/bar/src/compute.
package compute

import (
	"strconv"

	"github.com/geoffjay/templ-charts/charts/scales"
	d3scale "github.com/geoffjay/templ-charts/internal/d3/scale"
)

// GetIndexScale builds the band scale for the index axis (the axis that
// carries the index values, e.g. x for vertical layout). Mirrors
// @nivo/bar/src/compute/common.ts getIndexScale. padding is applied via
// scaleBand.padding(); the spec's Round is honored.
func GetIndexScale(
	data []map[string]any,
	getIndex func(d map[string]any) string,
	padding float64,
	indexSpec scales.ScaleBandSpec,
	size float64,
	axis scales.ScaleAxis,
) scales.Scale {
	r0, r1 := 0.0, size
	if axis == scales.ScaleAxisY {
		r0, r1 = size, 0
	}
	strs := make([]string, 0, len(data))
	for _, d := range data {
		strs = append(strs, getIndex(d))
	}
	b := d3scale.NewBand().
		SetRange(r0, r1).
		SetDomain(strs).
		SetRound(indexSpec.Round).
		SetPadding(padding)
	return bandScale{b: b}
}

// bandScale adapts a *d3scale.Band to the scales.Scale +
// scales.ScaleWithBandwidth interfaces, and exposes Domain() for the
// generators.
type bandScale struct{ b *d3scale.Band }

func (s bandScale) Type() scales.ScaleType { return scales.ScaleTypeBand }
func (s bandScale) Call(v any) float64     { return s.b.Call(toString(v)) }
func (s bandScale) Bandwidth() float64     { return s.b.Bandwidth() }
func (s bandScale) Step() float64          { return s.b.Step() }
func (s bandScale) Round() bool            { return s.b.Round() }
func (s bandScale) Domain() []string       { return s.b.Domain() }

// IndexScaleWithDomain is the interface the generators need: a band scale
// exposing Domain() (the ordered index values) plus the standard
// ScaleWithBandwidth surface.
type IndexScaleWithDomain interface {
	scales.ScaleWithBandwidth
	Domain() []string
}

// NormalizeData ensures every key is present on every datum (missing keys
// become nil). Mirrors @nivo/bar/src/compute/common.ts normalizeData.
func NormalizeData(data []map[string]any, keys []string) []map[string]any {
	out := make([]map[string]any, len(data))
	for i, item := range data {
		merged := map[string]any{}
		for _, k := range keys {
			merged[k] = nil
		}
		for k, v := range item {
			merged[k] = v
		}
		out[i] = merged
	}
	return out
}

// FilterNullValues returns a copy of `data` with nil/empty/zero-valued keys
// removed. Mirrors @nivo/bar/src/compute/common.ts filterNullValues.
func FilterNullValues(data map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range data {
		if v == nil {
			continue
		}
		if s, ok := v.(string); ok && s == "" {
			continue
		}
		if f, ok := v.(float64); ok && f == 0 {
			continue
		}
		out[k] = v
	}
	return out
}

// CoerceValue mirrors @nivo/bar coerceValue: returns (rawValue, Number(value)).
// A nil raw value yields (nil, 0). Non-numeric values yield (rawValue, 0).
func CoerceValue(v any) (any, float64) {
	if v == nil {
		return nil, 0
	}
	f, ok := toFloatOK(v)
	if !ok {
		return v, 0
	}
	return v, f
}

// BarLabelLayout is the label position/alignment for one bar. Mirrors nivo's
// BarLabelLayout.
type BarLabelLayout struct {
	LabelX     float64
	LabelY     float64
	TextAnchor string
}

// ComputedDatum is one key's value within an index, after normalization.
// Mirrors @nivo/bar ComputedDatum<D>. Value is nil when the raw value was
// null/missing (nivo: `value: rawValue === null ? rawValue : value`).
type ComputedDatum struct {
	ID             string
	Value          any // float64 | nil
	FormattedValue string
	Hidden         bool
	Index          int
	IndexValue     string
	Data           map[string]any
	Fill           string
}

// ComputedBarDatum is a ComputedDatum resolved to a bar rectangle (x/y/w/h)
// plus color and label. Mirrors @nivo/bar ComputedBarDatum<D>.
type ComputedBarDatum struct {
	Key           string
	Index         int
	Data          ComputedDatum
	X, Y          float64
	AbsX, AbsY    float64
	Width, Height float64
	Color         string
	Label         string
}

// LegendData is one legend entry. Mirrors @nivo/bar LegendData.
type LegendData struct {
	ID     string
	Label  string
	Hidden bool
	Color  string
}

// BarTotalData is one totals label. Mirrors @nivo/bar BarTotalsData.
type BarTotalData struct {
	Key             string
	X, Y            float64
	Value           float64
	FormattedValue  string
	AnimationOffset float64
}

// MarginLike is the subset of core.Margin used by the generators (absX/absY).
type MarginLike = struct {
	Top, Right, Bottom, Left float64
}

// ComputeLabelLayout builds a func(width, height) BarLabelLayout for the given
// layout/labelPosition/labelOffset, honoring valueScale.reverse. Mirrors nivo's
// useComputeLabelLayout.
func ComputeLabelLayout(layout string, reverse bool, labelPosition string, labelOffset float64) func(width, height float64) BarLabelLayout {
	return func(width, height float64) BarLabelLayout {
		offset := labelOffset
		if reverse {
			offset = -offset
		}
		if layout == "horizontal" {
			x := width / 2
			switch labelPosition {
			case "start":
				if reverse {
					x = width
				} else {
					x = 0
				}
			case "end":
				if reverse {
					x = 0
				} else {
					x = width
				}
			}
			anchor := "middle"
			if labelPosition != "middle" {
				if reverse {
					anchor = "end"
				} else {
					anchor = "start"
				}
			}
			return BarLabelLayout{LabelX: x + offset, LabelY: height / 2, TextAnchor: anchor}
		}
		// vertical
		y := height / 2
		switch labelPosition {
		case "start":
			if reverse {
				y = 0
			} else {
				y = height
			}
		case "end":
			if reverse {
				y = height
			} else {
				y = 0
			}
		}
		return BarLabelLayout{LabelX: width / 2, LabelY: y - offset, TextAnchor: "middle"}
	}
}

// --- helpers ---------------------------------------------------------------

func toString(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	}
	return formatAny(v)
}

func formatAny(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	}
	return ""
}

func toFloatOK(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case int32:
		return float64(x), true
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}
