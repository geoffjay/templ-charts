package legends_test

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
	"github.com/geoffjay/templ-charts/charts/legends"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return b.String()
}

func renderLegendSvg(t *testing.T, props legends.LegendProps) string {
	t.Helper()
	var b strings.Builder
	if err := legends.LegendSvg(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("LegendSvg.Render: %v", err)
	}
	return b.String()
}

func renderBoxLegend(t *testing.T, props legends.BoxLegendSvgProps) string {
	t.Helper()
	var b strings.Builder
	if err := legends.BoxLegendSvg(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("BoxLegendSvg.Render: %v", err)
	}
	return b.String()
}

func twoItems() []legends.Datum {
	return []legends.Datum{
		{ID: "a", Label: "Alpha", Color: "#ff0000"},
		{ID: "b", Label: "Beta", Color: "#00ff00"},
	}
}

func TestBoxLegendSvg_AnchorTranslate(t *testing.T) {
	tests := []struct {
		name      string
		anchor    legends.LegendAnchor
		transform string
	}{
		{"top-left", legends.LegendAnchorTopLeft, `transform="translate(0,0)"`},
		{"top", legends.LegendAnchorTop, `transform="translate(50,0)"`},
		{"top-right", legends.LegendAnchorTopRight, `transform="translate(100,0)"`},
		{"right", legends.LegendAnchorRight, `transform="translate(100,30)"`},
		{"bottom-right", legends.LegendAnchorBottomRight, `transform="translate(100,60)"`},
		{"bottom", legends.LegendAnchorBottom, `transform="translate(50,60)"`},
		{"bottom-left", legends.LegendAnchorBottomLeft, `transform="translate(0,60)"`},
		{"left", legends.LegendAnchorLeft, `transform="translate(0,30)"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Legend box is 100x40 (2 column items of 100x20) in a 200x100 chart.
			out := renderBoxLegend(t, legends.BoxLegendSvgProps{
				Props: legends.LegendProps{
					Anchor:    tc.anchor,
					Direction: legends.LegendDirectionColumn,
					Items:     twoItems(),
					ItemWidth: 100, ItemHeight: 20,
				},
				ChartWidth: 200, ChartHeight: 100,
			})
			if !strings.Contains(out, tc.transform) {
				t.Errorf("missing %s in output: %s", tc.transform, out)
			}
		})
	}
}

func TestBoxLegendSvg_TranslateOffsets(t *testing.T) {
	out := renderBoxLegend(t, legends.BoxLegendSvgProps{
		Props: legends.LegendProps{
			Anchor:    legends.LegendAnchorTopLeft,
			Items:     twoItems(),
			ItemWidth: 100, ItemHeight: 20,
			TranslateX: 12.5, TranslateY: -7,
		},
		ChartWidth: 200, ChartHeight: 100,
	})
	if !strings.Contains(out, `transform="translate(12.5,-7)"`) {
		t.Errorf("translate offsets not applied: %s", out)
	}
}

func TestBoxLegendSvg_ContainsItems(t *testing.T) {
	out := renderBoxLegend(t, legends.BoxLegendSvgProps{
		Props:      legends.LegendProps{Items: twoItems(), ItemWidth: 100, ItemHeight: 20},
		ChartWidth: 200, ChartHeight: 100,
	})
	if !strings.Contains(out, ">Alpha</text>") || !strings.Contains(out, ">Beta</text>") {
		t.Errorf("item labels missing: %s", out)
	}
	// Outer <g> + one <g> per item.
	if got := strings.Count(out, "<g"); got != 3 {
		t.Errorf("<g> count = %d, want 3", got)
	}
}

func TestLegendSvg_EmptyItems(t *testing.T) {
	out := renderLegendSvg(t, legends.LegendProps{})
	if out != "" {
		t.Errorf("empty legend rendered %q, want empty output", out)
	}
}

func TestLegendSvg_ColumnLayout(t *testing.T) {
	out := renderLegendSvg(t, legends.LegendProps{
		Direction: legends.LegendDirectionColumn,
		Items:     twoItems(),
		ItemWidth: 100, ItemHeight: 20,
	})
	if !strings.Contains(out, `transform="translate(0,0)"`) || !strings.Contains(out, `transform="translate(0,20)"`) {
		t.Errorf("column item translates wrong: %s", out)
	}
}

func TestLegendSvg_RowLayout(t *testing.T) {
	out := renderLegendSvg(t, legends.LegendProps{
		Direction: legends.LegendDirectionRow,
		Items:     twoItems(),
		ItemWidth: 100, ItemHeight: 20,
	})
	if !strings.Contains(out, `transform="translate(0,0)"`) || !strings.Contains(out, `transform="translate(100,0)"`) {
		t.Errorf("row item translates wrong: %s", out)
	}
}

func TestLegendSvgItem_Toggle(t *testing.T) {
	out := render(t, legends.LegendSvgItem(legends.LegendSvgItemProps{
		Props: legends.LegendProps{Toggle: "post", ChartID: "c1"},
		Datum: legends.Datum{ID: "s1", Label: "Series 1", Color: "#123456"},
	}))
	for _, want := range []string{
		`hx-post="/charts/c1/toggle?series=s1"`,
		`hx-target="#chart-c1"`,
		`hx-swap="innerHTML"`,
		`style="cursor: pointer"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestLegendSvgItem_NoToggleWithoutChartID(t *testing.T) {
	// Toggle set, but no ChartID → htmx attributes must be omitted.
	out := render(t, legends.LegendSvgItem(legends.LegendSvgItemProps{
		Props: legends.LegendProps{Toggle: "post"},
		Datum: legends.Datum{ID: "s1", Label: "Series 1", Color: "#123456"},
	}))
	if strings.Contains(out, "hx-post") {
		t.Errorf("hx-post emitted without ChartID: %s", out)
	}
}

func TestLegendSvgItem_Opacity(t *testing.T) {
	tests := []struct {
		name    string
		props   legends.LegendProps
		datum   legends.Datum
		want    string
		wantNot string
	}{
		{"hidden dims to 0.4", legends.LegendProps{}, legends.Datum{Label: "x", Hidden: true}, `opacity="0.4"`, ""},
		{"sub-1 item opacity", legends.LegendProps{ItemOpacity: 0.5}, legends.Datum{Label: "x"}, `opacity="0.5"`, ""},
		{"full opacity omitted", legends.LegendProps{ItemOpacity: 1}, legends.Datum{Label: "x"}, "", "opacity="},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := render(t, legends.LegendSvgItem(legends.LegendSvgItemProps{Props: tc.props, Datum: tc.datum}))
			if tc.want != "" && !strings.Contains(out, tc.want) {
				t.Errorf("missing %s in %s", tc.want, out)
			}
			if tc.wantNot != "" && strings.Contains(out, tc.wantNot) {
				t.Errorf("unexpected %s in %s", tc.wantNot, out)
			}
		})
	}
}

func TestLegendSvgItem_LabelPositionByDirection(t *testing.T) {
	// SymbolSize 18, SymbolSpacing 6, ItemHeight 20 (defaults made explicit).
	base := legends.LegendProps{SymbolSize: 18, SymbolSpacing: 6, ItemHeight: 20}
	tests := []struct {
		name string
		dir  legends.LegendItemDirection
		want []string
	}{
		{"left-to-right", legends.LegendItemDirectionLeftToRight, []string{`x="24"`, `y="10"`, `text-anchor="start"`}},
		{"right-to-left", legends.LegendItemDirectionRightToLeft, []string{`x="-6"`, `y="10"`, `text-anchor="end"`}},
		{"top-to-bottom", legends.LegendItemDirectionTopToBottom, []string{`x="9"`, `y="19"`, `text-anchor="middle"`}},
		{"bottom-to-top", legends.LegendItemDirectionBottomToTop, []string{`x="9"`, `y="1"`, `text-anchor="middle"`}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := base
			p.ItemDirection = tc.dir
			out := render(t, legends.LegendSvgItem(legends.LegendSvgItemProps{Props: p, Datum: legends.Datum{Label: "L"}}))
			for _, want := range tc.want {
				if !strings.Contains(out, want) {
					t.Errorf("missing %s in %s", want, out)
				}
			}
		})
	}
}

