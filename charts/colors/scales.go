package colors

import (
	"fmt"
	"math"

	d3array "github.com/geoffjay/templ-charts/internal/d3/array"
)

// DatumIdentityAccessor resolves a datum to its identity (string|number).
type DatumIdentityAccessor[D any] func(d D) string

// OrdinalColorScaleConfig is the union of nivo's ordinal color scale configs:
//   - static color string
//   - custom func(datum) string
//   - {scheme: id, size?: n} — predefined scheme
//   - []string — custom color array
//   - {datum: "path"} — color from datum property
//
// Mirrors @nivo/colors OrdinalColorScaleConfig.
type OrdinalColorScaleConfig struct {
	Type       OrdinalColorScaleType
	Static     string
	Func       func(any) string
	Scheme     string
	SchemeSize int // for diverging (3..11) / sequential (3..9)
	Colors     []string
	DatumPath  string
}

// OrdinalColorScaleType discriminates OrdinalColorScaleConfig variants.
type OrdinalColorScaleType int

const (
	OrdinalTypeStatic OrdinalColorScaleType = iota
	OrdinalTypeFunc
	OrdinalTypeScheme
	OrdinalTypeColors
	OrdinalTypeDatum
)

// ParseOrdinalColorScaleConfig accepts the loose forms nivo allows (string,
// func, []string, or map with scheme/datum/colors) and returns a typed
// OrdinalColorScaleConfig.
func ParseOrdinalColorScaleConfig(v any) (OrdinalColorScaleConfig, error) {
	switch x := v.(type) {
	case nil:
		return OrdinalColorScaleConfig{Type: OrdinalTypeStatic, Static: ""}, nil
	case string:
		// Could be a scheme id or a static color.
		if IsCategoricalColorScheme(x) || IsDivergingColorScheme(x) || IsSequentialColorScheme(x) {
			return OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: x}, nil
		}
		return OrdinalColorScaleConfig{Type: OrdinalTypeStatic, Static: x}, nil
	case func(any) string:
		return OrdinalColorScaleConfig{Type: OrdinalTypeFunc, Func: x}, nil
	case []string:
		return OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: x}, nil
	case []any:
		colors := make([]string, 0, len(x))
		for _, c := range x {
			if s, ok := c.(string); ok {
				colors = append(colors, s)
			}
		}
		return OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: colors}, nil
	case map[string]any:
		if path, ok := x["datum"].(string); ok {
			return OrdinalColorScaleConfig{Type: OrdinalTypeDatum, DatumPath: path}, nil
		}
		if scheme, ok := x["scheme"].(string); ok {
			cfg := OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: scheme}
			if sz, ok := x["size"]; ok {
				if n, ok := toInt(sz); ok {
					cfg.SchemeSize = n
				}
			}
			return cfg, nil
		}
		if rawColors, ok := x["colors"].([]any); ok {
			colors := make([]string, 0, len(rawColors))
			for _, c := range rawColors {
				if s, ok := c.(string); ok {
					colors = append(colors, s)
				}
			}
			return OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: colors}, nil
		}
		return OrdinalColorScaleConfig{}, fmt.Errorf("invalid ordinal color config: missing 'scheme', 'datum', or 'colors'")
	}
	return OrdinalColorScaleConfig{}, fmt.Errorf("invalid ordinal color config type: %T", v)
}

// OrdinalColorScale is the resolved color generator: given a datum, returns
// its color string.
type OrdinalColorScale[D any] func(d D) string

// GetOrdinalColorScale builds an OrdinalColorScale from a config + identity
// accessor. Mirrors @nivo/colors getOrdinalColorScale.
//
// identity is either a func(D) string or a string path resolved via
// getPathValue. If identity is nil/empty, the datum itself (as string) is
// used.
func GetOrdinalColorScale[D any](config OrdinalColorScaleConfig, identity any) OrdinalColorScale[D] {
	switch config.Type {
	case OrdinalTypeFunc:
		return func(d D) string { return config.Func(d) }
	case OrdinalTypeStatic:
		return func(D) string { return config.Static }
	case OrdinalTypeDatum:
		return func(d D) string { return getPathValue(d, config.DatumPath) }
	case OrdinalTypeColors:
		getID := resolveIdentity[D](identity)
		scale := newOrdinalScale(config.Colors)
		return func(d D) string { return scale(getID(d)) }
	case OrdinalTypeScheme:
		getID := resolveIdentity[D](identity)
		colors := resolveSchemeColors(config.Scheme, config.SchemeSize)
		scale := newOrdinalScale(colors)
		return func(d D) string { return scale(getID(d)) }
	}
	return func(D) string { return "" }
}

