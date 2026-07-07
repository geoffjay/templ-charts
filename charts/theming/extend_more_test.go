package theming

import "testing"

func sp(s string) *string   { return &s }
func fp(f float64) *float64 { return &f }

func TestExtendDefaultTheme_FullOverride(t *testing.T) {
	custom := &PartialTheme{
		Background: sp("#101010"),
		Text:       &TextStyle{FontFamily: "mono", FontSize: 14, Fill: "#eeeeee", OutlineWidth: 1, OutlineColor: "#000000", OutlineOpacity: 0.5, Extra: map[string]any{"letterSpacing": 1}},
		Axis: &PartialAxisTheme{
			Domain: &PartialAxisDomain{Line: &AxisDomainLine{Extra: map[string]any{"stroke": "#111111"}}},
			Ticks: &PartialAxisTicks{
				Line: &AxisTickLine{Extra: map[string]any{"stroke": "#222222"}},
				Text: &TextStyle{Fill: "#333333"},
			},
			Legend: &PartialAxisLegend{Text: &TextStyle{Fill: "#444444"}},
		},
		Grid:      &PartialGridTheme{Line: &GridLine{Extra: map[string]any{"stroke": "#555555"}}},
		Crosshair: &PartialCrosshairTheme{Line: &CrosshairLine{Stroke: "#666666", StrokeWidth: 2, StrokeOpacity: 0.8, StrokeDasharray: "4 4"}},
		Legends: &PartialLegendsTheme{
			Hidden: &PartialLegendsHidden{
				Symbol: &LegendsHiddenSymbol{Fill: "#777777", Opacity: 0.4},
				Text:   &TextStyle{Fill: "#888888"},
			},
			Text:  &TextStyle{Fill: "#999999"},
			Title: &PartialLegendsTitle{Text: &TextStyle{Fill: "#aaaaaa"}},
			Ticks: &PartialLegendsTicks{
				Line: &AxisTickLine{Extra: map[string]any{"stroke": "#bbbbbb"}},
				Text: &TextStyle{Fill: "#cccccc"},
			},
		},
		Labels:  &PartialLabelsTheme{Text: &TextStyle{Fill: "#dddddd"}},
		Markers: &PartialMarkersTheme{LineColor: sp("#eeeeee"), LineStrokeWidth: fp(3), Text: &TextStyle{Fill: "#f0f0f0"}},
		Dots:    &PartialDotsTheme{Text: &TextStyle{Fill: "#0a0a0a"}},
		Tooltip: &PartialTooltipTheme{
			Container:      map[string]any{"background": "#111"},
			Basic:          map[string]any{"whiteSpace": "pre"},
			Chip:           map[string]any{"width": 12},
			Table:          map[string]any{"width": "100%"},
			TableCell:      map[string]any{"padding": 3},
			TableCellValue: map[string]any{"fontWeight": "bold"},
		},
		Annotations: &PartialAnnotationsTheme{
			Text:    &TextStyle{Fill: "#0b0b0b"},
			Link:    &AnnotationLink{Stroke: "#0c0c0c", StrokeWidth: 2, OutlineWidth: 3, OutlineColor: "#0d0d0d", OutlineOpacity: 0.9, Extra: map[string]any{"x": 1}},
			Outline: &AnnotationOutline{Stroke: "#0e0e0e", StrokeWidth: 4, OutlineWidth: 5, OutlineColor: "#0f0f0f", OutlineOpacity: 0.7, Extra: map[string]any{"y": 2}},
			Symbol:  &AnnotationSymbol{Fill: "#1a1a1a", OutlineWidth: 6, OutlineColor: "#1b1b1b", OutlineOpacity: 0.6, Extra: map[string]any{"z": 3}},
		},
	}
	got := ExtendDefaultTheme(DefaultTheme, custom)

	if got.Background != "#101010" {
		t.Errorf("background = %q", got.Background)
	}
	if got.Text.FontFamily != "mono" || got.Text.Fill != "#eeeeee" || got.Text.OutlineWidth != 1 ||
		got.Text.OutlineColor != "#000000" || got.Text.OutlineOpacity != 0.5 {
		t.Errorf("text = %+v", got.Text)
	}
	if got.Text.Extra["letterSpacing"] != 1 {
		t.Errorf("text extra = %v", got.Text.Extra)
	}
	if got.Axis.Domain.Line.Extra["stroke"] != "#111111" {
		t.Errorf("axis domain line = %v", got.Axis.Domain.Line.Extra)
	}
	if got.Axis.Ticks.Line.Extra["stroke"] != "#222222" {
		t.Errorf("axis ticks line = %v", got.Axis.Ticks.Line.Extra)
	}
	if got.Axis.Ticks.Text.Fill != "#333333" {
		t.Errorf("axis ticks text fill = %q", got.Axis.Ticks.Text.Fill)
	}
	// Inheritance: unset ticks-text fields come from the merged root text.
	if got.Axis.Ticks.Text.FontFamily != "mono" {
		t.Errorf("axis ticks text fontFamily = %q, want inherited mono", got.Axis.Ticks.Text.FontFamily)
	}
	if got.Axis.Legend.Text.Fill != "#444444" {
		t.Errorf("axis legend text = %q", got.Axis.Legend.Text.Fill)
	}
	if got.Grid.Line.Extra["stroke"] != "#555555" {
		t.Errorf("grid = %v", got.Grid.Line.Extra)
	}
	ch := got.Crosshair.Line
	if ch.Stroke != "#666666" || ch.StrokeWidth != 2 || ch.StrokeOpacity != 0.8 || ch.StrokeDasharray != "4 4" {
		t.Errorf("crosshair = %+v", ch)
	}
	lg := got.Legends
	if lg.Hidden.Symbol.Fill != "#777777" || lg.Hidden.Symbol.Opacity != 0.4 {
		t.Errorf("legends hidden symbol = %+v", lg.Hidden.Symbol)
	}
	if lg.Hidden.Text.Fill != "#888888" || lg.Text.Fill != "#999999" || lg.Title.Text.Fill != "#aaaaaa" {
		t.Errorf("legends text/title = %+v", lg)
	}
	if lg.Ticks.Line.Extra["stroke"] != "#bbbbbb" || lg.Ticks.Text.Fill != "#cccccc" {
		t.Errorf("legends ticks = %+v", lg.Ticks)
	}
	if got.Labels.Text.Fill != "#dddddd" {
		t.Errorf("labels = %+v", got.Labels.Text)
	}
	if got.Markers.LineColor != "#eeeeee" || got.Markers.LineStrokeWidth != 3 || got.Markers.Text.Fill != "#f0f0f0" {
		t.Errorf("markers = %+v", got.Markers)
	}
	if got.Dots.Text.Fill != "#0a0a0a" {
		t.Errorf("dots = %+v", got.Dots.Text)
	}
	if got.Tooltip.Container["background"] != "#111" || got.Tooltip.Chip["width"] != 12 ||
		got.Tooltip.Basic["whiteSpace"] != "pre" || got.Tooltip.Table["width"] != "100%" ||
		got.Tooltip.TableCell["padding"] != 3 || got.Tooltip.TableCellValue["fontWeight"] != "bold" {
		t.Errorf("tooltip = %+v", got.Tooltip)
	}
	an := got.Annotations
	if an.Text.Fill != "#0b0b0b" {
		t.Errorf("annotations text = %+v", an.Text)
	}
	if an.Link.Stroke != "#0c0c0c" || an.Link.StrokeWidth != 2 || an.Link.OutlineWidth != 3 ||
		an.Link.OutlineColor != "#0d0d0d" || an.Link.OutlineOpacity != 0.9 || an.Link.Extra["x"] != 1 {
		t.Errorf("annotations link = %+v", an.Link)
	}
	if an.Outline.Stroke != "#0e0e0e" || an.Outline.StrokeWidth != 4 || an.Outline.OutlineWidth != 5 ||
		an.Outline.OutlineColor != "#0f0f0f" || an.Outline.OutlineOpacity != 0.7 || an.Outline.Extra["y"] != 2 {
		t.Errorf("annotations outline = %+v", an.Outline)
	}
	if an.Symbol.Fill != "#1a1a1a" || an.Symbol.OutlineWidth != 6 || an.Symbol.OutlineColor != "#1b1b1b" ||
		an.Symbol.OutlineOpacity != 0.6 || an.Symbol.Extra["z"] != 3 {
		t.Errorf("annotations symbol = %+v", an.Symbol)
	}
}

