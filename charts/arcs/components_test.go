package arcs_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/arcs"
)

func renderShape(t *testing.T, props arcs.ArcShapeProps) string {
	t.Helper()
	var b strings.Builder
	if err := arcs.ArcShape(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("ArcShape.Render: %v", err)
	}
	return b.String()
}

func TestArcShape_Minimal(t *testing.T) {
	out := renderShape(t, arcs.ArcShapeProps{Path: "M0,0L10,0Z", Fill: "#ff0000"})
	if !strings.Contains(out, `d="M0,0L10,0Z"`) {
		t.Errorf("missing path data, got %q", out)
	}
	if !strings.Contains(out, `fill="#ff0000"`) {
		t.Errorf("missing fill, got %q", out)
	}
	for _, attr := range []string{"stroke", "opacity", "aria-label", "data-tc-tooltip", "hx-get", "hx-trigger", "hx-swap", "hx-target", "<animate"} {
		if strings.Contains(out, attr) {
			t.Errorf("minimal arc must not emit %q, got %q", attr, out)
		}
	}
}

func TestArcShape_AllOptionalAttributes(t *testing.T) {
	out := renderShape(t, arcs.ArcShapeProps{
		Path:        "M0,0",
		Fill:        "#00ff00",
		Stroke:      "#000000",
		StrokeWidth: 1.5,
		Opacity:     0.75,
		HxGet:       "/hover?arc=0",
		HxTrigger:   "mouseenter",
		HxSwap:      "innerHTML",
		HxTarget:    "#tooltip",
		AriaLabel:   "slice A",
		DataTooltip: "<b>A</b>",
	})
	want := []string{
		`stroke="#000000"`,
		`stroke-width="1.5"`,
		`opacity="0.75"`,
		`aria-label="slice A"`,
		`data-tc-tooltip="&lt;b&gt;A&lt;/b&gt;"`,
		`hx-get="/hover?arc=0"`,
		`hx-trigger="mouseenter"`,
		`hx-swap="innerHTML"`,
		`hx-target="#tooltip"`,
	}
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("missing %q in %q", w, out)
		}
	}
}

func TestArcShape_StrokeRequiresBothFields(t *testing.T) {
	// StrokeWidth > 0 but empty Stroke → no stroke attributes.
	out := renderShape(t, arcs.ArcShapeProps{Path: "M0,0", Fill: "#fff", StrokeWidth: 2})
	if strings.Contains(out, "stroke") {
		t.Errorf("stroke-width without stroke color must not emit stroke attrs, got %q", out)
	}
	// Stroke set but StrokeWidth == 0 → no stroke attributes.
	out = renderShape(t, arcs.ArcShapeProps{Path: "M0,0", Fill: "#fff", Stroke: "#000"})
	if strings.Contains(out, "stroke") {
		t.Errorf("stroke color without width must not emit stroke attrs, got %q", out)
	}
}

func TestArcShape_OpacityBounds(t *testing.T) {
	// Opacity 0 and 1 are both outside the (0,1) open interval → no attribute.
	for _, o := range []float64{0, 1} {
		out := renderShape(t, arcs.ArcShapeProps{Path: "M0,0", Fill: "#fff", Opacity: o})
		if strings.Contains(out, "opacity") {
			t.Errorf("opacity=%v must not emit opacity attr, got %q", o, out)
		}
	}
	out := renderShape(t, arcs.ArcShapeProps{Path: "M0,0", Fill: "#fff", Opacity: 0.5})
	if !strings.Contains(out, `opacity="0.5"`) {
		t.Errorf("opacity=0.5 should emit opacity attr, got %q", out)
	}
}

func TestArcShape_AnimateFromPath(t *testing.T) {
	out := renderShape(t, arcs.ArcShapeProps{
		Path:            "M10,10",
		Fill:            "#fff",
		Animate:         true,
		AnimateFromPath: "M0,0",
	})
	if !strings.Contains(out, `<animate attributeName="d" from="M0,0" to="M10,10" begin="0s" dur="0.6s" fill="freeze">`) {
		t.Errorf("expected d-path grow animation, got %q", out)
	}
}

func TestArcShape_AnimateFadeIn(t *testing.T) {
	out := renderShape(t, arcs.ArcShapeProps{
		Path:         "M10,10",
		Fill:         "#fff",
		Animate:      true,
		AnimateBegin: "0.2s",
	})
	if !strings.Contains(out, `<animate attributeName="opacity" from="0" to="1" begin="0.2s" dur="0.6s" fill="freeze">`) {
		t.Errorf("expected opacity fade-in animation, got %q", out)
	}
	if strings.Contains(out, `attributeName="d"`) {
		t.Errorf("fade-in variant must not animate the d attribute, got %q", out)
	}
}