// resolveIdentity builds a func(D) string from identity (func or string path).
func resolveIdentity[D any](identity any) func(D) string {
	switch id := identity.(type) {
	case func(D) string:
		return id
	case func(any) string:
		return func(d D) string { return id(d) }
	case string:
		if id == "" {
			return func(d D) string { return fmt.Sprintf("%v", d) }
		}
		return func(d D) string { return getPathValue(d, id) }
	case nil:
		return func(d D) string { return fmt.Sprintf("%v", d) }
	}
	return func(d D) string { return fmt.Sprintf("%v", d) }
}

// resolveSchemeColors returns the []string for a scheme id, honoring size for
// diverging/sequential schemes.
func resolveSchemeColors(scheme string, size int) []string {
	if IsCategoricalColorScheme(scheme) {
		return CategoricalColorSchemes[scheme]
	}
	if IsDivergingColorScheme(scheme) {
		m := DivergingColorSchemes[scheme]
		if size < 3 || size > 11 {
			size = 11
		}
		return m[size]
	}
	if IsSequentialColorScheme(scheme) {
		m := SequentialColorSchemes[scheme]
		if size < 3 || size > 9 {
			size = 9
		}
		return m[size]
	}
	return CategoricalColorSchemes["nivo"]
}

// newOrdinalScale returns a func(key string) string that cycles through
// colors by key (d3-scale scaleOrdinal semantics). Unknown keys are assigned
// the next color in sequence and cached.
func newOrdinalScale(colors []string) func(string) string {
	if len(colors) == 0 {
		colors = []string{""}
	}
	idx := 0
	assigned := map[string]string{}
	return func(key string) string {
		if c, ok := assigned[key]; ok {
			return c
		}
		c := colors[idx%len(colors)]
		idx++
		assigned[key] = c
		return c
	}
}

func toInt(v any) (int, bool) {
	switch x := v.(type) {
	case int:
		return x, true
	case int64:
		return int(x), true
	case float64:
		return int(x), true
	}
	return 0, false
}

// --- Continuous color scales -----------------------------------------------

// SequentialColorScaleConfig mirrors @nivo/colors SequentialColorScaleConfig.
type SequentialColorScaleConfig struct {
	Type         string // "sequential"
	MinValue     *float64
	MaxValue     *float64
	Scheme       string // interpolator id
	Colors       [2]string
	Interpolator func(float64) string
}

// SequentialColorScaleValues holds the min/max the scale is bound to.
type SequentialColorScaleValues struct {
	Min, Max float64
}

// SequentialColorScaleDefaults mirrors nivo's sequentialColorScaleDefaults.
var SequentialColorScaleDefaults = struct {
	Scheme string
}{Scheme: "turbo"}

// GetSequentialColorScale builds a func(value float64) string for a
// sequential config + values. Mirrors @nivo/colors getSequentialColorScale.
func GetSequentialColorScale(config SequentialColorScaleConfig, values SequentialColorScaleValues) func(float64) string {
	min := values.Min
	if config.MinValue != nil {
		min = *config.MinValue
	}
	max := values.Max
	if config.MaxValue != nil {
		max = *config.MaxValue
	}
	var interp func(float64) string
	if config.Colors[0] != "" && config.Colors[1] != "" {
		interp = interpolateRgbBasis([]string{config.Colors[0], config.Colors[1]})
	} else if config.Interpolator != nil {
		interp = config.Interpolator
	} else {
		scheme := config.Scheme
		if scheme == "" {
			scheme = SequentialColorScaleDefaults.Scheme
		}
		interp = ColorInterpolators[scheme]
	}
	if interp == nil {
		interp = func(float64) string { return "#000" }
	}
	domain := min
	span := max - min
	if span == 0 {
		span = 1
	}
	return func(v float64) string {
		t := (v - domain) / span
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		return interp(t)
	}
}

// DivergingColorScaleConfig mirrors @nivo/colors DivergingColorScaleConfig.
type DivergingColorScaleConfig struct {
	Type         string // "diverging"
	MinValue     *float64
	MaxValue     *float64
	DivergeAt    *float64
	Scheme       string
	Colors       [3]string
	Interpolator func(float64) string
}

// DivergingColorScaleDefaults mirrors nivo's divergingColorScaleDefaults.
var DivergingColorScaleDefaults = struct {
	Scheme    string
	DivergeAt float64
}{Scheme: "red_yellow_blue", DivergeAt: 0.5}

