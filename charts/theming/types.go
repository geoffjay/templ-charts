// Package theming mirrors @nivo/theming: the Theme struct, defaultTheme
// (verbatim from defaults.ts), ExtendDefaultTheme (deep merge + 9-path text
// inheritance), ExtendAxisTheme, borderRadius helpers, the engine bridge
// (svgStyleAttributesMapping/convertStyleAttribute), and SanitizeSvgTextStyle
// / SanitizeHtmlTextStyle.
package theming

// TextStyle mirrors @nivo/theming TextStyle. Fill/FontFamily/FontSize etc.
// are the typed core; Extra carries any additional CSS properties.
type TextStyle struct {
	// FontFamily is the CSS font-family applied to the text (default "sans-serif").
	FontFamily string
	FontSize   any // number (px) or string (e.g. "11px")
	// Fill is the text color (default "#333333").
	Fill string
	// OutlineWidth is the width, in pixels, of the halo/outline stroked behind
	// the glyphs for legibility; 0 disables the outline (default 0).
	OutlineWidth float64
	// OutlineColor is the color of the text outline/halo (default "#ffffff").
	OutlineColor string
	// OutlineOpacity is the opacity of the text outline/halo, 0–1 (default 1).
	OutlineOpacity float64
	Extra          map[string]any // additional CSS properties (fontWeight, fontStyle, …)
}

// AxisDomainLine is the axis domain line style.
type AxisDomainLine struct {
	// Extra holds CSS/SVG style attributes for the axis domain (baseline) line,
	// e.g. stroke and strokeWidth (default stroke "transparent", strokeWidth 1).
	Extra map[string]any
}

// AxisTicks is the axis ticks style (line + text).
type AxisTicks struct {
	// Line is the style of the axis tick marks.
	Line AxisTickLine
	// Text is the style of the axis tick labels.
	Text TextStyle
}

// AxisTickLine is the axis tick line style.
type AxisTickLine struct {
	// Extra holds CSS/SVG style attributes for a tick line, e.g. stroke and
	// strokeWidth (default stroke "#777777", strokeWidth 1).
	Extra map[string]any
}

// AxisLegend is the axis legend style.
type AxisLegend struct {
	// Text is the style of the axis legend (title) text (default fontSize 12).
	Text TextStyle
}

// AxisTheme is the per-axis theme block.
type AxisTheme struct {
	// Domain styles the axis domain (baseline) line.
	Domain AxisDomain
	// Ticks styles the axis tick marks and their labels.
	Ticks AxisTicks
	// Legend styles the axis legend/title.
	Legend AxisLegend
}

// AxisDomain holds the axis domain line.
type AxisDomain struct {
	// Line is the axis domain (baseline) line style.
	Line AxisDomainLine
}

// GridTheme is the grid line style.
type GridTheme struct {
	// Line is the style of the background grid lines.
	Line GridLine
}

// GridLine is the grid line style.
type GridLine struct {
	// Extra holds CSS/SVG style attributes for the grid lines, e.g. stroke and
	// strokeWidth (default stroke "#dddddd", strokeWidth 1).
	Extra map[string]any
}

// CrosshairLine is the crosshair line style.
type CrosshairLine struct {
	// Stroke is the color of the crosshair line (default "#000000").
	Stroke string
	// StrokeWidth is the width, in pixels, of the crosshair line (default 1).
	StrokeWidth float64
	// StrokeOpacity is the opacity of the crosshair line, 0–1 (default 0.75).
	StrokeOpacity float64
	// StrokeDasharray is the SVG dash pattern for the crosshair line (default "6 6").
	StrokeDasharray string
}

// CrosshairTheme is the crosshair theme block.
type CrosshairTheme struct {
	// Line is the crosshair line style.
	Line CrosshairLine
}

// LegendsHidden is the legend hidden-item style.
type LegendsHidden struct {
	// Symbol is the style of the symbol shown for a hidden/toggled-off legend item.
	Symbol LegendsHiddenSymbol
	// Text is the style of the label for a hidden/toggled-off legend item.
	Text TextStyle
}

