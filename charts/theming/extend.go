package theming

// textPropsWithInheritance lists the dot paths of nested text styles that
// inherit from the root theme.text. Mirrors @nivo/theming extend.ts.
var textPropsWithInheritance = []struct {
	get func(*Theme) *TextStyle
}{
	{func(t *Theme) *TextStyle { return &t.Axis.Ticks.Text }},
	{func(t *Theme) *TextStyle { return &t.Axis.Legend.Text }},
	{func(t *Theme) *TextStyle { return &t.Legends.Title.Text }},
	{func(t *Theme) *TextStyle { return &t.Legends.Text }},
	{func(t *Theme) *TextStyle { return &t.Legends.Ticks.Text }},
	{func(t *Theme) *TextStyle { return &t.Labels.Text }},
	{func(t *Theme) *TextStyle { return &t.Dots.Text }},
	{func(t *Theme) *TextStyle { return &t.Markers.Text }},
	{func(t *Theme) *TextStyle { return &t.Annotations.Text }},
}

// inheritRootThemeText merges a partial text style over the root text style.
// Mirrors nivo's inheritRootThemeText: {...rootStyle, ...partialStyle}.
func inheritRootThemeText(partial TextStyle, root TextStyle) TextStyle {
	out := root
	out.Extra = mergeExtra(root.Extra, partial.Extra)
	if partial.FontFamily != "" {
		out.FontFamily = partial.FontFamily
	}
	if partial.FontSize != nil {
		out.FontSize = partial.FontSize
	}
	if partial.Fill != "" {
		out.Fill = partial.Fill
	}
	if partial.OutlineWidth != 0 {
		out.OutlineWidth = partial.OutlineWidth
	}
	if partial.OutlineColor != "" {
		out.OutlineColor = partial.OutlineColor
	}
	if partial.OutlineOpacity != 0 {
		out.OutlineOpacity = partial.OutlineOpacity
	}
	return out
}

// ExtendDefaultTheme deep-merges defaultTheme with customTheme, then applies
// 9-path text inheritance so each nested text style inherits unset fields
// from the root theme.text. Mirrors @nivo/theming extendDefaultTheme.
func ExtendDefaultTheme(defaultTheme Theme, customTheme *PartialTheme) Theme {
	theme := defaultTheme
	if customTheme == nil {
		applyTextInheritance(&theme)
		return theme
	}
	if customTheme.Background != nil {
		theme.Background = *customTheme.Background
	}
	if customTheme.Text != nil {
		theme.Text = mergeTextStyle(theme.Text, *customTheme.Text)
	}
	mergeAxis(&theme.Axis, customTheme.Axis)
	mergeGrid(&theme.Grid, customTheme.Grid)
	mergeCrosshair(&theme.Crosshair, customTheme.Crosshair)
	mergeLegends(&theme.Legends, customTheme.Legends)
	mergeLabels(&theme.Labels, customTheme.Labels)
	mergeMarkers(&theme.Markers, customTheme.Markers)
	mergeDots(&theme.Dots, customTheme.Dots)
	mergeTooltip(&theme.Tooltip, customTheme.Tooltip)
	mergeAnnotations(&theme.Annotations, customTheme.Annotations)

	applyTextInheritance(&theme)
	return theme
}

func applyTextInheritance(t *Theme) {
	for _, p := range textPropsWithInheritance {
		*p.get(t) = inheritRootThemeText(*p.get(t), t.Text)
	}
}

// ExtendAxisTheme extends a complete axis theme with overrides. Mirrors nivo's
// extendAxisTheme. Returns axisTheme unchanged if overrides is nil.
func ExtendAxisTheme(axisTheme AxisTheme, overrides *PartialAxisTheme) AxisTheme {
	if overrides == nil {
		return axisTheme
	}
	mergeAxis(&axisTheme, overrides)
	return axisTheme
}

// --- merge helpers ---------------------------------------------------------