// GetDivergingColorScale builds a func(value float64) string for a diverging
// config + values. Mirrors @nivo/colors getDivergingColorScale.
func GetDivergingColorScale(config DivergingColorScaleConfig, values SequentialColorScaleValues) func(float64) string {
	min := values.Min
	if config.MinValue != nil {
		min = *config.MinValue
	}
	max := values.Max
	if config.MaxValue != nil {
		max = *config.MaxValue
	}
	divergeAt := DivergingColorScaleDefaults.DivergeAt
	if config.DivergeAt != nil {
		divergeAt = *config.DivergeAt
	}
	offset := 0.5 - divergeAt

	var interp func(float64) string
	if config.Colors[0] != "" && config.Colors[1] != "" && config.Colors[2] != "" {
		interp = interpolateRgbBasis([]string{config.Colors[0], config.Colors[1], config.Colors[2]})
	} else if config.Interpolator != nil {
		interp = config.Interpolator
	} else {
		scheme := config.Scheme
		if scheme == "" {
			scheme = DivergingColorScaleDefaults.Scheme
		}
		interp = ColorInterpolators[scheme]
	}
	if interp == nil {
		interp = func(float64) string { return "#000" }
	}
	span := max - min
	if span == 0 {
		span = 1
	}
	return func(v float64) string {
		t := (v - min) / span
		if t < 0 {
			t = 0
		}
		if t > 1 {
			t = 1
		}
		return interp(t + offset)
	}
}

// QuantizeColorScaleConfig mirrors @nivo/colors QuantizeColorScaleConfig.
type QuantizeColorScaleConfig struct {
	Type   string // "quantize"
	Domain [2]float64
	Scheme string
	Steps  int
	Colors []string
}

// QuantizeColorScaleDefaults mirrors nivo's quantizeColorScaleDefaults.
var QuantizeColorScaleDefaults = struct {
	Scheme string
	Steps  int
}{Scheme: "turbo", Steps: 7}

// GetQuantizeColorScale builds a func(value float64) string for a quantize
// config + values. Mirrors @nivo/colors getQuantizeColorScale.
func GetQuantizeColorScale(config QuantizeColorScaleConfig, values SequentialColorScaleValues) func(float64) string {
	domain := config.Domain
	if domain[0] == 0 && domain[1] == 0 {
		domain = [2]float64{values.Min, values.Max}
	}
	// nivo calls .nice() on the domain; we approximate with d3array ticks.
	nice := niceDomain(domain[0], domain[1])

	var colors []string
	if len(config.Colors) > 0 {
		colors = config.Colors
	} else {
		scheme := config.Scheme
		if scheme == "" {
			scheme = QuantizeColorScaleDefaults.Scheme
		}
		steps := config.Steps
		if steps <= 0 {
			steps = QuantizeColorScaleDefaults.Steps
		}
		interp := ColorInterpolators[scheme]
		if interp == nil {
			interp = ColorInterpolators["turbo"]
		}
		colors = make([]string, steps)
		for i := 0; i < steps; i++ {
			t := float64(i) / float64(steps-1)
			if steps == 1 {
				t = 0
			}
			colors[i] = interp(t)
		}
	}
	if len(colors) == 0 {
		colors = []string{"#000"}
	}
	lo, hi := nice[0], nice[1]
	span := hi - lo
	if span == 0 {
		span = 1
	}
	n := len(colors)
	return func(v float64) string {
		if v <= lo {
			return colors[0]
		}
		if v >= hi {
			return colors[n-1]
		}
		idx := int(math.Floor((v - lo) / span * float64(n)))
		if idx >= n {
			idx = n - 1
		}
		if idx < 0 {
			idx = 0
		}
		return colors[idx]
	}
}

// niceDomain approximates d3-scale's .nice() for a [min,max] domain.
func niceDomain(min, max float64) [2]float64 {
	if min == max {
		return [2]float64{min, max}
	}
	ticks := d3array.Ticks(min, max, 10)
	if len(ticks) == 0 {
		return [2]float64{min, max}
	}
	return [2]float64{ticks[0], ticks[len(ticks)-1]}
}

// ContinuousColorScaleConfig is the union of sequential/diverging/quantize.
type ContinuousColorScaleConfig struct {
	Type       string // "sequential" | "diverging" | "quantize"
	Sequential *SequentialColorScaleConfig
	Diverging  *DivergingColorScaleConfig
	Quantize   *QuantizeColorScaleConfig
}

// GetContinuousColorScale dispatches to the right builder based on Type.
func GetContinuousColorScale(config ContinuousColorScaleConfig, values SequentialColorScaleValues) func(float64) string {
	switch config.Type {
	case "sequential":
		if config.Sequential != nil {
			return GetSequentialColorScale(*config.Sequential, values)
		}
	case "diverging":
		if config.Diverging != nil {
			return GetDivergingColorScale(*config.Diverging, values)
		}
	case "quantize":
		if config.Quantize != nil {
			return GetQuantizeColorScale(*config.Quantize, values)
		}
	}
	return func(float64) string { return "#000" }
}

// IsContinuousColorScale reports whether a config map's "type" indicates a
// continuous color scale. Mirrors @nivo/colors isContinuousColorScale.
func IsContinuousColorScale(config any) bool {
	if m, ok := config.(map[string]any); ok {
		t, _ := m["type"].(string)
		return t == "sequential" || t == "diverging" || t == "quantize"
	}
	return false
}