// LegendsHiddenSymbol is the hidden-legend symbol style.
type LegendsHiddenSymbol struct {
	// Fill is the fill color of the hidden-item symbol (default "#333333").
	Fill string
	// Opacity is the opacity of the hidden-item symbol, 0–1 (default 0.6).
	Opacity float64
}

// LegendsTicks is the legend ticks style.
type LegendsTicks struct {
	// Line is the style of legend tick lines (e.g. on continuous-color legends).
	Line AxisTickLine
	// Text is the style of legend tick labels (default fontSize 10).
	Text TextStyle
}

// LegendsTitle is the legend title style.
type LegendsTitle struct {
	// Text is the style of the legend title text.
	Text TextStyle
}

// LegendsTheme is the legends theme block.
type LegendsTheme struct {
	// Hidden styles hidden/toggled-off legend items.
	Hidden LegendsHidden
	// Text is the default style of legend item labels.
	Text TextStyle
	// Title styles the legend title.
	Title LegendsTitle
	// Ticks styles legend tick lines and labels.
	Ticks LegendsTicks
}

// LabelsTheme is the labels theme block.
type LabelsTheme struct {
	// Text is the style of data labels drawn on the chart.
	Text TextStyle
}

// MarkersTheme is the markers theme block.
type MarkersTheme struct {
	// LineColor is the color of reference-marker lines (default "#000000").
	LineColor string
	// LineStrokeWidth is the width, in pixels, of reference-marker lines (default 1).
	LineStrokeWidth float64
	// Text is the style of reference-marker labels.
	Text TextStyle
}

// DotsTheme is the dots theme block.
type DotsTheme struct {
	// Text is the style of dot/point labels.
	Text TextStyle
}

// TooltipTheme is the tooltip theme block (CSS-style values for HTML).
type TooltipTheme struct {
	// Container holds CSS style for the outer tooltip container (background,
	// padding, border radius, box shadow, …).
	Container map[string]any
	// Basic holds CSS style for the basic single-line tooltip layout.
	Basic map[string]any
	// Chip holds CSS style for the color chip/swatch shown in tooltips.
	Chip map[string]any
	// Table holds CSS style for the tooltip table wrapper.
	Table map[string]any
	// TableCell holds CSS style for a tooltip table cell.
	TableCell map[string]any
	// TableCellValue holds CSS style for the value cell in a tooltip table
	// (default fontWeight "bold").
	TableCellValue map[string]any
}

// AnnotationLink is the annotation link style.
type AnnotationLink struct {
	// Stroke is the color of the annotation link/connector line (default "#000000").
	Stroke string
	// StrokeWidth is the width, in pixels, of the link line (default 1).
	StrokeWidth float64
	// OutlineWidth is the width of the halo/outline behind the link line (default 2).
	OutlineWidth float64
	// OutlineColor is the color of the link outline/halo (default "#ffffff").
	OutlineColor string
	// OutlineOpacity is the opacity of the link outline/halo, 0–1 (default 1).
	OutlineOpacity float64
	// Extra holds additional CSS/SVG style attributes for the link.
	Extra map[string]any
}

// AnnotationOutline is the annotation outline style.
type AnnotationOutline struct {
	// Stroke is the color of the outline stroked around the annotated element
	// (default "#000000").
	Stroke string
	// StrokeWidth is the width, in pixels, of the outline stroke (default 2).
	StrokeWidth float64
	// OutlineWidth is the width of the halo behind the outline (default 2).
	OutlineWidth float64
	// OutlineColor is the color of the halo behind the outline (default "#ffffff").
	OutlineColor string
	// OutlineOpacity is the opacity of the halo behind the outline, 0–1 (default 1).
	OutlineOpacity float64
	// Extra holds additional CSS/SVG style attributes for the outline
	// (default fill "none").
	Extra map[string]any
}

