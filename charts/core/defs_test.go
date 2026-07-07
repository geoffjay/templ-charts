package core_test

import (
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
)

func TestDefFactories(t *testing.T) {
	stops := []core.GradientStop{{Offset: 0, Color: "#f00"}, {Offset: 100, Color: "#00f"}}
	g := core.LinearGradientDef("g", stops, map[string]any{"x1": 1})
	if g.ID != "g" || g.Type != core.DefTypeLinearGradient || len(g.Colors) != 2 || g.Raw["x1"] != 1 {
		t.Errorf("LinearGradientDef = %+v", g)
	}
	d := core.PatternDotsDef("d", nil)
	if d.ID != "d" || d.Type != core.DefTypePatternDots {
		t.Errorf("PatternDotsDef = %+v", d)
	}
	l := core.PatternLinesDef("l", nil)
	if l.ID != "l" || l.Type != core.DefTypePatternLines {
		t.Errorf("PatternLinesDef = %+v", l)
	}
	s := core.PatternSquaresDef("s", nil)
	if s.ID != "s" || s.Type != core.DefTypePatternSquares {
		t.Errorf("PatternSquaresDef = %+v", s)
	}
}

func TestBindDefs_EmptyDefsOrNodes(t *testing.T) {
	out := core.BindDefs(nil, []any{map[string]any{}}, nil, "", "", "")
	if len(out.Defs) != 0 || len(out.FillByNodeIndex) != 0 {
		t.Errorf("empty defs: %+v", out)
	}
	defs := []core.Def{core.PatternDotsDef("d", nil)}
	out = core.BindDefs(defs, nil, []core.DefRule{{ID: "d", Match: "*"}}, "", "", "")
	if len(out.Defs) != 1 || out.Defs[0].ID != "d" || len(out.FillByNodeIndex) != 0 {
		t.Errorf("empty nodes: %+v", out)
	}
}

func TestBindDefs_WildcardPatternNoInherit(t *testing.T) {
	defs := []core.Def{core.PatternDotsDef("dots", nil)}
	nodes := []any{
		map[string]any{"color": "#f00"},
		map[string]any{"color": "#0f0"},
	}
	out := core.BindDefs(defs, nodes, []core.DefRule{{ID: "dots", Match: "*"}}, "", "", "")
	if len(out.Defs) != 1 {
		t.Errorf("no-inherit pattern should not generate variants: %d defs", len(out.Defs))
	}
	for i := range nodes {
		if got := out.FillByNodeIndex[i]; got != "url(#dots)" {
			t.Errorf("node %d fill = %q, want url(#dots)", i, got)
		}
	}
}

func TestBindDefs_PatternInheritBackground(t *testing.T) {
	defs := []core.Def{{ID: "lines", Type: core.DefTypePatternLines, Background: "inherit", Color: "#fff"}}
	nodes := []any{
		map[string]any{"color": "#f00"},
		map[string]any{"color": "#f00"}, // same color → variant deduped
		map[string]any{"color": "#0f0"},
	}
	out := core.BindDefs(defs, nodes, []core.DefRule{{ID: "lines", Match: "*"}}, "", "", "")
	if got := out.FillByNodeIndex[0]; got != "url(#lines.bg.#f00)" {
		t.Errorf("node 0 fill = %q", got)
	}
	if got := out.FillByNodeIndex[2]; got != "url(#lines.bg.#0f0)" {
		t.Errorf("node 2 fill = %q", got)
	}
	if len(out.Defs) != 3 {
		t.Fatalf("defs = %d, want 3 (base + 2 color variants)", len(out.Defs))
	}
	v := out.Defs[1]
	if v.ID != "lines.bg.#f00" || v.Background != "#f00" || v.Color != "#fff" {
		t.Errorf("variant = %+v", v)
	}
}

func TestBindDefs_PatternInheritBoth(t *testing.T) {
	defs := []core.Def{{ID: "sq", Type: core.DefTypePatternSquares, Background: "inherit", Color: "inherit"}}
	nodes := []any{map[string]any{"color": "#123"}}
	out := core.BindDefs(defs, nodes, []core.DefRule{{ID: "sq", Match: "*"}}, "", "", "")
	want := "url(#sq.bg.#123.fg.#123)"
	if got := out.FillByNodeIndex[0]; got != want {
		t.Errorf("fill = %q, want %q", got, want)
	}
	if len(out.Defs) != 2 || out.Defs[1].Background != "#123" || out.Defs[1].Color != "#123" {
		t.Errorf("variant defs = %+v", out.Defs)
	}
}

