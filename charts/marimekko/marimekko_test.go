package marimekko_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/marimekko"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func sampleDims() []marimekko.MarimekkoDimension {
	return []marimekko.MarimekkoDimension{
		{ID: "agree", Key: "agree"},
		{ID: "disagree", Key: "disagree"},
		{ID: "neutral", Key: "neutral"},
	}
}

func sampleData() []marimekko.MarimekkoDatum {
	return []marimekko.MarimekkoDatum{
		{ID: "France", Value: 40, Dimensions: map[string]float64{"agree": 20, "disagree": 12, "neutral": 8}},
		{ID: "Japan", Value: 25, Dimensions: map[string]float64{"agree": 10, "disagree": 10, "neutral": 5}},
		{ID: "US", Value: 60, Dimensions: map[string]float64{"agree": 35, "disagree": 15, "neutral": 10}},
	}
}

func baseProps() marimekko.MarimekkoProps {
	return marimekko.MarimekkoProps{
		Width: 600, Height: 400,
		Margin:     core.Margin{Top: 20, Right: 30, Bottom: 50, Left: 60},
		Data:       sampleData(),
		Dimensions: sampleDims(),
	}
}

func render(t *testing.T, props marimekko.MarimekkoProps) string {
	t.Helper()
	var b strings.Builder
	if err := marimekko.Marimekko(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Marimekko.Render: %v", err)
	}
	return b.String()
}

func TestMarimekko_RendersSVG(t *testing.T) {
	out := render(t, baseProps())
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("expected <svg…>, got %q", out[:min(50, len(out))])
	}
	if !strings.Contains(out, "</svg>") {
		t.Errorf("missing closing </svg>")
	}
}

func TestMarimekko_SegmentCount(t *testing.T) {
	// 3 columns × 3 dimensions = 9 segment rects.
	out := render(t, baseProps())
	if got := strings.Count(out, `<rect x=`); got != 9 {
		t.Errorf("segment rect count = %d, want 9", got)
	}
}

func TestMarimekko_VariableWidths(t *testing.T) {
	// US (value 60) bar should be wider than Japan (value 25). Compare the
	// width of the first segment of each column via the computed model.
	res := marimekko.UseMarimekko(func() marimekko.MarimekkoProps {
		p := baseProps()
		p.Width, p.Height = 510, 330 // inner dims
		return p
	}())
	var franceW, japanW, usW float64
	for _, b := range res.Bars {
		switch b.Index {
		case "France":
			franceW = b.Width
		case "Japan":
			japanW = b.Width
		case "US":
			usW = b.Width
		}
	}
	if !(usW > franceW && franceW > japanW) {
		t.Errorf("expected US > France > Japan widths, got US=%.2f France=%.2f Japan=%.2f", usW, franceW, japanW)
	}
}

func TestMarimekko_ExpandOffset(t *testing.T) {
	p := baseProps()
	p.Offset = marimekko.OffsetExpand
	out := render(t, p)
	if !strings.Contains(out, "<rect x=") {
		t.Errorf("expand offset should still render segments")
	}
}

func TestMarimekko_HorizontalLayout(t *testing.T) {
	p := baseProps()
	p.Layout = marimekko.MarimekkoLayoutHorizontal
	out := render(t, p)
	if !strings.HasPrefix(out, "<svg") {
		t.Fatalf("horizontal layout did not render")
	}
}

func TestMarimekko_Golden(t *testing.T) {
	out := render(t, baseProps())
	golden.Assert(t, "marimekko-basic", out)
}

func TestMarimekko_A11yTitleDesc(t *testing.T) {
	p := baseProps()
	p.Title = "Survey by country"
	p.Desc = "Agree/disagree/neutral split, bar width by sample size."
	out := render(t, p)
	if !strings.Contains(out, "<title>Survey by country</title>") {
		t.Errorf("expected <title>")
	}
	if !strings.Contains(out, "<desc>Agree/disagree/neutral split, bar width by sample size.</desc>") {
		t.Errorf("expected <desc>")
	}
}

func TestMarimekko_InteractiveEmitsTooltip(t *testing.T) {
	p := baseProps()
	p.Interactive = true
	if !strings.Contains(render(t, p), "data-tc-tooltip") {
		t.Errorf("interactive marimekko should emit data-tc-tooltip")
	}
	if strings.Contains(render(t, baseProps()), "data-tc-tooltip") {
		t.Errorf("non-interactive marimekko must not emit data-tc-tooltip")
	}
}

func TestMarimekko_Legend(t *testing.T) {
	p := baseProps()
	p.Legends = []legends.LegendProps{
		{Anchor: legends.LegendAnchorTop, Direction: legends.LegendDirectionRow},
	}
	out := render(t, p)
	if !strings.Contains(out, "agree") {
		t.Errorf("expected legend to include dimension ids")
	}
}