func mergeTextStyle(base, override TextStyle) TextStyle {
	out := base
	out.Extra = mergeExtra(base.Extra, override.Extra)
	if override.FontFamily != "" {
		out.FontFamily = override.FontFamily
	}
	if override.FontSize != nil {
		out.FontSize = override.FontSize
	}
	if override.Fill != "" {
		out.Fill = override.Fill
	}
	if override.OutlineWidth != 0 {
		out.OutlineWidth = override.OutlineWidth
	}
	if override.OutlineColor != "" {
		out.OutlineColor = override.OutlineColor
	}
	if override.OutlineOpacity != 0 {
		out.OutlineOpacity = override.OutlineOpacity
	}
	return out
}

func mergeExtra(base, override map[string]any) map[string]any {
	if len(base) == 0 && len(override) == 0 {
		return nil
	}
	out := map[string]any{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		out[k] = v
	}
	return out
}

func mergeAxis(base *AxisTheme, ov *PartialAxisTheme) {
	if ov == nil {
		return
	}
	if ov.Domain != nil && ov.Domain.Line != nil {
		base.Domain.Line.Extra = mergeExtra(base.Domain.Line.Extra, ov.Domain.Line.Extra)
	}
	if ov.Ticks != nil {
		if ov.Ticks.Line != nil {
			base.Ticks.Line.Extra = mergeExtra(base.Ticks.Line.Extra, ov.Ticks.Line.Extra)
		}
		if ov.Ticks.Text != nil {
			base.Ticks.Text = mergeTextStyle(base.Ticks.Text, *ov.Ticks.Text)
		}
	}
	if ov.Legend != nil && ov.Legend.Text != nil {
		base.Legend.Text = mergeTextStyle(base.Legend.Text, *ov.Legend.Text)
	}
}

func mergeGrid(base *GridTheme, ov *PartialGridTheme) {
	if ov != nil && ov.Line != nil {
		base.Line.Extra = mergeExtra(base.Line.Extra, ov.Line.Extra)
	}
}

func mergeCrosshair(base *CrosshairTheme, ov *PartialCrosshairTheme) {
	if ov != nil && ov.Line != nil {
		if ov.Line.Stroke != "" {
			base.Line.Stroke = ov.Line.Stroke
		}
		if ov.Line.StrokeWidth != 0 {
			base.Line.StrokeWidth = ov.Line.StrokeWidth
		}
		if ov.Line.StrokeOpacity != 0 {
			base.Line.StrokeOpacity = ov.Line.StrokeOpacity
		}
		if ov.Line.StrokeDasharray != "" {
			base.Line.StrokeDasharray = ov.Line.StrokeDasharray
		}
	}
}

func mergeLegends(base *LegendsTheme, ov *PartialLegendsTheme) {
	if ov == nil {
		return
	}
	if ov.Hidden != nil {
		if ov.Hidden.Symbol != nil {
			if ov.Hidden.Symbol.Fill != "" {
				base.Hidden.Symbol.Fill = ov.Hidden.Symbol.Fill
			}
			if ov.Hidden.Symbol.Opacity != 0 {
				base.Hidden.Symbol.Opacity = ov.Hidden.Symbol.Opacity
			}
		}
		if ov.Hidden.Text != nil {
			base.Hidden.Text = mergeTextStyle(base.Hidden.Text, *ov.Hidden.Text)
		}
	}
	if ov.Text != nil {
		base.Text = mergeTextStyle(base.Text, *ov.Text)
	}
	if ov.Title != nil && ov.Title.Text != nil {
		base.Title.Text = mergeTextStyle(base.Title.Text, *ov.Title.Text)
	}
	if ov.Ticks != nil {
		if ov.Ticks.Line != nil {
			base.Ticks.Line.Extra = mergeExtra(base.Ticks.Line.Extra, ov.Ticks.Line.Extra)
		}
		if ov.Ticks.Text != nil {
			base.Ticks.Text = mergeTextStyle(base.Ticks.Text, *ov.Ticks.Text)
		}
	}
}

