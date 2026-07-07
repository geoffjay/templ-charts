package core

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func renderComponent(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatalf("render: %v", err)
	}
	return b.String()
}

// allComponents returns one instance of every templ component in the package,
// with props that exercise their main render path.
func allComponents() map[string]templ.Component {
	marker := CartesianMarker{Axis: "y", Value: 1, Legend: "l"}
	scale := func(any) float64 { return 10 }
	return map[string]templ.Component{
		"SvgWrapper":            SvgWrapper(SvgWrapperProps{Width: 10, Height: 10}, "<g></g>"),
		"Defs":                  Defs([]Def{PatternDotsDef("d", nil)}),
		"defOne":                defOne(PatternDotsDef("d", nil)),
		"linearGradientDef":     linearGradientDef(LinearGradientDef("g", []GradientStop{{Offset: 0, Color: "#f00"}}, nil)),
		"patternDotsDef":        patternDotsDef(PatternDotsDef("d", nil)),
		"patternLinesDef":       patternLinesDef(PatternLinesDef("l", nil)),
		"patternSquaresDef":     patternSquaresDef(PatternSquaresDef("s", nil)),
		"DotsItem":              DotsItem(DotsItemProps{Size: 4}),
		"CartesianMarkers":      CartesianMarkers(CartesianMarkersProps{Markers: []CartesianMarker{marker}, Width: 100, Height: 80, XScale: scale, YScale: scale}),
		"cartesianMarkersItem":  cartesianMarkersItem(marker, 100, 80, scale, scale),
		"cartesianMarkerLegend": cartesianMarkerLegend(marker, 100, 80),
	}
}

func TestDefOne_Dispatch(t *testing.T) {
	cases := []struct {
		def  Def
		want string
	}{
		{LinearGradientDef("g", []GradientStop{{Offset: 0, Color: "#f00"}}, nil), "<linearGradient"},
		{PatternDotsDef("d", nil), "<pattern"},
		{PatternLinesDef("l", nil), "<pattern"},
		{PatternSquaresDef("s", nil), "<pattern"},
	}
	for _, c := range cases {
		out := renderComponent(t, defOne(c.def))
		if !strings.Contains(out, c.want) {
			t.Errorf("defOne(%s) missing %q in %q", c.def.Type, c.want, out)
		}
	}
	if out := renderComponent(t, defOne(Def{ID: "u", Type: "unknown"})); out != "" {
		t.Errorf("defOne(unknown) = %q, want empty", out)
	}
}

func TestChildComponents_DirectRender(t *testing.T) {
	out := renderComponent(t, cartesianMarkersItem(
		CartesianMarker{Axis: "x", Value: 1, Legend: "cap"},
		100, 80,
		func(any) float64 { return 25 }, nil,
	))
	if !strings.Contains(out, `transform="translate(25,0)"`) || !strings.Contains(out, ">cap</text>") {
		t.Errorf("cartesianMarkersItem output: %q", out)
	}

	leg := renderComponent(t, cartesianMarkerLegend(CartesianMarker{Axis: "y", Legend: "cap"}, 100, 80))
	if !strings.Contains(leg, `dominant-baseline="central"`) || !strings.Contains(leg, ">cap</text>") {
		t.Errorf("cartesianMarkerLegend output: %q", leg)
	}

	grad := renderComponent(t, linearGradientDef(LinearGradientDef("g", []GradientStop{{Offset: 50, Color: "#0f0", Opacity: 0.25}}, nil)))
	if !strings.Contains(grad, `offset="50%"`) || !strings.Contains(grad, `stop-opacity="0.25"`) {
		t.Errorf("linearGradientDef output: %q", grad)
	}
}

func TestComponents_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for name, c := range allComponents() {
		var b strings.Builder
		if err := c.Render(ctx, &b); !errors.Is(err, context.Canceled) {
			t.Errorf("%s: Render with canceled ctx = %v, want context.Canceled", name, err)
		}
	}
}

func TestComponents_NilChildrenContext(t *testing.T) {
	// A nil children value in the context must be tolerated (replaced by the
	// runtime's no-op component).
	ctx := templ.WithChildren(context.Background(), nil)
	for name, c := range allComponents() {
		var b strings.Builder
		if err := c.Render(ctx, &b); err != nil {
			t.Errorf("%s: Render with nil children = %v", name, err)
		}
	}
}