func TestBindDefs_GradientInherit(t *testing.T) {
	defs := []core.Def{core.LinearGradientDef("grad", []core.GradientStop{
		{Offset: 0, Color: "inherit", Opacity: 1},
		{Offset: 100, Color: "#00f"},
	}, nil)}
	nodes := []any{
		map[string]any{"color": "#abc"},
		map[string]any{"color": "#abc"}, // deduped
	}
	out := core.BindDefs(defs, nodes, []core.DefRule{{ID: "grad", Match: "*"}}, "", "", "")
	want := "url(#grad.0.#abc)"
	if got := out.FillByNodeIndex[0]; got != want {
		t.Errorf("fill = %q, want %q", got, want)
	}
	if len(out.Defs) != 2 {
		t.Fatalf("defs = %d, want 2", len(out.Defs))
	}
	v := out.Defs[1]
	if v.Colors[0].Color != "#abc" || v.Colors[1].Color != "#00f" {
		t.Errorf("variant stops = %+v", v.Colors)
	}
}

func TestBindDefs_GradientNoInherit(t *testing.T) {
	defs := []core.Def{core.LinearGradientDef("grad", []core.GradientStop{{Offset: 0, Color: "#f00"}}, nil)}
	nodes := []any{map[string]any{"color": "#abc"}}
	out := core.BindDefs(defs, nodes, []core.DefRule{{ID: "grad", Match: "*"}}, "", "", "")
	if got := out.FillByNodeIndex[0]; got != "url(#grad)" {
		t.Errorf("fill = %q, want url(#grad)", got)
	}
	if len(out.Defs) != 1 {
		t.Errorf("defs = %d, want 1", len(out.Defs))
	}
}

func TestBindDefs_MatchVariants(t *testing.T) {
	defs := []core.Def{core.PatternDotsDef("d", nil)}
	node := map[string]any{"color": "#f00", "value": 5.0}

	// A func predicate.
	pred := func(n any) bool { return n.(map[string]any)["value"].(float64) > 3 }
	out := core.BindDefs(defs, []any{node}, []core.DefRule{{ID: "d", Match: pred}}, "", "", "")
	if out.FillByNodeIndex[0] != "url(#d)" {
		t.Errorf("func match fill = %q", out.FillByNodeIndex[0])
	}

	// A map matched against node data under dataKey.
	type wrapped struct{ Data map[string]any }
	nodes := []any{
		wrapped{Data: map[string]any{"kind": "a", "color": "#f00"}},
		wrapped{Data: map[string]any{"kind": "b", "color": "#0f0"}},
	}
	out = core.BindDefs(defs, nodes, []core.DefRule{{ID: "d", Match: map[string]any{"kind": "a"}}}, "Data", "", "")
	if out.FillByNodeIndex[0] != "url(#d)" {
		t.Errorf("map match fill = %q", out.FillByNodeIndex[0])
	}
	if _, ok := out.FillByNodeIndex[1]; ok {
		t.Errorf("map match should skip non-matching node")
	}

	// An empty key path compares the whole datum.
	out = core.BindDefs(defs, []any{"red-node"}, []core.DefRule{{ID: "d", Match: map[string]any{"": "red-node"}}}, "", "", "")
	if out.FillByNodeIndex[0] != "url(#d)" {
		t.Errorf("empty-path map match fill = %q", out.FillByNodeIndex[0])
	}

	// Non-wildcard strings, nil, and unsupported types never match.
	for _, m := range []any{"nope", nil, 42} {
		out = core.BindDefs(defs, []any{node}, []core.DefRule{{ID: "d", Match: m}}, "", "", "")
		if len(out.FillByNodeIndex) != 0 {
			t.Errorf("match %v should not apply", m)
		}
	}
}

func TestBindDefs_MissingRuleDef(t *testing.T) {
	defs := []core.Def{core.PatternDotsDef("d", nil)}
	out := core.BindDefs(defs, []any{map[string]any{}}, []core.DefRule{{ID: "ghost", Match: "*"}}, "", "", "")
	if len(out.FillByNodeIndex) != 0 || len(out.Defs) != 1 {
		t.Errorf("missing def id should be a no-op: %+v", out)
	}
}

func TestBindDefs_FirstMatchingRuleWins(t *testing.T) {
	defs := []core.Def{core.PatternDotsDef("a", nil), core.PatternDotsDef("b", nil)}
	rules := []core.DefRule{{ID: "a", Match: "*"}, {ID: "b", Match: "*"}}
	out := core.BindDefs(defs, []any{map[string]any{}}, rules, "", "", "")
	if got := out.FillByNodeIndex[0]; got != "url(#a)" {
		t.Errorf("fill = %q, want url(#a)", got)
	}
}

func TestBindDefs_CustomColorKey(t *testing.T) {
	defs := []core.Def{{ID: "p", Type: core.DefTypePatternDots, Color: "inherit"}}
	nodes := []any{map[string]any{"style": map[string]any{"stroke": "#777"}}}
	out := core.BindDefs(defs, nodes, []core.DefRule{{ID: "p", Match: "*"}}, "", "style.stroke", "")
	if got := out.FillByNodeIndex[0]; got != "url(#p.fg.#777)" {
		t.Errorf("fill = %q, want url(#p.fg.#777)", got)
	}
}