func mergeLabels(base *LabelsTheme, ov *PartialLabelsTheme) {
	if ov != nil && ov.Text != nil {
		base.Text = mergeTextStyle(base.Text, *ov.Text)
	}
}

func mergeMarkers(base *MarkersTheme, ov *PartialMarkersTheme) {
	if ov == nil {
		return
	}
	if ov.LineColor != nil {
		base.LineColor = *ov.LineColor
	}
	if ov.LineStrokeWidth != nil {
		base.LineStrokeWidth = *ov.LineStrokeWidth
	}
	if ov.Text != nil {
		base.Text = mergeTextStyle(base.Text, *ov.Text)
	}
}

func mergeDots(base *DotsTheme, ov *PartialDotsTheme) {
	if ov != nil && ov.Text != nil {
		base.Text = mergeTextStyle(base.Text, *ov.Text)
	}
}

func mergeTooltip(base *TooltipTheme, ov *PartialTooltipTheme) {
	if ov == nil {
		return
	}
	base.Container = mergeExtra(base.Container, ov.Container)
	base.Basic = mergeExtra(base.Basic, ov.Basic)
	base.Chip = mergeExtra(base.Chip, ov.Chip)
	base.Table = mergeExtra(base.Table, ov.Table)
	base.TableCell = mergeExtra(base.TableCell, ov.TableCell)
	base.TableCellValue = mergeExtra(base.TableCellValue, ov.TableCellValue)
}

func mergeAnnotations(base *AnnotationsTheme, ov *PartialAnnotationsTheme) {
	if ov == nil {
		return
	}
	if ov.Text != nil {
		base.Text = mergeTextStyle(base.Text, *ov.Text)
	}
	if ov.Link != nil {
		base.Link = mergeAnnotationLink(base.Link, *ov.Link)
	}
	if ov.Outline != nil {
		base.Outline = mergeAnnotationOutline(base.Outline, *ov.Outline)
	}
	if ov.Symbol != nil {
		base.Symbol = mergeAnnotationSymbol(base.Symbol, *ov.Symbol)
	}
}

func mergeAnnotationLink(base, ov AnnotationLink) AnnotationLink {
	out := base
	out.Extra = mergeExtra(base.Extra, ov.Extra)
	if ov.Stroke != "" {
		out.Stroke = ov.Stroke
	}
	if ov.StrokeWidth != 0 {
		out.StrokeWidth = ov.StrokeWidth
	}
	if ov.OutlineWidth != 0 {
		out.OutlineWidth = ov.OutlineWidth
	}
	if ov.OutlineColor != "" {
		out.OutlineColor = ov.OutlineColor
	}
	if ov.OutlineOpacity != 0 {
		out.OutlineOpacity = ov.OutlineOpacity
	}
	return out
}

func mergeAnnotationOutline(base, ov AnnotationOutline) AnnotationOutline {
	out := base
	out.Extra = mergeExtra(base.Extra, ov.Extra)
	if ov.Stroke != "" {
		out.Stroke = ov.Stroke
	}
	if ov.StrokeWidth != 0 {
		out.StrokeWidth = ov.StrokeWidth
	}
	if ov.OutlineWidth != 0 {
		out.OutlineWidth = ov.OutlineWidth
	}
	if ov.OutlineColor != "" {
		out.OutlineColor = ov.OutlineColor
	}
	if ov.OutlineOpacity != 0 {
		out.OutlineOpacity = ov.OutlineOpacity
	}
	return out
}

func mergeAnnotationSymbol(base, ov AnnotationSymbol) AnnotationSymbol {
	out := base
	out.Extra = mergeExtra(base.Extra, ov.Extra)
	if ov.Fill != "" {
		out.Fill = ov.Fill
	}
	if ov.OutlineWidth != 0 {
		out.OutlineWidth = ov.OutlineWidth
	}
	if ov.OutlineColor != "" {
		out.OutlineColor = ov.OutlineColor
	}
	if ov.OutlineOpacity != 0 {
		out.OutlineOpacity = ov.OutlineOpacity
	}
	return out
}
