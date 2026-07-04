package samples

import (
	"reflect"
	"testing"
)

// TestGeneratorsDeterministicAndNonEmpty exercises every sample generator: each
// must be deterministic (two calls produce equal data — the byte-stability
// guarantee) and return non-empty/non-zero data (so the data is actually
// renderable). Multi-value generators are wrapped so every returned value is
// checked.
func TestGeneratorsDeterministicAndNonEmpty(t *testing.T) {
	gens := map[string]func() []any{
		"Bar":                 func() []any { d, k := Bar(); return []any{d, k} },
		"Line":                func() []any { return []any{Line()} },
		"Pie":                 func() []any { return []any{Pie()} },
		"BoxPlot":             func() []any { return []any{BoxPlot()} },
		"Bullet":              func() []any { return []any{Bullet()} },
		"Bump":                func() []any { return []any{Bump()} },
		"Calendar":            func() []any { f, to, d := Calendar(); return []any{f, to, d} },
		"Chord":               func() []any { m, k := Chord(); return []any{m, k} },
		"CirclePacking":       func() []any { return []any{CirclePacking()} },
		"Funnel":              func() []any { return []any{Funnel()} },
		"Heatmap":             func() []any { return []any{Heatmap()} },
		"Icicle":              func() []any { return []any{Icicle()} },
		"Marimekko":           func() []any { d, dim := Marimekko(); return []any{d, dim} },
		"Network":             func() []any { n, l := Network(); return []any{n, l} },
		"ParallelCoordinates": func() []any { d, v := ParallelCoordinates(); return []any{d, v} },
		"PolarBar":            func() []any { d, k := PolarBar(); return []any{d, k} },
		"Radar":               func() []any { d, k := Radar(); return []any{d, k} },
		"RadialBar":           func() []any { return []any{RadialBar()} },
		"Sankey":              func() []any { n, l := Sankey(); return []any{n, l} },
		"ScatterPlot":         func() []any { return []any{ScatterPlot()} },
		"Stream":              func() []any { d, k := Stream(); return []any{d, k} },
		"Sunburst":            func() []any { return []any{Sunburst()} },
		"SwarmPlot":           func() []any { d, g := SwarmPlot(); return []any{d, g} },
		"Tree":                func() []any { return []any{Tree()} },
		"Treemap":             func() []any { return []any{Treemap()} },
		"Voronoi":             func() []any { return []any{Voronoi()} },
		"Waffle":              func() []any { return []any{Waffle()} },
		"Choropleth":          func() []any { f, d := Choropleth(); return []any{f, d} },
	}

	if len(gens) != 28 {
		t.Fatalf("expected 28 generators, listed %d", len(gens))
	}

	for name, gen := range gens {
		a, b := gen(), gen()
		if !reflect.DeepEqual(a, b) {
			t.Errorf("%s: generator is not deterministic across calls", name)
		}
		for i, v := range a {
			rv := reflect.ValueOf(v)
			switch rv.Kind() {
			case reflect.Slice, reflect.Map, reflect.Array:
				if rv.Len() == 0 {
					t.Errorf("%s: result #%d is empty", name, i)
				}
			default:
				if rv.IsZero() {
					t.Errorf("%s: result #%d is zero-valued", name, i)
				}
			}
		}
	}
}