// AnnotationSymbol is the annotation symbol style.
type AnnotationSymbol struct {
	// Fill is the fill color of the annotation symbol (default "#000000").
	Fill string
	// OutlineWidth is the width of the halo behind the symbol (default 2).
	OutlineWidth float64
	// OutlineColor is the color of the symbol halo (default "#ffffff").
	OutlineColor string
	// OutlineOpacity is the opacity of the symbol halo, 0–1 (default 1).
	OutlineOpacity float64
	// Extra holds additional CSS/SVG style attributes for the symbol.
	Extra map[string]any
}

// AnnotationsTheme is the annotations theme block.
type AnnotationsTheme struct {
	// Text is the style of annotation note text.
	Text TextStyle
	// Link styles the annotation connector line.
	Link AnnotationLink
	// Outline styles the outline drawn around the annotated element.
	Outline AnnotationOutline
	// Symbol styles the annotation symbol/dot.
	Symbol AnnotationSymbol
}

// Theme is the full theme after inheritance has been applied. Mirrors
// @nivo/theming Theme.
type Theme struct {
	// Background is the chart background color (default "transparent").
	Background string
	// Text is the base/global text style; nested text blocks that do not
	// override a property inherit it from here.
	Text TextStyle
	// Axis styles the axes (domain, ticks, legends).
	Axis AxisTheme
	// Grid styles the background grid lines.
	Grid GridTheme
	// Crosshair styles the crosshair.
	Crosshair CrosshairTheme
	// Legends styles the legends.
	Legends LegendsTheme
	// Labels styles data labels drawn on the chart.
	Labels LabelsTheme
	// Markers styles reference markers.
	Markers MarkersTheme
	// Dots styles dot/point labels.
	Dots DotsTheme
	// Tooltip styles the HTML tooltip.
	Tooltip TooltipTheme
	// Annotations styles annotations.
	Annotations AnnotationsTheme
}

// ThemeWithoutInheritance is the input form of a theme (nested text styles
// may omit inherited root text properties). Mirrors ThemeWithoutInheritance.
type ThemeWithoutInheritance = Theme

// PartialTheme is a theme where every field is optional. Used as input to
// ExtendDefaultTheme. Modeled as a pointer-bearing Theme so callers can
// express "unset" via nil pointers. In every Partial* struct a nil field
// leaves the corresponding default in place.
type PartialTheme struct {
	// Background overrides Theme.Background when non-nil.
	Background *string
	// Text overrides the base/global text style when non-nil.
	Text *TextStyle
	// Axis overrides axis styling when non-nil.
	Axis *PartialAxisTheme
	// Grid overrides grid-line styling when non-nil.
	Grid *PartialGridTheme
	// Crosshair overrides crosshair styling when non-nil.
	Crosshair *PartialCrosshairTheme
	// Legends overrides legend styling when non-nil.
	Legends *PartialLegendsTheme
	// Labels overrides data-label styling when non-nil.
	Labels *PartialLabelsTheme
	// Markers overrides reference-marker styling when non-nil.
	Markers *PartialMarkersTheme
	// Dots overrides dot/point-label styling when non-nil.
	Dots *PartialDotsTheme
	// Tooltip overrides tooltip styling when non-nil.
	Tooltip *PartialTooltipTheme
	// Annotations overrides annotation styling when non-nil.
	Annotations *PartialAnnotationsTheme
}

// PartialAxisTheme is the partial form of AxisTheme.
type PartialAxisTheme struct {
	// Domain overrides the axis domain-line style when non-nil.
	Domain *PartialAxisDomain
	// Ticks overrides the axis tick style when non-nil.
	Ticks *PartialAxisTicks
	// Legend overrides the axis legend style when non-nil.
	Legend *PartialAxisLegend
}

// PartialAxisDomain is the partial form of AxisDomain.
type PartialAxisDomain struct {
	// Line overrides the axis domain (baseline) line style when non-nil.
	Line *AxisDomainLine
}