func TestArcShape_NoAnimateByDefault(t *testing.T) {
	// Animate=false with AnimateFromPath set must still emit nothing.
	out := renderShape(t, arcs.ArcShapeProps{Path: "M10,10", Fill: "#fff", AnimateFromPath: "M0,0"})
	if strings.Contains(out, "<animate") {
		t.Errorf("non-animated arc must not emit <animate>, got %q", out)
	}
}

func TestArcsLayer(t *testing.T) {
	cases := []struct {
		name      string
		props     arcs.ArcsLayerProps
		wantPaths int
	}{
		{
			name:      "empty",
			props:     arcs.ArcsLayerProps{CenterX: 100, CenterY: 50.5},
			wantPaths: 0,
		},
		{
			name: "two-arcs",
			props: arcs.ArcsLayerProps{
				CenterX: 100, CenterY: 50.5,
				Items: []arcs.ArcLayerItem{
					{Props: arcs.ArcShapeProps{Path: "M0,0", Fill: "#111111"}},
					{Props: arcs.ArcShapeProps{Path: "M1,1", Fill: "#222222"}},
				},
			},
			wantPaths: 2,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var b strings.Builder
			if err := arcs.ArcsLayer(c.props).Render(context.Background(), &b); err != nil {
				t.Fatalf("ArcsLayer.Render: %v", err)
			}
			out := b.String()
			if !strings.Contains(out, `<g transform="translate(100,50.5)">`) {
				t.Errorf("missing centered <g>, got %q", out)
			}
			if got := strings.Count(out, "<path"); got != c.wantPaths {
				t.Errorf("path count = %d, want %d", got, c.wantPaths)
			}
			if c.wantPaths == 2 {
				if !strings.Contains(out, `fill="#111111"`) || !strings.Contains(out, `fill="#222222"`) {
					t.Errorf("missing per-arc fills, got %q", out)
				}
			}
		})
	}
}

func TestArcLabel_FullStyling(t *testing.T) {
	var b strings.Builder
	err := arcs.ArcLabel(arcs.ArcLabelProps{
		X: 12.345678, Y: -3.5,
		Label: "42%", Fill: "#333", FontSize: 11, FontFamily: "monospace",
	}).Render(context.Background(), &b)
	if err != nil {
		t.Fatalf("ArcLabel.Render: %v", err)
	}
	out := b.String()
	want := []string{
		`transform="translate(12.346,-3.5)"`, // fmtF rounds to 3 dp
		`text-anchor="middle"`,
		`dominant-baseline="central"`,
		`fill="#333"`,
		`font-size="11"`,
		`font-family="monospace"`,
		`>42%</text>`,
		`style="pointer-events: none"`,
	}
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("missing %q in %q", w, out)
		}
	}
}

func TestArcLabel_DefaultStyling(t *testing.T) {
	var b strings.Builder
	err := arcs.ArcLabel(arcs.ArcLabelProps{X: 0, Y: 0, Label: "x"}).Render(context.Background(), &b)
	if err != nil {
		t.Fatalf("ArcLabel.Render: %v", err)
	}
	out := b.String()
	for _, attr := range []string{"fill=", "font-size=", "font-family="} {
		if strings.Contains(out, attr) {
			t.Errorf("default label must not emit %q, got %q", attr, out)
		}
	}
	if !strings.Contains(out, `transform="translate(0,0)"`) {
		t.Errorf("missing origin transform, got %q", out)
	}
}

