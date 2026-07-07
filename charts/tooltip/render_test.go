package tooltip_test

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
	"github.com/geoffjay/templ-charts/charts/tooltip"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return b.String()
}

func TestBasicTooltip_WithChip(t *testing.T) {
	out := render(t, tooltip.BasicTooltip(tooltip.BasicTooltipProps{
		ID:             "series-a",
		FormattedValue: "42",
		Color:          "#ff0000",
		EnableChip:     true,
	}))
	for _, want := range []string{
		`class="nivo-tooltip-basic"`,
		`class="nivo-tooltip-chip"`,
		"background: #ff0000",
		"series-a: 42",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestBasicTooltip_ChipOmitted(t *testing.T) {
	tests := []struct {
		name  string
		props tooltip.BasicTooltipProps
	}{
		{"chip disabled", tooltip.BasicTooltipProps{ID: "a", FormattedValue: "1", Color: "#f00", EnableChip: false}},
		{"no color", tooltip.BasicTooltipProps{ID: "a", FormattedValue: "1", EnableChip: true}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := render(t, tooltip.BasicTooltip(tc.props))
			if strings.Contains(out, "nivo-tooltip-chip") {
				t.Errorf("chip rendered unexpectedly: %s", out)
			}
		})
	}
}

func TestBasicTooltip_LabelFallbacks(t *testing.T) {
	tests := []struct {
		name  string
		props tooltip.BasicTooltipProps
		want  string
	}{
		{"id + formatted", tooltip.BasicTooltipProps{ID: "x", FormattedValue: "9"}, "x: 9"},
		{"formatted only", tooltip.BasicTooltipProps{FormattedValue: "9"}, ">9</div>"},
		{"value only", tooltip.BasicTooltipProps{Value: "raw"}, ">raw</div>"},
		{"id only", tooltip.BasicTooltipProps{ID: "just-id"}, ">just-id</div>"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := render(t, tooltip.BasicTooltip(tc.props))
			if !strings.Contains(out, tc.want) {
				t.Errorf("missing %s in %s", tc.want, out)
			}
		})
	}
}