// PartialAxisTicks is the partial form of AxisTicks.
type PartialAxisTicks struct {
	// Line overrides the axis tick-mark style when non-nil.
	Line *AxisTickLine
	// Text overrides the axis tick-label style when non-nil.
	Text *TextStyle
}

// PartialAxisLegend is the partial form of AxisLegend.
type PartialAxisLegend struct {
	// Text overrides the axis legend (title) text style when non-nil.
	Text *TextStyle
}

// PartialGridTheme is the partial form of GridTheme.
type PartialGridTheme struct {
	// Line overrides the grid-line style when non-nil.
	Line *GridLine
}

// PartialCrosshairTheme is the partial form of CrosshairTheme.
type PartialCrosshairTheme struct {
	// Line overrides the crosshair line style when non-nil.
	Line *CrosshairLine
}

// PartialLegendsTheme is the partial form of LegendsTheme.
type PartialLegendsTheme struct {
	// Hidden overrides styling of hidden/toggled-off legend items when non-nil.
	Hidden *PartialLegendsHidden
	// Text overrides the default legend item label style when non-nil.
	Text *TextStyle
	// Title overrides the legend title style when non-nil.
	Title *PartialLegendsTitle
	// Ticks overrides legend tick styling when non-nil.
	Ticks *PartialLegendsTicks
}

// PartialLegendsHidden is the partial form of LegendsHidden.
type PartialLegendsHidden struct {
	// Symbol overrides the hidden-item symbol style when non-nil.
	Symbol *LegendsHiddenSymbol
	// Text overrides the hidden-item label style when non-nil.
	Text *TextStyle
}

// PartialLegendsTitle is the partial form of LegendsTitle.
type PartialLegendsTitle struct {
	// Text overrides the legend title text style when non-nil.
	Text *TextStyle
}

// PartialLegendsTicks is the partial form of LegendsTicks.
type PartialLegendsTicks struct {
	// Line overrides the legend tick-line style when non-nil.
	Line *AxisTickLine
	// Text overrides the legend tick-label style when non-nil.
	Text *TextStyle
}

// PartialLabelsTheme is the partial form of LabelsTheme.
type PartialLabelsTheme struct {
	// Text overrides the data-label text style when non-nil.
	Text *TextStyle
}

// PartialMarkersTheme is the partial form of MarkersTheme.
type PartialMarkersTheme struct {
	// LineColor overrides the reference-marker line color when non-nil.
	LineColor *string
	// LineStrokeWidth overrides the reference-marker line width when non-nil.
	LineStrokeWidth *float64
	// Text overrides the reference-marker label style when non-nil.
	Text *TextStyle
}

// PartialDotsTheme is the partial form of DotsTheme.
type PartialDotsTheme struct {
	// Text overrides the dot/point-label text style when non-nil.
	Text *TextStyle
}

// PartialTooltipTheme is the partial form of TooltipTheme.
type PartialTooltipTheme struct {
	// Container overrides CSS style for the outer tooltip container when non-nil.
	Container map[string]any
	// Basic overrides CSS style for the basic single-line tooltip layout when non-nil.
	Basic map[string]any
	// Chip overrides CSS style for the tooltip color chip/swatch when non-nil.
	Chip map[string]any
	// Table overrides CSS style for the tooltip table wrapper when non-nil.
	Table map[string]any
	// TableCell overrides CSS style for a tooltip table cell when non-nil.
	TableCell map[string]any
	// TableCellValue overrides CSS style for the tooltip table value cell when non-nil.
	TableCellValue map[string]any
}

// PartialAnnotationsTheme is the partial form of AnnotationsTheme.
type PartialAnnotationsTheme struct {
	// Text overrides annotation note text style when non-nil.
	Text *TextStyle
	// Link overrides annotation connector-line style when non-nil.
	Link *AnnotationLink
	// Outline overrides annotation outline style when non-nil.
	Outline *AnnotationOutline
	// Symbol overrides annotation symbol style when non-nil.
	Symbol *AnnotationSymbol
}