func TestArcLabelsLayer(t *testing.T) {
	props := arcs.ArcLabelsLayerProps{
		CenterX: 200, CenterY: 150,
		Items: []arcs.ArcLabelItem{
			{Props: arcs.ArcLabelProps{X: 1, Y: 2, Label: "a"}},
			{Props: arcs.ArcLabelProps{X: 3, Y: 4, Label: "b"}},
			{Props: arcs.ArcLabelProps{X: 5, Y: 6, Label: "c"}},
		},
	}
	var b strings.Builder
	if err := arcs.ArcLabelsLayer(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("ArcLabelsLayer.Render: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, `<g transform="translate(200,150)" style="pointer-events: none">`) {
		t.Errorf("missing layer <g>, got %q", out)
	}
	if got := strings.Count(out, "<text"); got != 3 {
		t.Errorf("text count = %d, want 3", got)
	}
	for _, l := range []string{">a</text>", ">b</text>", ">c</text>"} {
		if !strings.Contains(out, l) {
			t.Errorf("missing label %q", l)
		}
	}
}

func TestArcLabelsLayer_Empty(t *testing.T) {
	var b strings.Builder
	if err := arcs.ArcLabelsLayer(arcs.ArcLabelsLayerProps{CenterX: 1, CenterY: 2}).Render(context.Background(), &b); err != nil {
		t.Fatalf("ArcLabelsLayer.Render: %v", err)
	}
	if strings.Contains(b.String(), "<text") {
		t.Errorf("empty layer must not emit <text>, got %q", b.String())
	}
}

func TestArcLinkLabel_RightSide(t *testing.T) {
	props := arcs.ArcLinkLabelProps{
		Link: arcs.ArcLink{
			P0x: 10, P0y: -10,
			P1x: 20, P1y: -20,
			P2x: 44, P2y: -20,
			TextAnchor: "start", Side: "right",
		},
		Label: "Alpha", Thickness: 2, Color: "#999",
		TextOffset: 6,
	}
	var b strings.Builder
	if err := arcs.ArcLinkLabel(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("ArcLinkLabel.Render: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, `d="M10,-10 L20,-20 L44,-20"`) {
		t.Errorf("missing polyline path, got %q", out)
	}
	if !strings.Contains(out, `stroke="#999"`) || !strings.Contains(out, `stroke-width="2"`) || !strings.Contains(out, `fill="none"`) {
		t.Errorf("missing link stroke attrs, got %q", out)
	}
	// Right side: textX = P2x + offset = 50.
	if !strings.Contains(out, `x="50"`) || !strings.Contains(out, `y="-20"`) {
		t.Errorf("text should be at (50,-20), got %q", out)
	}
	if !strings.Contains(out, `text-anchor="start"`) || !strings.Contains(out, ">Alpha</text>") {
		t.Errorf("missing anchored label text, got %q", out)
	}
	// No optional font styling supplied.
	for _, attr := range []string{"font-size=", "font-family="} {
		if strings.Contains(out, attr) {
			t.Errorf("unstyled link label must not emit %q, got %q", attr, out)
		}
	}
}

func TestArcLinkLabel_LeftSideWithStyling(t *testing.T) {
	props := arcs.ArcLinkLabelProps{
		Link: arcs.ArcLink{
			P0x: -10, P0y: 10,
			P1x: -20, P1y: 20,
			P2x: -44, P2y: 20,
			TextAnchor: "end", Side: "left",
		},
		Label: "Beta", Thickness: 1, Color: "#000",
		TextOffset: 6, Fill: "#123456", FontSize: 10, FontFamily: "serif",
	}
	var b strings.Builder
	if err := arcs.ArcLinkLabel(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("ArcLinkLabel.Render: %v", err)
	}
	out := b.String()
	// Left side: textX = P2x - offset = -50.
	if !strings.Contains(out, `x="-50"`) {
		t.Errorf("left-side text should be at x=-50, got %q", out)
	}
	want := []string{`text-anchor="end"`, `fill="#123456"`, `font-size="10"`, `font-family="serif"`, ">Beta</text>"}
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("missing %q in %q", w, out)
		}
	}
}

func TestArcLinkLabelsLayer(t *testing.T) {
	props := arcs.ArcLinkLabelsLayerProps{
		CenterX: 250, CenterY: 175,
		Items: []arcs.ArcLinkLabelItem{
			{Props: arcs.ArcLinkLabelProps{Label: "one", Link: arcs.ArcLink{TextAnchor: "start", Side: "right"}}},
			{Props: arcs.ArcLinkLabelProps{Label: "two", Link: arcs.ArcLink{TextAnchor: "end", Side: "left"}}},
		},
	}
	var b strings.Builder
	if err := arcs.ArcLinkLabelsLayer(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("ArcLinkLabelsLayer.Render: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, `<g transform="translate(250,175)" style="pointer-events: none">`) {
		t.Errorf("missing layer <g>, got %q", out)
	}
	if got := strings.Count(out, "<path"); got != 2 {
		t.Errorf("link path count = %d, want 2", got)
	}
	if got := strings.Count(out, "<text"); got != 2 {
		t.Errorf("link text count = %d, want 2", got)
	}
}

func TestArcLinkLabelsLayer_Empty(t *testing.T) {
	var b strings.Builder
	if err := arcs.ArcLinkLabelsLayer(arcs.ArcLinkLabelsLayerProps{}).Render(context.Background(), &b); err != nil {
		t.Fatalf("ArcLinkLabelsLayer.Render: %v", err)
	}
	out := b.String()
	if strings.Contains(out, "<path") || strings.Contains(out, "<text") {
		t.Errorf("empty layer must emit no paths/texts, got %q", out)
	}
	if !strings.Contains(out, `<g transform="translate(0,0)"`) {
		t.Errorf("missing origin <g>, got %q", out)
	}
}