func TestLegendSymbol_Shapes(t *testing.T) {
	d := legends.Datum{Color: "#abcdef"}
	tests := []struct {
		shape legends.SymbolShape
		want  []string
	}{
		{legends.SymbolShapeCircle, []string{"<circle", `cx="9"`, `cy="10"`, `r="9"`, `fill="#abcdef"`}},
		{legends.SymbolShapeSquare, []string{"<rect", `x="0"`, `y="1"`, `width="18"`, `height="18"`, `fill="#abcdef"`}},
		{legends.SymbolShapeDiamond, []string{"<polygon", `points="9,1 18,10 9,19 0,10"`, `fill="#abcdef"`}},
		{legends.SymbolShapeTriangle, []string{"<polygon", `points="18,19 0,19 9,1"`, `fill="#abcdef"`}},
	}
	for _, tc := range tests {
		t.Run(string(tc.shape), func(t *testing.T) {
			out := render(t, legends.LegendSymbol(legends.LegendProps{SymbolShape: tc.shape}, d))
			for _, want := range tc.want {
				if !strings.Contains(out, want) {
					t.Errorf("shape %s missing %s in %s", tc.shape, want, out)
				}
			}
		})
	}
}

func TestLegendSymbol_DefaultsToCircle(t *testing.T) {
	out := render(t, legends.LegendSymbol(legends.LegendProps{}, legends.Datum{Color: "red"}))
	if !strings.Contains(out, "<circle") {
		t.Errorf("empty shape should default to circle: %s", out)
	}
}