func TestExtendDefaultTheme_EmptyPartialKeepsDefaults(t *testing.T) {
	// A non-nil but empty PartialTheme must leave everything at the default.
	got := ExtendDefaultTheme(DefaultTheme, &PartialTheme{})
	want := ExtendDefaultTheme(DefaultTheme, nil)
	if got.Background != want.Background || got.Text.Fill != want.Text.Fill {
		t.Errorf("empty partial changed the theme: %+v vs %+v", got.Text, want.Text)
	}
	if got.Crosshair.Line.Stroke != want.Crosshair.Line.Stroke {
		t.Errorf("crosshair changed: %+v", got.Crosshair)
	}
}

func TestExtendAxisTheme(t *testing.T) {
	base := DefaultTheme.Axis
	// Nil overrides return the input unchanged.
	if got := ExtendAxisTheme(base, nil); got.Ticks.Text.Fill != base.Ticks.Text.Fill {
		t.Errorf("nil overrides changed the axis theme")
	}
	got := ExtendAxisTheme(base, &PartialAxisTheme{
		Ticks: &PartialAxisTicks{Text: &TextStyle{Fill: "#ff00ff"}},
	})
	if got.Ticks.Text.Fill != "#ff00ff" {
		t.Errorf("ticks text fill = %q, want #ff00ff", got.Ticks.Text.Fill)
	}
	// Untouched parts keep the base values.
	if got.Legend.Text.Fill != base.Legend.Text.Fill {
		t.Errorf("legend text changed: %q", got.Legend.Text.Fill)
	}
}