func TestBasicTooltip_Styles(t *testing.T) {
	out := render(t, tooltip.BasicTooltip(tooltip.BasicTooltipProps{
		ID:             "s",
		FormattedValue: "1",
		Color:          "#0f0",
		EnableChip:     true,
		ContainerStyle: map[string]any{"backgroundColor": "white", "padding": "4px"},
		BasicStyle:     map[string]any{"whiteSpace": "pre"},
		ChipStyle:      map[string]any{"borderRadius": "2px"},
	}))
	for _, want := range []string{
		"background-color: white;",
		"padding: 4px;",
		"white-space: pre;",
		"border-radius: 2px;",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestChip(t *testing.T) {
	out := render(t, tooltip.Chip(tooltip.ChipProps{Color: "#123abc"}))
	for _, want := range []string{
		`class="nivo-tooltip-chip"`,
		"width: 12px",
		"height: 12px",
		"background: #123abc",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestTableTooltip_Rows(t *testing.T) {
	out := render(t, tooltip.TableTooltip(tooltip.TableTooltipProps{
		Rows: []tooltip.TableTooltipRow{
			{Label: "alpha", Value: "1", Color: "#f00"},
			{Label: "beta", Value: "2"},
		},
	}))
	if !strings.Contains(out, `class="nivo-tooltip-table"`) {
		t.Errorf("missing container class: %s", out)
	}
	if got := strings.Count(out, "<tr>"); got != 2 {
		t.Errorf("<tr> count = %d, want 2", got)
	}
	if got := strings.Count(out, "<td"); got != 4 {
		t.Errorf("<td> count = %d, want 4", got)
	}
	// Only the first row has a color → exactly one chip.
	if got := strings.Count(out, "nivo-tooltip-chip"); got != 1 {
		t.Errorf("chip count = %d, want 1", got)
	}
	for _, want := range []string{"alpha", ">1</td>", "beta", ">2</td>"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestTableTooltip_Empty(t *testing.T) {
	out := render(t, tooltip.TableTooltip(tooltip.TableTooltipProps{}))
	if !strings.Contains(out, "<table") || strings.Contains(out, "<tr>") {
		t.Errorf("empty table wrong: %s", out)
	}
}

func TestTableTooltip_Styles(t *testing.T) {
	out := render(t, tooltip.TableTooltip(tooltip.TableTooltipProps{
		Rows:           []tooltip.TableTooltipRow{{Label: "a", Value: "1"}},
		ContainerStyle: map[string]any{"background": "white"},
		TableStyle:     map[string]any{"borderCollapse": "collapse"},
		CellStyle:      map[string]any{"paddingRight": "8px"},
		ValueStyle:     map[string]any{"fontWeight": "bold"},
	}))
	for _, want := range []string{
		"background: white;",
		"border-collapse: collapse;",
		"padding-right: 8px;",
		"font-weight: bold;",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestCrosshair_CrossGeometry(t *testing.T) {
	out := render(t, tooltip.Crosshair(tooltip.CrosshairProps{
		Type:  tooltip.CrosshairTypeCross,
		Width: 100, Height: 50, X: 30, Y: 20,
	}))
	if !strings.Contains(out, `style="pointer-events: none"`) {
		t.Errorf("missing pointer-events group: %s", out)
	}
	if got := strings.Count(out, "<line"); got != 2 {
		t.Errorf("line count = %d, want 2", got)
	}
	for _, want := range []string{
		`x1="0"`, `y1="20"`, `x2="100"`, `y2="20"`, // horizontal
		`x1="30"`, `y1="0"`, `x2="30"`, `y2="50"`, // vertical
		`fill="none"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestCrosshair_LineCounts(t *testing.T) {
	tests := []struct {
		typ   tooltip.CrosshairType
		lines int
	}{
		{tooltip.CrosshairTypeTopLeft, 2},
		{tooltip.CrosshairTypeTopRight, 2},
		{tooltip.CrosshairTypeBottomLeft, 2},
		{tooltip.CrosshairTypeBottomRight, 2},
		{tooltip.CrosshairTypeMiddleLeft, 2},
		{tooltip.CrosshairTypeMiddleRight, 2},
		{tooltip.CrosshairTypeCross, 2},
		{tooltip.CrosshairTypeTop, 1},
		{tooltip.CrosshairTypeBottom, 1},
		{tooltip.CrosshairTypeMiddle, 1},
		{tooltip.CrosshairTypeLeft, 1},
		{tooltip.CrosshairTypeRight, 1},
		{tooltip.CrosshairType("bogus"), 0},
	}
	for _, tc := range tests {
		t.Run(string(tc.typ), func(t *testing.T) {
			out := render(t, tooltip.Crosshair(tooltip.CrosshairProps{
				Type:  tc.typ,
				Width: 100, Height: 50, X: 30, Y: 20,
			}))
			if got := strings.Count(out, "<line"); got != tc.lines {
				t.Errorf("%s line count = %d, want %d", tc.typ, got, tc.lines)
			}
		})
	}
}

func TestCrosshair_ThemeAttributes(t *testing.T) {
	out := render(t, tooltip.Crosshair(tooltip.CrosshairProps{
		Type:  tooltip.CrosshairTypeCross,
		Width: 100, Height: 50, X: 30, Y: 20,
		Theme: tooltip.CrosshairTheme{
			Stroke:          "#333333",
			StrokeWidth:     1.5,
			StrokeOpacity:   0.75,
			StrokeDasharray: "6 6",
		},
	}))
	for _, want := range []string{
		`stroke="#333333"`,
		`stroke-width="1.5"`,
		`stroke-opacity="0.75"`,
		`stroke-dasharray="6 6"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
	}
}

func TestCrosshair_ThemeAttributesOmitted(t *testing.T) {
	out := render(t, tooltip.Crosshair(tooltip.CrosshairProps{
		Type:  tooltip.CrosshairTypeCross,
		Width: 100, Height: 50, X: 30, Y: 20,
	}))
	for _, unwanted := range []string{"stroke=", "stroke-width=", "stroke-opacity=", "stroke-dasharray="} {
		if strings.Contains(out, unwanted) {
			t.Errorf("unexpected %s in %s", unwanted, out)
		}
	}
}

func TestTooltip_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var b strings.Builder
	if err := tooltip.BasicTooltip(tooltip.BasicTooltipProps{}).Render(ctx, &b); err == nil {
		t.Error("BasicTooltip: expected error with canceled context")
	}
	if err := tooltip.Chip(tooltip.ChipProps{}).Render(ctx, &b); err == nil {
		t.Error("Chip: expected error with canceled context")
	}
	if err := tooltip.TableTooltip(tooltip.TableTooltipProps{}).Render(ctx, &b); err == nil {
		t.Error("TableTooltip: expected error with canceled context")
	}
	if err := tooltip.Crosshair(tooltip.CrosshairProps{}).Render(ctx, &b); err == nil {
		t.Error("Crosshair: expected error with canceled context")
	}
}

func TestCrosshair_NonFiniteValuesRenderAsZero(t *testing.T) {
	// NaN/Inf coordinates render as 0; a tiny negative normalizes -0 to 0.
	out := render(t, tooltip.Crosshair(tooltip.CrosshairProps{
		Type:  tooltip.CrosshairTypeCross,
		Width: math.Inf(1), Height: math.NaN(),
		X: -1e-12, Y: math.Inf(-1),
	}))
	for _, want := range []string{`x1="0"`, `y1="0"`, `x2="0"`, `y2="0"`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in %s", want, out)
		}
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

// TestTooltip_WriteErrorsPropagate verifies that a failing writer surfaces an
// error from Render no matter which write fails (walking the failure point
// through every write in the component).
func TestTooltip_WriteErrorsPropagate(t *testing.T) {
	old := templruntime.DefaultBufferSize
	templruntime.DefaultBufferSize = 1
	defer func() { templruntime.DefaultBufferSize = old }()

	components := map[string]templ.Component{
		"BasicTooltip": tooltip.BasicTooltip(tooltip.BasicTooltipProps{
			ID: "a", FormattedValue: "1", Color: "#f00", EnableChip: true,
			ContainerStyle: map[string]any{"padding": "2px"},
			BasicStyle:     map[string]any{"color": "red"},
			ChipStyle:      map[string]any{"margin": "1px"},
		}),
		"Chip": tooltip.Chip(tooltip.ChipProps{Color: "#f00"}),
		"TableTooltip": tooltip.TableTooltip(tooltip.TableTooltipProps{
			Rows: []tooltip.TableTooltipRow{{Label: "a", Value: "1", Color: "#f00"}},
		}),
		"Crosshair": tooltip.Crosshair(tooltip.CrosshairProps{
			Type: tooltip.CrosshairTypeCross, Width: 100, Height: 50, X: 30, Y: 20,
			Theme: tooltip.CrosshairTheme{Stroke: "#333", StrokeWidth: 1, StrokeOpacity: 0.5, StrokeDasharray: "6 6"},
		}),
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