func TestLegendSymbol_Border(t *testing.T) {
	for _, shape := range []legends.SymbolShape{
		legends.SymbolShapeCircle,
		legends.SymbolShapeSquare,
		legends.SymbolShapeDiamond,
		legends.SymbolShapeTriangle,
	} {
		t.Run(string(shape)+" with border", func(t *testing.T) {
			out := render(t, legends.LegendSymbol(legends.LegendProps{
				SymbolShape:       shape,
				SymbolBorderWidth: 2,
				SymbolBorderColor: "#fff",
			}, legends.Datum{Color: "red"}))
			if !strings.Contains(out, `stroke="#fff"`) || !strings.Contains(out, `stroke-width="2"`) {
				t.Errorf("shape %s missing border attrs: %s", shape, out)
			}
		})
		t.Run(string(shape)+" without border", func(t *testing.T) {
			out := render(t, legends.LegendSymbol(legends.LegendProps{SymbolShape: shape}, legends.Datum{Color: "red"}))
			if strings.Contains(out, "stroke=") {
				t.Errorf("shape %s emitted stroke with zero border width: %s", shape, out)
			}
		})
	}
}

func grayscale(v float64) string {
	if v < 50 {
		return "#000000"
	}
	return "#ffffff"
}

func TestContinuousColorsLegendSvg_Basic(t *testing.T) {
	var b strings.Builder
	err := legends.ContinuousColorsLegendSvg(legends.ContinuousColorsLegendProps{
		Scale: grayscale,
		Min:   0, Max: 100,
		Width: 100, Height: 10,
		Samples: 2,
	}).Render(context.Background(), &b)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := b.String()
	for _, want := range []string{
		`<linearGradient id="nivo-continuous-legend"`,
		`<stop offset="0%" stop-color="#000000"`,
		`<stop offset="100%" stop-color="#ffffff"`,
		`<rect width="100" height="10" fill="url(#nivo-continuous-legend)"`,
		`>0</text>`,
		`>100</text>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
	if strings.Contains(out, `text-anchor="middle"`) {
		t.Errorf("title text rendered without a Title: %s", out)
	}
}

func TestContinuousColorsLegendSvg_Title(t *testing.T) {
	var b strings.Builder
	err := legends.ContinuousColorsLegendSvg(legends.ContinuousColorsLegendProps{
		Scale: grayscale,
		Min:   0, Max: 100,
		Samples: 2,
		Title:   "Heat",
	}).Render(context.Background(), &b)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := b.String()
	if !strings.Contains(out, ">Heat</text>") || !strings.Contains(out, `text-anchor="middle"`) {
		t.Errorf("title missing: %s", out)
	}
}

func TestContinuousColorsLegendSvg_DefaultSamplesAndSize(t *testing.T) {
	// Zero Samples/Width/Height fall back to defaults (32 samples, 100x10).
	var b strings.Builder
	err := legends.ContinuousColorsLegendSvg(legends.ContinuousColorsLegendProps{
		Scale: grayscale,
		Min:   0, Max: 100,
	}).Render(context.Background(), &b)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := b.String()
	if got := strings.Count(out, "<stop "); got != 32 {
		t.Errorf("stop count = %d, want 32 (default samples)", got)
	}
	if !strings.Contains(out, `<rect width="100" height="10"`) {
		t.Errorf("default bar size missing: %s", out)
	}
}

func TestContinuousColorsLegendSvg_AnchorPosition(t *testing.T) {
	// top-right anchor in a 300x200 chart with a 100-wide bar → x = 200.
	var b strings.Builder
	err := legends.ContinuousColorsLegendSvg(legends.ContinuousColorsLegendProps{
		Scale: grayscale,
		Min:   0, Max: 1,
		Anchor:     legends.LegendAnchorTopRight,
		ChartWidth: 300, ChartHeight: 200,
	}).Render(context.Background(), &b)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(b.String(), `transform="translate(200,0)"`) {
		t.Errorf("anchor translate wrong: %s", b.String())
	}
}

func TestLegends_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var b strings.Builder
	if err := legends.LegendSvg(legends.LegendProps{Items: twoItems()}).Render(ctx, &b); err == nil {
		t.Error("LegendSvg: expected error with canceled context")
	}
	if err := legends.BoxLegendSvg(legends.BoxLegendSvgProps{Props: legends.LegendProps{Items: twoItems()}}).Render(ctx, &b); err == nil {
		t.Error("BoxLegendSvg: expected error with canceled context")
	}
	if err := legends.LegendSvgItem(legends.LegendSvgItemProps{}).Render(ctx, &b); err == nil {
		t.Error("LegendSvgItem: expected error with canceled context")
	}
	if err := legends.LegendSymbol(legends.LegendProps{}, legends.Datum{}).Render(ctx, &b); err == nil {
		t.Error("LegendSymbol: expected error with canceled context")
	}
	if err := legends.ContinuousColorsLegendSvg(legends.ContinuousColorsLegendProps{Scale: grayscale}).Render(ctx, &b); err == nil {
		t.Error("ContinuousColorsLegendSvg: expected error with canceled context")
	}
}

// failingWriter fails every Write after the first `remaining` calls succeed.
type failingWriter struct{ remaining int }

func (w *failingWriter) Write(p []byte) (int, error) {
	if w.remaining <= 0 {
		return 0, errors.New("write failed")
	}
	w.remaining--
	return len(p), nil
}

// TestLegends_WriteErrorsPropagate verifies that a failing writer surfaces an
// error from Render no matter which write fails (walking the failure point
// through every write in the component).
func TestLegends_WriteErrorsPropagate(t *testing.T) {
	old := templruntime.DefaultBufferSize
	templruntime.DefaultBufferSize = 1
	defer func() { templruntime.DefaultBufferSize = old }()

	full := legends.LegendProps{
		Toggle: "post", ChartID: "c1",
		ItemOpacity:       0.5,
		SymbolBorderWidth: 2, SymbolBorderColor: "#fff",
		Items: []legends.Datum{
			{ID: "a", Label: "Alpha", Color: "#f00"},
			{ID: "b", Label: "Beta", Color: "#0f0", Hidden: true},
		},
	}
	shape := func(s legends.SymbolShape) legends.LegendProps {
		p := full
		p.SymbolShape = s
		return p
	}
	components := map[string]templ.Component{
		"BoxLegendSvg":    legends.BoxLegendSvg(legends.BoxLegendSvgProps{Props: full, ChartWidth: 200, ChartHeight: 100}),
		"circle symbol":   legends.LegendSymbol(shape(legends.SymbolShapeCircle), full.Items[0]),
		"square symbol":   legends.LegendSymbol(shape(legends.SymbolShapeSquare), full.Items[0]),
		"diamond symbol":  legends.LegendSymbol(shape(legends.SymbolShapeDiamond), full.Items[0]),
		"triangle symbol": legends.LegendSymbol(shape(legends.SymbolShapeTriangle), full.Items[0]),
		"continuous":      legends.ContinuousColorsLegendSvg(legends.ContinuousColorsLegendProps{Scale: grayscale, Min: 0, Max: 1, Samples: 2, Title: "T"}),
	}
	for name, c := range components {
		t.Run(name, func(t *testing.T) {
			for n := 0; n < 10000; n++ {
				buf := new(templruntime.Buffer)
				buf.Reset(&failingWriter{remaining: n})
				err := c.Render(context.Background(), buf)
				if err == nil {
					err = buf.Flush()
				}
				if err == nil {
					if n == 0 {
						t.Fatal("expected an error when every write fails")
					}
					return // failure point walked past the last write
				}
			}
			t.Fatal("render never succeeded")
		})
	}
}

func TestComputeDimensions_Defaults(t *testing.T) {
	// No items → zero box.
	if w, h := legends.ComputeDimensions(legends.LegendProps{}); w != 0 || h != 0 {
		t.Errorf("empty dims = (%v,%v), want (0,0)", w, h)
	}
	// Zero item sizes fall back to the defaults (100x20).
	w, h := legends.ComputeDimensions(legends.LegendProps{Items: twoItems()})
	if w != 100 || h != 40 {
		t.Errorf("default dims = (%v,%v), want (100,40)", w, h)
	}
}

func TestComputeContinuousColorsLegend_Clamps(t *testing.T) {
	// samples <= 0 falls back to the default of 32.
	if got := len(legends.ComputeContinuousColorsLegend(grayscale, 0, 1, 0)); got != 32 {
		t.Errorf("samples=0 stop count = %d, want 32", got)
	}
	// samples 1 clamps to 2; zero span is treated as 1.
	stops := legends.ComputeContinuousColorsLegend(grayscale, 5, 5, 1)
	if len(stops) != 2 {
		t.Fatalf("stop count = %d, want 2", len(stops))
	}
	if stops[0].Offset != 0 || stops[1].Offset != 1 {
		t.Errorf("offsets = %v,%v, want 0,1", stops[0].Offset, stops[1].Offset)
	}
}

func TestLegends_NonFiniteValuesRenderAsZero(t *testing.T) {
	out := render(t, legends.LegendSvgItem(legends.LegendSvgItemProps{
		Datum:  legends.Datum{Label: "x"},
		Layout: legends.ItemLayout{X: math.NaN(), Y: math.Inf(1)},
	}))
	if !strings.Contains(out, `transform="translate(0,0)"`) {
		t.Errorf("NaN/Inf should render as 0: %s", out)
	}
	// A tiny negative rounds to -0 and is normalized to 0.
	out = render(t, legends.LegendSvgItem(legends.LegendSvgItemProps{
		Datum:  legends.Datum{Label: "x"},
		Layout: legends.ItemLayout{X: -1e-12, Y: math.Inf(-1)},
	}))
	if !strings.Contains(out, `transform="translate(0,0)"`) {
		t.Errorf("-0 should normalize to 0: %s", out)
	}
}