func TestMergeExtra(t *testing.T) {
	if got := mergeExtra(nil, nil); got != nil {
		t.Errorf("mergeExtra(nil, nil) = %v, want nil", got)
	}
	got := mergeExtra(map[string]any{"a": 1, "b": 2}, map[string]any{"b": 3, "c": 4})
	if got["a"] != 1 || got["b"] != 3 || got["c"] != 4 {
		t.Errorf("mergeExtra = %v", got)
	}
	// Base-only and override-only both produce merged output.
	if got := mergeExtra(map[string]any{"a": 1}, nil); got["a"] != 1 {
		t.Errorf("base-only mergeExtra = %v", got)
	}
	if got := mergeExtra(nil, map[string]any{"z": 9}); got["z"] != 9 {
		t.Errorf("override-only mergeExtra = %v", got)
	}
}

func TestMergeTextStyle_ZeroOverrideKeepsBase(t *testing.T) {
	base := TextStyle{FontFamily: "sans", FontSize: 11, Fill: "#333333", OutlineWidth: 2, OutlineColor: "#fff", OutlineOpacity: 1}
	got := mergeTextStyle(base, TextStyle{})
	if got.FontFamily != base.FontFamily || got.FontSize != base.FontSize || got.Fill != base.Fill ||
		got.OutlineWidth != base.OutlineWidth || got.OutlineColor != base.OutlineColor || got.OutlineOpacity != base.OutlineOpacity {
		t.Errorf("zero override should keep base: %+v", got)
	}
	got = mergeTextStyle(base, TextStyle{Fill: "#ff0000", FontSize: 20})
	if got.Fill != "#ff0000" || got.FontSize != 20 || got.FontFamily != "sans" {
		t.Errorf("partial override = %+v", got)
	}
}

func TestInheritRootThemeText(t *testing.T) {
	root := TextStyle{FontFamily: "sans", FontSize: 11, Fill: "#333333"}
	// Empty partial: inherits everything.
	got := inheritRootThemeText(TextStyle{}, root)
	if got.FontFamily != "sans" || got.Fill != "#333333" {
		t.Errorf("empty partial inherit = %+v", got)
	}
	// Partial fill wins, family inherited.
	got = inheritRootThemeText(TextStyle{Fill: "#ff0000", OutlineWidth: 3, OutlineColor: "#000", OutlineOpacity: 0.5}, root)
	if got.Fill != "#ff0000" || got.FontFamily != "sans" || got.OutlineWidth != 3 || got.OutlineColor != "#000" || got.OutlineOpacity != 0.5 {
		t.Errorf("partial inherit = %+v", got)
	}
	// FontSize and FontFamily overrides win; extras are merged.
	got = inheritRootThemeText(TextStyle{FontFamily: "serif", FontSize: 20, Extra: map[string]any{"k": "v"}}, root)
	if got.FontFamily != "serif" || got.FontSize != 20 || got.Extra["k"] != "v" || got.Fill != "#333333" {
		t.Errorf("font override inherit = %+v", got)
	}
}
